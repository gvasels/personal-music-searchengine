#!/bin/bash
# Validates seed data in DynamoDB for hello-world local dev
set -e

LOCALSTACK_HOST="localhost"
AWS_REGION="us-east-1"
DYNAMODB_TABLE="MusicLibrary"
PASS=0
FAIL=0

assert_eq() {
    local desc=$1 expected=$2 actual=$3
    if [ "$expected" = "$actual" ]; then
        echo "  PASS: $desc"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $desc (expected=$expected, actual=$actual)"
        FAIL=$((FAIL + 1))
    fi
}

assert_ge() {
    local desc=$1 expected=$2 actual=$3
    if [ "$actual" -ge "$expected" ] 2>/dev/null; then
        echo "  PASS: $desc"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $desc (expected>=$expected, actual=$actual)"
        FAIL=$((FAIL + 1))
    fi
}

echo "=== Seed Data Validation ==="

# Test 1: Count items for seed-user
echo ""
echo "Test 1: Verify 20 items exist for seed-user"
COUNT=$(aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb query \
    --table-name ${DYNAMODB_TABLE} \
    --key-condition-expression "PK = :pk" \
    --expression-attribute-values '{":pk":{"S":"USER#seed-user"}}' \
    --select COUNT \
    --region ${AWS_REGION} 2>/dev/null | grep -o '"Count":[0-9]*' | grep -o '[0-9]*')
assert_eq "Table has 20 items for seed-user" "20" "$COUNT"

# Test 2: Verify required fields exist on first item
echo ""
echo "Test 2: Verify items have required fields"
ITEM=$(aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb get-item \
    --table-name ${DYNAMODB_TABLE} \
    --key '{"PK":{"S":"USER#seed-user"},"SK":{"S":"TRACK#track-001"}}' \
    --region ${AWS_REGION} 2>/dev/null)

for FIELD in title artist album genre year duration; do
    HAS_FIELD=$(echo "$ITEM" | grep -c "\"$FIELD\"" || true)
    assert_ge "Item has field: $FIELD" "1" "$HAS_FIELD"
done

# Test 3: Verify 5 distinct artists
echo ""
echo "Test 3: Verify 5 distinct artists"
ARTISTS=$(aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb query \
    --table-name ${DYNAMODB_TABLE} \
    --key-condition-expression "PK = :pk" \
    --expression-attribute-values '{":pk":{"S":"USER#seed-user"}}' \
    --projection-expression "artist" \
    --region ${AWS_REGION} 2>/dev/null | grep -o '"artist"' | wc -l | tr -d ' ')
UNIQUE_ARTISTS=$(aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb query \
    --table-name ${DYNAMODB_TABLE} \
    --key-condition-expression "PK = :pk" \
    --expression-attribute-values '{":pk":{"S":"USER#seed-user"}}' \
    --projection-expression "artist" \
    --region ${AWS_REGION} 2>/dev/null | grep -oP '"S":\s*"[^"]*"' | sort -u | wc -l | tr -d ' ')
assert_eq "5 distinct artists exist" "5" "$UNIQUE_ARTISTS"

# Test 4: Idempotency - run seed again and verify count unchanged
echo ""
echo "Test 4: Verify idempotency"
bash "$(dirname "$0")/init-seed-music.sh" > /dev/null 2>&1
COUNT_AFTER=$(aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb query \
    --table-name ${DYNAMODB_TABLE} \
    --key-condition-expression "PK = :pk" \
    --expression-attribute-values '{":pk":{"S":"USER#seed-user"}}' \
    --select COUNT \
    --region ${AWS_REGION} 2>/dev/null | grep -o '"Count":[0-9]*' | grep -o '[0-9]*')
assert_eq "Count unchanged after re-seed (idempotent)" "$COUNT" "$COUNT_AFTER"

# Summary
echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
