# LocalStack Initialization Scripts - CLAUDE.md

## Overview

Initialization scripts that run when LocalStack container starts. Creates AWS resources (DynamoDB, S3, Cognito) and seeds development data.

## File Descriptions

| File | Purpose |
|------|---------|
| `init-aws.sh` | Creates DynamoDB table (`MusicLibrary`) and S3 bucket (`music-library-local-media`) |
| `init-cognito.sh` | Creates Cognito user pool, app client, groups, and test users |
| `init-seed-music.sh` | Seeds 20 tracks for the Hello World feature (5 artists x 4 tracks) |

## init-seed-music.sh

Seeds DynamoDB with sample music data for local development testing.

### Data Created

**20 tracks total:** 5 artists with 4 tracks each

| Artist | Album(s) | Genre | Years |
|--------|----------|-------|-------|
| Aurora Waves | Dreamscape, Polar Nights | jazz | 2022-2023 |
| Electric Pulse | Circuit Breaker, Synth City | electronic | 2023-2024 |
| Stone Temple | Carved in Time, Seismic | rock | 2021-2022 |
| Velvet Soul | Silky Nights, Warm Embrace | soul | 2023-2024 |
| Ambient Dreams | Ethereal Skies, Nature Sounds | ambient | 2022-2023 |

### Track IDs

All tracks use stable IDs (`seed-t1` through `seed-t20`) for idempotent seeding.

### DynamoDB Schema

Tracks are stored with:
- `PK`: `USER#seed-user`
- `SK`: `TRACK#{id}` (e.g., `TRACK#seed-t1`)

### Key Functions

```bash
# put_track helper function
put_track() {
    local id="$1"
    local title="$2"
    local artist="$3"
    local album="$4"
    local genre="$5"
    local year="$6"
    local duration="$7"
    # Uses aws dynamodb put-item to LocalStack endpoint
}
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOCALSTACK_HOST` | `localhost` | LocalStack hostname |
| `AWS_ACCESS_KEY_ID` | `test` | Dummy credential for LocalStack |
| `AWS_SECRET_ACCESS_KEY` | `test` | Dummy credential for LocalStack |
| `AWS_DEFAULT_REGION` | `us-east-1` | AWS region |

### Usage

```bash
# Run directly
./docker/localstack-init/init-seed-music.sh

# Or via make
make local  # Runs all init scripts including seed
```

### Idempotency

The script is idempotent - running it multiple times overwrites existing data with the same IDs. No duplicate tracks are created.

## Dependencies

- AWS CLI installed
- LocalStack running on port 4566
- DynamoDB table `MusicLibrary` created (by `init-aws.sh`)

## Execution Order

Scripts run in alphabetical order when LocalStack starts:
1. `init-aws.sh` - Creates table and bucket
2. `init-cognito.sh` - Creates auth resources
3. `init-seed-music.sh` - Seeds music data
