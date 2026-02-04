# Validation Package - CLAUDE.md

## Overview

Input validation utilities for Lambda processors. Provides UUID validation, file size validation, and timeout constants to prevent security issues and resource exhaustion.

## File Descriptions

| File | Purpose |
|------|---------|
| `validation.go` | Core validation functions and constants |
| `validation_test.go` | Unit tests for validators |

## Constants

| Constant | Value | Purpose |
|----------|-------|---------|
| `MaxFileSizeBytes` | 100MB | Maximum allowed file size for audio uploads |
| `ProcessorTimeoutSeconds` | 55s | Context timeout for processor operations |

## Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `IsValidUUID` | `func IsValidUUID(s string) bool` | Returns true if string is valid UUID v4 |
| `ValidateUUID` | `func ValidateUUID(s, fieldName string) error` | Returns error with field name if invalid |
| `ValidateFileSize` | `func ValidateFileSize(ctx, client, bucket, key) error` | Checks S3 object size via HeadObject |

## Usage Examples

### UUID Validation
```go
import "github.com/gvasels/personal-music-searchengine/internal/validation"

func handleRequest(ctx context.Context, event Event) error {
    if err := validation.ValidateUUID(event.UserID, "userId"); err != nil {
        return err
    }
    if err := validation.ValidateUUID(event.UploadID, "uploadId"); err != nil {
        return err
    }
    // Continue processing...
}
```

### File Size Validation
```go
func handleRequest(ctx context.Context, event Event) error {
    if err := validation.ValidateFileSize(ctx, s3Client, event.BucketName, event.S3Key); err != nil {
        return fmt.Errorf("file validation failed: %w", err)
    }
    // Safe to download file...
}
```

### Context Timeout
```go
func handleRequest(ctx context.Context, event Event) error {
    ctx, cancel := context.WithTimeout(ctx, validation.ProcessorTimeoutSeconds*time.Second)
    defer cancel()

    // All operations now respect timeout...
}
```

## Security Considerations

- UUID validation prevents injection attacks via malformed IDs
- File size validation prevents OOM attacks via large file uploads
- Timeout prevents hung executions that consume Lambda resources

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/aws/aws-sdk-go-v2/service/s3` | S3 HeadObject interface |
