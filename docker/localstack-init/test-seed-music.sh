#!/bin/bash
# Test script for init-seed-music.sh
# Verifies that the seed script creates exactly 20 music tracks in DynamoDB
#
# Requirements tested:
# 1. 20 items exist in DynamoDB after running seed script
# 2. Items are under USER#seed-user partition key
# 3. Idempotency - running twice produces same 20 items, not 40
# 4. Each item has: Title, Artist, Album, Genre, Year, Duration

set -e

ENDPOINT="http://localhost:4566"
TABLE="MusicLibrary"
SEED_USER="USER#seed-user"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SEED_SCRIPT="${SCRIPT_DIR}/init-seed-music.sh"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "============================================"
echo "Testing init-seed-music.sh"
echo "============================================"
echo ""

# Check if seed script exists
echo -n "Checking if seed script exists... "
if [ ! -f "$SEED_SCRIPT" ]; then
    echo -e "${RED}FAIL${NC}"
    echo "Error: Seed script not found at: $SEED_SCRIPT"
    echo "This test expects init-seed-music.sh to exist."
    exit 1
fi
echo -e "${GREEN}OK${NC}"

# Check if seed script is executable
echo -n "Checking if seed script is executable... "
if [ ! -x "$SEED_SCRIPT" ]; then
    echo -e "${YELLOW}WARN${NC} - Making it executable"
    chmod +x "$SEED_SCRIPT"
fi
echo -e "${GREEN}OK${NC}"

# Run the seed script
echo ""
echo "Running seed script..."
"$SEED_SCRIPT"
if [ $? -ne 0 ]; then
    echo -e "${RED}FAIL${NC}: Seed script failed to execute"
    exit 1
fi
echo ""

# Function to query seed items count
query_seed_count() {
    aws --endpoint-url="$ENDPOINT" dynamodb query \
        --table-name "$TABLE" \
        --key-condition-expression "PK = :pk AND begins_with(SK, :sk)" \
        --expression-attribute-values '{":pk": {"S": "'"$SEED_USER"'"}, ":sk": {"S": "TRACK#"}}' \
        --select COUNT \
        --output json 2>/dev/null | jq -r '.Count // 0'
}

# Function to get all seed items
query_seed_items() {
    aws --endpoint-url="$ENDPOINT" dynamodb query \
        --table-name "$TABLE" \
        --key-condition-expression "PK = :pk AND begins_with(SK, :sk)" \
        --expression-attribute-values '{":pk": {"S": "'"$SEED_USER"'"}, ":sk": {"S": "TRACK#"}}' \
        --output json 2>/dev/null
}

# Test 1: Verify 20 items exist
echo "============================================"
echo "Test 1: Verify 20 items exist"
echo "============================================"
count=$(query_seed_count)
echo "Found $count items under $SEED_USER"
if [ "$count" -ne 20 ]; then
    echo -e "${RED}FAIL${NC}: Expected 20 items, found $count"
    exit 1
fi
echo -e "${GREEN}PASS${NC}: Found exactly 20 items"
echo ""

# Test 2: Verify items are under USER#seed-user partition key
echo "============================================"
echo "Test 2: Verify items are under USER#seed-user partition"
echo "============================================"
items_json=$(query_seed_items)
pk_check=$(echo "$items_json" | jq -r '.Items[].PK.S' | sort -u)
if [ "$pk_check" != "$SEED_USER" ]; then
    echo -e "${RED}FAIL${NC}: Items found with unexpected partition key"
    echo "Expected: $SEED_USER"
    echo "Found: $pk_check"
    exit 1
fi
echo -e "${GREEN}PASS${NC}: All items have PK=$SEED_USER"
echo ""

# Test 3: Verify each item has required attributes
echo "============================================"
echo "Test 3: Verify required attributes"
echo "============================================"
required_attrs=("Title" "Artist" "Album" "Genre" "Year" "Duration")
missing_attrs=0

for attr in "${required_attrs[@]}"; do
    # Check if all items have this attribute
    count_with_attr=$(echo "$items_json" | jq "[.Items[] | select(.${attr})] | length")
    if [ "$count_with_attr" -ne 20 ]; then
        echo -e "${RED}FAIL${NC}: Not all items have '$attr' attribute (found in $count_with_attr/20)"
        missing_attrs=$((missing_attrs + 1))
    else
        echo -e "${GREEN}PASS${NC}: All items have '$attr' attribute"
    fi
done

if [ $missing_attrs -gt 0 ]; then
    echo ""
    echo -e "${RED}FAIL${NC}: $missing_attrs required attributes missing from some items"
    exit 1
fi
echo ""

# Test 4: Verify idempotency - run seed script again
echo "============================================"
echo "Test 4: Verify idempotency"
echo "============================================"
echo "Running seed script again..."
"$SEED_SCRIPT"
if [ $? -ne 0 ]; then
    echo -e "${RED}FAIL${NC}: Seed script failed on second run"
    exit 1
fi

count_after=$(query_seed_count)
echo "Count after second run: $count_after"
if [ "$count_after" -ne 20 ]; then
    echo -e "${RED}FAIL${NC}: Expected 20 items after second run, found $count_after"
    echo "Seed script is not idempotent!"
    exit 1
fi
echo -e "${GREEN}PASS${NC}: Count remains 20 after second run (idempotent)"
echo ""

# Show sample data
echo "============================================"
echo "Sample seed data (first 3 items)"
echo "============================================"
echo "$items_json" | jq '.Items[:3][] | {
    SK: .SK.S,
    Title: .Title.S,
    Artist: .Artist.S,
    Album: .Album.S,
    Genre: .Genre.S,
    Year: .Year.N,
    Duration: .Duration.N
}'
echo ""

# All tests passed
echo "============================================"
echo -e "${GREEN}All tests passed!${NC}"
echo "============================================"
