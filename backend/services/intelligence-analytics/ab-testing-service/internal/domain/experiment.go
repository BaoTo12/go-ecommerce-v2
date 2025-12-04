package domain

import (
	"time"

	"github.com/google/uuid"
)

// ExperimentStatus represents the status of an A/B test
type ExperimentStatus string

const (
	ExperimentDraft   ExperimentStatus = "DRAFT"
	ExperimentRunning ExperimentStatus = "RUNNING"
	ExperimentPaused  ExperimentStatus = "PAUSED"
	ExperimentEnded   ExperimentStatus = "ENDED"
)

// Experiment represents an A/B test experiment
type Experiment struct {
	ExperimentID string
	Name         string
	Description  string
	Status       ExperimentStatus
	Variants     []Variant
	TrafficSplit map[string]int // variant_id -> percentage
	StartDate    *time.Time
	EndDate      *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Variant represents a test variant
type Variant struct {
	VariantID   string
	Name        string
	Description string
	Config      map[string]interface{}
	IsControl   bool
}

// Assignment represents a user's variant assignment
type Assignment struct {
	UserID       string
	ExperimentID string
	VariantID    string
	AssignedAt   time.Time
}

// NewExperiment creates a new experiment
func NewExperiment(name, description string) *Experiment {
	now := time.Now()
	return &Experiment{
		ExperimentID: uuid.New().String(),
		Name:         name,
		Description:  description,
		Status:       ExperimentDraft,
		Variants:     []Variant{},
		TrafficSplit: make(map[string]int),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddVariant adds a variant to the experiment
func (e *Experiment) AddVariant(name, description string, isControl bool, config map[string]interface{}) {
	variant := Variant{
		VariantID:   uuid.New().String(),
		Name:        name,
		Description: description,
		Config:      config,
		IsControl:   isControl,
	}
	e.Variants = append(e.Variants, variant)
}

// Start starts the experiment
func (e *Experiment) Start() {
	now := time.Now()
	e.StartDate = &now
	e.Status = ExperimentRunning
	e.UpdatedAt = now
}

// End ends the experiment
func (e *Experiment) End() {
	now := time.Now()
	e.EndDate = &now
	e.Status = ExperimentEnded
	e.UpdatedAt = now
}

