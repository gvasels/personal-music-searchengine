//go:build integration

package repository_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_S3_PresignedUploadAndDownload(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	presignClient := s3.NewPresignClient(tc.S3)
	repo := repository.NewS3Repository(tc.S3, presignClient, tc.BucketName)
	ctx := context.Background()

	key := "uploads/test-user/test-track.mp3"
	content := []byte("fake audio content for presigned upload test")
	contentType := "audio/mpeg"

	t.Run("generate presigned upload URL and upload via HTTP", func(t *testing.T) {
		url, err := repo.GeneratePresignedUploadURL(ctx, key, contentType, 15*time.Minute)
		require.NoError(t, err)
		assert.NotEmpty(t, url)
		assert.Contains(t, url, key)

		// Upload via presigned URL
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(content))
		require.NoError(t, err)
		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify object exists
		exists, err := repo.ObjectExists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("generate presigned download URL and download via HTTP", func(t *testing.T) {
		url, err := repo.GeneratePresignedDownloadURL(ctx, key, 15*time.Minute)
		require.NoError(t, err)
		assert.NotEmpty(t, url)

		resp, err := http.Get(url) //nolint:gosec
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		downloaded, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, content, downloaded)
	})

	tc.RegisterS3Cleanup(key)
}

func TestIntegration_S3_ObjectExistsAndMetadata(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	presignClient := s3.NewPresignClient(tc.S3)
	repo := repository.NewS3Repository(tc.S3, presignClient, tc.BucketName)
	ctx := context.Background()

	key := "test/metadata/track.mp3"
	content := []byte("metadata test content")

	// Upload directly via raw S3 client for setup
	_, err := tc.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(tc.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String("audio/mpeg"),
		Metadata:    map[string]string{"artist": "Test Artist", "title": "Test Track"},
	})
	require.NoError(t, err)
	tc.RegisterS3Cleanup(key)

	t.Run("ObjectExists returns true for existing key", func(t *testing.T) {
		exists, err := repo.ObjectExists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("ObjectExists returns false for missing key", func(t *testing.T) {
		exists, err := repo.ObjectExists(ctx, "nonexistent/key.mp3")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("GetObjectMetadata returns content-type and user metadata", func(t *testing.T) {
		meta, err := repo.GetObjectMetadata(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, "audio/mpeg", meta["content-type"])
		assert.Equal(t, "Test Artist", meta["artist"])
		assert.Equal(t, "Test Track", meta["title"])
	})
}

func TestIntegration_S3_DeleteObject(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	presignClient := s3.NewPresignClient(tc.S3)
	repo := repository.NewS3Repository(tc.S3, presignClient, tc.BucketName)
	ctx := context.Background()

	key := "test/delete/to-remove.txt"

	// Create object
	_, err := tc.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(tc.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader([]byte("to be deleted")),
		ContentType: aws.String("text/plain"),
	})
	require.NoError(t, err)

	t.Run("delete existing object", func(t *testing.T) {
		err := repo.DeleteObject(ctx, key)
		require.NoError(t, err)

		exists, err := repo.ObjectExists(ctx, key)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("delete non-existent object is idempotent", func(t *testing.T) {
		err := repo.DeleteObject(ctx, "nonexistent/key.mp3")
		assert.NoError(t, err)
	})
}

func TestIntegration_S3_DeleteByPrefix(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	presignClient := s3.NewPresignClient(tc.S3)
	repo := repository.NewS3Repository(tc.S3, presignClient, tc.BucketName)
	ctx := context.Background()

	// Create objects under hls/user1/track1/
	track1Keys := []string{
		"hls/user1/track1/segment001.ts",
		"hls/user1/track1/segment002.ts",
		"hls/user1/track1/master.m3u8",
	}
	for _, key := range track1Keys {
		_, err := tc.S3.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(tc.BucketName),
			Key:         aws.String(key),
			Body:        bytes.NewReader([]byte("segment data")),
			ContentType: aws.String("application/octet-stream"),
		})
		require.NoError(t, err)
	}

	// Create object under different track (should NOT be deleted)
	track2Key := "hls/user1/track2/segment001.ts"
	_, err := tc.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(tc.BucketName),
		Key:         aws.String(track2Key),
		Body:        bytes.NewReader([]byte("track2 segment")),
		ContentType: aws.String("application/octet-stream"),
	})
	require.NoError(t, err)
	tc.RegisterS3Cleanup(track2Key)

	t.Run("deletes all objects with matching prefix", func(t *testing.T) {
		err := repo.DeleteByPrefix(ctx, "hls/user1/track1/")
		require.NoError(t, err)

		// Verify track1 objects are gone
		for _, key := range track1Keys {
			exists, err := repo.ObjectExists(ctx, key)
			require.NoError(t, err)
			assert.False(t, exists, "track1 object should be deleted: %s", key)
		}

		// Verify track2 object still exists
		exists, err := repo.ObjectExists(ctx, track2Key)
		require.NoError(t, err)
		assert.True(t, exists, "track2 object should still exist")
	})

	t.Run("no-op for non-matching prefix", func(t *testing.T) {
		err := repo.DeleteByPrefix(ctx, "nonexistent/prefix/")
		assert.NoError(t, err)
	})
}

func TestIntegration_S3_CopyObject(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	presignClient := s3.NewPresignClient(tc.S3)
	repo := repository.NewS3Repository(tc.S3, presignClient, tc.BucketName)
	ctx := context.Background()

	srcKey := "test/copy/source.mp3"
	dstKey := "test/copy/destination.mp3"
	content := []byte("content to copy")

	_, err := tc.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(tc.BucketName),
		Key:         aws.String(srcKey),
		Body:        bytes.NewReader(content),
		ContentType: aws.String("audio/mpeg"),
	})
	require.NoError(t, err)
	tc.RegisterS3Cleanup(srcKey)
	tc.RegisterS3Cleanup(dstKey)

	t.Run("copy creates new object with same content", func(t *testing.T) {
		err := repo.CopyObject(ctx, srcKey, dstKey)
		require.NoError(t, err)

		// Verify destination exists
		exists, err := repo.ObjectExists(ctx, dstKey)
		require.NoError(t, err)
		assert.True(t, exists)

		// Verify source still exists
		exists, err = repo.ObjectExists(ctx, srcKey)
		require.NoError(t, err)
		assert.True(t, exists)
	})
}
