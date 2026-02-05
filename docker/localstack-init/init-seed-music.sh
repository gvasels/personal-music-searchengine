#!/bin/bash

# Seeds 20 mock tracks into LocalStack DynamoDB
# Idempotent: uses condition-expression to skip existing items
#
# Usage:
#   bash docker/localstack-init/init-seed-music.sh
#
# Prerequisites:
#   - LocalStack running on localhost:4566
#   - DynamoDB table "MusicLibrary" created (via init-aws.sh)
#
# Environment variables:
#   AWS_ENDPOINT  - LocalStack endpoint (default: http://localhost:4566)
#   AWS_DEFAULT_REGION - AWS region (default: us-east-1)

set -e

ENDPOINT="${AWS_ENDPOINT:-http://localhost:4566}"
TABLE="MusicLibrary"
USER_ID="seed-user"
REGION="${AWS_DEFAULT_REGION:-us-east-1}"

echo "Seeding music data into ${TABLE}..."

# Helper function to insert a single track (idempotent)
put_track() {
    local id="$1" title="$2" artist="$3" album="$4" genre="$5" year="$6" duration="$7"
    aws --endpoint-url="${ENDPOINT}" --region "${REGION}" dynamodb put-item \
        --table-name "${TABLE}" \
        --item "{
            \"PK\": {\"S\": \"USER#${USER_ID}\"},
            \"SK\": {\"S\": \"TRACK#${id}\"},
            \"TrackID\": {\"S\": \"${id}\"},
            \"UserID\": {\"S\": \"${USER_ID}\"},
            \"Title\": {\"S\": \"${title}\"},
            \"Artist\": {\"S\": \"${artist}\"},
            \"Album\": {\"S\": \"${album}\"},
            \"Genre\": {\"S\": \"${genre}\"},
            \"Year\": {\"N\": \"${year}\"},
            \"Duration\": {\"N\": \"${duration}\"},
            \"CreatedAt\": {\"S\": \"2024-01-01T00:00:00Z\"},
            \"UpdatedAt\": {\"S\": \"2024-01-01T00:00:00Z\"},
            \"Visibility\": {\"S\": \"public\"}
        }" \
        --condition-expression "attribute_not_exists(PK)" 2>/dev/null || true
}

# Aurora Waves - Neon Horizons (Synthwave, 2024)
put_track "hello-001" "Neon Dreams" "Aurora Waves" "Neon Horizons" "Synthwave" "2024" "234"
put_track "hello-002" "Electric Sunset" "Aurora Waves" "Neon Horizons" "Synthwave" "2024" "198"
put_track "hello-003" "Midnight Circuit" "Aurora Waves" "Neon Horizons" "Synthwave" "2024" "276"
put_track "hello-004" "Chrome Highways" "Aurora Waves" "Neon Horizons" "Synthwave" "2024" "312"

# The Midnight Echo - Velvet Thunder (Indie Rock, 2023)
put_track "hello-005" "Velvet Storm" "The Midnight Echo" "Velvet Thunder" "Indie Rock" "2023" "245"
put_track "hello-006" "Thunder Road" "The Midnight Echo" "Velvet Thunder" "Indie Rock" "2023" "267"
put_track "hello-007" "Echo Chamber" "The Midnight Echo" "Velvet Thunder" "Indie Rock" "2023" "189"
put_track "hello-008" "Midnight Rain" "The Midnight Echo" "Velvet Thunder" "Indie Rock" "2023" "301"

# Luna Noir - Shadow Dance (Electronic, 2024)
put_track "hello-009" "Shadow Pulse" "Luna Noir" "Shadow Dance" "Electronic" "2024" "223"
put_track "hello-010" "Dark Matter" "Luna Noir" "Shadow Dance" "Electronic" "2024" "256"
put_track "hello-011" "Lunar Eclipse" "Luna Noir" "Shadow Dance" "Electronic" "2024" "287"
put_track "hello-012" "Night Vision" "Luna Noir" "Shadow Dance" "Electronic" "2024" "198"

# Crimson Tide - Ocean Drive (Funk, 2023)
put_track "hello-013" "Tidal Wave" "Crimson Tide" "Ocean Drive" "Funk" "2023" "234"
put_track "hello-014" "Coral Reef" "Crimson Tide" "Ocean Drive" "Funk" "2023" "278"
put_track "hello-015" "Deep Current" "Crimson Tide" "Ocean Drive" "Funk" "2023" "312"
put_track "hello-016" "Sunset Harbor" "Crimson Tide" "Ocean Drive" "Funk" "2023" "245"

# Stellar Drift - Cosmic Bloom (Ambient, 2024)
put_track "hello-017" "Nebula Rising" "Stellar Drift" "Cosmic Bloom" "Ambient" "2024" "345"
put_track "hello-018" "Starfield" "Stellar Drift" "Cosmic Bloom" "Ambient" "2024" "412"
put_track "hello-019" "Galaxy Spin" "Stellar Drift" "Cosmic Bloom" "Ambient" "2024" "267"
put_track "hello-020" "Cosmic Dust" "Stellar Drift" "Cosmic Bloom" "Ambient" "2024" "289"

echo "Seed data complete! Inserted 20 tracks for USER#${USER_ID}"
