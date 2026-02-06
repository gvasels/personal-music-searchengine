#!/bin/bash

# Seed music data for local development
# Creates 20 tracks: 5 artists × 4 tracks each
# Idempotent: same IDs each run, put-item overwrites existing

set -e

# Dummy AWS credentials for LocalStack (overrides any SSO profile)
export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-test}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-test}"
export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-us-east-1}"

echo "Seeding music data..."

LOCALSTACK_HOST="${LOCALSTACK_HOST:-localhost}"
AWS_REGION="us-east-1"
DYNAMODB_TABLE="MusicLibrary"
ENDPOINT="http://${LOCALSTACK_HOST}:4566"

# Helper function to put a track item
put_track() {
    local id="$1"
    local title="$2"
    local artist="$3"
    local album="$4"
    local genre="$5"
    local year="$6"
    local duration="$7"

    aws --endpoint-url="$ENDPOINT" dynamodb put-item \
        --table-name "$DYNAMODB_TABLE" \
        --item '{
            "PK": {"S": "USER#seed-user"},
            "SK": {"S": "TRACK#'"$id"'"},
            "Title": {"S": "'"$title"'"},
            "Artist": {"S": "'"$artist"'"},
            "Album": {"S": "'"$album"'"},
            "Genre": {"S": "'"$genre"'"},
            "Year": {"N": "'"$year"'"},
            "Duration": {"N": "'"$duration"'"}
        }' \
        --region "$AWS_REGION" > /dev/null 2>&1
}

# Aurora Waves - 4 tracks (jazz)
echo "  Adding Aurora Waves tracks..."
put_track "seed-t1"  "Aurora Borealis"     "Aurora Waves"   "Dreamscape"      "jazz"       "2023" "245"
put_track "seed-t2"  "Midnight Glow"       "Aurora Waves"   "Dreamscape"      "jazz"       "2023" "312"
put_track "seed-t3"  "Northern Lights"     "Aurora Waves"   "Polar Nights"    "jazz"       "2022" "278"
put_track "seed-t4"  "Dawn Chorus"         "Aurora Waves"   "Polar Nights"    "jazz"       "2022" "195"

# Electric Pulse - 4 tracks (electronic)
echo "  Adding Electric Pulse tracks..."
put_track "seed-t5"  "Voltage Drop"        "Electric Pulse" "Circuit Breaker" "electronic" "2024" "203"
put_track "seed-t6"  "Neon Dreams"         "Electric Pulse" "Circuit Breaker" "electronic" "2024" "267"
put_track "seed-t7"  "Digital Wave"        "Electric Pulse" "Synth City"      "electronic" "2023" "189"
put_track "seed-t8"  "Power Surge"         "Electric Pulse" "Synth City"      "electronic" "2023" "224"

# Stone Temple - 4 tracks (rock)
echo "  Adding Stone Temple tracks..."
put_track "seed-t9"  "Ancient Ruins"       "Stone Temple"   "Carved in Time"  "rock"       "2021" "356"
put_track "seed-t10" "Granite Heart"       "Stone Temple"   "Carved in Time"  "rock"       "2021" "289"
put_track "seed-t11" "Earthquake"          "Stone Temple"   "Seismic"         "rock"       "2022" "341"
put_track "seed-t12" "Bedrock Blues"       "Stone Temple"   "Seismic"         "rock"       "2022" "278"

# Velvet Soul - 4 tracks (soul)
echo "  Adding Velvet Soul tracks..."
put_track "seed-t13" "Smooth Operator"     "Velvet Soul"    "Silky Nights"    "soul"       "2023" "298"
put_track "seed-t14" "Satin Dreams"        "Velvet Soul"    "Silky Nights"    "soul"       "2023" "256"
put_track "seed-t15" "Tender Touch"        "Velvet Soul"    "Warm Embrace"    "soul"       "2024" "312"
put_track "seed-t16" "Midnight Serenade"   "Velvet Soul"    "Warm Embrace"    "soul"       "2024" "267"

# Ambient Dreams - 4 tracks (ambient)
echo "  Adding Ambient Dreams tracks..."
put_track "seed-t17" "Floating Clouds"     "Ambient Dreams" "Ethereal Skies"  "ambient"    "2022" "445"
put_track "seed-t18" "Ocean Waves"         "Ambient Dreams" "Ethereal Skies"  "ambient"    "2022" "512"
put_track "seed-t19" "Forest Rain"         "Ambient Dreams" "Nature Sounds"   "ambient"    "2023" "389"
put_track "seed-t20" "Mountain Mist"       "Ambient Dreams" "Nature Sounds"   "ambient"    "2023" "423"

echo ""
echo "Seed data complete: 20 tracks (5 artists × 4 tracks)"
echo ""
