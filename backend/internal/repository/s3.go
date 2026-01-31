package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// S3Client interface for testability
type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	CopyObject(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	CreateMultipartUpload(ctx context.Context, params *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error)
	CompleteMultipartUpload(ctx context.Context, params *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error)
	AbortMultipartUpload(ctx context.Context, params *s3.AbortMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error)
}

// S3PresignClient interface for presigned URL operations
type S3PresignClient interface {
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
	PresignUploadPart(ctx context.Context, params *s3.UploadPartInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// S3RepositoryImpl implements S3Repository
type S3RepositoryImpl struct {
	client        S3Client
	presignClient S3PresignClient
	bucketName    string
}

// NewS3Repository creates a new S3 repository
func NewS3Repository(client S3Client, presignClient S3PresignClient, bucketName string) *S3RepositoryImpl {
	return &S3RepositoryImpl{
		client:        client,
		presignClient: presignClient,
		bucketName:    bucketName,
	}
}

// GeneratePresignedUploadURL generates a presigned URL for uploading a file
// Note: StorageClass is NOT set here to simplify browser uploads (fewer signed headers).
// Objects are automatically transitioned to INTELLIGENT_TIERING via S3 lifecycle rules.
func (r *S3RepositoryImpl) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	request, err := r.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return request.URL, nil
}

// GeneratePresignedDownloadURL generates a presigned URL for downloading a file
func (r *S3RepositoryImpl) GeneratePresignedDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	request, err := r.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	return request.URL, nil
}

// GeneratePresignedDownloadURLWithFilename generates a presigned URL with Content-Disposition header
// to force the browser to download the file with the specified filename
func (r *S3RepositoryImpl) GeneratePresignedDownloadURLWithFilename(ctx context.Context, key string, expiry time.Duration, filename string) (string, error) {
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", filename)
	request, err := r.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(r.bucketName),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(contentDisposition),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	return request.URL, nil
}

// InitiateMultipartUpload starts a multipart upload and returns the upload ID
func (r *S3RepositoryImpl) InitiateMultipartUpload(ctx context.Context, key, contentType string) (string, error) {
	result, err := r.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:       aws.String(r.bucketName),
		Key:          aws.String(key),
		ContentType:  aws.String(contentType),
		StorageClass: types.StorageClassIntelligentTiering,
	})
	if err != nil {
		return "", fmt.Errorf("failed to initiate multipart upload: %w", err)
	}

	return *result.UploadId, nil
}

// GenerateMultipartUploadURLs generates presigned URLs for each part of a multipart upload
func (r *S3RepositoryImpl) GenerateMultipartUploadURLs(ctx context.Context, key, uploadID string, numParts int, expiry time.Duration) ([]models.MultipartUploadPartURL, error) {
	partURLs := make([]models.MultipartUploadPartURL, 0, numParts)
	expiresAt := time.Now().Add(expiry)

	for partNumber := 1; partNumber <= numParts; partNumber++ {
		request, err := r.presignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(r.bucketName),
			Key:        aws.String(key),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int32(int32(partNumber)),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = expiry
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL for part %d: %w", partNumber, err)
		}

		partURLs = append(partURLs, models.MultipartUploadPartURL{
			PartNumber: partNumber,
			UploadURL:  request.URL,
			ExpiresAt:  expiresAt,
		})
	}

	return partURLs, nil
}

// CompleteMultipartUpload completes a multipart upload
func (r *S3RepositoryImpl) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []models.CompletedPartInfo) error {
	completedParts := make([]types.CompletedPart, 0, len(parts))
	for _, part := range parts {
		completedParts = append(completedParts, types.CompletedPart{
			PartNumber: aws.Int32(int32(part.PartNumber)),
			ETag:       aws.String(part.ETag),
		})
	}

	_, err := r.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(r.bucketName),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to complete multipart upload: %w", err)
	}

	return nil
}

// AbortMultipartUpload aborts a multipart upload
func (r *S3RepositoryImpl) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	_, err := r.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(r.bucketName),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	})
	if err != nil {
		return fmt.Errorf("failed to abort multipart upload: %w", err)
	}

	return nil
}

// DeleteObject deletes an object from S3
func (r *S3RepositoryImpl) DeleteObject(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// DeleteByPrefix deletes all objects with the given prefix from S3
func (r *S3RepositoryImpl) DeleteByPrefix(ctx context.Context, prefix string) error {
	if prefix == "" {
		return fmt.Errorf("prefix cannot be empty")
	}

	// List all objects with the prefix
	var continuationToken *string
	for {
		listResult, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(r.bucketName),
			Prefix:            aws.String(prefix),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return fmt.Errorf("failed to list objects with prefix %s: %w", prefix, err)
		}

		if len(listResult.Contents) == 0 {
			break
		}

		// Build list of objects to delete
		objectsToDelete := make([]types.ObjectIdentifier, 0, len(listResult.Contents))
		for _, obj := range listResult.Contents {
			objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{
				Key: obj.Key,
			})
		}

		// Delete the batch of objects
		_, err = r.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(r.bucketName),
			Delete: &types.Delete{
				Objects: objectsToDelete,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete objects with prefix %s: %w", prefix, err)
		}

		// Check if there are more objects
		if !*listResult.IsTruncated {
			break
		}
		continuationToken = listResult.NextContinuationToken
	}

	return nil
}

// CopyObject copies an object within S3
func (r *S3RepositoryImpl) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	_, err := r.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:       aws.String(r.bucketName),
		Key:          aws.String(destKey),
		CopySource:   aws.String(fmt.Sprintf("%s/%s", r.bucketName, sourceKey)),
		StorageClass: types.StorageClassIntelligentTiering,
	})
	if err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	return nil
}

// GetObjectMetadata retrieves metadata for an S3 object
func (r *S3RepositoryImpl) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}

	metadata := make(map[string]string)
	if result.ContentType != nil {
		metadata["content-type"] = *result.ContentType
	}
	if result.ContentLength != nil {
		metadata["content-length"] = fmt.Sprintf("%d", *result.ContentLength)
	}
	if result.ETag != nil {
		metadata["etag"] = *result.ETag
	}
	if result.LastModified != nil {
		metadata["last-modified"] = result.LastModified.Format(time.RFC3339)
	}

	// Include user metadata
	for k, v := range result.Metadata {
		metadata[k] = v
	}

	return metadata, nil
}

// ObjectExists checks if an object exists in S3
func (r *S3RepositoryImpl) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if it's a "not found" error
		// AWS SDK v2 doesn't expose typed errors the same way, so we check the error message
		if isNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

// isNotFoundError checks if an error is a "not found" error from S3.
// Uses errors.As to unwrap the AWS SDK error chain (e.g. *smithyhttp.ResponseError).
func isNotFoundError(err error) bool {
	var notFound *types.NotFound
	if errors.As(err, &notFound) {
		return true
	}
	var noSuchKey *types.NoSuchKey
	if errors.As(err, &noSuchKey) {
		return true
	}
	// Fallback: check HTTP status code for wrapped 404 responses
	var respErr interface{ HTTPStatusCode() int }
	if errors.As(err, &respErr) && respErr.HTTPStatusCode() == 404 {
		return true
	}
	return false
}
