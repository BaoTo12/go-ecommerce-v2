package sharding

import (
	"context"
	"errors"
	"hash/fnv"
	"sync"
	"time"
)

/*
SHARDING STRATEGY - Horizontal Database Partitioning

Implements consistent hashing and range-based sharding strategies
for distributing data across multiple database instances.

Supports:
- Consistent hashing (minimal key redistribution on rebalance)
- Range-based partitioning
- Virtual nodes for better distribution
- Dynamic shard addition/removal
*/

var (
	ErrNoShards     = errors.New("no shards available")
	ErrShardNotFound = errors.New("shard not found")
)

// Shard represents a database shard
type Shard struct {
	ID       string
	Host     string
	Port     int
	Weight   int
	IsActive bool
	conn     interface{} // Database connection
}

// ShardingStrategy defines how keys are mapped to shards
type ShardingStrategy interface {
	GetShard(key string) (*Shard, error)
	AddShard(shard *Shard) error
	RemoveShard(shardID string) error
	Rebalance() error
}

// ConsistentHashRing implements consistent hashing
type ConsistentHashRing struct {
	ring         map[uint32]*Shard
	sortedHashes []uint32
	shards       map[string]*Shard
	virtualNodes int
	mu           sync.RWMutex
}

// NewConsistentHashRing creates a consistent hash ring
func NewConsistentHashRing(virtualNodes int) *ConsistentHashRing {
	if virtualNodes <= 0 {
		virtualNodes = 150 // Default virtual nodes per shard
	}
	return &ConsistentHashRing{
		ring:         make(map[uint32]*Shard),
		sortedHashes: make([]uint32, 0),
		shards:       make(map[string]*Shard),
		virtualNodes: virtualNodes,
	}
}

func (r *ConsistentHashRing) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// AddShard adds a shard to the ring
func (r *ConsistentHashRing) AddShard(shard *Shard) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.shards[shard.ID]; exists {
		return errors.New("shard already exists")
	}

	r.shards[shard.ID] = shard

	// Add virtual nodes
	for i := 0; i < r.virtualNodes; i++ {
		virtualKey := shard.ID + "_" + string(rune(i))
		hash := r.hash(virtualKey)
		r.ring[hash] = shard
		r.sortedHashes = append(r.sortedHashes, hash)
	}

	// Sort hashes
	r.sortHashes()

	return nil
}

// RemoveShard removes a shard from the ring
func (r *ConsistentHashRing) RemoveShard(shardID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	shard, exists := r.shards[shardID]
	if !exists {
		return ErrShardNotFound
	}

	// Remove virtual nodes
	for i := 0; i < r.virtualNodes; i++ {
		virtualKey := shard.ID + "_" + string(rune(i))
		hash := r.hash(virtualKey)
		delete(r.ring, hash)
	}

	// Rebuild sorted hashes
	r.sortedHashes = make([]uint32, 0, len(r.ring))
	for h := range r.ring {
		r.sortedHashes = append(r.sortedHashes, h)
	}
	r.sortHashes()

	delete(r.shards, shardID)
	return nil
}

// GetShard returns the shard for a given key
func (r *ConsistentHashRing) GetShard(key string) (*Shard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.sortedHashes) == 0 {
		return nil, ErrNoShards
	}

	hash := r.hash(key)

	// Binary search for the first hash >= key hash
	idx := r.binarySearch(hash)
	if idx >= len(r.sortedHashes) {
		idx = 0 // Wrap around
	}

	return r.ring[r.sortedHashes[idx]], nil
}

