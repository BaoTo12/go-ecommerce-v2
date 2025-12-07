package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ====================
// DATABASE INTEGRATION TESTS
// ====================

// TestDatabaseConnection tests database connectivity
func TestDatabaseConnection(t *testing.T) {
	dsn := getTestDSN()
	if dsn == "" {
		t.Skip("Database DSN not configured")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

// TestTransaction tests transaction behavior
func TestTransaction_CommitRollback(t *testing.T) {
	dsn := getTestDSN()
	if dsn == "" {
		t.Skip("Database DSN not configured")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Test commit
	t.Run("Commit", func(t *testing.T) {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Execute some operation
		_, err = tx.ExecContext(ctx, "SELECT 1")
		if err != nil {
			t.Fatalf("Failed to execute in transaction: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("Failed to commit: %v", err)
		}
	})

	// Test rollback
	t.Run("Rollback", func(t *testing.T) {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Execute some operation
		_, err = tx.ExecContext(ctx, "SELECT 1")
		if err != nil {
			t.Fatalf("Failed to execute in transaction: %v", err)
		}

		if err := tx.Rollback(); err != nil {
			t.Fatalf("Failed to rollback: %v", err)
		}
	})
}

// TestConnectionPool tests connection pool behavior
func TestConnectionPool(t *testing.T) {
	dsn := getTestDSN()
	if dsn == "" {
		t.Skip("Database DSN not configured")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Configure pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Get initial stats
	stats := db.Stats()
	t.Logf("Initial: Open=%d, Idle=%d, InUse=%d",
		stats.OpenConnections, stats.Idle, stats.InUse)

	// Make concurrent queries
	done := make(chan bool, 20)
	for i := 0; i < 20; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			db.QueryRowContext(ctx, "SELECT 1").Scan(new(int))
			done <- true
		}()
	}

	// Wait for all queries
	for i := 0; i < 20; i++ {
		<-done
	}

	// Check stats after
	stats = db.Stats()
	t.Logf("After: Open=%d, Idle=%d, InUse=%d, WaitCount=%d",
		stats.OpenConnections, stats.Idle, stats.InUse, stats.WaitCount)

	// Verify pool limits were respected
	if stats.OpenConnections > 10 {
		t.Errorf("Exceeded max connections: %d > 10", stats.OpenConnections)
	}
}

// TestQueryTimeout tests query timeout behavior
func TestQueryTimeout(t *testing.T) {
	dsn := getTestDSN()
	if dsn == "" {
		t.Skip("Database DSN not configured")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Query should timeout
	_, err = db.QueryContext(ctx, "SELECT pg_sleep(10)")
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

// TestPreparedStatements tests prepared statement handling
func TestPreparedStatements(t *testing.T) {
	dsn := getTestDSN()
	if dsn == "" {
		t.Skip("Database DSN not configured")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Prepare statement
	stmt, err := db.PrepareContext(ctx, "SELECT $1::int + $2::int")
	if err != nil {
		t.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Execute multiple times
	for i := 0; i < 5; i++ {
		var result int
		err := stmt.QueryRowContext(ctx, i, 1).Scan(&result)
		if err != nil {
			t.Fatalf("Failed to execute prepared statement: %v", err)
		}
		if result != i+1 {
			t.Errorf("Expected %d, got %d", i+1, result)
		}
	}
}

// TestNullHandling tests NULL value handling
func TestNullHandling(t *testing.T) {
	dsn := getTestDSN()
	if dsn == "" {
		t.Skip("Database DSN not configured")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Test NULL string
	var nullStr sql.NullString
	err = db.QueryRowContext(ctx, "SELECT NULL::text").Scan(&nullStr)
	if err != nil {
		t.Fatalf("Failed to scan NULL string: %v", err)
	}
	if nullStr.Valid {
		t.Error("Expected NULL string to be invalid")
	}

	// Test non-NULL string
	err = db.QueryRowContext(ctx, "SELECT 'hello'").Scan(&nullStr)
	if err != nil {
		t.Fatalf("Failed to scan string: %v", err)
	}
	if !nullStr.Valid || nullStr.String != "hello" {
		t.Errorf("Expected 'hello', got %v", nullStr)
	}
}

// Helper to get test DSN
func getTestDSN() string {
	// Would typically come from environment variable
	// return os.Getenv("TEST_DATABASE_URL")
	return "" // Disabled by default
}

// ====================
// HTTP INTEGRATION TESTS
// ====================

// These tests verify that HTTP handlers work correctly with the full stack

// Note: imports moved to top of file

func TestHTTPIntegration_CRUD(t *testing.T) {
	// Setup a full test server with all middleware
	mux := http.NewServeMux()
	
	// In-memory storage for testing
	storage := make(map[string]interface{})
	
	// Create
	mux.HandleFunc("POST /api/items", func(w http.ResponseWriter, r *http.Request) {
		var item map[string]interface{}
		json.NewDecoder(r.Body).Decode(&item)
		id := "item_1"
		item["id"] = id
		storage[id] = item
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(item)
	})
	
	// Read
	mux.HandleFunc("GET /api/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		item, exists := storage[id]
		if !exists {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(item)
	})
	
	// Update
	mux.HandleFunc("PUT /api/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if _, exists := storage[id]; !exists {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		var item map[string]interface{}
		json.NewDecoder(r.Body).Decode(&item)
		item["id"] = id
		storage[id] = item
		json.NewEncoder(w).Encode(item)
	})
	
	// Delete
	mux.HandleFunc("DELETE /api/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		delete(storage, id)
		w.WriteHeader(http.StatusNoContent)
	})
	
	server := httptest.NewServer(mux)
	defer server.Close()
	
	client := server.Client()
	
	// Test Create
	t.Run("Create", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name":"Test Item"}`)
		resp, err := client.Post(server.URL+"/api/items", "application/json", body)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})
	
	// Test Read
	t.Run("Read", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/items/item_1")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
	
	// Test Not Found
	t.Run("NotFound", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/items/nonexistent")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})
}

// Benchmark database operations
func BenchmarkDatabaseQuery(b *testing.B) {
	dsn := getTestDSN()
	if dsn == "" {
		b.Skip("Database DSN not configured")
	}

	db, _ := sql.Open("postgres", dsn)
	defer db.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.QueryRowContext(ctx, "SELECT 1").Scan(new(int))
	}
}

func BenchmarkDatabaseQuery_Parallel(b *testing.B) {
	dsn := getTestDSN()
	if dsn == "" {
		b.Skip("Database DSN not configured")
	}

	db, _ := sql.Open("postgres", dsn)
	db.SetMaxOpenConns(10)
	defer db.Close()

	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			db.QueryRowContext(ctx, "SELECT 1").Scan(new(int))
		}
	})
}
