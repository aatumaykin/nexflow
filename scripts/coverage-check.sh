#!/bin/bash
# Coverage checker script
# Verifies that code coverage meets the minimum threshold (60%)

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Running tests with coverage...${NC}"

# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Check total coverage
TOTAL=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')

echo ""
echo -e "${YELLOW}=== Coverage Report ===${NC}"
echo -e "Total coverage: ${TOTAL}%"

MIN_COVERAGE=60

# Compare with threshold
if (( $(echo "$TOTAL < $MIN_COVERAGE" | bc -l) )); then
    echo -e "${RED}❌ Coverage (${TOTAL}%) is below threshold (${MIN_COVERAGE}%)${NC}"
    echo ""
    echo -e "${YELLOW}Top 10 least covered packages:${NC}"
    go tool cover -func=coverage.out | grep "^github.com" | \
        awk '{print $1, $3}' | \
        sort -t' ' -k2 -n | \
        head -10 | \
        awk '{printf "  %-50s %s\n", $1, $2}'
    exit 1
else
    echo -e "${GREEN}✅ Coverage (${TOTAL}%) meets threshold (${MIN_COVERAGE}%)${NC}"
    echo ""
    echo -e "${YELLOW}Per-package coverage:${NC}"
    go tool cover -func=coverage.out | grep "^github.com" | \
        awk '{
            pkg = $1
            gsub(/\/github.com\/atumaikin\/nexflow\//, "", pkg)
            cov = $3
            printf "  %-50s %s\n", pkg, cov
        }' | sort
fi

echo ""
echo -e "${YELLOW}=== Generating HTML report ===${NC}"
go tool cover -html=coverage.out -o coverage.html
echo -e "${GREEN}HTML report generated: coverage.html${NC}"
