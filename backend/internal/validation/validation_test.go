package validation

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestIsValidUUID_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"standard format", "550e8400-e29b-41d4-a716-446655440000"},
		{"uppercase", "550E8400-E29B-41D4-A716-446655440000"},
		{"mixed case", "550e8400-E29B-41d4-A716-446655440000"},
		{"uuid v1", "550e8400-e29b-11d4-a716-446655440000"},
		{"uuid v5", "550e8400-e29b-51d4-a716-446655440000"},
		{"cognito style", "14c88468-d0f1-700d-646e-0e2e934f67b0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, IsValidUUID(tt.input))
		})
	}
}

func TestIsValidUUID_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"too short", "550e8400-e29b-41d4-a716"},
		{"too long", "550e8400-e29b-41d4-a716-446655440000-extra"},
		{"invalid characters", "550e8400-e29b-41d4-a716-44665544000g"},
		{"sql injection attempt", "'; DROP TABLE users;--"},
		{"path traversal", "../../../etc/passwd"},
		{"spaces", "550e8400 e29b 41d4 a716 446655440000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.False(t, IsValidUUID(tt.input))
		})
	}
}

func TestValidateUUID_Error(t *testing.T) {
	err := ValidateUUID("invalid", "userId")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userId")
	assert.Contains(t, err.Error(), "valid UUID")
}

func TestValidateUUID_Success(t *testing.T) {
	err := ValidateUUID("550e8400-e29b-41d4-a716-446655440000", "userId")
	assert.NoError(t, err)
}

// mockS3Client implements S3HeadObjectAPI for testing
type mockS3Client struct {
	contentLength int64
	err           error
}

func (m *mockS3Client) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &s3.HeadObjectOutput{
		ContentLength: aws.Int64(m.contentLength),
	}, nil
}

func TestValidateFileSize_Under100MB(t *testing.T) {
	client := &mockS3Client{contentLength: 50 * 1024 * 1024} // 50MB

	err := ValidateFileSize(context.Background(), client, "bucket", "key")
	assert.NoError(t, err)
}

func TestValidateFileSize_Exactly100MB(t *testing.T) {
	client := &mockS3Client{contentLength: MaxFileSizeBytes}

	err := ValidateFileSize(context.Background(), client, "bucket", "key")
	assert.NoError(t, err)
}

func TestValidateFileSize_Over100MB(t *testing.T) {
	client := &mockS3Client{contentLength: MaxFileSizeBytes + 1}

	err := ValidateFileSize(context.Background(), client, "bucket", "key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum")
}

func TestValidateFileSize_S3Error(t *testing.T) {
	client := &mockS3Client{err: errors.New("access denied")}

	err := ValidateFileSize(context.Background(), client, "bucket", "key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get file metadata")
}
