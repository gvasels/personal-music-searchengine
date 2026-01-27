# Personal Music Search Engine

A personal music library application for uploading, organizing, searching, and streaming your audio files. Built with a serverless architecture on AWS.

## Features

### Core Features (Implemented)

| Feature | Description |
|---------|-------------|
| **Audio Upload** | Upload MP3, FLAC, WAV, AAC, OGG files with automatic metadata extraction |
| **Metadata Extraction** | Extracts title, artist, album, year, genre, BPM, key from audio files |
| **Cover Art** | Extracts embedded cover art and displays in library and player |
| **Full-Text Search** | Search tracks by title, artist, album with fuzzy matching |
| **Playlist Management** | Create, edit, delete playlists; add/remove tracks; reorder via drag-and-drop |
| **Tag System** | Add custom tags to tracks for organization |
| **Audio Streaming** | Secure streaming via CloudFront signed URLs |
| **Audio Player** | Full-featured player with queue, shuffle, repeat, volume control |
| **BPM/Key Detection** | Automatic BPM and musical key detection during upload |
| **Mobile Responsive** | Works on desktop and mobile devices |
| **Dark/Light Theme** | Toggle between dark and light themes |

### Global User Type (Phase 2 - Complete)

| Feature | Description | Status |
|---------|-------------|--------|
| **User Roles** | Role-based access (admin, artist, subscriber, guest) via Cognito groups | ✅ Complete |
| **Public Playlists** | Visibility options (private/unlisted/public) with discovery page | ✅ Complete |
| **Artist Profiles** | Artist profile management with bio, social links, catalog linking | ✅ Complete |
| **Follow System** | Follow/unfollow artists with follower/following counts | ✅ Complete |
| **Admin Panel** | User management with search, role changes, enable/disable users | ✅ Complete |

### Creator Studio Features

| Feature | Description | Status |
|---------|-------------|--------|
| **Feature Flags** | Role-based feature gating | ✅ Backend + Frontend |
| **DJ Crates** | Organize tracks into DJ crates/folders | ✅ Backend + Frontend |
| **Hot Cues** | Set up to 8 cue points per track with labels and colors | ✅ Backend + Frontend |
| **BPM/Key Matching** | Find compatible tracks using Camelot wheel | ✅ Backend + Frontend |
| **Waveform Display** | Visual waveform in player (UI ready, data pending) | 🔄 UI Complete |
| **Creator Dashboard** | Stats and module cards for creator features | ✅ Frontend |

### Planned Features

See [ROADMAP.md](./ROADMAP.md) for the full roadmap including:
- AI Chatbot & Agent System
- Rights Management System (geographic scopes)
- Advanced DJ features (Mix Recorder, Beat Grid Editor, Live Streaming)
- Podcaster & Producer Studio features
- Video/Marengo integration

---

## Technology Stack

| Layer | Technologies |
|-------|--------------|
| **Backend** | Go 1.22+, Echo v4, AWS Lambda (ARM64), dhowden/tag |
| **Frontend** | React 18, TanStack Router/Query, Tailwind CSS v4, DaisyUI 5, Zustand, Vite |
| **Infrastructure** | AWS (Lambda, DynamoDB, S3, CloudFront, API Gateway, Cognito, Step Functions) |
| **Search** | Nixiesearch (Lucene-based) |
| **IaC** | OpenTofu 1.8+ |

---

## Repository Structure

```
├── backend/                 # Go Lambda services
│   ├── cmd/
│   │   ├── api/             # Main API Lambda
│   │   └── processor/       # Upload processing Lambdas
│   └── internal/
│       ├── handlers/        # HTTP handlers (Echo)
│       ├── models/          # Domain models
│       ├── repository/      # DynamoDB + S3 data access
│       ├── service/         # Business logic
│       ├── metadata/        # Audio metadata extraction
│       └── search/          # Nixiesearch client
├── frontend/                # React SPA
│   └── src/
│       ├── components/      # Reusable React components
│       ├── hooks/           # TanStack Query hooks
│       ├── lib/             # API client, stores
│       └── routes/          # TanStack Router pages
├── infrastructure/          # OpenTofu IaC
│   ├── global/              # Route53, ACM
│   ├── shared/              # Cognito, shared resources
│   ├── backend/             # Lambda, API Gateway, DynamoDB, S3
│   └── frontend/            # S3, CloudFront
└── .spec-workflow/          # Feature specifications
```

