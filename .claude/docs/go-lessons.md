# Go Lessons Learned

Troubleshooting patterns and solutions discovered during development.

## DynamoDB Struct Tags

### Problem: Empty fields when unmarshaling DynamoDB items

**Symptom**: Struct fields are empty after `attributevalue.UnmarshalMap()` even though data exists in DynamoDB.

**Root Cause**: Field names in DynamoDB don't match Go struct tags.

**Solution**: Use lowercase `dynamodbav` tags that match DynamoDB attribute names:

```go
// WRONG - uses uppercase field names in DynamoDB
type RouteRecord struct {
    Product string `dynamodbav:"Product"`  // DynamoDB has "product"
}

// CORRECT - matches DynamoDB attribute names exactly
type RouteRecord struct {
    Product string `dynamodbav:"product"`
}
```

**Debugging**: Print the raw DynamoDB item to see actual attribute names:
```go
item, _ := client.GetItem(ctx, &dynamodb.GetItemInput{...})
fmt.Printf("Raw item: %+v\n", item.Item)
```

---

## Nested Structs in DynamoDB

### Problem: Nested objects not marshaling/unmarshaling correctly

**Symptom**: Nested structs come back as empty or nil.

**Solution**: Ensure nested structs also have proper `dynamodbav` tags:

```go
type Versions struct {
    Stable *VersionInfo      `dynamodbav:"stable"`
    Beta   *VersionInfo      `dynamodbav:"beta"`
    Canary *CanaryVersionInfo `dynamodbav:"canary"`
}

type VersionInfo struct {
    Version string `dynamodbav:"version"`
    S3URL   string `dynamodbav:"s3Url"`  // Note: camelCase to match DynamoDB
}
```

**DynamoDB JSON structure**:
```json
{
  "versions": {
    "stable": {
      "version": "1.0.0",
      "s3Url": "https://..."
    }
  }
}
```

---

## API Gateway Stage-Prefixed Routes

### Problem: Lambda returns 404 for all routes via API Gateway

**Symptom**: Health check works locally but returns 404 when called via API Gateway.

**Root Cause**: API Gateway HTTP API with named stages (e.g., "prod") includes the stage name in the request path. A request to `/prod/health` arrives at Lambda as `/prod/health`, not `/health`.

**Solution**: Add both direct and stage-prefixed routes:

```go
// Direct paths (for local testing and $default stage)
e.GET("/api/manifest", handler.GetManifest)
e.GET("/health", handler.HealthCheck)

// Stage-prefixed paths (for named stages like "prod", "staging")
e.GET("/:stage/api/manifest", handler.GetManifest)
e.GET("/:stage/health", handler.HealthCheck)
```

**Debugging**: Check Lambda CloudWatch logs for the actual request URI:
```bash
aws logs tail /aws/lambda/function-name --follow
# Look for: "uri":"/prod/health"
```

---

## Lambda Build Flags

### Problem: Lambda function fails to start or has large binary size

**Solution**: Use proper build flags for Lambda ARM64:

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
  go build -tags lambda.norpc \
  -ldflags="-s -w" \
  -o bootstrap cmd/api/main.go
```

| Flag | Purpose |
|------|---------|
| `GOOS=linux` | Target Linux (Lambda runs on Amazon Linux) |
| `GOARCH=arm64` | Target ARM64 (Graviton2 processors) |
| `CGO_ENABLED=0` | Disable CGO for static binary |
| `-tags lambda.norpc` | Exclude RPC code not needed for Lambda |
| `-ldflags="-s -w"` | Strip debug info for smaller binary |

---

## Testing with Embedded Structs

### Problem: Type mismatch when initializing structs with embedded types

**Symptom**:
```
cannot use &models.VersionInfo{} as *models.CanaryVersionInfo
```

**Solution**: Initialize embedded structs properly:

```go
// WRONG
Canary: &models.VersionInfo{
    Version: "2.0-canary",
    S3URL:   "https://...",
}

// CORRECT - CanaryVersionInfo embeds VersionInfo
Canary: &models.CanaryVersionInfo{
    VersionInfo: models.VersionInfo{
        Version: "2.0-canary",
        S3URL:   "https://...",
    },
    TrafficPercent: 10,
}
```

---

## Integration Tests with Build Tags

### Problem: Integration tests run during unit test phase

**Solution**: Use build tags to separate integration tests:

```go
//go:build integration

package tests

import "testing"

func TestIntegration(t *testing.T) {
    // This only runs with: go test -tags integration
}
```

**Run unit tests only**:
```bash
go test ./...
```

**Run integration tests**:
```bash
go test -tags integration ./tests/...
```

---

## Mock Interface Pattern

### Problem: Can't test code that uses AWS SDK clients

**Solution**: Define interfaces and inject dependencies:

```go
// Define interface matching AWS SDK methods you use
type DynamoDBAPI interface {
    Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
    GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

// Use interface in your code
type Service struct {
    db        DynamoDBAPI
    tableName string
}

// In tests, create a mock
type mockDynamoDB struct {
    QueryFunc   func(...) (*dynamodb.QueryOutput, error)
    GetItemFunc func(...) (*dynamodb.GetItemOutput, error)
}

func (m *mockDynamoDB) Query(...) (*dynamodb.QueryOutput, error) {
    return m.QueryFunc(...)
}
```

---

## Echo Context in Tests

### Problem: Need to test Echo handlers without starting a server

**Solution**: Use `httptest` with Echo's test utilities:

```go
func TestHandler(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/api/manifest", nil)
    req.Header.Set("x-canary", "true")
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    handler := NewManifestHandler(mockService, mockOverride)
    err := handler.GetManifest(c)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

---

## Go Module Replace Directives

### Problem: Need to use local version of a module during development

**Solution**: Use replace directive in go.mod:

```go
module github.com/org/project

require github.com/org/shared-lib v1.0.0

// During development, use local version
replace github.com/org/shared-lib => ../shared-lib
```

**Remove before committing** or use a separate `go.mod.local`.
