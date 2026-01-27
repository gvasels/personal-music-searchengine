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

func (s *trackService) GetTrack(ctx context.Context, requesterID, trackID string, hasGlobal bool) (*models.TrackResponse, error) {
	var track *models.Track
	var err error
	var isOwner bool

	// First, try to get as owner (most common case)
	track, err = s.repo.GetTrack(ctx, requesterID, trackID)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	if track != nil {
		isOwner = true
	}

	// If not found as owner, check if requester has global access or track is public
	if track == nil {
		// Use GetTrackByID to find the track regardless of owner
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
			// Unlisted tracks can be accessed via direct link (treat as accessible)
		} else {
			// Private track - return 403 Forbidden
			return nil, models.NewForbiddenError("you do not have permission to access this track")
		}
	}

	coverArtURL := ""
	if track.CoverArtKey != "" {
		// Generate signed URL for cover art
		url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
		if err == nil {
			coverArtURL = url
		}
	}

	// For admin view of other users' tracks, populate owner display name
	if hasGlobal && !isOwner && track.UserID != "" {
		name, err := s.repo.GetUserDisplayName(ctx, track.UserID)
		if err == nil && name != "" {
			track.OwnerDisplayName = name
		} else {
			track.OwnerDisplayName = track.UserID
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

func (s *trackService) DeleteTrack(ctx context.Context, userID, trackID string, hasGlobal bool) error {
	var track *models.Track
	var err error
	var ownerID string

	// Try to get track as owner first
	track, err = s.repo.GetTrack(ctx, userID, trackID)
	if err != nil && err != repository.ErrNotFound {
		return err
	}

	if track != nil {
		// User owns this track
		ownerID = userID
	} else if hasGlobal {
		// Admin trying to delete another user's track - look up the track by ID
		track, err = s.repo.GetTrackByID(ctx, trackID)
		if err != nil {
			if err == repository.ErrNotFound {
				return models.NewNotFoundError("Track", trackID)
			}
			return err
		}
		ownerID = track.UserID
	} else {
		// Regular user trying to delete a track they don't own
		return models.NewNotFoundError("Track", trackID)
	}

	// Delete from repository using the actual owner's ID
	if err := s.repo.DeleteTrack(ctx, ownerID, trackID); err != nil {
		return err
	}

	// Delete files from S3 (best effort - don't fail if S3 delete fails)
	if track.S3Key != "" {
		_ = s.s3Repo.DeleteObject(ctx, track.S3Key)
	}
	if track.CoverArtKey != "" {
		_ = s.s3Repo.DeleteObject(ctx, track.CoverArtKey)
	}

	// Delete HLS transcoded files if they exist (best effort)
	// HLS files are stored at hls/{userID}/{trackID}/
	if track.HLSPlaylistKey != "" {
		hlsPrefix := "hls/" + ownerID + "/" + trackID + "/"
		_ = s.s3Repo.DeleteByPrefix(ctx, hlsPrefix)
	}

	return nil
}

func (s *trackService) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error) {
	// For admin users (GlobalScope=true), get all tracks
	if filter.GlobalScope {
		return s.listAllTracks(ctx, userID, filter)
	}

	// For regular users: get own tracks + public tracks from others
	return s.listTracksForRegularUser(ctx, userID, filter)
}

// listAllTracks returns all tracks for admin users, including owner display names
func (s *trackService) listAllTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error) {
	result, err := s.repo.ListTracks(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Collect unique user IDs to look up display names
	userIDSet := make(map[string]bool)
	for _, track := range result.Items {
		if track.UserID != "" {
			userIDSet[track.UserID] = true
		}
	}

	// Look up display names for all unique users
	displayNames := make(map[string]string)
	for uid := range userIDSet {
		name, err := s.repo.GetUserDisplayName(ctx, uid)
		if err == nil && name != "" {
			displayNames[uid] = name
		} else {
			// Fallback to user ID if display name not found
			displayNames[uid] = uid
		}
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
		// Set owner display name for admin view
		// Show "You" for current user's tracks, otherwise show the owner's display name
		if track.UserID == userID {
			track.OwnerDisplayName = "You"
		} else {
			track.OwnerDisplayName = displayNames[track.UserID]
		}
		responses = append(responses, track.ToResponse(coverArtURL))
	}

	return &repository.PaginatedResult[models.TrackResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

// listTracksForRegularUser returns own tracks + public tracks for non-admin users
func (s *trackService) listTracksForRegularUser(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error) {
	// Get user's own tracks
	ownResult, err := s.repo.ListTracks(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Track IDs we've seen for deduplication
	seenIDs := make(map[string]bool)
	responses := make([]models.TrackResponse, 0)

	// Process own tracks
	for _, track := range ownResult.Items {
		seenIDs[track.ID] = true
		coverArtURL := ""
		if track.CoverArtKey != "" {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
			if err == nil {
				coverArtURL = url
			}
		}
		responses = append(responses, track.ToResponse(coverArtURL))
	}

	// Also fetch public tracks from other users
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}
	publicResult, err := s.repo.ListPublicTracks(ctx, limit, "")
	if err == nil {
		// Add public tracks not already in results (avoid duplicates for user's own public tracks)
		for _, track := range publicResult.Items {
			if seenIDs[track.ID] {
				continue
			}
			seenIDs[track.ID] = true

			coverArtURL := ""
			if track.CoverArtKey != "" {
				url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
				if err == nil {
					coverArtURL = url
				}
			}
			responses = append(responses, track.ToResponse(coverArtURL))
		}
	}

	return &repository.PaginatedResult[models.TrackResponse]{
		Items:      responses,
		NextCursor: ownResult.NextCursor,
		HasMore:    ownResult.HasMore || (publicResult != nil && publicResult.HasMore),
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

// UpdateVisibility updates the visibility of a track.
// Only the track owner can update visibility.
func (s *trackService) UpdateVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	// Validate visibility value
	if !visibility.IsValid() {
		return models.NewValidationError("invalid visibility value")
	}

	// Verify track exists and belongs to user
	_, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("Track", trackID)
		}
		return err
	}

	// Update visibility in repository (this also updates GSI3 keys for public discovery)
	return s.repo.UpdateTrackVisibility(ctx, userID, trackID, visibility)
}

// GetLibraryStats returns aggregated library statistics based on scope
func (s *trackService) GetLibraryStats(ctx context.Context, userID string, scope StatsScope, hasGlobal bool) (*LibraryStats, error) {
	var tracks []models.Track

	switch scope {
	case StatsScopeAll:
		// Admin scope - get all tracks across all users
		if !hasGlobal {
			return nil, models.NewForbiddenError("admin access required for global stats")
		}
		filter := models.TrackFilter{Limit: 10000, GlobalScope: true}
		result, err := s.repo.ListTracks(ctx, userID, filter)
		if err != nil {
			return nil, err
		}
		tracks = result.Items

	case StatsScopePublic:
		// Public scope - only public tracks (for subscriber simulation)
		result, err := s.repo.ListPublicTracks(ctx, 10000, "")
		if err != nil {
			return nil, err
		}
		tracks = result.Items

	case StatsScopeOwn:
		fallthrough
	default:
		// Own scope - user's own tracks + public tracks
		filter := models.TrackFilter{Limit: 10000, GlobalScope: false}
		result, err := s.repo.ListTracks(ctx, userID, filter)
		if err != nil {
			return nil, err
		}
		tracks = result.Items

		// Also add public tracks from others
		publicResult, err := s.repo.ListPublicTracks(ctx, 10000, "")
		if err == nil {
			seenIDs := make(map[string]bool)
			for _, t := range tracks {
				seenIDs[t.ID] = true
			}
			for _, t := range publicResult.Items {
				if !seenIDs[t.ID] {
					tracks = append(tracks, t)
				}
			}
		}
	}

	// Aggregate stats
	albums := make(map[string]bool)
	artists := make(map[string]bool)
	totalDuration := 0

	for _, track := range tracks {
		if track.Album != "" {
			albums[track.Album] = true
		}
		if track.Artist != "" {
			artists[track.Artist] = true
		}
		totalDuration += track.Duration
	}

	return &LibraryStats{
		TotalTracks:   len(tracks),
		TotalAlbums:   len(albums),
		TotalArtists:  len(artists),
		TotalDuration: totalDuration,
	}, nil
}
