package service

import (
	"context"
	"os"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// TrackService handles track-related operations
type TrackService struct {
	repo        repository.Repository
	mediaBucket string
}

// NewTrackService creates a new TrackService
func NewTrackService(repo repository.Repository) *TrackService {
	mediaBucket := os.Getenv("MEDIA_BUCKET")
	if mediaBucket == "" {
		mediaBucket = "music-library-media"
	}

	return &TrackService{
		repo:        repo,
		mediaBucket: mediaBucket,
	}
}

// GetTrack retrieves a track by ID
func (s *TrackService) GetTrack(ctx context.Context, userID, trackID string) (*models.TrackResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(track.CoverArtKey)
	response := track.ToResponse(coverArtURL)
	return &response, nil
}

// ListTracks lists tracks with filtering
func (s *TrackService) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResponse[models.TrackResponse], error) {
	result, err := s.repo.ListTracks(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.TrackResponse, len(result.Items))
	for i, track := range result.Items {
		coverArtURL := s.getCoverArtURL(track.CoverArtKey)
		responses[i] = track.ToResponse(coverArtURL)
	}

	return &models.PaginatedResponse[models.TrackResponse]{
		Items:      responses,
		Pagination: result.Pagination,
	}, nil
}

// UpdateTrack updates a track's metadata
func (s *TrackService) UpdateTrack(ctx context.Context, userID, trackID string, req models.UpdateTrackRequest) (*models.TrackResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
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

	track.UpdatedAt = time.Now()

	if err := s.repo.UpdateTrack(ctx, track); err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(track.CoverArtKey)
	response := track.ToResponse(coverArtURL)
	return &response, nil
}

// DeleteTrack deletes a track
func (s *TrackService) DeleteTrack(ctx context.Context, userID, trackID string) error {
	return s.repo.DeleteTrack(ctx, userID, trackID)
}

// AddTagsToTrack adds tags to a track
func (s *TrackService) AddTagsToTrack(ctx context.Context, userID, trackID string, tags []string) (*models.TrackResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}

	// Add new tags (avoid duplicates)
	existingTags := make(map[string]bool)
	for _, t := range track.Tags {
		existingTags[t] = true
	}

	for _, tag := range tags {
		if !existingTags[tag] {
			track.Tags = append(track.Tags, tag)

			// Create track-tag association
			tt := &models.TrackTag{
				UserID:  userID,
				TrackID: trackID,
				TagName: tag,
				AddedAt: time.Now(),
			}
			if err := s.repo.AddTrackTag(ctx, tt); err != nil {
				return nil, err
			}
		}
	}

	track.UpdatedAt = time.Now()
	if err := s.repo.UpdateTrack(ctx, track); err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(track.CoverArtKey)
	response := track.ToResponse(coverArtURL)
	return &response, nil
}

// RemoveTagFromTrack removes a tag from a track
func (s *TrackService) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) (*models.TrackResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}

	// Remove tag from list
	newTags := make([]string, 0, len(track.Tags))
	for _, t := range track.Tags {
		if t != tagName {
			newTags = append(newTags, t)
		}
	}
	track.Tags = newTags

	// Remove track-tag association
	if err := s.repo.RemoveTrackTag(ctx, userID, trackID, tagName); err != nil {
		return nil, err
	}

	track.UpdatedAt = time.Now()
	if err := s.repo.UpdateTrack(ctx, track); err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(track.CoverArtKey)
	response := track.ToResponse(coverArtURL)
	return &response, nil
}

func (s *TrackService) getCoverArtURL(coverArtKey string) string {
	if coverArtKey == "" {
		return ""
	}
	// Return CloudFront URL for cover art
	cloudFrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")
	if cloudFrontDomain == "" {
		return ""
	}
	return "https://" + cloudFrontDomain + "/" + coverArtKey
}
