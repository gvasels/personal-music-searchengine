// Package validation provides input validation utilities for Lambda processors.
package validation

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// MaxFileSizeBytes is the maximum allowed file size for audio uploads (100MB).
const MaxFileSizeBytes int64 = 100 * 1024 * 1024

// ProcessorTimeoutSeconds is the timeout for processor Lambda operations.
// Set to 5 seconds less than Lambda timeout to allow graceful shutdown.
const ProcessorTimeoutSeconds = 55

// uuidRegex matches UUID v4 format (with or without hyphens).
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-?[0-9a-fA-F]{4}-?4[0-9a-fA-F]{3}-?[89abAB][0-9a-fA-F]{3}-?[0-9a-fA-F]{12}$`)

// IsValidUUID returns true if the string is a valid UUID v4 format.
func IsValidUUID(s string) bool {
	if s == "" {
		return false
	}
	return uuidRegex.MatchString(s)
}

// ValidateUUID returns an error if the string is not a valid UUID.
func ValidateUUID(s, fieldName string) error {
	if !IsValidUUID(s) {
		return fmt.Errorf("invalid %s: must be a valid UUID", fieldName)
	}
	return nil
}

// S3HeadObjectAPI defines the interface for S3 HeadObject operation.
type S3HeadObjectAPI interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

// ValidateFileSize checks if the S3 object is within the allowed size limit.
// Returns an error if the file exceeds MaxFileSizeBytes.
func ValidateFileSize(ctx context.Context, client S3HeadObjectAPI, bucket, key string) error {
	result, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to get file metadata: %w", err)
	}

	if result.ContentLength != nil && *result.ContentLength > MaxFileSizeBytes {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes (100MB)",
			*result.ContentLength, MaxFileSizeBytes)
	}

	return nil
}

// FileSizeError is returned when a file exceeds the maximum allowed size.
type FileSizeError struct {
	Size    int64
	MaxSize int64
}

func (e *FileSizeError) Error() string {
	return fmt.Sprintf("file size %d bytes exceeds maximum allowed size of %d bytes",
		e.Size, e.MaxSize)
}
