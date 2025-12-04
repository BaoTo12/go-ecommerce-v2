package domain

import "context"

type ExperimentRepository interface {
	SaveExperiment(ctx context.Context, exp *Experiment) error
	GetExperiment(ctx context.Context, experimentID string) (*Experiment, error)
	UpdateExperiment(ctx context.Context, exp *Experiment) error
	GetActiveExperiments(ctx context.Context) ([]*Experiment, error)
	SaveAssignment(ctx context.Context, assignment *Assignment) error
	GetAssignment(ctx context.Context, userID, experimentID string) (*Assignment, error)
}

