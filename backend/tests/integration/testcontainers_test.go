// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer manages a PostgreSQL testcontainer
type PostgresContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
	DSN       string
}

func StartPostgres(ctx context.Context, t *testing.T) *PostgresContainer {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	return &PostgresContainer{
		Container: container,
		Host:      host,
		Port:      port.Port(),
		DSN:       "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable",
	}
}

func (c *PostgresContainer) Stop(ctx context.Context) {
	c.Container.Terminate(ctx)
}

// RedisContainer manages a Redis testcontainer
type RedisContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
	Addr      string
}

func StartRedis(ctx context.Context, t *testing.T) *RedisContainer {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "6379")

	return &RedisContainer{
		Container: container,
		Host:      host,
		Port:      port.Port(),
		Addr:      host + ":" + port.Port(),
	}
}

func (c *RedisContainer) Stop(ctx context.Context) {
	c.Container.Terminate(ctx)
}

// MongoContainer manages a MongoDB testcontainer
type MongoContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
	URI       string
}

func StartMongo(ctx context.Context, t *testing.T) *MongoContainer {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "27017")

	return &MongoContainer{
		Container: container,
		Host:      host,
		Port:      port.Port(),
		URI:       "mongodb://" + host + ":" + port.Port(),
	}
}

func (c *MongoContainer) Stop(ctx context.Context) {
	c.Container.Terminate(ctx)
}

func TestOrderServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start PostgreSQL container
	postgres := StartPostgres(ctx, t)
	defer postgres.Stop(ctx)

	t.Log("PostgreSQL started at:", postgres.DSN)

	// Test order creation with real database
	// In production, you would:
	// 1. Run migrations
	// 2. Initialize the real repository
	// 3. Run tests against the real database

	assert.NotEmpty(t, postgres.DSN)
}

func TestCartServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start Redis container
	redis := StartRedis(ctx, t)
	defer redis.Stop(ctx)

	t.Log("Redis started at:", redis.Addr)

	// Test cart operations with real Redis
	// In production, you would:
	// 1. Initialize the real Redis repository
	// 2. Run tests against real Redis

	assert.NotEmpty(t, redis.Addr)
}

func TestChatServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start MongoDB container
	mongo := StartMongo(ctx, t)
	defer mongo.Stop(ctx)

	t.Log("MongoDB started at:", mongo.URI)

	// Test chat operations with real MongoDB
	// In production, you would:
	// 1. Initialize the real MongoDB repository
	// 2. Run tests against real MongoDB

	assert.NotEmpty(t, mongo.URI)
}
