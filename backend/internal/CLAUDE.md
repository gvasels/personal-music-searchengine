# Backend Internal Packages - CLAUDE.md

## Overview

Internal packages for the backend services. These packages are not exported and can only be imported within the `backend` module. Contains domain models, data access, business logic, and HTTP handlers.

## Directory Structure

```
internal/
├── handlers/       # HTTP request handlers (Echo)
├── metadata/       # Audio metadata extraction utilities
├── models/         # Domain models, DTOs, and constants
├── repository/     # Data access layer (DynamoDB, S3)
├── search/         # Nixiesearch client
└── service/        # Business logic layer
```

## Package Descriptions

| Package | Purpose | Key Types |
|---------|---------|-----------|
| `handlers` | HTTP request/response handling | `Handlers`, handler methods |
| `metadata` | Audio file metadata extraction | `Extractor`, `Metadata` |
| `models` | Domain models and data structures | `Track`, `Album`, `User`, etc. |
| `repository` | DynamoDB and S3 operations | `Repository`, `DynamoDBRepository` |
| `search` | Full-text search integration | `SearchClient`, `SearchResult` |
| `service` | Business logic and orchestration | `*Service` types |

## Dependency Flow

```
handlers → service → repository → models
              ↓
           metadata
              ↓
           search
```

**Rules:**
- `handlers` depends on `service` and `models`
- `service` depends on `repository`, `metadata`, `search`, and `models`
- `repository` depends on `models` and AWS SDK
- `metadata` depends on `models` and dhowden/tag
- `search` depends on `models`
- `models` has no internal dependencies

## Testing

Each package has corresponding `*_test.go` files:
- `models/` - Pure unit tests, no mocks needed
- `repository/` - Integration tests with DynamoDB Local
- `service/` - Unit tests with mocked repository
- `handlers/` - HTTP tests with mocked service
