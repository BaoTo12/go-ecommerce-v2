package scheduler

import (
	"context"
	"errors"
	"sync"
	"time"
)

/*
JOB SCHEDULER - Distributed Task Scheduling

Implements:
- Cron-like scheduling
- One-time jobs
- Recurring jobs
- Job dependencies
- Retry with backoff
- Distributed locking
- Job queue priorities
*/

// JobStatus represents job status
type JobStatus string

const (
	JobStatusPending   JobStatus = "PENDING"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
	JobStatusRetrying  JobStatus = "RETRYING"
)

// Job represents a scheduled job
type Job struct {
	ID          string
	Name        string
	Handler     string // Handler name
	Payload     []byte
	Schedule    string // Cron expression or interval
	NextRun     time.Time
	LastRun     *time.Time
	Status      JobStatus
	Retries     int
	MaxRetries  int
	RetryDelay  time.Duration
	Priority    int
	Dependencies []string
	Timeout     time.Duration
	Metadata    map[string]string
}

// JobHandler executes a job
type JobHandler func(ctx context.Context, payload []byte) error

// JobStore stores jobs
type JobStore interface {
	Save(ctx context.Context, job *Job) error
	Get(ctx context.Context, id string) (*Job, error)
	GetPending(ctx context.Context, limit int) ([]*Job, error)
	UpdateStatus(ctx context.Context, id string, status JobStatus) error
	Delete(ctx context.Context, id string) error
}

// DistributedLock for distributed locking
type DistributedLock interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string) error
	Extend(ctx context.Context, key string, ttl time.Duration) error
}

// Scheduler manages job scheduling
type Scheduler struct {
	store      JobStore
	handlers   map[string]JobHandler
	lock       DistributedLock
	workerPool int
	pollInterval time.Duration
	
	running bool
	stopCh  chan struct{}
	wg      sync.WaitGroup
	mu      sync.RWMutex
}

// SchedulerConfig configures the scheduler
type SchedulerConfig struct {
	WorkerPool   int
	PollInterval time.Duration
}

// DefaultSchedulerConfig returns default config
func DefaultSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		WorkerPool:   10,
		PollInterval: 1 * time.Second,
	}
}

// NewScheduler creates a job scheduler
func NewScheduler(store JobStore, lock DistributedLock, config SchedulerConfig) *Scheduler {
	return &Scheduler{
		store:        store,
		handlers:     make(map[string]JobHandler),
		lock:         lock,
		workerPool:   config.WorkerPool,
		pollInterval: config.PollInterval,
		stopCh:       make(chan struct{}),
	}
}

// Register registers a job handler
func (s *Scheduler) Register(name string, handler JobHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[name] = handler
}

// Schedule schedules a new job
func (s *Scheduler) Schedule(ctx context.Context, job *Job) error {
	if job.ID == "" {
		job.ID = generateJobID()
	}
	if job.Status == "" {
		job.Status = JobStatusPending
	}
	if job.NextRun.IsZero() {
		job.NextRun = time.Now()
	}
	return s.store.Save(ctx, job)
}

// ScheduleIn schedules a job to run after a delay
func (s *Scheduler) ScheduleIn(ctx context.Context, name string, payload []byte, delay time.Duration) error {
	job := &Job{
		ID:       generateJobID(),
		Name:     name,
		Handler:  name,
		Payload:  payload,
		NextRun:  time.Now().Add(delay),
		Status:   JobStatusPending,
	}
	return s.store.Save(ctx, job)
}

// ScheduleRecurring schedules a recurring job
func (s *Scheduler) ScheduleRecurring(ctx context.Context, name string, payload []byte, interval time.Duration) error {
	job := &Job{
		ID:       generateJobID(),
		Name:     name,
		Handler:  name,
		Payload:  payload,
		Schedule: interval.String(),
		NextRun:  time.Now(),
		Status:   JobStatusPending,
	}
	return s.store.Save(ctx, job)
}

// Start starts the scheduler
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return errors.New("scheduler already running")
	}
	s.running = true
	s.mu.Unlock()

	// Start worker pool
	jobs := make(chan *Job, s.workerPool*2)
	for i := 0; i < s.workerPool; i++ {
		s.wg.Add(1)
		go s.worker(ctx, jobs)
	}

	// Start poll loop
	s.wg.Add(1)
	go s.pollLoop(ctx, jobs)

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopCh)
	s.wg.Wait()
}

func (s *Scheduler) pollLoop(ctx context.Context, jobs chan<- *Job) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.fetchJobs(ctx, jobs)
		}
	}
}

