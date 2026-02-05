#!/bin/bash

# Seed DynamoDB with 20 sample music tracks for the hello-world search demo.
# Uses put-item for idempotency — re-running this script will overwrite
# existing items with the same keys rather than creating duplicates.

set -e

echo "Seeding DynamoDB with sample music tracks..."

LOCALSTACK_HOST="localhost"
AWS_REGION="us-east-1"
TABLE="MusicLibrary"

put_track() {
    local id="$1" title="$2" artist="$3" album="$4" genre="$5" year="$6" duration="$7"
    aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb put-item \
        --table-name "${TABLE}" \
        --item '{
            "PK": {"S": "USER#seed-user"},
            "SK": {"S": "TRACK#seed-'"${id}"'"},
            "Title": {"S": "'"${title}"'"},
            "Artist": {"S": "'"${artist}"'"},
            "Album": {"S": "'"${album}"'"},
            "Genre": {"S": "'"${genre}"'"},
            "Year": {"N": "'"${year}"'"},
            "Duration": {"N": "'"${duration}"'"}
        }' \
        --region ${AWS_REGION}
}

put_track "t1"  "Aurora Borealis"      "Midnight Echo"          "Northern Lights"      "Electronic" "2023" "245"
put_track "t2"  "Velvet Dreams"        "Luna Wave"              "Dreamscape"           "Jazz"       "2022" "312"
put_track "t3"  "Crimson Tide"         "Storm Riders"           "Ocean Storm"          "Rock"       "2021" "198"
put_track "t4"  "Silver Lining"        "Midnight Echo"          "Northern Lights"      "Electronic" "2023" "267"
put_track "t5"  "Neon Pulse"           "Cyber Drift"            "Digital Dawn"         "Synthwave"  "2024" "220"
put_track "t6"  "Whisper in the Wind"  "Folk Tales"             "Countryside"          "Folk"       "2020" "185"
put_track "t7"  "Jazz Nocturne"        "Blue Note Trio"         "Late Night Sessions"  "Jazz"       "2019" "340"
put_track "t8"  "Thunderstrike"        "Storm Riders"           "Ocean Storm"          "Rock"       "2021" "210"
put_track "t9"  "Celestial Aurora"     "Stargazer"              "Cosmos"               "Ambient"    "2023" "420"
put_track "t10" "Midnight Serenade"    "Luna Wave"              "Dreamscape"           "Jazz"       "2022" "290"
put_track "t11" "Pixel Storm"          "Cyber Drift"            "Digital Dawn"         "Synthwave"  "2024" "195"
put_track "t12" "Desert Mirage"        "Sand Nomads"            "Sahara"               "World"      "2022" "275"
put_track "t13" "Frozen Lake"          "Arctic Sound"           "Polar"                "Ambient"    "2023" "360"
put_track "t14" "Funky Grooves"        "Bass Kitchen"           "Cookout"              "Funk"       "2021" "230"
put_track "t15" "Rainy Day Blues"      "Blue Note Trio"         "Late Night Sessions"  "Jazz"       "2019" "305"
put_track "t16" "Electric Sunrise"     "Neon Flux"              "Voltage"              "Electronic" "2024" "215"
put_track "t17" "Mountain Echo"        "Folk Tales"             "Countryside"          "Folk"       "2020" "200"
put_track "t18" "Deep Blue Ocean"      "Aqua Sound"             "Depths"               "Ambient"    "2022" "380"
put_track "t19" "City Lights"          "Urban Jazz Collective"  "Metropolitan"         "Jazz"       "2023" "255"
put_track "t20" "Solar Flare"          "Stargazer"              "Cosmos"               "Electronic" "2023" "235"

echo "Seeded 20 tracks into ${TABLE} (PK=USER#seed-user)."
