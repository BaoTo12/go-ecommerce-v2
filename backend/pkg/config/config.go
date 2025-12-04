package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServiceName      string
	ServicePort      int
	CellID           string
	LogLevel         string
	DatabaseURL      string
	DatabaseMaxConns int
	RedisAddr        string
	RedisPassword    string
	KafkaBrokers     []string
	MongoURI         string
	MongoDatabase    string
	GRPCPort         int
	HTTPPort         int
	JWTSecret        string
	JWTAccessExpiry  int
}

func Load() (*Config, error) {
	cfg := &Config{
		ServiceName:      getEnv("SERVICE_NAME", "unknown-service"),
		ServicePort:      getEnvInt("SERVICE_PORT", 8080),
		CellID:           getEnv("CELL_ID", "cell-000"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		DatabaseMaxConns: getEnvInt("DATABASE_MAX_CONNS", 25),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		KafkaBrokers:     []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:    getEnv("MONGO_DATABASE", "titan_commerce"),
		GRPCPort:         getEnvInt("GRPC_PORT", 9000),
		HTTPPort:         getEnvInt("HTTP_PORT", 8080),
		JWTSecret:        getEnv("JWT_SECRET", "changeme-in-production"),
		JWTAccessExpiry:  getEnvInt("JWT_ACCESS_EXPIRY", 15),
	}

	if cfg.ServiceName == "unknown-service" {
		return nil, fmt.Errorf("SERVICE_NAME environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
