package pool

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

// ConnectionPool manages database connections with pooling
type ConnectionPool struct {
	db          *sql.DB
	maxConns    int
	maxIdleTime time.Duration
	mu          sync.RWMutex
	stats       ConnectionStats
}

// ConnectionStats tracks pool statistics
type ConnectionStats struct {
	OpenConnections   int32
	InUseConnections  int32
	IdleConnections   int32
	WaitCount         int64
	WaitDuration      time.Duration
	MaxOpenReached    int32
}

// ConnectionPoolConfig configures the connection pool
type ConnectionPoolConfig struct {
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
	ConnMaxIdleTime  time.Duration
	HealthCheckFreq  time.Duration
}

// DefaultPoolConfig returns sensible defaults for production
func DefaultPoolConfig() ConnectionPoolConfig {
	return ConnectionPoolConfig{
		MaxOpenConns:    100,
		MaxIdleConns:    25,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		HealthCheckFreq: 1 * time.Minute,
	}
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(dsn string, cfg ConnectionPoolConfig) (*ConnectionPool, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	pool := &ConnectionPool{
		db:          db,
		maxConns:    cfg.MaxOpenConns,
		maxIdleTime: cfg.ConnMaxIdleTime,
	}

	// Start health check goroutine
	go pool.healthCheck(cfg.HealthCheckFreq)

	return pool, nil
}

// GetDB returns the underlying *sql.DB
func (p *ConnectionPool) GetDB() *sql.DB {
	return p.db
}

// Exec executes a query without returning rows
func (p *ConnectionPool) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

// Query executes a query that returns rows
func (p *ConnectionPool) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that returns at most one row
func (p *ConnectionPool) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

// Stats returns pool statistics
func (p *ConnectionPool) Stats() ConnectionStats {
	stats := p.db.Stats()
	return ConnectionStats{
		OpenConnections:  int32(stats.OpenConnections),
		InUseConnections: int32(stats.InUse),
		IdleConnections:  int32(stats.Idle),
		WaitCount:        stats.WaitCount,
		WaitDuration:     stats.WaitDuration,
		MaxOpenReached:   int32(stats.MaxOpenConnections),
	}
}

// healthCheck periodically checks connection health
func (p *ConnectionPool) healthCheck(freq time.Duration) {
	ticker := time.NewTicker(freq)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := p.db.PingContext(ctx); err != nil {
			// Log error, could trigger alerting
		}
		cancel()
	}
}

// Close closes the connection pool
func (p *ConnectionPool) Close() error {
	return p.db.Close()
}

// Transaction executes a function within a transaction
func (p *ConnectionPool) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
