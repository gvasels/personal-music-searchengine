package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// StreamService handles streaming-related operations
type StreamService struct {
	repo             repository.Repository
	s3Client         *s3.Client
	mediaBucket      string
	cloudFrontDomain string
}

// NewStreamService creates a new StreamService
func NewStreamService(repo repository.Repository, s3Client *s3.Client) *StreamService {
	mediaBucket := os.Getenv("MEDIA_BUCKET")
	if mediaBucket == "" {
		mediaBucket = "music-library-media"
	}

	cloudFrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")

	return &StreamService{
		repo:             repo,
		s3Client:         s3Client,
		mediaBucket:      mediaBucket,
		cloudFrontDomain: cloudFrontDomain,
	}
}

// GetStreamURL returns a signed URL for streaming a track
func (s *StreamService) GetStreamURL(ctx context.Context, userID, trackID, quality string) (*models.StreamResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}

	// Generate presigned URL for streaming
	expiresAt := time.Now().Add(4 * time.Hour)

	var streamURL string
	if s.cloudFrontDomain != "" {
		// Use CloudFront signed URL
		streamURL = fmt.Sprintf("https://%s/%s", s.cloudFrontDomain, track.S3Key)
	} else {
		// Fallback to S3 presigned URL
		s3Client := s.s3Client
		presignClient := s3.NewPresignClient(s3Client)

		presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.mediaBucket),
			Key:    aws.String(track.S3Key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 4 * time.Hour
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create presigned URL: %w", err)
		}
		streamURL = presignedReq.URL
	}

	return &models.StreamResponse{
		TrackID:   trackID,
		StreamURL: streamURL,
		ExpiresAt: expiresAt,
		Format:    string(track.Format),
		Bitrate:   track.Bitrate,
	}, nil
}

// GetDownloadURL returns a signed URL for downloading a track
func (s *StreamService) GetDownloadURL(ctx context.Context, userID, trackID string) (*models.DownloadResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}

	// Generate presigned URL for download
	s3Client := s.s3Client
	presignClient := s3.NewPresignClient(s3Client)

	fileName := fmt.Sprintf("%s - %s%s", track.Artist, track.Title, filepath.Ext(track.S3Key))

	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(s.mediaBucket),
		Key:                        aws.String(track.S3Key),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", fileName)),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 1 * time.Hour
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create download URL: %w", err)
	}

	return &models.DownloadResponse{
		TrackID:     trackID,
		DownloadURL: presignedReq.URL,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		FileName:    fileName,
		FileSize:    track.FileSize,
		Format:      string(track.Format),
	}, nil
}

// RecordPlayback records a playback event
func (s *StreamService) RecordPlayback(ctx context.Context, userID string, req models.RecordPlayRequest) error {
	track, err := s.repo.GetTrack(ctx, userID, req.TrackID)
	if err != nil {
		return err
	}

	// Only count as a play if listened for more than 30 seconds
	if req.Duration >= 30 || req.Completed {
		track.PlayCount++
		now := time.Now()
		track.LastPlayed = &now
		track.UpdatedAt = now

		if err := s.repo.UpdateTrack(ctx, track); err != nil {
			return err
		}
	}

	return nil
}

// GetQueue returns the user's play queue
func (s *StreamService) GetQueue(ctx context.Context, userID string) (*models.PlayQueue, error) {
	// For simplicity, return empty queue
	// In production, this would be stored in DynamoDB or Redis
	return &models.PlayQueue{
		UserID:       userID,
		CurrentIndex: 0,
		TrackIDs:     []string{},
		ShuffleMode:  false,
		RepeatMode:   "none",
		UpdatedAt:    time.Now(),
	}, nil
}

// UpdateQueue updates the user's play queue
func (s *StreamService) UpdateQueue(ctx context.Context, userID string, req models.UpdateQueueRequest) (*models.PlayQueue, error) {
	queue, _ := s.GetQueue(ctx, userID)

	if req.TrackIDs != nil {
		queue.TrackIDs = req.TrackIDs
	}
	if req.CurrentIndex != nil {
		queue.CurrentIndex = *req.CurrentIndex
	}
	if req.ShuffleMode != nil {
		queue.ShuffleMode = *req.ShuffleMode
	}
	if req.RepeatMode != nil {
		queue.RepeatMode = *req.RepeatMode
	}
	queue.UpdatedAt = time.Now()

	return queue, nil
}

// QueueAction performs an action on the play queue
func (s *StreamService) QueueAction(ctx context.Context, userID string, req models.QueueActionRequest) (*models.PlayQueue, error) {
	queue, _ := s.GetQueue(ctx, userID)

	switch req.Action {
	case models.QueueActionAddNext:
		if len(req.TrackIDs) > 0 {
			// Insert after current position
			insertPos := queue.CurrentIndex + 1
			if insertPos > len(queue.TrackIDs) {
				insertPos = len(queue.TrackIDs)
			}
			newTrackIDs := make([]string, 0, len(queue.TrackIDs)+len(req.TrackIDs))
			newTrackIDs = append(newTrackIDs, queue.TrackIDs[:insertPos]...)
			newTrackIDs = append(newTrackIDs, req.TrackIDs...)
			newTrackIDs = append(newTrackIDs, queue.TrackIDs[insertPos:]...)
			queue.TrackIDs = newTrackIDs
		}
	case models.QueueActionAddLast:
		queue.TrackIDs = append(queue.TrackIDs, req.TrackIDs...)
	case models.QueueActionRemove:
		// Remove specified tracks
		newTrackIDs := make([]string, 0, len(queue.TrackIDs))
		removeSet := make(map[string]bool)
		for _, id := range req.TrackIDs {
			removeSet[id] = true
		}
		for _, id := range queue.TrackIDs {
			if !removeSet[id] {
				newTrackIDs = append(newTrackIDs, id)
			}
		}
		queue.TrackIDs = newTrackIDs
	case models.QueueActionClear:
		queue.TrackIDs = []string{}
		queue.CurrentIndex = 0
	case models.QueueActionShuffle:
		queue.ShuffleMode = true
	case models.QueueActionUnshuffle:
		queue.ShuffleMode = false
	}

	queue.UpdatedAt = time.Now()
	return queue, nil
}
