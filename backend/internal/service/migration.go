package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// MigrationService defines artist migration operations
type MigrationService interface {
	MigrateArtists(ctx context.Context, userID string) (*MigrationResult, error)
	GetMigrationStatus(ctx context.Context, userID string) (*MigrationStatus, error)
}

// MigrationResult represents the result of an artist migration
type MigrationResult struct {
	UserID          string    `json:"userId"`
	ArtistsCreated  int       `json:"artistsCreated"`
	TracksUpdated   int       `json:"tracksUpdated"`
	TracksSkipped   int       `json:"tracksSkipped"`
	Errors          []string  `json:"errors,omitempty"`
	StartedAt       time.Time `json:"startedAt"`
	CompletedAt     time.Time `json:"completedAt"`
	DurationSeconds float64   `json:"durationSeconds"`
}

// MigrationStatus represents the current status of a migration
type MigrationStatus struct {
	UserID           string    `json:"userId"`
	TotalTracks      int       `json:"totalTracks"`
	MigratedTracks   int       `json:"migratedTracks"`
	UnmigratedTracks int       `json:"unmigratedTracks"`
	TotalArtists     int       `json:"totalArtists"`
	LastMigration    *time.Time `json:"lastMigration,omitempty"`
	Status           string    `json:"status"` // "not_started", "in_progress", "completed", "partial"
}

// MigrationRepository defines the repository interface for migration operations
type MigrationRepository interface {
	CreateArtist(ctx context.Context, artist models.Artist) error
	GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error)
	GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error)
	ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error)
	UpdateTrack(ctx context.Context, track models.Track) error
	ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error)
}

// migrationService implements MigrationService
type migrationService struct {
	repo MigrationRepository
}

// NewMigrationService creates a new MigrationService
func NewMigrationService(repo MigrationRepository) MigrationService {
	return &migrationService{
		repo: repo,
	}
}

// MigrateArtists migrates existing string-based artists to entity model
// This is idempotent - safe to run multiple times
func (s *migrationService) MigrateArtists(ctx context.Context, userID string) (*MigrationResult, error) {
	startedAt := time.Now()
	result := &MigrationResult{
		UserID:    userID,
		StartedAt: startedAt,
		Errors:    []string{},
	}

	// Track unique artists encountered
	artistMap := make(map[string]*models.Artist) // name -> artist

	// Scan all tracks for this user
	var lastKey string
	for {
		filter := models.TrackFilter{
			Limit:   100,
			LastKey: lastKey,
		}

		tracksResult, err := s.repo.ListTracks(ctx, userID, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to list tracks: %w", err)
		}

		for _, track := range tracksResult.Items {
			// Skip tracks that already have artistId
			if track.ArtistID != "" {
				result.TracksSkipped++
				continue
			}

			// Skip tracks with no artist name
			if track.Artist == "" {
				result.TracksSkipped++
				continue
			}

			// Get or create artist
			artist, err := s.getOrCreateArtist(ctx, userID, track.Artist, artistMap)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to create artist '%s': %v", track.Artist, err))
				continue
			}

			// Update track with artistId
			track.ArtistLegacy = track.Artist // Backup original artist string
			track.ArtistID = artist.ID
			track.Artists = []models.ArtistContribution{
				{
					ArtistID:   artist.ID,
					ArtistName: artist.Name,
					Role:       models.RoleMain,
				},
			}

			if err := s.repo.UpdateTrack(ctx, track); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to update track '%s': %v", track.ID, err))
				continue
			}

			result.TracksUpdated++
		}

		if !tracksResult.HasMore {
			break
		}
		lastKey = tracksResult.NextCursor
	}

	// Count artists created
	for _, artist := range artistMap {
		if artist != nil {
			result.ArtistsCreated++
		}
	}

	result.CompletedAt = time.Now()
	result.DurationSeconds = result.CompletedAt.Sub(startedAt).Seconds()

	return result, nil
}

// getOrCreateArtist gets an existing artist or creates a new one
func (s *migrationService) getOrCreateArtist(ctx context.Context, userID, artistName string, cache map[string]*models.Artist) (*models.Artist, error) {
	// Check cache first
	if artist, ok := cache[artistName]; ok {
		return artist, nil
	}

	// Check if artist already exists in database
	existingArtists, err := s.repo.GetArtistByName(ctx, userID, artistName)
	if err == nil && len(existingArtists) > 0 {
		cache[artistName] = existingArtists[0]
		return existingArtists[0], nil
	}

	// Create new artist with deterministic ID
	artistID := GenerateDeterministicArtistID(userID, artistName)

	// Check if artist with this ID already exists (idempotency check)
	existingArtist, err := s.repo.GetArtist(ctx, userID, artistID)
	if err == nil && existingArtist != nil {
		cache[artistName] = existingArtist
		return existingArtist, nil
	}

	// Create new artist
	now := time.Now()
	artist := models.Artist{
		ID:        artistID,
		UserID:    userID,
		Name:      artistName,
		SortName:  models.GenerateSortName(artistName),
		IsActive:  true,
	}
	artist.CreatedAt = now
	artist.UpdatedAt = now

	if err := s.repo.CreateArtist(ctx, artist); err != nil {
		// If creation fails due to conflict, try to get the existing one
		existingArtist, getErr := s.repo.GetArtist(ctx, userID, artistID)
		if getErr == nil && existingArtist != nil {
			cache[artistName] = existingArtist
			return existingArtist, nil
		}
		return nil, fmt.Errorf("failed to create artist: %w", err)
	}

	cache[artistName] = &artist
	return &artist, nil
}

// GetMigrationStatus returns the current migration status for a user
func (s *migrationService) GetMigrationStatus(ctx context.Context, userID string) (*MigrationStatus, error) {
	status := &MigrationStatus{
		UserID: userID,
		Status: "not_started",
	}

	// Count total tracks
	totalTracks := 0
	migratedTracks := 0
	var lastKey string

	for {
		filter := models.TrackFilter{
			Limit:   100,
			LastKey: lastKey,
		}

		tracksResult, err := s.repo.ListTracks(ctx, userID, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to list tracks: %w", err)
		}

		for _, track := range tracksResult.Items {
			totalTracks++
			if track.ArtistID != "" {
				migratedTracks++
			}
		}

		if !tracksResult.HasMore {
			break
		}
		lastKey = tracksResult.NextCursor
	}

	// Count total artists
	artistFilter := models.ArtistFilter{Limit: 1000}
	artistsResult, err := s.repo.ListArtists(ctx, userID, artistFilter)
	if err == nil {
		status.TotalArtists = len(artistsResult.Items)
	}

	status.TotalTracks = totalTracks
	status.MigratedTracks = migratedTracks
	status.UnmigratedTracks = totalTracks - migratedTracks

	// Determine status
	if migratedTracks == 0 && totalTracks > 0 {
		status.Status = "not_started"
	} else if migratedTracks == totalTracks {
		status.Status = "completed"
	} else if migratedTracks > 0 {
		status.Status = "partial"
	}

	return status, nil
}
