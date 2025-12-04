package main

import (
	"fmt"
	"os"

	"github.com/titan-commerce/backend/pkg/config"
	"github.com/titan-commerce/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(logger.Config{
		Level:       cfg.LogLevel,
		ServiceName: cfg.ServiceName,
		CellID:      cfg.CellID,
		Pretty:      true,
	})

	log.Info("🔐 Auth Service starting...")
	
	// TODO: Implement authentication
	// - JWT access tokens (15 min expiry) + refresh tokens (30 days)
	// - MFA support (TOTP, SMS)  
	// - OAuth2/OIDC integration (Google, Facebook login)
	// - Redis for token blacklist (logout)
	// - Rate limiting on login attempts
	// - Password hashing (bcrypt)
	
	select {}
}