---

## Deployment

### Prerequisites

| Requirement | Version | Purpose |
|-------------|---------|---------|
| AWS CLI | v2+ | AWS resource management |
| OpenTofu | 1.8+ | Infrastructure as code (NOT Terraform) |
| Go | 1.22+ | Backend compilation |
| Node.js | 20+ | Frontend build |
| Make | - | Build automation |

**Current Deployment:**
- Account: `887395463840`
- Region: `us-east-1`
- AWS Profile: `gvasels-muza`

---

### Fresh AWS Account Deployment

Follow these steps to deploy to a new AWS account from scratch.

#### Step 1: Configure AWS Credentials

```bash
# Create AWS profile for the new account
aws configure --profile your-profile-name
# Enter: Access Key ID, Secret Access Key, Region (us-east-1), Output format (json)

# Verify access
aws sts get-caller-identity --profile your-profile-name
```

#### Step 2: Update Provider Configuration

Edit `infrastructure/global/main.tf` and update the profile:

```hcl
provider "aws" {
  region  = "us-east-1"
  profile = "your-profile-name"  # Change this
}
```

Update the same profile in:
- `infrastructure/shared/main.tf`
- `infrastructure/backend/main.tf`
- `infrastructure/frontend/main.tf`

#### Step 3: Bootstrap Global Infrastructure (State Backend)

```bash
cd infrastructure/global

# Initialize without backend (first time only - no state bucket yet)
tofu init

# Create state bucket, DynamoDB lock table, ECR repos, base IAM
tofu apply

# After successful apply, the S3 backend config in main.tf can be uncommented
# Then migrate local state to S3:
tofu init -migrate-state
```

**Resources created:**
| Resource | Name | Purpose |
|----------|------|---------|
| S3 Bucket | `music-library-prod-tofu-state` | Terraform state storage |
| DynamoDB Table | `music-library-prod-tofu-lock` | State locking |
| ECR Repositories | `music-library-prod-{api,processor,indexer}` | Lambda container images |
| IAM Role | `music-library-prod-lambda-execution` | Base Lambda execution |

#### Step 4: Deploy Shared Infrastructure

```bash
cd infrastructure/shared

tofu init
tofu apply
```

**Resources created:**
- Cognito User Pool and App Client
- DynamoDB table (`MusicLibrary`)
- S3 buckets (media, upload)
- CloudFront key pair for signed URLs

#### Step 5: Build and Push Lambda Code

```bash
cd backend

# Build all Lambda functions
make build

# Package as ZIP files
make package
```

#### Step 6: Deploy Backend Infrastructure

```bash
cd infrastructure/backend

tofu init
tofu apply
```

**Resources created:**
- API Gateway HTTP API with routes
- Lambda functions (API, processors)
- Step Functions state machine
- CloudFront distribution for media

**Note:** First deploy may require manual import if resources already exist:
```bash
# Example: Import existing API Gateway route
tofu import aws_apigatewayv2_route.route_name api-id/route-key
```

#### Step 7: Deploy Lambda Code

```bash
cd backend

# Deploy API Lambda
aws lambda update-function-code \
  --function-name music-library-prod-api \
  --zip-file fileb://api.zip \
  --region us-east-1 \
  --profile your-profile-name

# Deploy processor Lambdas
for fn in metadata cover-art-processor track-creator file-mover search-indexer upload-status-updater; do
  aws lambda update-function-code \
    --function-name music-library-prod-${fn} \
    --zip-file fileb://cmd/processor/${fn}.zip \
    --region us-east-1 \
    --profile your-profile-name
done
```

#### Step 8: Build and Deploy Frontend

