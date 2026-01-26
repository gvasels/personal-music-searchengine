#!/bin/bash
# migrate-playlist-visibility.sh - Add default visibility to existing playlists
#
# Usage: ./migrate-playlist-visibility.sh [--dry-run]
#
# This migration:
# - Scans all playlists in DynamoDB
# - Sets Visibility="private" for playlists without the attribute
# - Is idempotent (safe to run multiple times)

set -euo pipefail

# Configuration
AWS_PROFILE="${AWS_PROFILE:-gvasels-muza}"
AWS_REGION="${AWS_REGION:-us-east-1}"
DYNAMODB_TABLE="MusicLibrary"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TOTAL=0
MIGRATED=0
SKIPPED=0
ERRORS=0

# Parse arguments
DRY_RUN=false
for arg in "$@"; do
    case $arg in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
    esac
done

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

migrate_playlist() {
    local pk="$1"
    local sk="$2"

    TOTAL=$((TOTAL + 1))

    # Check if Visibility already exists
    local existing=$(aws dynamodb get-item \
        --table-name "$DYNAMODB_TABLE" \
        --profile "$AWS_PROFILE" \
        --region "$AWS_REGION" \
        --key "{\"PK\": {\"S\": \"$pk\"}, \"SK\": {\"S\": \"$sk\"}}" \
        --projection-expression "Visibility" \
        --query "Item.Visibility.S" \
        --output text 2>/dev/null || echo "")

    if [ -n "$existing" ] && [ "$existing" != "None" ] && [ "$existing" != "null" ]; then
        log_debug "Skipping $pk/$sk - already has Visibility: $existing"
        SKIPPED=$((SKIPPED + 1))
        return 0
    fi

    if [ "$DRY_RUN" = true ]; then
        log_info "[DRY-RUN] Would set Visibility=private for $pk/$sk"
        MIGRATED=$((MIGRATED + 1))
        return 0
    fi

    # Update the playlist with default visibility
    local current_time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    if aws dynamodb update-item \
        --table-name "$DYNAMODB_TABLE" \
        --profile "$AWS_PROFILE" \
        --region "$AWS_REGION" \
        --key "{\"PK\": {\"S\": \"$pk\"}, \"SK\": {\"S\": \"$sk\"}}" \
        --update-expression "SET Visibility = :vis, UpdatedAt = :updated" \
        --condition-expression "attribute_not_exists(Visibility)" \
        --expression-attribute-values "{\":vis\": {\"S\": \"private\"}, \":updated\": {\"S\": \"$current_time\"}}" \
        2>/dev/null; then
        log_info "Migrated $pk/$sk -> Visibility=private"
        MIGRATED=$((MIGRATED + 1))
    else
        # Condition failed (attribute exists now) - skip
        log_debug "Skipping $pk/$sk - condition check failed (already migrated?)"
        SKIPPED=$((SKIPPED + 1))
    fi
}

# Main
echo ""
log_info "=== Playlist Visibility Migration ==="
log_info "Table: $DYNAMODB_TABLE"
log_info "Region: $AWS_REGION"
log_info "Profile: $AWS_PROFILE"
if [ "$DRY_RUN" = true ]; then
    log_warn "DRY RUN MODE - No changes will be made"
fi
echo ""

# Scan for all playlists
log_info "Scanning for playlists..."

# Use pagination to handle large tables
last_key=""
while true; do
    if [ -z "$last_key" ]; then
        result=$(aws dynamodb scan \
            --table-name "$DYNAMODB_TABLE" \
            --profile "$AWS_PROFILE" \
            --region "$AWS_REGION" \
            --filter-expression "begins_with(SK, :sk_prefix)" \
            --expression-attribute-values "{\":sk_prefix\": {\"S\": \"PLAYLIST#\"}}" \
            --projection-expression "PK, SK" \
            --output json)
    else
        result=$(aws dynamodb scan \
            --table-name "$DYNAMODB_TABLE" \
            --profile "$AWS_PROFILE" \
            --region "$AWS_REGION" \
            --filter-expression "begins_with(SK, :sk_prefix)" \
            --expression-attribute-values "{\":sk_prefix\": {\"S\": \"PLAYLIST#\"}}" \
            --projection-expression "PK, SK" \
            --exclusive-start-key "$last_key" \
            --output json)
    fi

    # Process items
    items=$(echo "$result" | jq -c '.Items[]' 2>/dev/null || echo "")

    while IFS= read -r item; do
        [ -z "$item" ] && continue
        pk=$(echo "$item" | jq -r '.PK.S')
        sk=$(echo "$item" | jq -r '.SK.S')
        migrate_playlist "$pk" "$sk"
    done <<< "$items"

    # Check for more pages
    last_key=$(echo "$result" | jq -c '.LastEvaluatedKey // empty')
    if [ -z "$last_key" ]; then
        break
    fi
done

# Summary
echo ""
log_info "=== Migration Summary ==="
log_info "Total playlists found: $TOTAL"
log_info "Migrated: $MIGRATED"
log_info "Skipped (already had visibility): $SKIPPED"
if [ "$ERRORS" -gt 0 ]; then
    log_error "Errors: $ERRORS"
fi
echo ""

if [ "$DRY_RUN" = true ]; then
    log_warn "This was a DRY RUN. Run without --dry-run to apply changes."
else
    log_info "Migration complete!"
fi
