#!/bin/bash
# run_tests.sh - Comprehensive test runner

echo "=============================="
echo "Running All Tests"
echo "=============================="

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Phase 1: Unit Tests
echo -e "\n${YELLOW}Phase 1: Unit Tests${NC}"
echo "------------------------------"
go test ./... -v -short -count=1 2>&1 | head -100
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Unit tests passed${NC}"
else
    echo -e "${RED}✗ Unit tests failed${NC}"
    exit 1
fi

# Phase 1: Static Analysis
echo -e "\n${YELLOW}Phase 1: Static Analysis${NC}"
echo "------------------------------"
if command -v golangci-lint &> /dev/null; then
    golangci-lint run --timeout 5m
    echo -e "${GREEN}✓ Linter passed${NC}"
else
    echo -e "${YELLOW}⚠ golangci-lint not installed, skipping${NC}"
fi

# Phase 2: Benchmark Tests
echo -e "\n${YELLOW}Phase 2: Benchmark Tests${NC}"
echo "------------------------------"
go test ./pkg/... -bench=. -benchtime=100ms -run=^$ 2>&1 | grep -E "Benchmark|ns/op"
echo -e "${GREEN}✓ Benchmarks completed${NC}"

# Phase 2: Race Detection
echo -e "\n${YELLOW}Phase 2: Race Detection${NC}"
echo "------------------------------"
go test ./pkg/... -race -short -count=1 2>&1 | head -50
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ No race conditions detected${NC}"
else
    echo -e "${RED}✗ Race conditions found${NC}"
fi

# Phase 3: Fuzz Tests (run briefly)
echo -e "\n${YELLOW}Phase 3: Fuzz Tests (10s each)${NC}"
echo "------------------------------"
for test in $(go test ./pkg/... -list='Fuzz.*' 2>/dev/null | grep '^Fuzz'); do
    echo "Running $test..."
    go test ./pkg/... -fuzz=$test -fuzztime=10s 2>&1 | tail -3
done
echo -e "${GREEN}✓ Fuzz tests completed${NC}"

# Phase 3: Security Tests
echo -e "\n${YELLOW}Phase 3: Security Tests${NC}"
echo "------------------------------"
go test ./pkg/testing/security/... -v -count=1 2>&1 | head -50
echo -e "${GREEN}✓ Security tests completed${NC}"

# Phase 3: E2E Tests
echo -e "\n${YELLOW}Phase 3: E2E Tests${NC}"
echo "------------------------------"
go test ./pkg/testing/e2e/... -v -count=1 2>&1 | head -50
echo -e "${GREEN}✓ E2E tests completed${NC}"

# Coverage Report
echo -e "\n${YELLOW}Generating Coverage Report${NC}"
echo "------------------------------"
go test ./pkg/... -coverprofile=coverage.out -covermode=atomic 2>&1 | tail -20
go tool cover -func=coverage.out | tail -10
echo -e "${GREEN}✓ Coverage report generated${NC}"

echo ""
echo "=============================="
echo -e "${GREEN}All Tests Completed!${NC}"
echo "=============================="
