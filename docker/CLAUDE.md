# Docker - CLAUDE.md

## Overview

Docker configuration for local development using LocalStack to emulate AWS services (DynamoDB, S3).

## Directory Structure

```
docker/
├── docker-compose.yml       # LocalStack service definition
├── localstack-init/         # Initialization scripts
│   └── init-aws.sh          # Creates DynamoDB table and S3 bucket
└── CLAUDE.md                # This file
```

## File Descriptions

| File | Purpose |
|------|---------|
| `docker-compose.yml` | Defines LocalStack container with DynamoDB and S3 services |
| `localstack-init/init-aws.sh` | Initialization script that creates required AWS resources |

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

## Resources Created by init-aws.sh

### DynamoDB Table: `MusicLibrary`
- Primary Key: `PK` (Hash), `SK` (Range)
- GSI1: `GSI1PK`, `GSI1SK` (for artist/tag queries)
- Billing: Pay-per-request

### S3 Bucket: `music-library-local-media`
- Folders: `uploads/`, `media/`, `covers/`
- CORS configured for localhost:5173 and localhost:3000

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
```
