#!/bin/bash

# TDD Red Phase: This test validates seed data from init-seed-music.sh.
# Expected to FAIL because init-seed-music.sh does not exist yet.
# Once init-seed-music.sh is implemented (Green phase), all tests should pass.

set -e

################################################################################
# Constants
################################################################################
LOCALSTACK_HOST="localhost"
AWS_REGION="us-east-1"
TABLE="MusicLibrary"
SEED_USER_PK="USER#seed-user"
EXPECTED_COUNT=20

ENDPOINT_URL="http://${LOCALSTACK_HOST}:4566"
AWS_CMD="aws --endpoint-url=${ENDPOINT_URL} --region ${AWS_REGION}"

PASS_COUNT=0
FAIL_COUNT=0
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

################################################################################
# Helpers
################################################################################
pass() {
    echo "  PASS: $1"
    PASS_COUNT=$((PASS_COUNT + 1))
}

fail() {
    echo "  FAIL: $1"
    FAIL_COUNT=$((FAIL_COUNT + 1))
}

query_seed_items() {
    ${AWS_CMD} dynamodb query \
        --table-name "${TABLE}" \
        --key-condition-expression "PK = :pk AND begins_with(SK, :sk)" \
        --expression-attribute-values '{
            ":pk": {"S": "'"${SEED_USER_PK}"'"},
            ":sk": {"S": "TRACK#"}
        }' \
        --output json
}

################################################################################
# Pre-flight: Verify LocalStack is reachable
################################################################################
echo ""
echo "============================================"
echo "Seed Data Tests"
echo "============================================"
echo ""

echo "Checking LocalStack connectivity..."
if ! ${AWS_CMD} dynamodb list-tables > /dev/null 2>&1; then
    echo "FAIL: Cannot connect to LocalStack at ${ENDPOINT_URL}"
    echo "Make sure LocalStack is running (docker-compose up -d)"
    exit 1
fi
echo "LocalStack is reachable."
echo ""

################################################################################
# Test 1: Seed data count is exactly 20
################################################################################
echo "Test 1: Verify seed data count is ${EXPECTED_COUNT}"

QUERY_RESULT=$(query_seed_items)
ACTUAL_COUNT=$(echo "${QUERY_RESULT}" | python3 -c "import sys,json; print(json.load(sys.stdin)['Count'])")

if [ "${ACTUAL_COUNT}" -eq "${EXPECTED_COUNT}" ]; then
    pass "Found ${ACTUAL_COUNT} seed tracks (expected ${EXPECTED_COUNT})"
else
    fail "Expected ${EXPECTED_COUNT} seed tracks, found ${ACTUAL_COUNT}"
fi

echo ""

################################################################################
# Test 2: Each item has required attributes
################################################################################
echo "Test 2: Verify required attributes on every seed item"

REQUIRED_ATTRS=("Title" "Artist" "Album" "Genre" "Year" "Duration")

ITEMS_JSON=$(echo "${QUERY_RESULT}" | python3 -c "
import sys, json
data = json.load(sys.stdin)
items = data.get('Items', [])
for i, item in enumerate(items):
    attrs = list(item.keys())
    print(json.dumps({'index': i, 'attrs': attrs, 'sk': item.get('SK', {}).get('S', 'unknown')}))
")

ALL_ATTRS_OK=true
while IFS= read -r line; do
    ITEM_INDEX=$(echo "${line}" | python3 -c "import sys,json; print(json.load(sys.stdin)['index'])")
    ITEM_SK=$(echo "${line}" | python3 -c "import sys,json; print(json.load(sys.stdin)['sk'])")
    ITEM_ATTRS=$(echo "${line}" | python3 -c "import sys,json; print(' '.join(json.load(sys.stdin)['attrs']))")

    for attr in "${REQUIRED_ATTRS[@]}"; do
        if ! echo "${ITEM_ATTRS}" | grep -qw "${attr}"; then
            fail "Item ${ITEM_SK} is missing required attribute: ${attr}"
            ALL_ATTRS_OK=false
        fi
    done
done <<< "${ITEMS_JSON}"

if [ "${ALL_ATTRS_OK}" = true ]; then
    pass "All ${ACTUAL_COUNT} items have required attributes: ${REQUIRED_ATTRS[*]}"
fi

echo ""

################################################################################
# Test 3: Idempotency - running seed script again keeps count at 20
################################################################################
echo "Test 3: Verify idempotency (re-running seed script keeps count at ${EXPECTED_COUNT})"

SEED_SCRIPT="${SCRIPT_DIR}/init-seed-music.sh"

if [ ! -f "${SEED_SCRIPT}" ]; then
    fail "Seed script not found at ${SEED_SCRIPT} (expected for Red phase)"
else
    echo "  Running seed script again..."
    bash "${SEED_SCRIPT}"

    QUERY_RESULT_AFTER=$(query_seed_items)
    COUNT_AFTER=$(echo "${QUERY_RESULT_AFTER}" | python3 -c "import sys,json; print(json.load(sys.stdin)['Count'])")

    if [ "${COUNT_AFTER}" -eq "${EXPECTED_COUNT}" ]; then
        pass "After re-run, count is still ${COUNT_AFTER} (idempotent)"
    else
        fail "After re-run, count is ${COUNT_AFTER} (expected ${EXPECTED_COUNT}, not idempotent)"
    fi
fi

echo ""

################################################################################
# Summary
################################################################################
echo "============================================"
echo "Results: ${PASS_COUNT} passed, ${FAIL_COUNT} failed"
echo "============================================"
echo ""

if [ "${FAIL_COUNT}" -gt 0 ]; then
    echo "FAILED"
    exit 1
fi

echo "ALL TESTS PASSED"
exit 0
