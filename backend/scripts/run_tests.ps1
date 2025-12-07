# run_tests.ps1 - Comprehensive test runner for Windows

Write-Host "==============================" -ForegroundColor Cyan
Write-Host "Running All Tests" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan

$ErrorActionPreference = "Continue"

# Phase 1: Unit Tests
Write-Host "`nPhase 1: Unit Tests" -ForegroundColor Yellow
Write-Host "------------------------------"
$result = go test ./... -v -short -count=1 2>&1
$result | Select-Object -First 50
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Unit tests passed" -ForegroundColor Green
} else {
    Write-Host "✗ Unit tests failed" -ForegroundColor Red
}

# Phase 1: Build Check
Write-Host "`nPhase 1: Build Check" -ForegroundColor Yellow
Write-Host "------------------------------"
go build ./...
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful" -ForegroundColor Green
} else {
    Write-Host "✗ Build failed" -ForegroundColor Red
}

# Phase 1: go vet
Write-Host "`nPhase 1: go vet" -ForegroundColor Yellow
Write-Host "------------------------------"
go vet ./...
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ go vet passed" -ForegroundColor Green
} else {
    Write-Host "✗ go vet found issues" -ForegroundColor Red
}

# Phase 2: Benchmark Tests
Write-Host "`nPhase 2: Benchmark Tests" -ForegroundColor Yellow
Write-Host "------------------------------"
$benchResult = go test ./pkg/... -bench=. -benchtime=100ms -run='^$' 2>&1
$benchResult | Select-String -Pattern "Benchmark|ns/op"
Write-Host "✓ Benchmarks completed" -ForegroundColor Green

# Phase 2: Race Detection
Write-Host "`nPhase 2: Race Detection" -ForegroundColor Yellow
Write-Host "------------------------------"
$raceResult = go test ./pkg/... -race -short -count=1 2>&1
$raceResult | Select-Object -First 30
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ No race conditions detected" -ForegroundColor Green
} else {
    Write-Host "✗ Race conditions found" -ForegroundColor Red
}

# Phase 3: Security Tests
Write-Host "`nPhase 3: Security Tests" -ForegroundColor Yellow
Write-Host "------------------------------"
$secResult = go test ./pkg/testing/security/... -v -count=1 2>&1
$secResult | Select-Object -First 30
Write-Host "✓ Security tests completed" -ForegroundColor Green

# Phase 3: E2E Tests
Write-Host "`nPhase 3: E2E Tests" -ForegroundColor Yellow
Write-Host "------------------------------"
$e2eResult = go test ./pkg/testing/e2e/... -v -count=1 2>&1
$e2eResult | Select-Object -First 30
Write-Host "✓ E2E tests completed" -ForegroundColor Green

# Coverage Report
Write-Host "`nGenerating Coverage Report" -ForegroundColor Yellow
Write-Host "------------------------------"
go test ./pkg/... -coverprofile=coverage.out -covermode=atomic 2>&1 | Select-Object -Last 10
go tool cover -func=coverage.out | Select-Object -Last 10
Write-Host "✓ Coverage report generated" -ForegroundColor Green

Write-Host ""
Write-Host "==============================" -ForegroundColor Cyan
Write-Host "All Tests Completed!" -ForegroundColor Green
Write-Host "==============================" -ForegroundColor Cyan
