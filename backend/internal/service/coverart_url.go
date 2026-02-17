package service

import (
	"context"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// generateCoverArtURL generates a signed URL for cover art, preferring CloudFront over S3.
func generateCoverArtURL(ctx context.Context, cf repository.CloudFrontSigner, s3Repo repository.S3Repository, key string) string {
	if key == "" {
		return ""
	}
	if cf != nil {
		url, err := cf.GenerateSignedURL(ctx, key, 24*time.Hour)
		if err == nil {
			return url
		}
	}
	url, err := s3Repo.GeneratePresignedDownloadURL(ctx, key, 24*time.Hour)
	if err == nil {
		return url
	}
	return ""
}