```bash
cd frontend

# Install dependencies
npm install

# Update environment variables in .env.production with new values:
# - VITE_API_URL (from API Gateway output)
# - VITE_COGNITO_USER_POOL_ID (from Cognito output)
# - VITE_COGNITO_CLIENT_ID (from Cognito output)

# Build for production
npm run build
```

#### Step 9: Deploy Frontend Infrastructure

```bash
cd infrastructure/frontend

tofu init
tofu apply
```

**Resources created:**
- S3 bucket for static hosting
- CloudFront distribution

#### Step 10: Upload Frontend Assets

```bash
# Get bucket name from tofu output
FRONTEND_BUCKET=$(cd infrastructure/frontend && tofu output -raw frontend_bucket_name)

# Sync built assets
aws s3 sync frontend/dist/ s3://${FRONTEND_BUCKET} \
  --delete \
  --profile your-profile-name

# Invalidate CloudFront cache
DISTRIBUTION_ID=$(cd infrastructure/frontend && tofu output -raw cloudfront_distribution_id)
aws cloudfront create-invalidation \
  --distribution-id ${DISTRIBUTION_ID} \
  --paths "/*" \
  --profile your-profile-name
```

---

### Updating Existing Deployment

#### Infrastructure Changes

When OpenTofu files are modified:

```bash
# Navigate to the changed layer
cd infrastructure/{global|shared|backend|frontend}

# Review changes
tofu plan

# Apply changes
tofu apply
```

**Layer dependencies:** Deploy in order if cross-layer changes:
```
global → shared → backend → frontend
```

#### Backend Code Changes

```bash
cd backend

# Build and package
make build && make package

# Deploy all Lambdas
aws lambda update-function-code \
  --function-name music-library-prod-api \
  --zip-file fileb://api.zip \
  --profile gvasels-muza

for fn in metadata cover-art-processor track-creator file-mover search-indexer upload-status-updater; do
  aws lambda update-function-code \
    --function-name music-library-prod-${fn} \
    --zip-file fileb://cmd/processor/${fn}.zip \
    --profile gvasels-muza
done
```

#### Frontend Code Changes

```bash
cd frontend

# Build
npm run build

# Deploy to S3
aws s3 sync dist/ s3://music-library-prod-frontend \
  --delete \
  --profile gvasels-muza

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id EWDTYVQMDOQK2 \
  --paths "/*" \
  --profile gvasels-muza
```

---

### Quick Deploy Commands

```bash
# Full backend deploy (build + upload)
cd backend && make build && make package && \
  aws lambda update-function-code --function-name music-library-prod-api --zip-file fileb://api.zip --profile gvasels-muza

# Full frontend deploy (build + sync + invalidate)
cd frontend && npm run build && \
  aws s3 sync dist/ s3://music-library-prod-frontend --delete --profile gvasels-muza && \
  aws cloudfront create-invalidation --distribution-id EWDTYVQMDOQK2 --paths "/*" --profile gvasels-muza
```

---

### Lambda Functions Reference

| Function | Purpose |
|----------|---------|
| `music-library-prod-api` | Main API (Echo + Lambda adapter) |
| `music-library-prod-metadata` | Metadata extraction from uploaded files |
| `music-library-prod-cover-art-processor` | Cover art extraction and resizing |
| `music-library-prod-track-creator` | Create track records in DynamoDB |
| `music-library-prod-file-mover` | Move files from upload to media bucket |
| `music-library-prod-search-indexer` | Index tracks in Nixiesearch |
| `music-library-prod-upload-status-updater` | Update upload status in DynamoDB |

### Environment Variables

**Backend (Lambda):**
| Variable | Value |
|----------|-------|
| `DYNAMODB_TABLE_NAME` | `MusicLibrary` |
| `MEDIA_BUCKET` | `music-library-prod-media` |
| `CLOUDFRONT_DOMAIN` | `d8wn3lkytn5qe.cloudfront.net` |
| `STEP_FUNCTIONS_ARN` | `arn:aws:states:us-east-1:887395463840:stateMachine:...` |