func (r *ConsistentHashRing) binarySearch(hash uint32) int {
	low, high := 0, len(r.sortedHashes)-1
	for low < high {
		mid := (low + high) / 2
		if r.sortedHashes[mid] < hash {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low
}

func (r *ConsistentHashRing) sortHashes() {
	// Simple insertion sort for small arrays
	for i := 1; i < len(r.sortedHashes); i++ {
		for j := i; j > 0 && r.sortedHashes[j] < r.sortedHashes[j-1]; j-- {
			r.sortedHashes[j], r.sortedHashes[j-1] = r.sortedHashes[j-1], r.sortedHashes[j]
		}
	}
}

func (r *ConsistentHashRing) Rebalance() error {
	// No-op for consistent hashing - already balanced
	return nil
}

// RangeSharding implements range-based partitioning
type RangeSharding struct {
	ranges []ShardRange
	shards map[string]*Shard
	mu     sync.RWMutex
}

// ShardRange defines a key range for a shard
type ShardRange struct {
	ShardID  string
	StartKey string
	EndKey   string
}

// NewRangeSharding creates range-based sharding
func NewRangeSharding() *RangeSharding {
	return &RangeSharding{
		ranges: make([]ShardRange, 0),
		shards: make(map[string]*Shard),
	}
}

// AddRange adds a key range for a shard
func (s *RangeSharding) AddRange(shardID, startKey, endKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	shard, exists := s.shards[shardID]
	if !exists {
		return ErrShardNotFound
	}

	_ = shard // Use shard if needed

	s.ranges = append(s.ranges, ShardRange{
		ShardID:  shardID,
		StartKey: startKey,
		EndKey:   endKey,
	})

	return nil
}

func (s *RangeSharding) AddShard(shard *Shard) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shards[shard.ID] = shard
	return nil
}

func (s *RangeSharding) RemoveShard(shardID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.shards, shardID)
	return nil
}

func (s *RangeSharding) GetShard(key string) (*Shard, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, r := range s.ranges {
		if key >= r.StartKey && key < r.EndKey {
			return s.shards[r.ShardID], nil
		}
	}

	return nil, ErrShardNotFound
}

func (s *RangeSharding) Rebalance() error {
	// Range sharding typically requires manual rebalancing
	return nil
}

// ShardRouter routes queries to appropriate shards
type ShardRouter struct {
	strategy ShardingStrategy
	readReplicas map[string][]*Shard
	mu sync.RWMutex
}

// NewShardRouter creates a shard router
func NewShardRouter(strategy ShardingStrategy) *ShardRouter {
	return &ShardRouter{
		strategy:     strategy,
		readReplicas: make(map[string][]*Shard),
	}
}

// Execute executes a query on the appropriate shard
func (r *ShardRouter) Execute(ctx context.Context, key string, query func(shard *Shard) error) error {
	shard, err := r.strategy.GetShard(key)
	if err != nil {
		return err
	}

	return query(shard)
}

// ExecuteRead executes a read query, potentially on a replica
func (r *ShardRouter) ExecuteRead(ctx context.Context, key string, query func(shard *Shard) error) error {
	shard, err := r.strategy.GetShard(key)
	if err != nil {
		return err
	}

	r.mu.RLock()
	replicas := r.readReplicas[shard.ID]
	r.mu.RUnlock()

	// Use replica if available
	if len(replicas) > 0 {
		// Simple round-robin (could be improved with health checks)
		replica := replicas[time.Now().UnixNano()%int64(len(replicas))]
		return query(replica)
	}

	return query(shard)
}

// AddReadReplica adds a read replica for a shard
func (r *ShardRouter) AddReadReplica(primaryID string, replica *Shard) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.readReplicas[primaryID] = append(r.readReplicas[primaryID], replica)
}

// ScatterGather executes a query on all shards and gathers results
func (r *ShardRouter) ScatterGather(ctx context.Context, query func(shard *Shard) (interface{}, error)) ([]interface{}, error) {
	r.mu.RLock()
	// Get all unique shards from the ring
	ring, ok := r.strategy.(*ConsistentHashRing)
	if !ok {
		r.mu.RUnlock()
		return nil, errors.New("scatter-gather not supported for this strategy")
	}
	
	shards := make([]*Shard, 0, len(ring.shards))
	for _, shard := range ring.shards {
		shards = append(shards, shard)
	}
	r.mu.RUnlock()

	// Execute in parallel
	results := make([]interface{}, len(shards))
	errs := make([]error, len(shards))
	var wg sync.WaitGroup

	for i, shard := range shards {
		wg.Add(1)
		go func(idx int, s *Shard) {
			defer wg.Done()
			results[idx], errs[idx] = query(s)
		}(i, shard)
	}

	wg.Wait()

	// Check for errors
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
