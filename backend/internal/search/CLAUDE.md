# Search Package - CLAUDE.md

## Overview

Nixiesearch client package for full-text search functionality. Provides a client that communicates with a Nixiesearch Lambda function for search, index, and delete operations.

## File Descriptions

| File | Purpose |
|------|---------|
| `types.go` | Search request/response types and document schema |
| `client.go` | Nixiesearch Lambda client implementation |
| `client_test.go` | Unit tests with mock Lambda client |

## Key Types

### Document
Represents a searchable track in the index.
```go
type Document struct {
    ID        string    // Track UUID
    UserID    string    // Owner's user UUID
    Title     string    // Track title
    Artist    string    // Artist name
    Album     string    // Album name
    Genre     string    // Genre
    Year      int       // Release year
    Duration  int       // Duration in seconds
    Filename  string    // Original filename
    IndexedAt time.Time // Index timestamp
}
```

### SearchQuery
```go
type SearchQuery struct {
    Query   string        // Full-text search query
    Filters SearchFilters // Optional filters
    Sort    *SortOption   // Optional sorting
    Limit   int           // Page size (default 20, max 100)
    Cursor  string        // Pagination cursor
}
```

## Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NewClient` | `func NewClient(lambda LambdaInvoker, fn string) *Client` | Creates search client |
| `Search` | `func (c *Client) Search(ctx, userID, query) (*SearchResponse, error)` | Executes search query |
| `Index` | `func (c *Client) Index(ctx, doc) (*IndexResponse, error)` | Indexes a document |
| `Delete` | `func (c *Client) Delete(ctx, docID) (*DeleteResponse, error)` | Deletes a document |
| `BulkIndex` | `func (c *Client) BulkIndex(ctx, docs) (*BulkIndexResponse, error)` | Bulk index documents |

## Usage Example

```go
import (
    "github.com/aws/aws-sdk-go-v2/service/lambda"
    "github.com/gvasels/personal-music-searchengine/internal/search"
)

func main() {
    lambdaClient := lambda.NewFromConfig(cfg)
    searchClient := search.NewClient(lambdaClient, "nixiesearch-lambda")

    // Search for tracks
    resp, err := searchClient.Search(ctx, "user-123", search.SearchQuery{
        Query: "beatles",
        Filters: search.SearchFilters{
            Genre: "Rock",
        },
        Limit: 20,
    })

    // Index a new track
    _, err = searchClient.Index(ctx, search.Document{
        ID:     "track-uuid",
        UserID: "user-123",
        Title:  "Hey Jude",
        Artist: "The Beatles",
    })
}
```

## Architecture

```
Client ──► Lambda.Invoke() ──► Nixiesearch Lambda ──► EFS Index
```

The client invokes a dedicated Nixiesearch Lambda function that:
1. Receives JSON-encoded requests
2. Performs operations on the search index (stored in EFS)
3. Returns JSON-encoded responses

## Security

- All search queries are automatically scoped to the authenticated user
- The `UserID` filter is added by the client before sending to Lambda
- Documents are isolated per user in the index

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/aws/aws-sdk-go-v2/service/lambda` | Lambda invocation |
