# Local Development Guide

This guide covers setting up and running the Personal Music Search Engine locally using LocalStack to emulate AWS services.

## Prerequisites

| Tool | Minimum Version | Installation |
|------|-----------------|--------------|
| Docker | 20.10+ | [docker.com](https://docs.docker.com/get-docker/) |
| Go | 1.22+ | [go.dev](https://go.dev/dl/) |
| Node.js | 18+ | [nodejs.org](https://nodejs.org/) |
| AWS CLI | 2.0+ | [aws.amazon.com](https://aws.amazon.com/cli/) |

Verify installation:
```bash
docker --version
go version
node --version
aws --version
```

## Quick Start

### One-Command Setup

```bash
# Using Make (recommended)
make local

# Or using shell script
./scripts/local-dev.sh start
```

This starts:
- **LocalStack** at http://localhost:4566
- **Backend API** at http://localhost:8080
- **Frontend** at http://localhost:5173

### Stop Everything

```bash
make local-stop
# or
./scripts/local-dev.sh stop
```

## Manual Setup

### Step 1: Start LocalStack

```bash
# Start LocalStack container
docker-compose -f docker/docker-compose.yml up -d

# Wait for services to be healthy
./scripts/wait-for-localstack.sh 60

# Initialize AWS resources
./docker/localstack-init/init-aws.sh
./docker/localstack-init/init-cognito.sh
```

### Step 2: Start Backend

```bash
cd backend
AWS_ENDPOINT=http://localhost:4566 \
DYNAMODB_TABLE_NAME=MusicLibrary \
MEDIA_BUCKET=music-library-local-media \
go run ./cmd/api
```

Backend runs at http://localhost:8080

### Step 3: Configure Frontend

Copy the environment template:
```bash
cp frontend/.env.local.example frontend/.env.local
```

Update `frontend/.env.local` with values from Cognito init output:
```env
VITE_LOCAL_STACK=true
VITE_API_URL=http://localhost:8080
VITE_COGNITO_USER_POOL_ID=us-east-1_xxxxxxxxx  # From init-cognito.sh
VITE_COGNITO_CLIENT_ID=xxxxxxxxxx              # From init-cognito.sh
VITE_COGNITO_REGION=us-east-1
VITE_COGNITO_ENDPOINT=http://localhost:4566
```

### Step 4: Start Frontend

```bash
cd frontend
npm run dev:local
```

Frontend runs at http://localhost:5173

## Test Users

All test users have the same password: `LocalTest123!`

| Email | Role | Permissions |
|-------|------|-------------|
| `admin@local.test` | Admin | Full system access |
| `subscriber@local.test` | Subscriber | Create playlists, follow artists |
| `artist@local.test` | Artist | Upload tracks, manage profile |

## Service Endpoints

| Service | URL | Health Check |
|---------|-----|--------------|
| LocalStack | http://localhost:4566 | `/_localstack/health` |
| Backend API | http://localhost:8080 | `/health` |
| Frontend | http://localhost:5173 | Direct access |

## Running Integration Tests

Integration tests run against LocalStack and require the `integration` build tag.

### Using Make

```bash
make test-integration
```

### Manual

```bash
# Ensure LocalStack is running
docker-compose -f docker/docker-compose.yml up -d
./scripts/wait-for-localstack.sh 60
./docker/localstack-init/init-aws.sh

# Run tests
cd backend
go test -tags=integration -v ./...
```

### Test Utilities

The `backend/internal/testutil` package provides helpers:

```go
import "github.com/gvasels/personal-music-searchengine/backend/internal/testutil"

func TestSomething(t *testing.T) {
    // Setup LocalStack connection
    ctx, cleanup := testutil.SetupLocalStack(t)
    defer cleanup()

    // Create test data
    trackID := ctx.CreateTestTrack(t, "user123")

    // Get auth token for test user
    token := ctx.GetTestUserToken(t, "admin")

    // Your test logic here...
}
```

## AWS CLI with LocalStack

Use the `--endpoint-url` flag for all AWS CLI commands:

```bash
# DynamoDB
aws --endpoint-url=http://localhost:4566 dynamodb list-tables
aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name MusicLibrary

# S3
aws --endpoint-url=http://localhost:4566 s3 ls
aws --endpoint-url=http://localhost:4566 s3 ls s3://music-library-local-media/

# Cognito
aws --endpoint-url=http://localhost:4566 cognito-idp list-user-pools --max-results 10
```

## Data Persistence

LocalStack data is stored in a Docker volume. To reset all data:

```bash
make local-reset
# or
./scripts/local-dev.sh reset
```

This removes:
- All DynamoDB data
- All S3 objects
- All Cognito users and pools

After reset, run `make local` to recreate resources.

## Troubleshooting

### LocalStack Not Starting

Check Docker is running:
```bash
docker ps
```

Check LocalStack logs:
```bash
docker-compose -f docker/docker-compose.yml logs localstack
```

### Backend Can't Connect to LocalStack

Verify LocalStack health:
```bash
curl http://localhost:4566/_localstack/health
```

Ensure environment variables are set:
```bash
AWS_ENDPOINT=http://localhost:4566 go run ./cmd/api
```

### Frontend Auth Not Working

1. Check `.env.local` has correct Cognito values
2. Verify Cognito init completed:
   ```bash
   aws --endpoint-url=http://localhost:4566 cognito-idp list-user-pools --max-results 10
   ```
3. Re-run Cognito init:
   ```bash
   ./docker/localstack-init/init-cognito.sh
   ```

### Tests Skipped

Integration tests require LocalStack. Check for:
```
=== SKIP: TestIntegration_...
    localstack.go:XX: LocalStack not running, skipping integration tests
```

Start LocalStack and re-run tests.

### Service Status

Check all services:
```bash
./scripts/local-dev.sh status
```

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    Frontend     │────▶│    Backend      │────▶│   LocalStack    │
│  localhost:5173 │     │  localhost:8080 │     │  localhost:4566 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                                        │
                              ┌─────────────────────────┼─────────────────────────┐
                              │                         │                         │
                        ┌─────▼─────┐           ┌──────▼──────┐          ┌───────▼───────┐
                        │ DynamoDB  │           │     S3      │          │   Cognito     │
                        │MusicLibrary│          │local-media  │          │local-pool     │
                        └───────────┘           └─────────────┘          └───────────────┘
```

## Differences from Production

| Aspect | Local | Production |
|--------|-------|------------|
| AWS Services | LocalStack emulation | Real AWS |
| Authentication | LocalStack Cognito | AWS Cognito |
| Database | LocalStack DynamoDB | AWS DynamoDB |
| Storage | LocalStack S3 | AWS S3 |
| API | Local Go server | AWS Lambda + API Gateway |
| CDN | None | CloudFront |

## Best Practices

1. **Use test utilities** - Don't create manual test data, use `testutil` package
2. **Clean up after tests** - Use `defer cleanup()` pattern
3. **Check for LocalStack** - Tests gracefully skip if LocalStack is unavailable
4. **Separate concerns** - Unit tests don't need LocalStack, integration tests do
5. **Reset regularly** - Use `make local-reset` when data gets corrupted

## Command Reference

| Command | Purpose |
|---------|---------|
| `make local` | Start full local environment |
| `make local-stop` | Stop all local services |
| `make local-services` | Start LocalStack only |
| `make local-backend` | Start backend only (requires LocalStack) |
| `make local-frontend` | Start frontend only |
| `make local-reset` | Reset LocalStack data |
| `make test-integration` | Run integration tests |
| `make test` | Run unit tests |
| `./scripts/local-dev.sh status` | Check service status |
