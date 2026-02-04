#!/bin/bash
# migrate-track-visibility.sh - Add default visibility to existing tracks
#
# Usage: ./migrate-track-visibility.sh [--dry-run]
#
# This script sets Visibility="private" for all existing tracks that don't have
# a visibility field set. Safe to run multiple times (idempotent).
#
# Prerequisites:
# - AWS CLI v2 configured with appropriate permissions
# - jq installed for JSON processing
# - AWS_PROFILE or credentials configured

set -euo pipefail

# Configuration
TABLE_NAME="${DYNAMODB_TABLE_NAME:-MusicLibrary}"
AWS_REGION="${AWS_REGION:-us-east-1}"
DEFAULT_VISIBILITY="private"
BATCH_SIZE=25
DRY_RUN=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    --table)
      TABLE_NAME="$2"
      shift 2
      ;;
    --region)
      AWS_REGION="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--dry-run] [--table TABLE_NAME] [--region REGION]"
      exit 1
      ;;
  esac
done

echo "=== Track Visibility Migration ==="
echo "Table: $TABLE_NAME"
echo "Region: $AWS_REGION"
echo "Default Visibility: $DEFAULT_VISIBILITY"
echo "Dry Run: $DRY_RUN"
echo ""

# Check for required tools
if ! command -v aws &> /dev/null; then
  echo "Error: AWS CLI is required but not installed."
  exit 1
fi

if ! command -v jq &> /dev/null; then
  echo "Error: jq is required but not installed."
  exit 1
fi

# Counters
TOTAL_SCANNED=0
TOTAL_UPDATED=0
TOTAL_SKIPPED=0

# Function to update a single track
update_track() {
  local pk="$1"
  local sk="$2"

  if [ "$DRY_RUN" = true ]; then
    echo "  [DRY RUN] Would update: PK=$pk, SK=$sk"
    return 0
  fi

  aws dynamodb update-item \
    --table-name "$TABLE_NAME" \
    --region "$AWS_REGION" \
    --key "{\"PK\": {\"S\": \"$pk\"}, \"SK\": {\"S\": \"$sk\"}}" \
    --update-expression "SET #vis = :vis" \
    --condition-expression "attribute_not_exists(#vis)" \
    --expression-attribute-names '{"#vis": "Visibility"}' \
    --expression-attribute-values "{\":vis\": {\"S\": \"$DEFAULT_VISIBILITY\"}}" \
    2>/dev/null || return 1

  return 0
}

# Scan for tracks and update them
echo "Scanning for tracks..."
LAST_EVALUATED_KEY=""

while true; do
  # Build scan command
  SCAN_CMD="aws dynamodb scan \
    --table-name $TABLE_NAME \
    --region $AWS_REGION \
    --filter-expression 'begins_with(SK, :track_prefix)' \
    --expression-attribute-values '{\":track_prefix\": {\"S\": \"TRACK#\"}}' \
    --projection-expression 'PK, SK, Visibility' \
    --limit 100"

  if [ -n "$LAST_EVALUATED_KEY" ]; then
    SCAN_CMD="$SCAN_CMD --exclusive-start-key '$LAST_EVALUATED_KEY'"
  fi

  # Execute scan
  RESULT=$(eval $SCAN_CMD)

  # Process items
  ITEMS=$(echo "$RESULT" | jq -c '.Items[]' 2>/dev/null || echo "")

  if [ -z "$ITEMS" ]; then
    # Check if there are more pages
    LAST_EVALUATED_KEY=$(echo "$RESULT" | jq -r '.LastEvaluatedKey // empty')
    if [ -z "$LAST_EVALUATED_KEY" ]; then
      break
    fi
    continue
  fi

  # Process each track
  while IFS= read -r item; do
    ((TOTAL_SCANNED++))

    PK=$(echo "$item" | jq -r '.PK.S')
    SK=$(echo "$item" | jq -r '.SK.S')
    VISIBILITY=$(echo "$item" | jq -r '.Visibility.S // empty')

    if [ -n "$VISIBILITY" ]; then
      ((TOTAL_SKIPPED++))
      continue
    fi

    echo "Updating track: PK=$PK, SK=$SK"

    if update_track "$PK" "$SK"; then
      ((TOTAL_UPDATED++))
    else
      echo "  Warning: Failed to update or already has visibility"
      ((TOTAL_SKIPPED++))
    fi

    # Rate limiting
    if [ $((TOTAL_UPDATED % BATCH_SIZE)) -eq 0 ] && [ "$DRY_RUN" = false ]; then
      echo "  Processed $TOTAL_UPDATED tracks, pausing..."
      sleep 1
    fi

  done <<< "$ITEMS"

  # Check for more pages
  LAST_EVALUATED_KEY=$(echo "$RESULT" | jq -r '.LastEvaluatedKey // empty')
  if [ -z "$LAST_EVALUATED_KEY" ]; then
    break
  fi

  echo "Fetching next page..."
done

echo ""
echo "=== Migration Complete ==="
echo "Total Scanned: $TOTAL_SCANNED"
echo "Total Updated: $TOTAL_UPDATED"
echo "Total Skipped: $TOTAL_SKIPPED"

if [ "$DRY_RUN" = true ]; then
  echo ""
  echo "This was a dry run. No changes were made."
  echo "Run without --dry-run to apply changes."
fi
