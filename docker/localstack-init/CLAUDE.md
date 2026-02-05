# LocalStack Init Scripts - CLAUDE.md

## Overview

Shell scripts that initialize LocalStack AWS resources for local development. Run automatically by Docker Compose on container startup, or manually via `make local`. Scripts are idempotent and safe to re-run.

## File Descriptions

| File | Purpose |
|------|---------|
| `init-aws.sh` | Creates DynamoDB table (`MusicLibrary`) and S3 bucket (`music-library-local-media`) with CORS and folder structure |
| `init-cognito.sh` | Creates Cognito user pool, app client, groups (admin/artist/subscriber), and three test users |
| `init-seed-music.sh` | Seeds DynamoDB with 20 sample tracks under `USER#seed-user` for the hello-world demo |
| `test-seed-music.sh` | Validates seed data: count = 20, required attributes present, idempotency check |

## Resources Created

### init-aws.sh

| Resource | Name | Details |
|----------|------|---------|
| DynamoDB Table | `MusicLibrary` | PK/SK + GSI1, GSI2, GSI3; pay-per-request |
| S3 Bucket | `music-library-local-media` | CORS for localhost:5173/3000; folders: `uploads/`, `media/`, `covers/` |

### init-cognito.sh

| Resource | Name | Details |
|----------|------|---------|
| User Pool | `music-library-local-pool` | Email sign-in, password policy (8+ chars, mixed case, numbers) |
| App Client | `music-library-local-client` | No secret (SPA), USER_PASSWORD_AUTH + USER_SRP_AUTH |
| Groups | `admin`, `artist`, `subscriber` | Role-based access groups |
| Test Users | `admin@local.test`, `subscriber@local.test`, `artist@local.test` | Password: `LocalTest123!` |

Config is saved to `/tmp/localstack-cognito-config.env` with `USER_POOL_ID` and `CLIENT_ID`.

### init-seed-music.sh

Seeds 20 tracks with attributes: `PK`, `SK`, `Title`, `Artist`, `Album`, `Genre`, `Year`, `Duration`. All tracks use `PK=USER#seed-user` and `SK=TRACK#seed-{id}`. Uses `put-item` for idempotency.

## Key Functions

| Function | File | Description |
|----------|------|-------------|
| `put_track(id, title, artist, album, genre, year, duration)` | `init-seed-music.sh` | Inserts a single track into DynamoDB |
| `create_test_user(email, group)` | `init-cognito.sh` | Creates a Cognito user and adds to a group (idempotent) |
| `query_seed_items()` | `test-seed-music.sh` | Queries DynamoDB for all seed tracks |
| `pass(msg)` / `fail(msg)` | `test-seed-music.sh` | Test result reporting helpers |

## Dependencies

- AWS CLI v2 (available in LocalStack container)
- `python3` (used by test-seed-music.sh for JSON parsing)
- LocalStack running on `localhost:4566`
- `AWS_REQUEST_CHECKSUM_CALCULATION=WHEN_REQUIRED` to work around AWS CLI v2 checksum trailers

## Usage

```bash
# Run automatically via Docker Compose
docker-compose up -d

# Run manually
bash docker/localstack-init/init-aws.sh
bash docker/localstack-init/init-cognito.sh
bash docker/localstack-init/init-seed-music.sh

# Validate seed data
bash docker/localstack-init/test-seed-music.sh
```