**Frontend (.env.production):**
```env
VITE_API_URL=https://r1simytb2i.execute-api.us-east-1.amazonaws.com
VITE_COGNITO_USER_POOL_ID=us-east-1_XXXXXXXX
VITE_COGNITO_CLIENT_ID=XXXXXXXXXXXXXXXXXXXXXXXXXX
VITE_COGNITO_REGION=us-east-1
```

---

## Local Development

The project includes a complete **LocalStack-based development environment** that emulates AWS services locally (DynamoDB, S3, Cognito).

### Quick Start (LocalStack)

```bash
# Start full local environment (LocalStack + Backend + Frontend)
make local

# Or using shell script
./scripts/local-dev.sh start
```

This starts:
| Service | URL |
|---------|-----|
| **LocalStack** | http://localhost:4566 |
| **Backend API** | http://localhost:8080 |
| **Frontend** | http://localhost:5173 |

### Test Users

| Email | Password | Role |
|-------|----------|------|
| `admin@local.test` | `LocalTest123!` | Admin |
| `subscriber@local.test` | `LocalTest123!` | Subscriber |
| `artist@local.test` | `LocalTest123!` | Artist |

### Local Commands

| Command | Purpose |
|---------|---------|
| `make local` | Start full environment |
| `make local-stop` | Stop all services |
| `make test-integration` | Run integration tests against LocalStack |
| `make local-reset` | Reset LocalStack data |
| `./scripts/local-dev.sh status` | Check service status |

See [LOCAL_DEV.md](./LOCAL_DEV.md) for complete documentation.

### Backend (without LocalStack)

```bash
cd backend

# Run unit tests
go test -short ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build locally
go build -o bootstrap ./cmd/api
```

### Frontend (without LocalStack)

```bash
cd frontend

# Install dependencies
npm install

# Start dev server (against production API)
npm run dev

# Start dev server (local mode)
npm run dev:local

# Run tests
npm test

# Type check
npm run typecheck

# Lint
npm run lint
```

---

## API Endpoints

All endpoints require authentication via Cognito JWT.

### User
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/features` | Get user's role-based feature flags |
| GET | `/api/v1/users/me` | Get current user's profile |
| PUT | `/api/v1/users/me` | Update current user's profile |
| GET | `/api/v1/users/me/settings` | Get user settings |
| PATCH | `/api/v1/users/me/settings` | Update user settings |
| GET | `/api/v1/users/me/following` | Get artists user is following |

### Tracks
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/tracks` | List tracks with pagination |
| GET | `/api/v1/tracks/:id` | Get track by ID |
| PUT | `/api/v1/tracks/:id` | Update track metadata |
| DELETE | `/api/v1/tracks/:id` | Delete track |
| POST | `/api/v1/tracks/:id/tags` | Add tags to track |
| DELETE | `/api/v1/tracks/:id/tags/:tag` | Remove tag from track |

### Albums
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/albums` | List albums |
| GET | `/api/v1/albums/:id` | Get album with tracks |

### Artists
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/artists` | List artists |
| GET | `/api/v1/artists/:name/tracks` | Get artist's tracks |
| GET | `/api/v1/artists/:name/albums` | Get artist's albums |

### Playlists
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/playlists` | List playlists |
| POST | `/api/v1/playlists` | Create playlist |
| GET | `/api/v1/playlists/:id` | Get playlist with tracks |
| PUT | `/api/v1/playlists/:id` | Update playlist |
| DELETE | `/api/v1/playlists/:id` | Delete playlist |
| POST | `/api/v1/playlists/:id/tracks` | Add tracks |
| DELETE | `/api/v1/playlists/:id/tracks` | Remove tracks |
| PUT | `/api/v1/playlists/:id/tracks/reorder` | Reorder tracks (full array) |

### Tags
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/tags` | List tags |
| GET | `/api/v1/tags/:name/tracks` | Get tracks by tag |

