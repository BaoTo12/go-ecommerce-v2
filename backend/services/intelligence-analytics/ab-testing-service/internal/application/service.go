package application

import (
	"context"
	"hash/fnv"

	"github.com/titan-commerce/backend/ab-testing-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type ABTestingService struct {
	repo   domain.ExperimentRepository
	logger *logger.Logger
}

func NewABTestingService(repo domain.ExperimentRepository, logger *logger.Logger) *ABTestingService {
	return &ABTestingService{repo: repo, logger: logger}
}

// CreateExperiment creates a new A/B test
func (s *ABTestingService) CreateExperiment(
	ctx context.Context,
	name, description string,
) (*domain.Experiment, error) {
	exp := domain.NewExperiment(name, description)

	if err := s.repo.SaveExperiment(ctx, exp); err != nil {
		s.logger.Error(err, "failed to create experiment")
		return nil, err
	}

	s.logger.Infof("Experiment created: %s", name)
	return exp, nil
}

// AssignVariant assigns a user to a variant
func (s *ABTestingService) AssignVariant(
	ctx context.Context,
	userID, experimentID string,
) (string, error) {
	// Check existing assignment
	existing, _ := s.repo.GetAssignment(ctx, userID, experimentID)
	if existing != nil {
		return existing.VariantID, nil
	}

	// Get experiment
	exp, err := s.repo.GetExperiment(ctx, experimentID)
	if err != nil {
		return "", err
	}

	// Assign based on hash
	variantID := s.hashAssignment(userID, experimentID, exp.Variants)

	// Save assignment
	assignment := &domain.Assignment{
		UserID:       userID,
		ExperimentID: experimentID,
		VariantID:    variantID,
	}
	s.repo.SaveAssignment(ctx, assignment)

	return variantID, nil
}

func (s *ABTestingService) hashAssignment(userID, experimentID string, variants []domain.Variant) string {
	h := fnv.New32a()
	h.Write([]byte(userID + experimentID))
	idx := int(h.Sum32()) % len(variants)
	return variants[idx].VariantID
}

