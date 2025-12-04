package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"
	"github.com/titan-commerce/backend/campaign-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CampaignRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewCampaignRepository(databaseURL string, logger *logger.Logger) (*CampaignRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Campaign PostgreSQL repository initialized")
	return &CampaignRepository{db: db, logger: logger}, nil
}

func (r *CampaignRepository) SaveCampaign(ctx context.Context, campaign *domain.Campaign) error {
	contentJSON, _ := json.Marshal(campaign.Content)
	scheduleJSON, _ := json.Marshal(campaign.Schedule)
	audienceJSON, _ := json.Marshal(campaign.TargetAudience)
	metricsJSON, _ := json.Marshal(campaign.Metrics)

	query := `
		INSERT INTO campaigns (campaign_id, name, description, type, content, schedule, target_audience, metrics, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		campaign.CampaignID, campaign.Name, campaign.Description, campaign.Type,
		contentJSON, scheduleJSON, audienceJSON, metricsJSON, campaign.Status,
		campaign.CreatedAt, campaign.UpdatedAt,
	)
	return err
}

func (r *CampaignRepository) GetCampaign(ctx context.Context, campaignID string) (*domain.Campaign, error) {
	query := `SELECT campaign_id, name, description, type, content, schedule, target_audience, metrics, status, created_at, updated_at FROM campaigns WHERE campaign_id = $1`

	var campaign domain.Campaign
	var contentJSON, scheduleJSON, audienceJSON, metricsJSON []byte

	err := r.db.QueryRowContext(ctx, query, campaignID).Scan(
		&campaign.CampaignID, &campaign.Name, &campaign.Description, &campaign.Type,
		&contentJSON, &scheduleJSON, &audienceJSON, &metricsJSON, &campaign.Status,
		&campaign.CreatedAt, &campaign.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "campaign not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(contentJSON, &campaign.Content)
	json.Unmarshal(scheduleJSON, &campaign.Schedule)
	json.Unmarshal(audienceJSON, &campaign.TargetAudience)
	json.Unmarshal(metricsJSON, &campaign.Metrics)
	return &campaign, nil
}

func (r *CampaignRepository) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	contentJSON, _ := json.Marshal(campaign.Content)
	scheduleJSON, _ := json.Marshal(campaign.Schedule)
	audienceJSON, _ := json.Marshal(campaign.TargetAudience)
	metricsJSON, _ := json.Marshal(campaign.Metrics)

	query := `
		UPDATE campaigns 
		SET name = $2, description = $3, type = $4, content = $5, schedule = $6, target_audience = $7, metrics = $8, status = $9, updated_at = $10
		WHERE campaign_id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		campaign.CampaignID, campaign.Name, campaign.Description, campaign.Type,
		contentJSON, scheduleJSON, audienceJSON, metricsJSON, campaign.Status, campaign.UpdatedAt,
	)
	return err
}

func (r *CampaignRepository) ListActiveCampaigns(ctx context.Context) ([]*domain.Campaign, error) {
	query := `SELECT campaign_id, name, description, type, content, schedule, target_audience, metrics, status, created_at, updated_at 
			  FROM campaigns WHERE status = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, domain.CampaignActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*domain.Campaign
	for rows.Next() {
		var campaign domain.Campaign
		var contentJSON, scheduleJSON, audienceJSON, metricsJSON []byte

		if err := rows.Scan(&campaign.CampaignID, &campaign.Name, &campaign.Description, &campaign.Type,
			&contentJSON, &scheduleJSON, &audienceJSON, &metricsJSON, &campaign.Status,
			&campaign.CreatedAt, &campaign.UpdatedAt); err != nil {
			return nil, err
		}

		json.Unmarshal(contentJSON, &campaign.Content)
		json.Unmarshal(scheduleJSON, &campaign.Schedule)
		json.Unmarshal(audienceJSON, &campaign.TargetAudience)
		json.Unmarshal(metricsJSON, &campaign.Metrics)
		campaigns = append(campaigns, &campaign)
	}
	return campaigns, nil
}

func (r *CampaignRepository) Close() error {
	return r.db.Close()
}
