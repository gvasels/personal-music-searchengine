package service

import (
	"context"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// trackService implements TrackService
type trackService struct {
	repo   repository.Repository
	s3Repo repository.S3Repository
}

// NewTrackService creates a new track service
func NewTrackService(repo repository.Repository, s3Repo repository.S3Repository) TrackService {
	return &trackService{
		repo:   repo,
		s3Repo: s3Repo,
	}
}

func (s *trackService) GetTrack(ctx context.Context, userID, trackID string) (*models.TrackResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Track", trackID)
		}
		return nil, err
	}

	coverArtURL := ""
	if track.CoverArtKey != "" {
		// Generate signed URL for cover art
		url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
		if err == nil {
			coverArtURL = url
		}
	}

	response := track.ToResponse(coverArtURL)
	return &response, nil
}

func (s *trackService) UpdateTrack(ctx context.Context, userID, trackID string, req models.UpdateTrackRequest) (*models.TrackResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Track", trackID)
		}
		return nil, err
	}

	// Apply updates
	if req.Title != nil {
		track.Title = *req.Title
	}
	if req.Artist != nil {
		track.Artist = *req.Artist
	}
	if req.AlbumArtist != nil {
		track.AlbumArtist = *req.AlbumArtist
	}
	if req.Album != nil {
		track.Album = *req.Album
	}
	if req.Genre != nil {
		track.Genre = *req.Genre
	}
	if req.Year != nil {
		track.Year = *req.Year
	}
	if req.TrackNumber != nil {
		track.TrackNumber = *req.TrackNumber
	}
	if req.DiscNumber != nil {
		track.DiscNumber = *req.DiscNumber
	}
	if req.Lyrics != nil {
		track.Lyrics = *req.Lyrics
	}
	if req.Comment != nil {
		track.Comment = *req.Comment
	}
	if req.Tags != nil {
		track.Tags = req.Tags
	}

	if err := s.repo.UpdateTrack(ctx, *track); err != nil {
		return nil, err
	}

	coverArtURL := ""
	if track.CoverArtKey != "" {
		url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
		if err == nil {
			coverArtURL = url
		}
	}

	response := track.ToResponse(coverArtURL)
	return &response, nil
}

func (s *trackService) DeleteTrack(ctx context.Context, userID, trackID string) error {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("Track", trackID)
		}
		return err
	}

	// Delete from repository
	if err := s.repo.DeleteTrack(ctx, userID, trackID); err != nil {
		return err
	}

	// Delete files from S3 (best effort - don't fail if S3 delete fails)
	if track.S3Key != "" {
		_ = s.s3Repo.DeleteObject(ctx, track.S3Key)
	}
	if track.CoverArtKey != "" {
		_ = s.s3Repo.DeleteObject(ctx, track.CoverArtKey)
	}

	return nil
}

func (s *trackService) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error) {
	result, err := s.repo.ListTracks(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.TrackResponse, 0, len(result.Items))
	for _, track := range result.Items {
		coverArtURL := ""
		if track.CoverArtKey != "" {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
			if err == nil {
				coverArtURL = url
			}
		}
		responses = append(responses, track.ToResponse(coverArtURL))
	}

	return &repository.PaginatedResult[models.TrackResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (s *trackService) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.TrackResponse, error) {
	tracks, err := s.repo.ListTracksByArtist(ctx, userID, artist)
	if err != nil {
		return nil, err
	}

	responses := make([]models.TrackResponse, 0, len(tracks))
	for _, track := range tracks {
		coverArtURL := ""
		if track.CoverArtKey != "" {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
			if err == nil {
				coverArtURL = url
			}
		}
		responses = append(responses, track.ToResponse(coverArtURL))
	}

	return responses, nil
}

func (s *trackService) IncrementPlayCount(ctx context.Context, userID, trackID string) error {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("Track", trackID)
		}
		return err
	}

	track.PlayCount++
	now := time.Now()
	track.LastPlayed = &now

	return s.repo.UpdateTrack(ctx, *track)
}
