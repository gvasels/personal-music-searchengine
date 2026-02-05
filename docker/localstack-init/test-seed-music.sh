#!/bin/bash

# Test script for init-seed-music.sh
# Validates seed data in LocalStack DynamoDB
# TDD Red Phase: MUST FAIL if init-seed-music.sh has not been created yet
#
# Usage:
#   bash docker/localstack-init/test-seed-music.sh
#
# Prerequisites:
#   - LocalStack running on localhost:4566
#   - DynamoDB table "MusicLibrary" created (via init-aws.sh)
#
# Environment variables:
#   AWS_ENDPOINT  - LocalStack endpoint (default: http://localhost:4566)
#   AWS_REGION    - AWS region (default: us-east-1)

set -e

ENDPOINT="${AWS_ENDPOINT:-http://localhost:4566}"
AWS_REGION="${AWS_REGION:-us-east-1}"
TABLE="MusicLibrary"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SEED_SCRIPT="${SCRIPT_DIR}/init-seed-music.sh"
SEED_USER="USER#seed-user"
EXPECTED_COUNT=20
ERRORS=0

echo "=== Testing Seed Music Data ==="
echo "Endpoint: ${ENDPOINT}"
echo "Table:    ${TABLE}"
echo ""

# ---------------------------------------------------------------
# Test 1: Seed script exists
# ---------------------------------------------------------------
echo "Test 1: Checking seed script exists..."
if [ ! -f "${SEED_SCRIPT}" ]; then
    echo "FAIL: ${SEED_SCRIPT} does not exist"
    echo ""
    echo "=== 1 test(s) FAILED (seed script missing) ==="
    exit 1
fi

if [ ! -x "${SEED_SCRIPT}" ]; then
    echo "FAIL: ${SEED_SCRIPT} is not executable"
    echo ""
    echo "=== 1 test(s) FAILED (seed script not executable) ==="
    exit 1
fi

echo "PASS: Seed script exists and is executable"

# ---------------------------------------------------------------
# Test 2: Run seed script and verify item count = 20
# ---------------------------------------------------------------
echo ""
echo "Test 2: Running seed script and verifying item count..."
bash "${SEED_SCRIPT}"

COUNT=$(aws --endpoint-url="${ENDPOINT}" --region "${AWS_REGION}" \
    dynamodb scan \
    --table-name "${TABLE}" \
    --filter-expression "begins_with(PK, :pk)" \
    --expression-attribute-values '{":pk":{"S":"'"${SEED_USER}"'"}}' \
    --select COUNT \
    --query "Count" \
    --output text 2>/dev/null)

if [ "${COUNT}" != "${EXPECTED_COUNT}" ]; then
    echo "FAIL: Expected ${EXPECTED_COUNT} items with PK=${SEED_USER}, got ${COUNT}"
    ERRORS=$((ERRORS + 1))
else
    echo "PASS: Item count is ${EXPECTED_COUNT}"
fi

# ---------------------------------------------------------------
# Test 3: Idempotency - run seed again, count unchanged
# ---------------------------------------------------------------
echo ""
echo "Test 3: Testing idempotency (running seed script a second time)..."
bash "${SEED_SCRIPT}"

COUNT_AFTER=$(aws --endpoint-url="${ENDPOINT}" --region "${AWS_REGION}" \
    dynamodb scan \
    --table-name "${TABLE}" \
    --filter-expression "begins_with(PK, :pk)" \
    --expression-attribute-values '{":pk":{"S":"'"${SEED_USER}"'"}}' \
    --select COUNT \
    --query "Count" \
    --output text 2>/dev/null)

if [ "${COUNT_AFTER}" != "${EXPECTED_COUNT}" ]; then
    echo "FAIL: After second run, expected ${EXPECTED_COUNT} items, got ${COUNT_AFTER}"
    ERRORS=$((ERRORS + 1))
else
    echo "PASS: Idempotency verified (still ${EXPECTED_COUNT} items after re-run)"
fi

# ---------------------------------------------------------------
# Test 4: Verify a track has all expected fields
# ---------------------------------------------------------------
echo ""
echo "Test 4: Verifying track fields on TRACK#hello-001..."

ITEM=$(aws --endpoint-url="${ENDPOINT}" --region "${AWS_REGION}" \
    dynamodb get-item \
    --table-name "${TABLE}" \
    --key '{"PK":{"S":"'"${SEED_USER}"'"},"SK":{"S":"TRACK#hello-001"}}' \
    --output json 2>/dev/null)

if [ -z "${ITEM}" ] || [ "${ITEM}" = "{}" ]; then
    echo "FAIL: Track hello-001 not found in table"
    ERRORS=$((ERRORS + 1))
else
    for field in Title Artist Album Genre Year Duration; do
        if ! echo "${ITEM}" | grep -q "\"${field}\""; then
            echo "FAIL: Track hello-001 missing required field: ${field}"
            ERRORS=$((ERRORS + 1))
        else
            echo "  OK: Field '${field}' present"
        fi
    done
fi

# ---------------------------------------------------------------
# Results
# ---------------------------------------------------------------
echo ""
if [ ${ERRORS} -eq 0 ]; then
    echo "=== All seed data tests PASSED ==="
    exit 0
else
    echo "=== ${ERRORS} test(s) FAILED ==="
    exit 1
fi
