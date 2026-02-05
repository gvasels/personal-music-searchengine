# LocalStack Init Scripts - CLAUDE.md

## Overview

Bash initialization scripts that run when the LocalStack Docker container starts. They create AWS resources (DynamoDB, S3, Cognito) and seed test data for local development. Scripts are idempotent and can be re-run safely.

## Files

| File | Purpose |
|------|---------|
| `init-aws.sh` | Creates DynamoDB table (`MusicLibrary`) with 3 GSIs, S3 bucket (`music-library-local-media`) with CORS and folder structure |
| `init-cognito.sh` | Creates Cognito user pool, app client, groups (admin/artist/subscriber), and 3 test users |
| `init-seed-music.sh` | Seeds 20 mock tracks into DynamoDB for hello-world search validation |
| `test-seed-music.sh` | Bash test script that validates seed data (count, idempotency, field presence) |

## Key Functions/Exports

### init-aws.sh
- Creates DynamoDB table `MusicLibrary` with PK/SK and GSI1-GSI3
- Creates S3 bucket `music-library-local-media` with CORS for `localhost:5173` and `localhost:3000`
- Creates `uploads/`, `media/`, `covers/` folder placeholders in S3
- Disables AWS CLI v2 checksum trailers (`AWS_REQUEST_CHECKSUM_CALCULATION=WHEN_REQUIRED`) for LocalStack compatibility

### init-cognito.sh
- `create_test_user(email, group)` - Creates a Cognito user and adds to specified group (idempotent)
- Creates user pool `music-library-local-pool` with email-based sign-in
- Creates app client `music-library-local-client` (no secret, for SPA)
- Writes config to `/tmp/localstack-cognito-config.env` (USER_POOL_ID, CLIENT_ID)

### init-seed-music.sh
- `put_track(id, title, artist, album, genre, year, duration)` - Inserts a single track with `condition-expression "attribute_not_exists(PK)"` for idempotency
- Seeds 20 tracks across 5 artists/albums:
  - Aurora Waves / Neon Horizons (Synthwave, 2024) - 4 tracks
  - The Midnight Echo / Velvet Thunder (Indie Rock, 2023) - 4 tracks
  - Luna Noir / Shadow Dance (Electronic, 2024) - 4 tracks
  - Crimson Tide / Ocean Drive (Funk, 2023) - 4 tracks
  - Stellar Drift / Cosmic Bloom (Ambient, 2024) - 4 tracks
- All tracks use `PK=USER#seed-user`, `SK=TRACK#hello-{001..020}`, `Visibility=public`

### test-seed-music.sh
- **Test 1**: Verifies `init-seed-music.sh` exists and is executable
- **Test 2**: Runs seed script, verifies exactly 20 items with `PK=USER#seed-user`
- **Test 3**: Runs seed script again, verifies count unchanged (idempotency)
- **Test 4**: Fetches `TRACK#hello-001` and validates all required fields (Title, Artist, Album, Genre, Year, Duration)

## Dependencies

- AWS CLI v2 (available in LocalStack container and host)
- LocalStack running on `localhost:4566` (or `$AWS_ENDPOINT`)
- DynamoDB table `MusicLibrary` must exist before running seed script (created by `init-aws.sh`)

## Environment Variables

| Variable | Default | Used By |
|----------|---------|---------|
| `AWS_ENDPOINT` | `http://localhost:4566` | `init-seed-music.sh`, `test-seed-music.sh` |
| `AWS_DEFAULT_REGION` | `us-east-1` | `init-seed-music.sh` |
| `AWS_REGION` | `us-east-1` | `test-seed-music.sh` |
| `LOCALSTACK_HOST` | `localhost` | `init-aws.sh`, `init-cognito.sh` |

## Usage Examples

```bash
# Run all init scripts in order
bash docker/localstack-init/init-aws.sh
bash docker/localstack-init/init-cognito.sh
bash docker/localstack-init/init-seed-music.sh

# Validate seed data
bash docker/localstack-init/test-seed-music.sh

# Run with custom endpoint
AWS_ENDPOINT=http://localstack:4566 bash docker/localstack-init/init-seed-music.sh

# Verify seed data manually
aws --endpoint-url=http://localhost:4566 dynamodb scan \
  --table-name MusicLibrary \
  --filter-expression "begins_with(PK, :pk)" \
  --expression-attribute-values '{":pk":{"S":"USER#seed-user"}}' \
  --select COUNT
```
