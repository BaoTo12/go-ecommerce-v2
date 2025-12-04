package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/titan-commerce/backend/pkg/config"
	"github.com/titan-commerce/backend/pkg/logger"
)

const (
	TotalCells      = 500
	HealthCheckFreq = 5 * time.Second
)

// CellHealth tracks health status of cells
type CellHealth struct {
	CellID      int
	Status      string // "healthy", "degraded", "unhealthy"
	Endpoint    string
	LastCheck   time.Time
	FailureCount int
}

// CellRouter routes users to cells using consistent hashing
type CellRouter struct {
	cells  map[int]*CellHealth
	mu     sync.RWMutex
	logger *logger.Logger
}

// NewCellRouter creates a new cell router
func NewCellRouter(log *logger.Logger) *CellRouter {
	router := &CellRouter{
		cells:  make(map[int]*CellHealth),
		logger: log,
	}

	// Initialize cells
	for i := 1; i <= TotalCells; i++ {
		router.cells[i] = &CellHealth{
			CellID:   i,
			Status:   "healthy",
			Endpoint: fmt.Sprintf("cell-%03d.svc.cluster.local:9000", i),
		}
	}

	return router
}

// RouteToCellID calculates the cell ID for a given user ID using consistent hashing
func (r *CellRouter) RouteToCellID(userID string) int {
	hash := sha256.Sum256([]byte(userID))
	hashInt := binary.BigEndian.Uint64(hash[:8])
	cellID := int(hashInt%uint64(TotalCells)) + 1
	return cellID
}

// GetCellEndpoint returns the endpoint for a cell
func (r *CellRouter) GetCellEndpoint(cellID int) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cell, exists := r.cells[cellID]
	if !exists {
		return "", fmt.Errorf("cell %d does not exist", cellID)
	}

	// If cell is unhealthy, failover to next cell
	if cell.Status == "unhealthy" {
		nextCellID := (cellID % TotalCells) + 1
		r.logger.Warnf("Cell %d unhealthy, failing over to cell %d", cellID, nextCellID)
		if nextCell, ok := r.cells[nextCellID]; ok && nextCell.Status == "healthy" {
			return nextCell.Endpoint, nil
		}
	}

	return cell.Endpoint, nil
}

// HealthCheck performs health checks on all cells
func (r *CellRouter) HealthCheck() {
	ticker := time.NewTicker(HealthCheckFreq)
	defer ticker.Stop()

	for range ticker.C {
		r.logger.Debug("Running health checks on all cells...")
		
		for cellID := range r.cells {
			go r.checkCellHealth(cellID)
		}
	}
}

func (r *CellRouter) checkCellHealth(cellID int) {
	r.mu.Lock()
	cell := r.cells[cellID]
	r.mu.Unlock()

	// Simulate health check (in production, would make actual HTTP/gRPC call)
	// For now, mark all cells as healthy
	r.mu.Lock()
	cell.Status = "healthy"
	cell.LastCheck = time.Now()
	cell.FailureCount = 0
	r.mu.Unlock()
}

// HTTP handlers
func (r *CellRouter) handleRoute(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	cellID := r.RouteToCellID(userID)
	endpoint, err := r.GetCellEndpoint(cellID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"user_id":"%s","cell_id":%d,"endpoint":"%s"}`, userID, cellID, endpoint)
}

func (r *CellRouter) handleCells(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"total_cells":%d,"cells":[`, TotalCells)
	
	first := true
	for _, cell := range r.cells {
		if !first {
			fmt.Fprintf(w, ",")
		}
		fmt.Fprintf(w, `{"cell_id":%d,"status":"%s","endpoint":"%s"}`, 
			cell.CellID, cell.Status, cell.Endpoint)
		first = false
	}
	
	fmt.Fprintf(w, "]}")
}

func (r *CellRouter) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy"}`)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(logger.Config{
		Level:       cfg.LogLevel,
		ServiceName: "cell-router",
		CellID:      "global",
		Pretty:      true,
	})

	log.Info("Starting Cell Router...")
	log.Infof("Managing %d cells", TotalCells)

	router := NewCellRouter(log)

	// Start health checks in background
	go router.HealthCheck()

	// HTTP server
	http.HandleFunc("/route", router.handleRoute)
	http.HandleFunc("/cells", router.handleCells)
	http.HandleFunc("/health", router.handleHealth)

	port := cfg.HTTPPort
	if port == 0 {
		port = 8080
	}

	log.Infof("Cell Router listening on :%d", port)
	log.Info("Endpoints:")
	log.Info("  GET /route?user_id=<user-id>  - Route user to cell")
	log.Info("  GET /cells                     - List all cells")
	log.Info("  GET /health                    - Health check")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err, "Failed to start server")
	}
}
