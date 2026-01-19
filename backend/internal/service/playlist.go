package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// playlistService implements PlaylistService
type playlistService struct {
	repo   repository.Repository
	s3Repo repository.S3Repository
}

// NewPlaylistService creates a new playlist service
func NewPlaylistService(repo repository.Repository, s3Repo repository.S3Repository) PlaylistService {
	return &playlistService{
		repo:   repo,
		s3Repo: s3Repo,
	}
}

func (s *playlistService) CreatePlaylist(ctx context.Context, userID string, req models.CreatePlaylistRequest) (*models.PlaylistResponse, error) {
	now := time.Now()
	playlist := models.Playlist{
		ID:            uuid.New().String(),
		UserID:        userID,
		Name:          req.Name,
		Description:   req.Description,
		TrackCount:    0,
		TotalDuration: 0,
		IsPublic:      req.IsPublic,
	}
	playlist.CreatedAt = now
	playlist.UpdatedAt = now

	if err := s.repo.CreatePlaylist(ctx, playlist); err != nil {
		return nil, err
	}

	response := playlist.ToResponse("")
	return &response, nil
}

func (s *playlistService) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.PlaylistWithTracks, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Playlist", playlistID)
		}
		return nil, err
	}

	coverArtURL := ""
	if playlist.CoverArtKey != "" {
		url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, playlist.CoverArtKey, 24*time.Hour)
		if err == nil {
			coverArtURL = url
		}
	}

	// Get playlist tracks
	playlistTracks, err := s.repo.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	// Get full track details for each playlist track
	tracks := make([]models.TrackResponse, 0, len(playlistTracks))
	for _, pt := range playlistTracks {
		track, err := s.repo.GetTrack(ctx, userID, pt.TrackID)
		if err != nil {
			if err == repository.ErrNotFound {
				continue // Skip deleted tracks
			}
			return nil, err
		}

		trackCoverURL := ""
		if track.CoverArtKey != "" {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
			if err == nil {
				trackCoverURL = url
			}
		}
		tracks = append(tracks, track.ToResponse(trackCoverURL))
	}

	return &models.PlaylistWithTracks{
		Playlist: playlist.ToResponse(coverArtURL),
		Tracks:   tracks,
	}, nil
}

func (s *playlistService) UpdatePlaylist(ctx context.Context, userID, playlistID string, req models.UpdatePlaylistRequest) (*models.PlaylistResponse, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Playlist", playlistID)
		}
		return nil, err
	}

	// Apply updates
	if req.Name != nil {
		playlist.Name = *req.Name
	}
	if req.Description != nil {
		playlist.Description = *req.Description
	}
	if req.IsPublic != nil {
		playlist.IsPublic = *req.IsPublic
	}

	if err := s.repo.UpdatePlaylist(ctx, *playlist); err != nil {
		return nil, err
	}

	coverArtURL := ""
	if playlist.CoverArtKey != "" {
		url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, playlist.CoverArtKey, 24*time.Hour)
		if err == nil {
			coverArtURL = url
		}
	}

	response := playlist.ToResponse(coverArtURL)
	return &response, nil
}

func (s *playlistService) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	_, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("Playlist", playlistID)
		}
		return err
	}

	return s.repo.DeletePlaylist(ctx, userID, playlistID)
}

func (s *playlistService) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.PlaylistResponse], error) {
	result, err := s.repo.ListPlaylists(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.PlaylistResponse, 0, len(result.Items))
	for _, playlist := range result.Items {
		coverArtURL := ""
		if playlist.CoverArtKey != "" {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, playlist.CoverArtKey, 24*time.Hour)
			if err == nil {
				coverArtURL = url
			}
		}
		responses = append(responses, playlist.ToResponse(coverArtURL))
	}

	return &repository.PaginatedResult[models.PlaylistResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (s *playlistService) AddTracks(ctx context.Context, userID, playlistID string, req models.AddTracksToPlaylistRequest) (*models.PlaylistResponse, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Playlist", playlistID)
		}
		return nil, err
	}

	// Determine position for new tracks
	position := playlist.TrackCount
	if req.Position != nil {
		position = *req.Position
	}

	// Validate that all tracks exist and calculate duration
	var totalDuration int
	for _, trackID := range req.TrackIDs {
		track, err := s.repo.GetTrack(ctx, userID, trackID)
		if err != nil {
			if err == repository.ErrNotFound {
				return nil, models.NewNotFoundError("Track", trackID)
			}
			return nil, err
		}
		totalDuration += track.Duration
	}

	// Add tracks to playlist
	if err := s.repo.AddTracksToPlaylist(ctx, playlistID, req.TrackIDs, position); err != nil {
		return nil, err
	}

	// Update playlist stats
	playlist.TrackCount += len(req.TrackIDs)
	playlist.TotalDuration += totalDuration
	if err := s.repo.UpdatePlaylist(ctx, *playlist); err != nil {
		return nil, err
	}

	coverArtURL := ""
	if playlist.CoverArtKey != "" {
		url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, playlist.CoverArtKey, 24*time.Hour)
		if err == nil {
			coverArtURL = url
		}
	}

	response := playlist.ToResponse(coverArtURL)
	return &response, nil
}

func (s *playlistService) RemoveTracks(ctx context.Context, userID, playlistID string, req models.RemoveTracksFromPlaylistRequest) (*models.PlaylistResponse, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Playlist", playlistID)
		}
		return nil, err
	}

	// Calculate duration to subtract
	var totalDuration int
	for _, trackID := range req.TrackIDs {
		track, err := s.repo.GetTrack(ctx, userID, trackID)
		if err == nil {
			totalDuration += track.Duration
		}
	}

	// Remove tracks from playlist
	if err := s.repo.RemoveTracksFromPlaylist(ctx, playlistID, req.TrackIDs); err != nil {
		return nil, err
	}

	// Update playlist stats
	playlist.TrackCount -= len(req.TrackIDs)
	if playlist.TrackCount < 0 {
		playlist.TrackCount = 0
	}
	playlist.TotalDuration -= totalDuration
	if playlist.TotalDuration < 0 {
		playlist.TotalDuration = 0
	}
	if err := s.repo.UpdatePlaylist(ctx, *playlist); err != nil {
		return nil, err
	}

	coverArtURL := ""
	if playlist.CoverArtKey != "" {
		url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, playlist.CoverArtKey, 24*time.Hour)
		if err == nil {
			coverArtURL = url
		}
	}

	response := playlist.ToResponse(coverArtURL)
	return &response, nil
}
