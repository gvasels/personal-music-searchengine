package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

const (
	// URL expiration times
	streamURLExpiry   = 4 * time.Hour  // Shorter for streaming
	downloadURLExpiry = 24 * time.Hour // Longer for downloads
	coverArtURLExpiry = 24 * time.Hour
)

// streamService implements StreamService
type streamService struct {
	repo       repository.Repository
	cloudfront repository.CloudFrontSigner
	s3Repo     repository.S3Repository
}

// NewStreamService creates a new stream service
func NewStreamService(repo repository.Repository, cloudfront repository.CloudFrontSigner, s3Repo repository.S3Repository) StreamService {
	return &streamService{
		repo:       repo,
		cloudfront: cloudfront,
		s3Repo:     s3Repo,
	}
}

func (s *streamService) GetStreamURL(ctx context.Context, userID, trackID string) (*models.StreamResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Track", trackID)
		}
		return nil, err
	}

	var hlsURL, fallbackURL string

	// Generate HLS URL if available
	if track.HLSStatus == models.HLSStatusReady && track.HLSPlaylistKey != "" {
		if s.cloudfront != nil {
			hlsURL, err = s.cloudfront.GenerateSignedURL(ctx, track.HLSPlaylistKey, streamURLExpiry)
			if err != nil {
				// Log error but continue with fallback
				fmt.Printf("Warning: failed to generate HLS URL: %v\n", err)
			}
		}
	}

	// Generate fallback URL (direct audio file)
	if s.cloudfront != nil {
		fallbackURL, err = s.cloudfront.GenerateSignedURL(ctx, track.S3Key, streamURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate stream URL: %w", err)
		}
	} else {
		fallbackURL, err = s.s3Repo.GeneratePresignedDownloadURL(ctx, track.S3Key, streamURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate stream URL: %w", err)
		}
	}

	// Use HLS URL as primary if available, otherwise fallback
	streamURL := hlsURL
	if streamURL == "" {
		streamURL = fallbackURL
	}

	// Increment play count asynchronously (best effort)
	go func() {
		bgCtx := context.Background()
		track.PlayCount++
		now := time.Now()
		track.LastPlayed = &now
		_ = s.repo.UpdateTrack(bgCtx, *track)
	}()

	return &models.StreamResponse{
		TrackID:     trackID,
		StreamURL:   streamURL,
		HLSURL:      hlsURL,
		FallbackURL: fallbackURL,
		HLSReady:    track.HLSStatus == models.HLSStatusReady,
		ExpiresAt:   time.Now().Add(streamURLExpiry),
		Format:      string(track.Format),
		Bitrate:     track.Bitrate,
	}, nil
}

func (s *streamService) GetDownloadURL(ctx context.Context, userID, trackID string) (*models.DownloadResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Track", trackID)
		}
		return nil, err
	}

	// Generate download URL
	var downloadURL string
	if s.cloudfront != nil {
		// Use CloudFront for download
		downloadURL, err = s.cloudfront.GenerateSignedURL(ctx, track.S3Key, downloadURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate download URL: %w", err)
		}
	} else {
		// Fallback to S3 presigned URL
		downloadURL, err = s.s3Repo.GeneratePresignedDownloadURL(ctx, track.S3Key, downloadURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate download URL: %w", err)
		}
	}

	// Generate friendly filename
	fileName := fmt.Sprintf("%s - %s%s", track.Artist, track.Title, getExtensionFromFormat(track.Format))

	return &models.DownloadResponse{
		TrackID:     trackID,
		DownloadURL: downloadURL,
		ExpiresAt:   time.Now().Add(downloadURLExpiry),
		FileName:    fileName,
		FileSize:    track.FileSize,
		Format:      string(track.Format),
	}, nil
}

func (s *streamService) GetCoverArtURL(ctx context.Context, userID, trackID string) (string, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return "", models.NewNotFoundError("Track", trackID)
		}
		return "", err
	}

	if track.CoverArtKey == "" {
		return "", nil // No cover art
	}

	// Generate signed URL for cover art
	if s.cloudfront != nil {
		return s.cloudfront.GenerateSignedURL(ctx, track.CoverArtKey, coverArtURLExpiry)
	}
	return s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, coverArtURLExpiry)
}

// getExtensionFromFormat returns the file extension for an audio format
func getExtensionFromFormat(format models.AudioFormat) string {
	switch format {
	case models.AudioFormatMP3:
		return ".mp3"
	case models.AudioFormatFLAC:
		return ".flac"
	case models.AudioFormatWAV:
		return ".wav"
	case models.AudioFormatAAC:
		return ".m4a"
	case models.AudioFormatOGG:
		return ".ogg"
	default:
		return filepath.Ext(string(format))
	}
}
