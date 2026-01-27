# Docker - CLAUDE.md

## Overview

Docker configuration for local development using LocalStack to emulate AWS services (DynamoDB, S3, Cognito). Provides a complete local AWS environment for integration testing without requiring real AWS credentials.

## Directory Structure

```
docker/
├── docker-compose.yml       # LocalStack service definition
├── localstack-init/         # Initialization scripts
│   ├── init-aws.sh          # Creates DynamoDB table and S3 bucket
│   └── init-cognito.sh      # Creates Cognito user pool and test users
└── CLAUDE.md                # This file
```

## File Descriptions

| File | Purpose |
|------|---------|
| `docker-compose.yml` | Defines LocalStack container with DynamoDB, S3, and Cognito services |
| `localstack-init/init-aws.sh` | Initialization script that creates DynamoDB table and S3 bucket |
| `localstack-init/init-cognito.sh` | Creates Cognito user pool, app client, groups, and test users |

## Quick Start

```bash
# Start LocalStack
docker-compose up -d

# Check health
curl http://localhost:4566/_localstack/health

# View logs
docker-compose logs -f localstack

# Stop LocalStack
docker-compose down
```

## Services Configured

| Service | Port | Purpose |
|---------|------|---------|
| DynamoDB | 4566 | Database for tracks, albums, users, playlists |
| S3 | 4566 | Media storage for audio files and cover art |
| STS | 4566 | Security Token Service (for credentials) |
| IAM | 4566 | Identity and Access Management |
| Cognito | 4566 | User authentication and authorization |

## Resources Created by init-aws.sh

### DynamoDB Table: `MusicLibrary`
- Primary Key: `PK` (Hash), `SK` (Range)
- GSI1: `GSI1PK`, `GSI1SK` (for artist/tag queries)
- Billing: Pay-per-request

### S3 Bucket: `music-library-local-media`
- Folders: `uploads/`, `media/`, `covers/`
- CORS configured for localhost:5173 and localhost:3000

## Resources Created by init-cognito.sh

### Cognito User Pool: `music-library-local-pool`
- Email-based sign-in (no username)
- Password policy: 8+ characters, mixed case, numbers, symbols
- Email verification enabled (auto-verified in LocalStack)

### App Client: `music-library-local-client`
- No client secret (public SPA client)
- Auth flows: USER_PASSWORD_AUTH, USER_SRP_AUTH

### User Groups

| Group | Role | Purpose |
|-------|------|---------|
| `admin` | Admin | Full system access |
| `artist` | Artist | Can upload tracks, manage profile |
| `subscriber` | Subscriber | Can create playlists, follow artists |

### Test Users

| Email | Password | Role/Group |
|-------|----------|------------|
| `admin@local.test` | `LocalTest123!` | admin |
| `subscriber@local.test` | `LocalTest123!` | subscriber |
| `artist@local.test` | `LocalTest123!` | artist |

**Note**: All test users are pre-confirmed with verified email addresses.

## AWS CLI Usage

```bash
# List DynamoDB tables
aws --endpoint-url=http://localhost:4566 dynamodb list-tables

# Scan table
aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name MusicLibrary

# List S3 buckets
aws --endpoint-url=http://localhost:4566 s3 ls

# List bucket contents
aws --endpoint-url=http://localhost:4566 s3 ls s3://music-library-local-media/

# List Cognito user pools
aws --endpoint-url=http://localhost:4566 cognito-idp list-user-pools --max-results 10

# List users in pool
aws --endpoint-url=http://localhost:4566 cognito-idp list-users --user-pool-id <pool-id>

# List groups in pool
aws --endpoint-url=http://localhost:4566 cognito-idp list-groups --user-pool-id <pool-id>

# Get user info
aws --endpoint-url=http://localhost:4566 cognito-idp admin-get-user \
    --user-pool-id <pool-id> --username admin@local.test
```

## Data Persistence

LocalStack data is persisted to a Docker volume (`music-library-localstack-data`). To reset all data:

```bash
docker-compose down -v  # Removes volumes
docker-compose up -d    # Recreates resources
```

## Environment Variables

The backend connects to LocalStack when these environment variables are set:

```bash
LOCAL_MODE=true
LOCALSTACK_ENDPOINT=http://localhost:4566
AWS_REGION=us-east-1
DYNAMODB_TABLE_NAME=MusicLibrary
MEDIA_BUCKET=music-library-local-media
```

## Troubleshooting

**LocalStack not ready:**
```bash
# Wait for health check
until curl -s http://localhost:4566/_localstack/health | grep -q '"dynamodb": "running"'; do
    echo "Waiting for LocalStack..."
    sleep 2
done
```

**Resources not created:**
```bash
# Manually run init script
bash docker/localstack-init/init-aws.sh
```

**Permission denied:**
```bash
# Fix init script permissions
chmod +x docker/localstack-init/init-aws.sh
chmod +x docker/localstack-init/init-cognito.sh
```

**Cognito pool not created:**
```bash
# Manually run Cognito init script
bash docker/localstack-init/init-cognito.sh
```

**Test user login fails:**
```bash
# Verify user exists and is confirmed
aws --endpoint-url=http://localhost:4566 cognito-idp admin-get-user \
    --user-pool-id <pool-id> --username admin@local.test

# Re-run init to recreate users
docker-compose -f docker/docker-compose.yml down -v
docker-compose -f docker/docker-compose.yml up -d
./scripts/wait-for-localstack.sh 60
./docker/localstack-init/init-aws.sh
./docker/localstack-init/init-cognito.sh
```

## One-Command Setup

For convenience, use the root Makefile or local-dev.sh script:

```bash
# Using Make
make local           # Start full environment
make local-stop      # Stop all services
make test-integration # Run integration tests

# Using shell script
./scripts/local-dev.sh start   # Start full environment
./scripts/local-dev.sh stop    # Stop all services
./scripts/local-dev.sh test    # Run integration tests
./scripts/local-dev.sh status  # Check service status
./scripts/local-dev.sh reset   # Reset LocalStack data
```