func (s *Scheduler) fetchJobs(ctx context.Context, jobs chan<- *Job) {
	pending, err := s.store.GetPending(ctx, s.workerPool*2)
	if err != nil {
		return
	}

	for _, job := range pending {
		if time.Now().After(job.NextRun) {
			select {
			case jobs <- job:
			default:
				// Channel full, skip
			}
		}
	}
}

func (s *Scheduler) worker(ctx context.Context, jobs <-chan *Job) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case job := <-jobs:
			s.executeJob(ctx, job)
		}
	}
}

func (s *Scheduler) executeJob(ctx context.Context, job *Job) {
	// Try to acquire lock
	lockKey := "job_lock_" + job.ID
	if s.lock != nil {
		acquired, err := s.lock.Acquire(ctx, lockKey, job.Timeout+time.Minute)
		if err != nil || !acquired {
			return // Another worker is handling this job
		}
		defer s.lock.Release(ctx, lockKey)
	}

	// Update status
	s.store.UpdateStatus(ctx, job.ID, JobStatusRunning)

	// Get handler
	s.mu.RLock()
	handler, exists := s.handlers[job.Handler]
	s.mu.RUnlock()

	if !exists {
		s.store.UpdateStatus(ctx, job.ID, JobStatusFailed)
		return
	}

	// Set timeout
	execCtx := ctx
	if job.Timeout > 0 {
		var cancel context.CancelFunc
		execCtx, cancel = context.WithTimeout(ctx, job.Timeout)
		defer cancel()
	}

	// Execute
	err := handler(execCtx, job.Payload)
	if err != nil {
		// Handle retry
		if job.Retries < job.MaxRetries {
			job.Retries++
			job.NextRun = time.Now().Add(job.RetryDelay * time.Duration(job.Retries))
			job.Status = JobStatusRetrying
			s.store.Save(ctx, job)
			return
		}
		s.store.UpdateStatus(ctx, job.ID, JobStatusFailed)
		return
	}

	// Handle recurring
	if job.Schedule != "" {
		interval, err := time.ParseDuration(job.Schedule)
		if err == nil {
			now := time.Now()
			job.LastRun = &now
			job.NextRun = time.Now().Add(interval)
			job.Status = JobStatusPending
			job.Retries = 0
			s.store.Save(ctx, job)
			return
		}
	}

	s.store.UpdateStatus(ctx, job.ID, JobStatusCompleted)
}

// InMemoryJobStore is an in-memory implementation
type InMemoryJobStore struct {
	jobs map[string]*Job
	mu   sync.RWMutex
}

func NewInMemoryJobStore() *InMemoryJobStore {
	return &InMemoryJobStore{
		jobs: make(map[string]*Job),
	}
}

func (s *InMemoryJobStore) Save(ctx context.Context, job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
	return nil
}

func (s *InMemoryJobStore) Get(ctx context.Context, id string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobs[id], nil
}

func (s *InMemoryJobStore) GetPending(ctx context.Context, limit int) ([]*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pending := make([]*Job, 0, limit)
	for _, job := range s.jobs {
		if job.Status == JobStatusPending || job.Status == JobStatusRetrying {
			pending = append(pending, job)
			if len(pending) >= limit {
				break
			}
		}
	}

	// Sort by priority (simple bubble sort)
	for i := 0; i < len(pending); i++ {
		for j := i + 1; j < len(pending); j++ {
			if pending[j].Priority > pending[i].Priority {
				pending[i], pending[j] = pending[j], pending[i]
			}
		}
	}

	return pending, nil
}

func (s *InMemoryJobStore) UpdateStatus(ctx context.Context, id string, status JobStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[id]; ok {
		job.Status = status
	}
	return nil
}

func (s *InMemoryJobStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, id)
	return nil
}

// InMemoryLock is an in-memory distributed lock
type InMemoryLock struct {
	locks map[string]time.Time
	mu    sync.Mutex
}

func NewInMemoryLock() *InMemoryLock {
	return &InMemoryLock{
		locks: make(map[string]time.Time),
	}
}

func (l *InMemoryLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if expires, ok := l.locks[key]; ok {
		if time.Now().Before(expires) {
			return false, nil // Already locked
		}
	}

	l.locks[key] = time.Now().Add(ttl)
	return true, nil
}

func (l *InMemoryLock) Release(ctx context.Context, key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.locks, key)
	return nil
}

func (l *InMemoryLock) Extend(ctx context.Context, key string, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.locks[key] = time.Now().Add(ttl)
	return nil
}

func generateJobID() string {
	return "job_" + time.Now().Format("20060102150405.000000")
}
