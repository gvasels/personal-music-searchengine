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

func (s *streamService) GetStreamURL(ctx context.Context, userID, trackID string, hasGlobal bool) (*models.StreamResponse, error) {
	var track *models.Track
	var err error

	// First try to get as owner
	track, err = s.repo.GetTrack(ctx, userID, trackID)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	// If not found as owner, check if requester has global access or track is public
	if track == nil {
		track, err = s.repo.GetTrackByID(ctx, trackID)
		if err != nil {
			if err == repository.ErrNotFound {
				return nil, models.NewNotFoundError("Track", trackID)
			}
			return nil, err
		}

		// Track exists but requester doesn't own it - check access
		if hasGlobal {
			// Admins can access any track
		} else if track.Visibility == models.VisibilityPublic {
			// Public tracks can be accessed by anyone
		} else if track.Visibility == models.VisibilityUnlisted {
			// Unlisted tracks can be accessed via direct link
		} else {
			// Private track - return 403 Forbidden
			return nil, models.NewForbiddenError("you do not have permission to stream this track")
		}
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

func (s *streamService) GetDownloadURL(ctx context.Context, userID, trackID string, hasGlobal bool) (*models.DownloadResponse, error) {
	var track *models.Track
	var err error

	// First try to get as owner
	track, err = s.repo.GetTrack(ctx, userID, trackID)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	// If not found as owner, check if requester has global access or track is public
	if track == nil {
		track, err = s.repo.GetTrackByID(ctx, trackID)
		if err != nil {
			if err == repository.ErrNotFound {
				return nil, models.NewNotFoundError("Track", trackID)
			}
			return nil, err
		}

		// Track exists but requester doesn't own it - check access
		if hasGlobal {
			// Admins can access any track
		} else if track.Visibility == models.VisibilityPublic {
			// Public tracks can be accessed by anyone
		} else if track.Visibility == models.VisibilityUnlisted {
			// Unlisted tracks can be accessed via direct link
		} else {
			// Private track - return 403 Forbidden
			return nil, models.NewForbiddenError("you do not have permission to download this track")
		}
	}

	// Generate friendly filename
	fileName := fmt.Sprintf("%s - %s%s", track.Artist, track.Title, getExtensionFromFormat(track.Format))

	// Use S3 presigned URL for downloads - it supports Content-Disposition header natively
	// CloudFront would require query string forwarding configuration to support this
	downloadURL, err := s.s3Repo.GeneratePresignedDownloadURLWithFilename(ctx, track.S3Key, downloadURLExpiry, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate download URL: %w", err)
	}

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
