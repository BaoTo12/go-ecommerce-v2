package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/titan-commerce/backend/fraud-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type FraudRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewFraudRepository(databaseURL string, logger *logger.Logger) (*FraudRepository, error) {
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

	repo := &FraudRepository{db: db, logger: logger}
	if err := repo.createTables(ctx); err != nil {
		return nil, err
	}

	logger.Info("Fraud PostgreSQL repository initialized")
	return repo, nil
}

func (r *FraudRepository) createTables(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS fraud_checks (
			id VARCHAR(64) PRIMARY KEY,
			transaction_id VARCHAR(64) NOT NULL,
			user_id VARCHAR(64) NOT NULL,
			amount DECIMAL(12,2) NOT NULL,
			currency VARCHAR(3) NOT NULL,
			ip VARCHAR(45),
			device_id VARCHAR(64),
			user_agent TEXT,
			score DECIMAL(5,4) NOT NULL,
			risk_level VARCHAR(20) NOT NULL,
			decision VARCHAR(20) NOT NULL,
			reasons TEXT[],
			processing_time BIGINT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_fraud_txn ON fraud_checks(transaction_id)`,
		`CREATE INDEX IF NOT EXISTS idx_fraud_user ON fraud_checks(user_id, created_at DESC)`,
		`CREATE TABLE IF NOT EXISTS fraud_features (
			fraud_check_id VARCHAR(64) PRIMARY KEY REFERENCES fraud_checks(id),
			account_age_days INT,
			total_orders INT,
			avg_order_value DECIMAL(12,2),
			orders_last_24h INT,
			orders_last_7d INT,
			failed_payments_7d INT,
			unique_devices INT,
			unique_ips INT,
			new_device BOOLEAN,
			new_ip BOOLEAN,
			velocity_score DECIMAL(5,4),
			amount_deviation DECIMAL(5,2)
		)`,
		`CREATE TABLE IF NOT EXISTS fraud_rules (
			id VARCHAR(64) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			condition TEXT NOT NULL,
			action VARCHAR(20) NOT NULL,
			score_weight DECIMAL(5,2),
			is_active BOOLEAN DEFAULT true,
			priority INT DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS fraud_alerts (
			id VARCHAR(64) PRIMARY KEY,
			fraud_check_id VARCHAR(64) REFERENCES fraud_checks(id),
			alert_type VARCHAR(30) NOT NULL,
			severity VARCHAR(20) NOT NULL,
			message TEXT,
			acknowledged BOOLEAN DEFAULT false,
			resolved_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_unresolved ON fraud_alerts(acknowledged, created_at DESC)`,
		`CREATE TABLE IF NOT EXISTS user_devices (
			user_id VARCHAR(64) NOT NULL,
			device_id VARCHAR(64) NOT NULL,
			first_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, device_id)
		)`,
		`CREATE TABLE IF NOT EXISTS user_ips (
			user_id VARCHAR(64) NOT NULL,
			ip VARCHAR(45) NOT NULL,
			first_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, ip)
		)`,
		`CREATE TABLE IF NOT EXISTS user_stats (
			user_id VARCHAR(64) PRIMARY KEY,
			account_created_at TIMESTAMP,
			total_orders INT DEFAULT 0,
			avg_order_value DECIMAL(12,2) DEFAULT 0,
			failed_payments INT DEFAULT 0,
			last_order_at TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
	}

	for _, query := range queries {
		if _, err := r.db.ExecContext(ctx, query); err != nil {
			return errors.Wrap(errors.ErrInternal, "failed to create table", err)
		}
	}
	return nil
}

func (r *FraudRepository) Save(ctx context.Context, check *domain.FraudCheck) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO fraud_checks 
			  (id, transaction_id, user_id, amount, currency, ip, device_id, user_agent, 
			   score, risk_level, decision, reasons, processing_time, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err = tx.ExecContext(ctx, query,
		check.ID, check.TransactionID, check.UserID, check.Amount, check.Currency,
		check.IP, check.DeviceID, check.UserAgent, check.Score, check.RiskLevel,
		check.Decision, check.Reasons, check.ProcessingTime, check.CreatedAt)
	if err != nil {
		return err
	}

	featuresQuery := `INSERT INTO fraud_features 
					  (fraud_check_id, account_age_days, total_orders, avg_order_value,
					   orders_last_24h, orders_last_7d, failed_payments_7d, unique_devices,
					   unique_ips, new_device, new_ip, velocity_score, amount_deviation)
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = tx.ExecContext(ctx, featuresQuery,
		check.ID, check.Features.AccountAgeDays, check.Features.TotalOrders,
		check.Features.AvgOrderValue, check.Features.OrdersLast24h, check.Features.OrdersLast7d,
		check.Features.FailedPaymentsLast7d, check.Features.UniqueDevices, check.Features.UniqueIPs,
		check.Features.NewDevice, check.Features.NewIP, check.Features.VelocityScore,
		check.Features.AmountDeviation)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *FraudRepository) FindByID(ctx context.Context, checkID string) (*domain.FraudCheck, error) {
	query := `SELECT c.id, c.transaction_id, c.user_id, c.amount, c.currency, c.ip, c.device_id,
			  c.user_agent, c.score, c.risk_level, c.decision, c.reasons, c.processing_time, c.created_at,
			  f.account_age_days, f.total_orders, f.avg_order_value, f.orders_last_24h, f.orders_last_7d,
			  f.failed_payments_7d, f.unique_devices, f.unique_ips, f.new_device, f.new_ip,
			  f.velocity_score, f.amount_deviation
			  FROM fraud_checks c
			  LEFT JOIN fraud_features f ON c.id = f.fraud_check_id
			  WHERE c.id = $1`

	var check domain.FraudCheck
	var reasons []string
	err := r.db.QueryRowContext(ctx, query, checkID).Scan(
		&check.ID, &check.TransactionID, &check.UserID, &check.Amount, &check.Currency,
		&check.IP, &check.DeviceID, &check.UserAgent, &check.Score, &check.RiskLevel,
		&check.Decision, &reasons, &check.ProcessingTime, &check.CreatedAt,
		&check.Features.AccountAgeDays, &check.Features.TotalOrders, &check.Features.AvgOrderValue,
		&check.Features.OrdersLast24h, &check.Features.OrdersLast7d, &check.Features.FailedPaymentsLast7d,
		&check.Features.UniqueDevices, &check.Features.UniqueIPs, &check.Features.NewDevice,
		&check.Features.NewIP, &check.Features.VelocityScore, &check.Features.AmountDeviation,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "fraud check not found")
	}
	check.Reasons = reasons
	return &check, err
}

func (r *FraudRepository) FindByTransaction(ctx context.Context, txnID string) (*domain.FraudCheck, error) {
	query := `SELECT id FROM fraud_checks WHERE transaction_id = $1 LIMIT 1`
	var checkID string
	err := r.db.QueryRowContext(ctx, query, txnID).Scan(&checkID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, checkID)
}

func (r *FraudRepository) FindByUser(ctx context.Context, userID string, limit int) ([]*domain.FraudCheck, error) {
	query := `SELECT id, transaction_id, user_id, amount, currency, score, risk_level, 
			  decision, created_at FROM fraud_checks 
			  WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []*domain.FraudCheck
	for rows.Next() {
		var c domain.FraudCheck
		if err := rows.Scan(&c.ID, &c.TransactionID, &c.UserID, &c.Amount,
			&c.Currency, &c.Score, &c.RiskLevel, &c.Decision, &c.CreatedAt); err != nil {
			return nil, err
		}
		checks = append(checks, &c)
	}
	return checks, nil
}

func (r *FraudRepository) GetActiveRules(ctx context.Context) ([]*domain.FraudRule, error) {
	query := `SELECT id, name, description, condition, action, score_weight, priority
			  FROM fraud_rules WHERE is_active = true ORDER BY priority DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.FraudRule
	for rows.Next() {
		var rule domain.FraudRule
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.Description, &rule.Condition,
			&rule.Action, &rule.ScoreWeight, &rule.Priority); err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}
	return rules, nil
}

func (r *FraudRepository) SaveAlert(ctx context.Context, alert *domain.FraudAlert) error {
	query := `INSERT INTO fraud_alerts (id, fraud_check_id, alert_type, severity, message, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		alert.ID, alert.FraudCheckID, alert.AlertType, alert.Severity, alert.Message, alert.CreatedAt)
	return err
}

func (r *FraudRepository) GetUnresolvedAlerts(ctx context.Context) ([]*domain.FraudAlert, error) {
	query := `SELECT id, fraud_check_id, alert_type, severity, message, acknowledged, created_at
			  FROM fraud_alerts WHERE acknowledged = false ORDER BY created_at DESC LIMIT 100`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*domain.FraudAlert
	for rows.Next() {
		var a domain.FraudAlert
		if err := rows.Scan(&a.ID, &a.FraudCheckID, &a.AlertType, &a.Severity,
			&a.Message, &a.Acknowledged, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, &a)
	}
	return alerts, nil
}

func (r *FraudRepository) GetUserStats(ctx context.Context, userID string) (*domain.FraudFeatures, error) {
	query := `SELECT total_orders, avg_order_value, failed_payments,
			  EXTRACT(DAY FROM NOW() - account_created_at)::INT as age_days
			  FROM user_stats WHERE user_id = $1`

	var features domain.FraudFeatures
	var ageDays sql.NullInt32
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&features.TotalOrders, &features.AvgOrderValue, &features.FailedPaymentsLast7d, &ageDays)
	if err == sql.ErrNoRows {
		return &domain.FraudFeatures{AccountAgeDays: 0}, nil
	}
	if ageDays.Valid {
		features.AccountAgeDays = int(ageDays.Int32)
	}
	return &features, nil
}

func (r *FraudRepository) RecordDevice(ctx context.Context, userID, deviceID string) error {
	query := `INSERT INTO user_devices (user_id, device_id) VALUES ($1, $2)
			  ON CONFLICT (user_id, device_id) DO UPDATE SET last_seen = NOW()`
	_, err := r.db.ExecContext(ctx, query, userID, deviceID)
	return err
}

func (r *FraudRepository) RecordIP(ctx context.Context, userID, ip string) error {
	query := `INSERT INTO user_ips (user_id, ip) VALUES ($1, $2)
			  ON CONFLICT (user_id, ip) DO UPDATE SET last_seen = NOW()`
	_, err := r.db.ExecContext(ctx, query, userID, ip)
	return err
}

func (r *FraudRepository) IsNewDevice(ctx context.Context, userID, deviceID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_devices WHERE user_id = $1 AND device_id = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, deviceID).Scan(&exists)
	return !exists, err
}

func (r *FraudRepository) IsNewIP(ctx context.Context, userID, ip string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_ips WHERE user_id = $1 AND ip = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, ip).Scan(&exists)
	return !exists, err
}

func (r *FraudRepository) Close() error {
	return r.db.Close()
}
