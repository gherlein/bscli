#!/bin/bash

# BrightSign CLI Integration Test Runner
# This script runs comprehensive integration tests against a real BrightSign player

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}BrightSign CLI Integration Test Runner${NC}"
echo "======================================"

# Check if environment variables are set
if [[ -z "${BSCLI_TEST_HOST}" ]]; then
    echo -e "${RED}Error: BSCLI_TEST_HOST environment variable is required${NC}"
    echo "Example: export BSCLI_TEST_HOST=192.168.1.100"
    exit 1
fi

if [[ -z "${BSCLI_TEST_PASSWORD}" ]]; then
    echo -e "${RED}Error: BSCLI_TEST_PASSWORD environment variable is required${NC}"
    echo "Example: export BSCLI_TEST_PASSWORD=yourpassword"
    exit 1
fi

echo -e "${YELLOW}Test Configuration:${NC}"
echo "  Host: ${BSCLI_TEST_HOST}"
echo "  Username: ${BSCLI_TEST_USERNAME:-admin}"
echo "  Password: [REDACTED]"
echo ""

# Warning about running tests
echo -e "${YELLOW}WARNING:${NC} These tests will:"
echo "  - Connect to the BrightSign player"
echo "  - Create and delete test files on the SD card"
echo "  - Modify registry keys (temporary test keys)"
echo "  - Run network diagnostics"
echo ""

read -p "Continue with integration tests? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Tests cancelled."
    exit 0
fi

echo -e "${BLUE}Running Integration Tests...${NC}"
echo ""

# Change to test directory
cd "$(dirname "$0")"

# Run the Go integration tests
echo -e "${YELLOW}Starting Go integration tests...${NC}"

if go test -v -timeout 10m .; then
    echo ""
    echo -e "${GREEN}✅ All integration tests passed!${NC}"
else
    echo ""
    echo -e "${RED}❌ Some integration tests failed!${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}Integration test summary:${NC}"
echo "  - Tested all major CLI commands"
echo "  - Verified JSON output consistency"
echo "  - Tested file upload/download operations"
echo "  - Verified error handling"
echo "  - Tested registry operations"
echo ""
echo -e "${GREEN}Integration tests completed successfully!${NC}"