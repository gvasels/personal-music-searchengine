#!/bin/bash
# Seeds mock music data into DynamoDB for hello-world local dev
set -e

LOCALSTACK_HOST="localhost"
AWS_REGION="us-east-1"
DYNAMODB_TABLE="MusicLibrary"
SEED_USER="seed-user"

# Check if seed data already exists
EXISTING=$(aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb query \
    --table-name ${DYNAMODB_TABLE} \
    --key-condition-expression "PK = :pk" \
    --expression-attribute-values '{":pk":{"S":"USER#seed-user"}}' \
    --select COUNT \
    --region ${AWS_REGION} 2>/dev/null | grep -o '"Count":[0-9]*' | grep -o '[0-9]*')

if [ "${EXISTING}" -ge 20 ] 2>/dev/null; then
    echo "Seed data already exists (${EXISTING} items). Skipping..."
    exit 0
fi

echo "Seeding mock music data..."

# Function to seed a track
seed_track() {
    local ID=$1 TITLE=$2 ARTIST=$3 ALBUM=$4 GENRE=$5 YEAR=$6 DURATION=$7

    aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb put-item \
        --table-name ${DYNAMODB_TABLE} \
        --item "{
            \"PK\": {\"S\": \"USER#${SEED_USER}\"},
            \"SK\": {\"S\": \"TRACK#${ID}\"},
            \"Type\": {\"S\": \"TRACK\"},
            \"id\": {\"S\": \"${ID}\"},
            \"userId\": {\"S\": \"${SEED_USER}\"},
            \"title\": {\"S\": \"${TITLE}\"},
            \"artist\": {\"S\": \"${ARTIST}\"},
            \"album\": {\"S\": \"${ALBUM}\"},
            \"genre\": {\"S\": \"${GENRE}\"},
            \"year\": {\"N\": \"${YEAR}\"},
            \"duration\": {\"N\": \"${DURATION}\"},
            \"format\": {\"S\": \"MP3\"},
            \"fileSize\": {\"N\": \"5242880\"},
            \"s3Key\": {\"S\": \"uploads/${SEED_USER}/${ID}.mp3\"},
            \"Visibility\": {\"S\": \"public\"},
            \"playCount\": {\"N\": \"0\"},
            \"createdAt\": {\"S\": \"2024-01-15T00:00:00Z\"},
            \"updatedAt\": {\"S\": \"2024-01-15T00:00:00Z\"}
        }" \
        --region ${AWS_REGION}
}

# Luna Waves - Electronic
seed_track "track-001" "Midnight Drift" "Luna Waves" "Waveforms" "Electronic" "2024" "240"
seed_track "track-002" "Solar Wind" "Luna Waves" "Waveforms" "Electronic" "2024" "195"
seed_track "track-003" "Neon Pulse" "Luna Waves" "Waveforms" "Electronic" "2024" "210"
seed_track "track-004" "Crystal Shore" "Luna Waves" "Waveforms" "Electronic" "2024" "268"

# The Ember Collective - Indie Rock
seed_track "track-005" "Burning Ground" "The Ember Collective" "Aftermath" "Indie Rock" "2023" "198"
seed_track "track-006" "Ashfall" "The Ember Collective" "Aftermath" "Indie Rock" "2023" "225"
seed_track "track-007" "Smoke Signals" "The Ember Collective" "Aftermath" "Indie Rock" "2023" "187"
seed_track "track-008" "Wildfire" "The Ember Collective" "Aftermath" "Indie Rock" "2023" "312"

# DJ Phantom - House
seed_track "track-009" "Ghost Protocol" "DJ Phantom" "Spectre" "House" "2025" "330"
seed_track "track-010" "Phantom Zone" "DJ Phantom" "Spectre" "House" "2025" "285"
seed_track "track-011" "Shadow Drop" "DJ Phantom" "Spectre" "House" "2025" "300"
seed_track "track-012" "Dark Frequency" "DJ Phantom" "Spectre" "House" "2025" "345"

# Aria Chen - Classical Crossover
seed_track "track-013" "Silk Road" "Aria Chen" "Eastern Wind" "Classical Crossover" "2022" "264"
seed_track "track-014" "Paper Crane" "Aria Chen" "Eastern Wind" "Classical Crossover" "2022" "198"
seed_track "track-015" "Jade Garden" "Aria Chen" "Eastern Wind" "Classical Crossover" "2022" "230"
seed_track "track-016" "Moonlight Sonata" "Aria Chen" "Eastern Wind" "Classical Crossover" "2022" "356"

# Voltage - Techno
seed_track "track-017" "Circuit Break" "Voltage" "Overclocked" "Techno" "2024" "290"
seed_track "track-018" "Overload" "Voltage" "Overclocked" "Techno" "2024" "310"
seed_track "track-019" "Wired" "Voltage" "Overclocked" "Techno" "2024" "275"
seed_track "track-020" "Blackout" "Voltage" "Overclocked" "Techno" "2024" "340"

echo "Seeded 20 mock tracks across 5 artists."
