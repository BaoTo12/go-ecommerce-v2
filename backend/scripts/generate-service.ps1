# Service Generation Script (PowerShell)

param(
    [Parameter(Mandatory=$true)]
    [string]$ServiceName,
    
    [Parameter(Mandatory=$true)]
    [ValidateSet('transaction-core', 'catalog-discovery', 'user-social', 'communication', 'logistics-fulfillment', 'marketing-engagement', 'intelligence-analytics')]
    [string]$Category,
    
    [Parameter(Mandatory=$false)]
    [string]$Description = "Microservice for Titan Commerce Platform"
)

$ServiceDir = "services\$Category\$ServiceName"

Write-Host "Generating $ServiceName in $Category..." -ForegroundColor Green

# Create directory structure
New-Item -ItemType Directory -Force -Path "$ServiceDir\cmd\server" | Out-Null
New-Item -ItemType Directory -Force -Path "$ServiceDir\internal\domain" | Out-Null
New-Item -ItemType Directory -Force -Path "$ServiceDir\internal\application" | Out-Null
New-Item -ItemType Directory -Force -Path "$ServiceDir\internal\infrastructure" | Out-Null
New-Item -ItemType Directory -Force -Path "$ServiceDir\internal\interfaces" | Out-Null
New-Item -ItemType Directory -Force -Path "$ServiceDir\proto\$ServiceName\v1" | Out-Null

# Create go.mod
@"
module github.com/titan-commerce/backend/$ServiceName

go 1.23

require (
	github.com/google/uuid v1.5.0
	github.com/titan-commerce/backend/pkg v0.0.0
	google.golang.org/grpc v1.60.1
)

replace github.com/titan-commerce/backend/pkg => ../../pkg
"@ | Out-File -FilePath "$ServiceDir\go.mod" -Encoding UTF8

# Create main.go
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

	log.Info("$ServiceName starting...")
	log.Infof("gRPC server would listen on :%d", cfg.GRPCPort)
	
	select {}
}
"@ | Out-File -FilePath "$ServiceDir\cmd\server\main.go" -Encoding UTF8

# Create README.md
@"
# $ServiceName

$Description

## Quick Start

``````bash
export SERVICE_NAME=$ServiceName
export CELL_ID=cell-001
go run cmd/server/main.go
``````

## Status

🚧 **Under Development** - Skeleton structure created
"@ | Out-File -FilePath "$ServiceDir\README.md" -Encoding UTF8

Write-Host "✓ $ServiceName skeleton created at $ServiceDir" -ForegroundColor Green
