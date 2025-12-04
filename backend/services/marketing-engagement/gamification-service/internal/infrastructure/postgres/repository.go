package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/titan-commerce/backend/gamification-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type GamificationRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewGamificationRepository(databaseURL string, logger *logger.Logger) (*GamificationRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to open database", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	repo := &GamificationRepository{db: db, logger: logger}
	if err := repo.createTables(ctx); err != nil {
		return nil, err
	}

	logger.Info("Gamification PostgreSQL repository initialized")
	return repo, nil
}

func (r *GamificationRepository) createTables(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS coin_wallets (
			user_id VARCHAR(64) PRIMARY KEY,
			balance INT NOT NULL DEFAULT 0,
			lifetime INT NOT NULL DEFAULT 0,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS coin_transactions (
			id VARCHAR(64) PRIMARY KEY,
			user_id VARCHAR(64) NOT NULL,
			amount INT NOT NULL,
			type VARCHAR(20) NOT NULL,
			source VARCHAR(50),
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_coin_txn_user ON coin_transactions(user_id, created_at DESC)`,
		`CREATE TABLE IF NOT EXISTS daily_checkins (
			user_id VARCHAR(64) PRIMARY KEY,
			last_checkin TIMESTAMP NOT NULL,
			current_streak INT NOT NULL DEFAULT 1,
			longest_streak INT NOT NULL DEFAULT 1,
			total_checkins INT NOT NULL DEFAULT 1
		)`,
		`CREATE TABLE IF NOT EXISTS missions (
			id VARCHAR(64) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			type VARCHAR(30) NOT NULL,
			target INT NOT NULL,
			reward INT NOT NULL,
			start_time TIMESTAMP,
			end_time TIMESTAMP,
			is_active BOOLEAN DEFAULT true
		)`,
		`CREATE TABLE IF NOT EXISTS user_missions (
			id VARCHAR(64) PRIMARY KEY,
			user_id VARCHAR(64) NOT NULL,
			mission_id VARCHAR(64) NOT NULL,
			progress INT NOT NULL DEFAULT 0,
			completed BOOLEAN NOT NULL DEFAULT false,
			claimed_at TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, mission_id)
		)`,
		`CREATE TABLE IF NOT EXISTS lucky_draw_prizes (
			id VARCHAR(64) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			type VARCHAR(30) NOT NULL,
			value INT NOT NULL,
			probability DECIMAL(5,2) NOT NULL,
			image_url TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS lucky_draw_results (
			id VARCHAR(64) PRIMARY KEY,
			user_id VARCHAR(64) NOT NULL,
			prize_id VARCHAR(64),
			spun_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
	}

	for _, query := range queries {
		if _, err := r.db.ExecContext(ctx, query); err != nil {
			return errors.Wrap(errors.ErrInternal, "failed to create table", err)
		}
	}
	return nil
}

func (r *GamificationRepository) GetWallet(ctx context.Context, userID string) (*domain.CoinWallet, error) {
	query := `SELECT user_id, balance, lifetime, updated_at FROM coin_wallets WHERE user_id = $1`

	var wallet domain.CoinWallet
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&wallet.UserID, &wallet.Balance, &wallet.Lifetime, &wallet.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		wallet = *domain.NewCoinWallet(userID)
		if err := r.SaveWallet(ctx, &wallet); err != nil {
			return nil, err
		}
		return &wallet, nil
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to get wallet", err)
	}
	return &wallet, nil
}

func (r *GamificationRepository) SaveWallet(ctx context.Context, wallet *domain.CoinWallet) error {
	query := `INSERT INTO coin_wallets (user_id, balance, lifetime, updated_at)
			  VALUES ($1, $2, $3, $4)
			  ON CONFLICT (user_id) DO UPDATE SET
			  balance = EXCLUDED.balance, lifetime = EXCLUDED.lifetime, updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		wallet.UserID, wallet.Balance, wallet.Lifetime, wallet.UpdatedAt)
	return err
}

func (r *GamificationRepository) SaveTransaction(ctx context.Context, txn *domain.CoinTransaction) error {
	query := `INSERT INTO coin_transactions (id, user_id, amount, type, source, description, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		txn.ID, txn.UserID, txn.Amount, txn.Type, txn.Source, txn.Description, txn.CreatedAt)
	return err
}

func (r *GamificationRepository) GetTransactions(ctx context.Context, userID string, limit int) ([]*domain.CoinTransaction, error) {
	query := `SELECT id, user_id, amount, type, source, description, created_at
			  FROM coin_transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*domain.CoinTransaction
	for rows.Next() {
		var txn domain.CoinTransaction
		var source, desc sql.NullString
		if err := rows.Scan(&txn.ID, &txn.UserID, &txn.Amount, &txn.Type,
			&source, &desc, &txn.CreatedAt); err != nil {
			return nil, err
		}
		txn.Source = source.String
		txn.Description = desc.String
		txns = append(txns, &txn)
	}
	return txns, nil
}

func (r *GamificationRepository) GetCheckIn(ctx context.Context, userID string) (*domain.DailyCheckIn, error) {
	query := `SELECT user_id, last_checkin, current_streak, longest_streak, total_checkins 
			  FROM daily_checkins WHERE user_id = $1`

	var checkIn domain.DailyCheckIn
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&checkIn.UserID, &checkIn.LastCheckIn, &checkIn.CurrentStreak, 
		&checkIn.LongestStreak, &checkIn.TotalCheckIns)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (r *GamificationRepository) SaveCheckIn(ctx context.Context, checkIn *domain.DailyCheckIn) error {
	query := `INSERT INTO daily_checkins (user_id, last_checkin, current_streak, longest_streak, total_checkins)
			  VALUES ($1, $2, $3, $4, $5)
			  ON CONFLICT (user_id) DO UPDATE SET
			  last_checkin = EXCLUDED.last_checkin, current_streak = EXCLUDED.current_streak,
			  longest_streak = EXCLUDED.longest_streak, total_checkins = EXCLUDED.total_checkins`

	_, err := r.db.ExecContext(ctx, query,
		checkIn.UserID, checkIn.LastCheckIn, checkIn.CurrentStreak, 
		checkIn.LongestStreak, checkIn.TotalCheckIns)
	return err
}

func (r *GamificationRepository) GetActiveMissions(ctx context.Context) ([]*domain.Mission, error) {
	query := `SELECT id, name, description, type, target, reward, start_time, end_time, is_active
			  FROM missions WHERE is_active = true AND (end_time IS NULL OR end_time > NOW())`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*domain.Mission
	for rows.Next() {
		var m domain.Mission
		var startTime, endTime sql.NullTime
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Type, &m.Target, &m.Reward, 
			&startTime, &endTime, &m.IsActive); err != nil {
			return nil, err
		}
		if startTime.Valid {
			m.StartTime = startTime.Time
		}
		if endTime.Valid {
			m.EndTime = endTime.Time
		}
		missions = append(missions, &m)
	}
	return missions, nil
}

func (r *GamificationRepository) GetUserMission(ctx context.Context, userID, missionID string) (*domain.UserMission, error) {
	query := `SELECT id, user_id, mission_id, progress, completed, claimed_at, updated_at 
			  FROM user_missions WHERE user_id = $1 AND mission_id = $2`

	var um domain.UserMission
	var claimedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, userID, missionID).Scan(
		&um.ID, &um.UserID, &um.MissionID, &um.Progress, &um.Completed, &claimedAt, &um.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if claimedAt.Valid {
		um.ClaimedAt = &claimedAt.Time
	}
	return &um, err
}

func (r *GamificationRepository) SaveUserMission(ctx context.Context, um *domain.UserMission) error {
	query := `INSERT INTO user_missions (id, user_id, mission_id, progress, completed, claimed_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  ON CONFLICT (user_id, mission_id) DO UPDATE SET
			  progress = EXCLUDED.progress, completed = EXCLUDED.completed,
			  claimed_at = EXCLUDED.claimed_at, updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		um.ID, um.UserID, um.MissionID, um.Progress, um.Completed, um.ClaimedAt, um.UpdatedAt)
	return err
}

func (r *GamificationRepository) GetUserMissions(ctx context.Context, userID string) ([]*domain.UserMission, error) {
	query := `SELECT id, user_id, mission_id, progress, completed, claimed_at, updated_at 
			  FROM user_missions WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*domain.UserMission
	for rows.Next() {
		var um domain.UserMission
		var claimedAt sql.NullTime
		if err := rows.Scan(&um.ID, &um.UserID, &um.MissionID, &um.Progress, 
			&um.Completed, &claimedAt, &um.UpdatedAt); err != nil {
			return nil, err
		}
		if claimedAt.Valid {
			um.ClaimedAt = &claimedAt.Time
		}
		missions = append(missions, &um)
	}
	return missions, nil
}

func (r *GamificationRepository) GetPrizes(ctx context.Context) ([]*domain.LuckyDrawPrize, error) {
	query := `SELECT id, name, type, value, probability, image_url FROM lucky_draw_prizes`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prizes []*domain.LuckyDrawPrize
	for rows.Next() {
		var p domain.LuckyDrawPrize
		var imageURL sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Value, &p.Probability, &imageURL); err != nil {
			return nil, err
		}
		p.ImageURL = imageURL.String
		prizes = append(prizes, &p)
	}

	if len(prizes) == 0 {
		return []*domain.LuckyDrawPrize{
			{ID: "1", Name: "100 Coins", Type: "coins", Value: 100, Probability: 50},
			{ID: "2", Name: "500 Coins", Type: "coins", Value: 500, Probability: 30},
			{ID: "3", Name: "1000 Coins", Type: "coins", Value: 1000, Probability: 15},
			{ID: "4", Name: "No Prize", Type: "nothing", Value: 0, Probability: 5},
		}, nil
	}
	return prizes, nil
}

func (r *GamificationRepository) SaveDrawResult(ctx context.Context, result *domain.LuckyDrawResult) error {
	query := `INSERT INTO lucky_draw_results (id, user_id, prize_id, spun_at)
			  VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query,
		result.ID, result.UserID, result.PrizeID, result.SpunAt)
	return err
}

func (r *GamificationRepository) Close() error {
	return r.db.Close()
}
