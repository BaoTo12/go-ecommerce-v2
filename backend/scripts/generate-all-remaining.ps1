#!/usr/bin/env pwsh
# Batch Service Generator - Creates all remaining 20 services

$services = @(
    @{Name = "wallet-service"; Category = "transaction-core"; Desc = "Digital wallet with escrow" },
    @{Name = "refund-service"; Category = "transaction-core"; Desc = "Automated refund processing" },
    @{Name = "voucher-service"; Category = "transaction-core"; Desc = "Discount code management" },
    @{Name = "search-service"; Category = "catalog-discovery"; Desc = "Elasticsearch full-text search" },
    @{Name = "recommendation-service"; Category = "catalog-discovery"; Desc = "ML-based product recommendations" },
    @{Name = "category-service"; Category = "catalog-discovery"; Desc = "Category tree management" },
    @{Name = "seller-service"; Category = "catalog-discovery"; Desc = "Seller profiles and KYC" },
    @{Name = "review-service"; Category = "catalog-discovery"; Desc = "Product reviews with spam detection" },
    @{Name = "user-service"; Category = "user-social"; Desc = "User profile management" },
    @{Name = "social-service"; Category = "user-social"; Desc = "Social graph - following/followers" },
    @{Name = "feed-service"; Category = "user-social"; Desc = "TikTok-style activity feed" },
    @{Name = "notification-service"; Category = "user-social"; Desc = "Multi-channel notifications" },
    @{Name = "chat-service"; Category = "communication"; Desc = "WebSocket real-time chat" },
    @{Name = "videocall-service"; Category = "communication"; Desc = "WebRTC video calls" },
    @{Name = "shipping-service"; Category = "logistics-fulfillment"; Desc = "Multi-carrier shipping" },
    @{Name = "tracking-service"; Category = "logistics-fulfillment"; Desc = "Real-time package tracking" },
    @{Name = "warehouse-service"; Category = "logistics-fulfillment"; Desc = "Warehouse management" },
    @{Name = "campaign-service"; Category = "marketing-engagement"; Desc = "Marketing campaign orchestration" },
    @{Name = "coupon-service"; Category = "marketing-engagement"; Desc = "Coupon distribution and validation" },
    @{Name = "pricing-service"; Category = "intelligence-analytics"; Desc = "Dynamic ML-based pricing" },
    @{Name = "fraud-service"; Category = "intelligence-analytics"; Desc = "Real-time fraud detection" },
    @{Name = "analytics-service"; Category = "intelligence-analytics"; Desc = "Business analytics with ClickHouse" },
    @{Name = "ab-testing-service"; Category = "intelligence-analytics"; Desc = "A/B testing and experimentation" }
)

foreach ($svc in $services) {
    Write-Host "Creating $($svc.Name)..." -ForegroundColor Green
    
    $dir = "services\$($svc.Category)\$($svc.Name)"
    
    # Create directories
    New-Item -ItemType Directory -Force -Path "$dir\cmd\server" | Out-Null
    New-Item -ItemType Directory -Force -Path "$dir\internal" | Out-Null
    
    # go.mod
    @"
module github.com/titan-commerce/backend/$($svc.Name)

go 1.23

require (
	github.com/google/uuid v1.5.0
	github.com/titan-commerce/backend/pkg v0.0.0
	google.golang.org/grpc v1.60.1
)

replace github.com/titan-commerce/backend/pkg => ../../pkg
"@ | Out-File -FilePath "$dir\go.mod" -Encoding UTF8
    
    # main.go
    @"
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

	log.Info("$($svc.Name) starting...")
	// TODO: Implement service
	select {}
}
"@ | Out-File -FilePath "$dir\cmd\server\main.go" -Encoding UTF8
    
    # README.md
    @"
# $($svc.Name)

$($svc.Desc)

## Status

🚧 **Under Development** - Skeleton structure created
"@ | Out-File -FilePath "$dir\README.md" -Encoding UTF8
    
    Write-Host "  ✓ Created $($svc.Name)" -ForegroundColor Gray
}

Write-Host "`n✅ All 20 services created successfully!" -ForegroundColor Green
