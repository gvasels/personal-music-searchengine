# Personal Music Search Engine

A personal music library application for uploading, organizing, searching, and streaming your audio files. Built with a serverless architecture on AWS.

## Live Demo

- **Frontend**: https://d8wn3lkytn5qe.cloudfront.net
- **API**: https://r1simytb2i.execute-api.us-east-1.amazonaws.com

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

### Creator Studio Features (Phase 2 - Foundation Complete)

| Feature | Description | Status |
|---------|-------------|--------|
| **Feature Flags** | Tier-based feature gating (free/creator/pro) | ✅ Backend + Frontend |
| **Subscription Tiers** | Free, Creator ($9.99), Pro ($19.99), Studio ($49.99) | ✅ Backend + Frontend |
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

- AWS CLI configured with profile `gvasels-muza`
- Go 1.22+
- Node.js 20+
- Account: 887395463840, Region: us-east-1

### Backend Deployment

```bash
cd backend

# Build all Lambda functions (ARM64)
make build

# Package as ZIP files
make package

# Deploy API Lambda
aws lambda update-function-code \
  --function-name music-library-prod-api \
  --zip-file fileb://api.zip \
  --region us-east-1 \
  --profile gvasels-muza

# Deploy processor Lambdas
for fn in metadata cover-art-processor track-creator file-mover search-indexer upload-status-updater; do
  aws lambda update-function-code \
    --function-name music-library-prod-${fn} \
    --zip-file fileb://cmd/processor/${fn}.zip \
    --region us-east-1 \
    --profile gvasels-muza
done
```

**Lambda Functions:**

| Function | Purpose |
|----------|---------|
| `music-library-prod-api` | Main API (Echo + Lambda adapter) |
| `music-library-prod-metadata` | Metadata extraction from uploaded files |
| `music-library-prod-cover-art-processor` | Cover art extraction and resizing |
| `music-library-prod-track-creator` | Create track records in DynamoDB |
| `music-library-prod-file-mover` | Move files from upload to media bucket |
| `music-library-prod-search-indexer` | Index tracks in Nixiesearch |
| `music-library-prod-upload-status-updater` | Update upload status in DynamoDB |

### Frontend Deployment

```bash
cd frontend

# Install dependencies
npm install

# Build for production
npm run build

# Deploy to S3
aws s3 sync dist/ s3://music-library-prod-frontend \
  --delete \
  --profile gvasels-muza

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id E2XXXXXXXXXX \
  --paths "/*" \
  --profile gvasels-muza
```

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

### Backend

```bash
cd backend

# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build locally
go build -o bootstrap ./cmd/api
```

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev

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

---

## License

Private project - all rights reserved.

---

## Contributing

This is a personal project. See [CLAUDE.md](./CLAUDE.md) for development guidelines and SDLC workflow.
