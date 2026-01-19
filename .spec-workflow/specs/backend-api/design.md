# Design - Backend API (Epic 2)

## Architecture Overview

```
┌────────────────────────────────────────────────────────────────────────────┐
│                         API Gateway (HTTP API v2)                          │
│                        Cognito JWT Authorizer                              │
└─────────────────────────────────┬──────────────────────────────────────────┘
                                  │
                                  ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                           API Lambda (Go/ARM64)                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Echo Framework + aws-lambda-go-api-proxy                           │   │
│  │  handlers.go → service layer → repository layer                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────┘
                                  │
                                  │ (ConfirmUpload triggers)
                                  ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                    Step Functions: Upload Processor                        │
│                                                                            │
│  ExtractMetadata → ProcessCoverArt → CreateTrack → MoveFile → Index       │
│        │                 │               │            │          │         │
│        ▼                 ▼               ▼            ▼          ▼         │
│   ┌─────────┐      ┌──────────┐    ┌─────────┐  ┌─────────┐ ┌─────────┐   │
│   │Metadata │      │CoverArt  │    │Track    │  │File     │ │Search   │   │
│   │Extractor│      │Processor │    │Creator  │  │Mover    │ │Indexer  │   │
│   │ Lambda  │      │ Lambda   │    │ Lambda  │  │ Lambda  │ │ Lambda  │   │
│   └─────────┘      └──────────┘    └─────────┘  └─────────┘ └─────────┘   │
│                                                                            │
│                                  ▼                                         │
│                         UploadStatusUpdater Lambda                         │
└────────────────────────────────────────────────────────────────────────────┘
```

## Design Decisions

### DD-1: Lambda Architecture - ARM64 (Graviton2)
**Decision**: Use ARM64 architecture for all Lambdas
**Rationale**: 20% cost reduction, better Go performance on Graviton2

### DD-2: Cover Art Processing - No Resize
**Decision**: Extract and store cover art as-is without resizing
**Rationale**: Simplifies processing, CloudFront can handle resize via query params if needed later

### DD-3: Local Development - LocalStack
**Decision**: Use LocalStack via docker-compose for local development
**Rationale**: Full AWS mock environment, matches production behavior

### DD-4: API Gateway Authorization
**Decision**: Use API Gateway's built-in Cognito JWT authorizer
**Rationale**: No custom middleware needed, user ID from JWT claims passed to Lambda

### DD-5: Metadata Extraction Library
**Decision**: Use `dhowden/tag` for audio metadata extraction
**Rationale**: Pure Go, no CGO dependencies, supports MP3/FLAC/WAV/OGG/M4A

### DD-6: Lambda Entry Modes
**Decision**: Single binary supports both Lambda and HTTP server modes
**Rationale**: Same code for local development and production

### DD-7: Missing Metadata Handling
**Decision**: Use filename as title, "Unknown Artist" for artist when metadata missing
**Rationale**: Allow raw WAV and files without tags, user can edit metadata later

### DD-8: Duplicate Upload Handling
**Decision**: Allow duplicate uploads, each creates a new track
**Rationale**: Simple implementation, user can delete duplicates manually

---

## Component Design

### 1. API Lambda (`cmd/api/main.go`)

```go
func main() {
    // Initialize dependencies
    cfg := loadConfig()
    repo := repository.NewDynamoDBRepository(cfg.DynamoClient, cfg.TableName)
    s3Repo := repository.NewS3Repository(cfg.S3Client, cfg.MediaBucket)
    services := service.NewServices(repo, s3Repo, cloudfront, cfg.MediaBucket, cfg.StepFunctionsARN)
    h := handlers.NewHandlers(services)

    // Create Echo instance
    e := echo.New()
    e.Validator = NewValidator()
    h.RegisterRoutes(e)

    // Run in Lambda or HTTP mode
    if isLambda() {
        adapter := echoadapter.New(e)
        lambda.Start(adapter.ProxyWithContext)
    } else {
        e.Start(":8080")
    }
}
```

**Functions:**
| Function | Description |
|----------|-------------|
| `main()` | Entry point - initializes deps, chooses Lambda or HTTP mode |
| `loadConfig()` | Loads config from environment variables |
| `isLambda()` | Detects Lambda runtime via AWS_LAMBDA_FUNCTION_NAME env |
| `NewValidator()` | Creates validator for request binding |

### 2. Metadata Extractor (`internal/metadata/extractor.go`)

```go
type Extractor interface {
    Extract(ctx context.Context, reader io.ReadSeeker) (*models.UploadMetadata, error)
    ExtractCoverArt(ctx context.Context, reader io.ReadSeeker) ([]byte, string, error)
}
```

**Functions:**
| Function | Description |
|----------|-------------|
| `NewExtractor()` | Creates metadata extractor instance |
| `Extract()` | Extracts all metadata from audio file |
| `ExtractCoverArt()` | Extracts embedded cover art image bytes and mime type |
| `detectFormat()` | Detects audio format from file header |
| `parseDuration()` | Calculates duration from metadata or estimates from file size |

### 3. Processor Lambdas (`cmd/processor/*/main.go`)

#### 3.1 Metadata Extractor Lambda
**Input:**
```json
{
  "uploadId": "uuid",
  "userId": "uuid",
  "s3Key": "uploads/{uploadId}/{fileName}",
  "fileName": "song.mp3",
  "bucketName": "music-library-prod-media"
}
```