### Search
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/search?q=query` | Simple search (accepts `q` or `query`) |
| POST | `/api/v1/search` | Advanced search with filters |
| GET | `/api/v1/autocomplete?q=query` | Search suggestions |

### Upload
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/upload/presigned` | Get presigned upload URL |
| POST | `/api/v1/upload/confirm` | Confirm upload completion |
| GET | `/api/v1/uploads` | List upload history |
| GET | `/api/v1/uploads/:id` | Get upload status |

### Streaming
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/stream/:trackId` | Get signed streaming URL |
| GET | `/api/v1/download/:trackId` | Get signed download URL |

### Artist Profiles
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/artists/entity` | List artist profiles |
| POST | `/api/v1/artists/entity` | Create artist profile (artist role) |
| GET | `/api/v1/artists/entity/:id` | Get artist profile |
| PUT | `/api/v1/artists/entity/:id` | Update artist profile (owner) |
| POST | `/api/v1/artists/entity/:id/follow` | Follow artist |
| DELETE | `/api/v1/artists/entity/:id/follow` | Unfollow artist |
| GET | `/api/v1/artists/entity/:id/followers` | Get followers |

### Admin (Admin role required)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/admin/users` | Search users by email |
| GET | `/api/v1/admin/users/:id` | Get user details |
| PUT | `/api/v1/admin/users/:id/role` | Update user role |
| PUT | `/api/v1/admin/users/:id/status` | Enable/disable user |

---

## Testing

### Backend Tests

```bash
cd backend
go test ./...                    # All tests
go test -short ./...             # Unit tests only
go test ./internal/service/...   # Specific package
```

**Coverage**: 134 tests in service layer, targeting 80%+ coverage.

### Frontend Tests

```bash
cd frontend
npm test                         # All tests (Vitest)
npm test -- --coverage           # With coverage
npm test -- --watch              # Watch mode
```

**Coverage**: 400+ tests across hooks, components, routes.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        CloudFront                                │
│                  (Static assets + API proxy)                     │
└─────────────────────────────┬───────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│   S3 Static   │     │  API Gateway  │     │   S3 Media    │
│   (Frontend)  │     │               │     │   (Audio)     │
└───────────────┘     └───────┬───────┘     └───────────────┘
                              │
                      ┌───────▼───────┐
                      │    Lambda     │
                      │   (Go API)    │
                      └───────┬───────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│   DynamoDB    │     │  Nixiesearch  │     │    Cognito    │
│ (Single-table)│     │   (Search)    │     │    (Auth)     │
└───────────────┘     └───────────────┘     └───────────────┘
```

### Upload Processing Flow

```
Upload → S3 (upload bucket) → Step Functions → [Metadata → CoverArt → Track → Move → Index → Status]
```

---

## DynamoDB Schema (Single-Table Design)

| Entity | PK | SK |
|--------|----|----|
| User | `USER#{userId}` | `PROFILE` |
| Track | `USER#{userId}` | `TRACK#{trackId}` |
| Album | `USER#{userId}` | `ALBUM#{albumId}` |
| Playlist | `USER#{userId}` | `PLAYLIST#{playlistId}` |
| PlaylistTrack | `USER#{userId}` | `PLAYLIST#{playlistId}#TRACK#{position}` |
| Tag | `USER#{userId}` | `TAG#{tagName}` |
| Upload | `USER#{userId}` | `UPLOAD#{uploadId}` |
| Crate | `USER#{userId}` | `CRATE#{crateId}` |
| HotCue | `USER#{userId}` | `HOTCUE#{trackId}#{slot}` |
| ArtistProfile | `ARTIST_PROFILE#{id}` | `PROFILE` |
| Follow | `USER#{userId}` | `FOLLOW#{artistProfileId}` |

### GSI Access Patterns
| GSI | Purpose |
|-----|---------|
| GSI1 | User → ArtistProfile lookup, Artist followers |
| GSI2 | Public playlist discovery, Linked artist uniqueness |

---

## License

Private project - all rights reserved.

---

## Contributing

This is a personal project. See [CLAUDE.md](./CLAUDE.md) for development guidelines and SDLC workflow.
