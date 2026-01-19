package service

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// PlaylistService handles playlist-related operations
type PlaylistService struct {
	repo repository.Repository
}

// NewPlaylistService creates a new PlaylistService
func NewPlaylistService(repo repository.Repository) *PlaylistService {
	return &PlaylistService{repo: repo}
}

// ListPlaylists lists playlists with filtering
func (s *PlaylistService) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*models.PaginatedResponse[models.PlaylistResponse], error) {
	result, err := s.repo.ListPlaylists(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.PlaylistResponse, len(result.Items))
	for i, playlist := range result.Items {
		coverArtURL := s.getCoverArtURL(playlist.CoverArtKey)
		responses[i] = playlist.ToResponse(coverArtURL)
	}

	return &models.PaginatedResponse[models.PlaylistResponse]{
		Items:      responses,
		Pagination: result.Pagination,
	}, nil
}

// CreatePlaylist creates a new playlist
func (s *PlaylistService) CreatePlaylist(ctx context.Context, userID string, req models.CreatePlaylistRequest) (*models.PlaylistResponse, error) {
	now := time.Now()
	playlist := &models.Playlist{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err := s.repo.CreatePlaylist(ctx, playlist); err != nil {
		return nil, err
	}

	response := playlist.ToResponse("")
	return &response, nil
}

// GetPlaylistWithTracks retrieves a playlist with its tracks
func (s *PlaylistService) GetPlaylistWithTracks(ctx context.Context, userID, playlistID string) (*models.PlaylistWithTracks, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}

	playlistTracks, err := s.repo.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	// Get full track details
	tracks := make([]models.TrackResponse, 0, len(playlistTracks))
	for _, pt := range playlistTracks {
		track, err := s.repo.GetTrack(ctx, userID, pt.TrackID)
		if err != nil {
			continue // Skip tracks that can't be found
		}
		coverArtURL := s.getCoverArtURL(track.CoverArtKey)
		tracks = append(tracks, track.ToResponse(coverArtURL))
	}

	coverArtURL := s.getCoverArtURL(playlist.CoverArtKey)
	return &models.PlaylistWithTracks{
		Playlist: playlist.ToResponse(coverArtURL),
		Tracks:   tracks,
	}, nil
}

// UpdatePlaylist updates a playlist
func (s *PlaylistService) UpdatePlaylist(ctx context.Context, userID, playlistID string, req models.UpdatePlaylistRequest) (*models.PlaylistResponse, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		playlist.Name = *req.Name
	}
	if req.Description != nil {
		playlist.Description = *req.Description
	}
	if req.IsPublic != nil {
		playlist.IsPublic = *req.IsPublic
	}

	playlist.UpdatedAt = time.Now()

	if err := s.repo.UpdatePlaylist(ctx, playlist); err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(playlist.CoverArtKey)
	response := playlist.ToResponse(coverArtURL)
	return &response, nil
}

// DeletePlaylist deletes a playlist
func (s *PlaylistService) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	return s.repo.DeletePlaylist(ctx, userID, playlistID)
}

// AddTracks adds tracks to a playlist
func (s *PlaylistService) AddTracks(ctx context.Context, userID, playlistID string, req models.AddTracksToPlaylistRequest) (*models.PlaylistResponse, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}

	existingTracks, err := s.repo.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	// Determine starting position
	startPosition := len(existingTracks)
	if req.Position != nil && *req.Position <= startPosition {
		startPosition = *req.Position
	}

	now := time.Now()
	totalDuration := playlist.TotalDuration

	for i, trackID := range req.TrackIDs {
		// Verify track exists
		track, err := s.repo.GetTrack(ctx, userID, trackID)
		if err != nil {
			continue
		}

		pt := &models.PlaylistTrack{
			PlaylistID: playlistID,
			TrackID:    trackID,
			Position:   startPosition + i,
			AddedAt:    now,
		}

		if err := s.repo.AddPlaylistTrack(ctx, pt); err != nil {
			return nil, err
		}

		totalDuration += track.Duration
	}

	// Update playlist stats
	playlist.TrackCount = len(existingTracks) + len(req.TrackIDs)
	playlist.TotalDuration = totalDuration
	playlist.UpdatedAt = now

	if err := s.repo.UpdatePlaylist(ctx, playlist); err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(playlist.CoverArtKey)
	response := playlist.ToResponse(coverArtURL)
	return &response, nil
}

// RemoveTracks removes tracks from a playlist
func (s *PlaylistService) RemoveTracks(ctx context.Context, userID, playlistID string, trackIDs []string) (*models.PlaylistResponse, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}

	for _, trackID := range trackIDs {
		if err := s.repo.RemovePlaylistTrack(ctx, playlistID, trackID); err != nil {
			return nil, err
		}
	}

	// Recalculate stats
	remainingTracks, err := s.repo.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	totalDuration := 0
	for _, pt := range remainingTracks {
		track, err := s.repo.GetTrack(ctx, userID, pt.TrackID)
		if err == nil {
			totalDuration += track.Duration
		}
	}

	playlist.TrackCount = len(remainingTracks)
	playlist.TotalDuration = totalDuration
	playlist.UpdatedAt = time.Now()

	if err := s.repo.UpdatePlaylist(ctx, playlist); err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(playlist.CoverArtKey)
	response := playlist.ToResponse(coverArtURL)
	return &response, nil
}

// ReorderTracks reorders tracks in a playlist
func (s *PlaylistService) ReorderTracks(ctx context.Context, userID, playlistID string, req models.ReorderPlaylistTracksRequest) (*models.PlaylistResponse, error) {
	playlist, err := s.repo.GetPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}

	existingTracks, err := s.repo.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	// Build new order
	trackIDs := make([]string, len(existingTracks))
	for i, pt := range existingTracks {
		trackIDs[i] = pt.TrackID
	}

	// Find and remove the track to be moved
	var movedTrackID string
	newTrackIDs := make([]string, 0, len(trackIDs))
	for _, id := range trackIDs {
		if id == req.TrackID {
			movedTrackID = id
		} else {
			newTrackIDs = append(newTrackIDs, id)
		}
	}

	if movedTrackID == "" {
		return nil, models.NewNotFoundError("Track in playlist", req.TrackID)
	}

	// Insert at new position
	position := req.NewPosition
	if position > len(newTrackIDs) {
		position = len(newTrackIDs)
	}

	finalTrackIDs := make([]string, 0, len(trackIDs))
	finalTrackIDs = append(finalTrackIDs, newTrackIDs[:position]...)
	finalTrackIDs = append(finalTrackIDs, movedTrackID)
	finalTrackIDs = append(finalTrackIDs, newTrackIDs[position:]...)

	if err := s.repo.ReorderPlaylistTracks(ctx, playlistID, finalTrackIDs); err != nil {
		return nil, err
	}

	playlist.UpdatedAt = time.Now()
	if err := s.repo.UpdatePlaylist(ctx, playlist); err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(playlist.CoverArtKey)
	response := playlist.ToResponse(coverArtURL)
	return &response, nil
}

func (s *PlaylistService) getCoverArtURL(coverArtKey string) string {
	if coverArtKey == "" {
		return ""
	}
	cloudFrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")
	if cloudFrontDomain == "" {
		return ""
	}
	return "https://" + cloudFrontDomain + "/" + coverArtKey
}