**Output:**
```json
{
  "title": "Song Title",
  "artist": "Artist Name",
  "album": "Album Name",
  "duration": 180,
  "format": "MP3",
  "hasCoverArt": true
}
```

**Functions:**
| Function | Description |
|----------|-------------|
| `handleRequest()` | Lambda handler - downloads file, extracts metadata |
| `downloadFromS3()` | Downloads file to memory for processing |

#### 3.2 Cover Art Processor Lambda
**Input:** Step Functions state with metadata
**Output:** `{ "coverArtKey": "covers/{userId}/{trackId}.jpg" }` or empty

**Functions:**
| Function | Description |
|----------|-------------|
| `handleRequest()` | Lambda handler - extracts and uploads cover art |
| `uploadCoverArt()` | Uploads cover art bytes to S3 |

#### 3.3 Track Creator Lambda
**Input:** Step Functions state with metadata and cover art
**Output:** `{ "trackId": "uuid", "albumId": "uuid" }`

**Functions:**
| Function | Description |
|----------|-------------|
| `handleRequest()` | Creates Track and Album (if needed) in DynamoDB |
| `createTrack()` | Creates Track record with all metadata |
| `getOrCreateAlbum()` | Gets existing or creates new Album |

#### 3.4 File Mover Lambda
**Input:** Step Functions state with trackId
**Output:** `{ "finalKey": "media/{userId}/{trackId}.mp3" }`

**Functions:**
| Function | Description |
|----------|-------------|
| `handleRequest()` | Moves file from uploads/ to media/ prefix |
| `copyAndDelete()` | S3 copy then delete original |

#### 3.5 Search Indexer Lambda (Stub)
**Input:** Track metadata
**Output:** `{ "indexed": false, "reason": "not_implemented" }`

**Functions:**
| Function | Description |
|----------|-------------|
| `handleRequest()` | Returns stub response (Epic 3 implementation) |

#### 3.6 Upload Status Updater Lambda
**Input:** uploadId, status, optional error
**Output:** Updated upload record

**Functions:**
| Function | Description |
|----------|-------------|
| `handleRequest()` | Updates Upload status in DynamoDB |
| `updateStatus()` | Sets status, completedAt, errorMsg fields |

---

## Infrastructure Design

### API Gateway (`api-gateway.tf`)

```hcl
resource "aws_apigatewayv2_api" "api" {
  name          = "${local.name_prefix}-api"
  protocol_type = "HTTP"
  cors_configuration {
    allow_origins = ["http://localhost:5173", "https://${cloudfront_domain}"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["Authorization", "Content-Type"]
    max_age       = 86400
  }
}

resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.api.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "cognito"
  jwt_configuration {
    audience = [cognito_client_id]
    issuer   = "https://cognito-idp.${region}.amazonaws.com/${cognito_user_pool_id}"
  }
}
```

### Lambda Definitions (`lambda-api.tf`, `lambda-processors.tf`)

All Lambdas use:
- **Runtime**: `provided.al2023` (custom Go runtime)
- **Architecture**: `arm64`
- **Handler**: `bootstrap`
- **Memory**: 256MB (API), 512MB (processors)
- **Timeout**: 30s (API), 60s (processors)

---

## Local Development

### docker-compose.yml

```yaml
version: '3.8'
services:
  localstack:
    image: localstack/localstack:latest
    ports:
      - "4566:4566"
    environment:
      - SERVICES=s3,dynamodb,stepfunctions
      - DEFAULT_REGION=us-east-1
    volumes:
      - "./docker/localstack-init:/etc/localstack/init/ready.d"

  api:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - AWS_ENDPOINT=http://localstack:4566
      - DYNAMODB_TABLE_NAME=MusicLibrary
      - MEDIA_BUCKET=music-library-dev-media
    depends_on:
      - localstack
```

### LocalStack Init Script

```bash
#!/bin/bash
awslocal dynamodb create-table --table-name MusicLibrary ...
awslocal s3 mb s3://music-library-dev-media
```

---

## File Structure

```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go                    # API Lambda entrypoint
│   └── processor/
│       ├── metadata/main.go           # Metadata extractor
│       ├── coverart/main.go           # Cover art processor
│       ├── track/main.go              # Track creator
│       ├── mover/main.go              # File mover
│       ├── indexer/main.go            # Search indexer (stub)
│       └── status/main.go             # Status updater
├── internal/
│   ├── metadata/
│   │   ├── extractor.go               # Metadata extraction
│   │   └── extractor_test.go          # Tests
│   └── ... (existing packages)
└── go.mod                             # Add dhowden/tag dependency

infrastructure/backend/
├── api-gateway.tf                     # HTTP API + authorizer
├── lambda-api.tf                      # API Lambda
├── lambda-processors.tf               # Processor Lambdas
└── ... (existing files)

docker/
├── docker-compose.yml                 # LocalStack setup
├── localstack-init/
│   └── init-aws.sh                    # Create tables/buckets
└── Dockerfile.dev                     # Local dev container
```

---

## Testing Strategy

### Unit Tests
| Test | Coverage |
|------|----------|
| `metadata/extractor_test.go` | MP3, FLAC, WAV metadata extraction |
| `cmd/api/main_test.go` | Config loading, validator |

### Integration Tests
| Test | Coverage |
|------|----------|
| LocalStack + API | Full request flow |
| Step Functions mock | Processor Lambda chain |

### Test Fixtures
- Sample MP3 with metadata and cover art
- Sample FLAC with metadata
- Sample WAV (no metadata)
