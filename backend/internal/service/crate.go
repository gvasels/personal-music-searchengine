package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// CrateRepository defines the repository interface for crates
type CrateRepository interface {
	CreateCrate(ctx context.Context, crate models.Crate) error
	GetCrate(ctx context.Context, userID, crateID string) (*models.Crate, error)
	UpdateCrate(ctx context.Context, crate models.Crate) error
	DeleteCrate(ctx context.Context, userID, crateID string) error
	ListCrates(ctx context.Context, userID string, filter models.CrateFilter) ([]models.Crate, error)
	CountUserCrates(ctx context.Context, userID string) (int, error)
	GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error)
}

// CrateService handles crate operations
type CrateService struct {
	repo       CrateRepository
	featureSvc *FeatureService
}

// NewCrateService creates a new crate service
func NewCrateService(repo CrateRepository, featureSvc *FeatureService) *CrateService {
	return &CrateService{
		repo:       repo,
		featureSvc: featureSvc,
	}
}

// CreateCrate creates a new crate
func (s *CrateService) CreateCrate(ctx context.Context, userID string, req models.CreateCrateRequest) (*models.Crate, error) {
	// Check feature access
	enabled, err := s.featureSvc.IsEnabled(ctx, userID, models.FeatureCrates)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, fmt.Errorf("crates feature is not enabled for your subscription tier")
	}

	// Check crate limit
	count, err := s.repo.CountUserCrates(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= models.MaxCratesPerUser {
		return nil, fmt.Errorf("maximum number of crates (%d) reached", models.MaxCratesPerUser)
	}

	now := time.Now()
	crate := models.Crate{
		ID:            uuid.New().String(),
		UserID:        userID,
		Name:          req.Name,
		Description:   req.Description,
		Color:         req.Color,
		TrackIDs:      []string{},
		TrackCount:    0,
		SortOrder:     models.CrateSortCustom,
		IsSmartCrate:  req.IsSmartCrate,
		SmartCriteria: req.SmartCriteria,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err := s.repo.CreateCrate(ctx, crate); err != nil {
		return nil, err
	}

	return &crate, nil
}

// GetCrate retrieves a crate by ID
func (s *CrateService) GetCrate(ctx context.Context, userID, crateID string) (*models.Crate, error) {
	crate, err := s.repo.GetCrate(ctx, userID, crateID)
	if err != nil {
		return nil, err
	}
	if crate == nil {
		return nil, fmt.Errorf("crate not found")
	}
	return crate, nil
}

// UpdateCrate updates an existing crate
func (s *CrateService) UpdateCrate(ctx context.Context, userID, crateID string, req models.UpdateCrateRequest) (*models.Crate, error) {
	crate, err := s.repo.GetCrate(ctx, userID, crateID)
	if err != nil {
		return nil, err
	}
	if crate == nil {
		return nil, fmt.Errorf("crate not found")
	}

	// Apply updates
	if req.Name != nil {
		crate.Name = *req.Name
	}
	if req.Description != nil {
		crate.Description = *req.Description
	}
	if req.Color != nil {
		crate.Color = *req.Color
	}
	if req.SortOrder != nil {
		crate.SortOrder = *req.SortOrder
		// Resort tracks if sort order changed
		if *req.SortOrder != models.CrateSortCustom {
			s.sortCrateTracks(ctx, userID, crate)
		}
	}
	if req.SmartCriteria != nil {
		crate.SmartCriteria = req.SmartCriteria
	}

	crate.UpdatedAt = time.Now()

	if err := s.repo.UpdateCrate(ctx, *crate); err != nil {
		return nil, err
	}

	return crate, nil
}

// DeleteCrate deletes a crate
func (s *CrateService) DeleteCrate(ctx context.Context, userID, crateID string) error {
	crate, err := s.repo.GetCrate(ctx, userID, crateID)
	if err != nil {
		return err
	}
	if crate == nil {
		return fmt.Errorf("crate not found")
	}

	return s.repo.DeleteCrate(ctx, userID, crateID)
}

// ListCrates returns all crates for a user
func (s *CrateService) ListCrates(ctx context.Context, userID string, filter models.CrateFilter) ([]models.CrateResponse, error) {
	crates, err := s.repo.ListCrates(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.CrateResponse, len(crates))
	for i, crate := range crates {
		responses[i] = crate.ToResponse()
	}

	return responses, nil
}

// AddTracks adds tracks to a crate
func (s *CrateService) AddTracks(ctx context.Context, userID, crateID string, req models.AddTracksToCrateRequest) error {
	crate, err := s.repo.GetCrate(ctx, userID, crateID)
	if err != nil {
		return err
	}
	if crate == nil {
		return fmt.Errorf("crate not found")
	}

	// Check track limit
	newCount := crate.TrackCount + len(req.TrackIDs)
	if newCount > models.MaxTracksPerCrate {
		return fmt.Errorf("adding these tracks would exceed the maximum of %d tracks per crate", models.MaxTracksPerCrate)
	}

	// Verify all tracks exist and belong to user
	for _, trackID := range req.TrackIDs {
		track, err := s.repo.GetTrack(ctx, userID, trackID)
		if err != nil {
			return err
		}
		if track == nil {
			return fmt.Errorf("track %s not found", trackID)
		}
	}

	// Remove duplicates
	existingTrackMap := make(map[string]bool)
	for _, id := range crate.TrackIDs {
		existingTrackMap[id] = true
	}

	newTracks := make([]string, 0)
	for _, trackID := range req.TrackIDs {
		if !existingTrackMap[trackID] {
			newTracks = append(newTracks, trackID)
			existingTrackMap[trackID] = true
		}
	}

	if len(newTracks) == 0 {
		return nil // All tracks already in crate
	}

	// Insert at position
	position := req.Position
	if position < 0 || position >= len(crate.TrackIDs) {
		// Append at end
		crate.TrackIDs = append(crate.TrackIDs, newTracks...)
	} else {
		// Insert at position
		newTrackIDs := make([]string, 0, len(crate.TrackIDs)+len(newTracks))
		newTrackIDs = append(newTrackIDs, crate.TrackIDs[:position]...)
		newTrackIDs = append(newTrackIDs, newTracks...)
		newTrackIDs = append(newTrackIDs, crate.TrackIDs[position:]...)
		crate.TrackIDs = newTrackIDs
	}

	crate.TrackCount = len(crate.TrackIDs)
	crate.UpdatedAt = time.Now()

	return s.repo.UpdateCrate(ctx, *crate)
}

// RemoveTracks removes tracks from a crate
func (s *CrateService) RemoveTracks(ctx context.Context, userID, crateID string, req models.RemoveTracksFromCrateRequest) error {
	crate, err := s.repo.GetCrate(ctx, userID, crateID)
	if err != nil {
		return err
	}
	if crate == nil {
		return fmt.Errorf("crate not found")
	}

	// Build removal set
	removeSet := make(map[string]bool)
	for _, id := range req.TrackIDs {
		removeSet[id] = true
	}

	// Filter out removed tracks
	newTrackIDs := make([]string, 0, len(crate.TrackIDs))
	for _, id := range crate.TrackIDs {
		if !removeSet[id] {
			newTrackIDs = append(newTrackIDs, id)
		}
	}

	crate.TrackIDs = newTrackIDs
	crate.TrackCount = len(newTrackIDs)
	crate.UpdatedAt = time.Now()

	return s.repo.UpdateCrate(ctx, *crate)
}

// ReorderTracks sets a new order for tracks in a crate
func (s *CrateService) ReorderTracks(ctx context.Context, userID, crateID string, req models.ReorderTracksRequest) error {
	crate, err := s.repo.GetCrate(ctx, userID, crateID)
	if err != nil {
		return err
	}
	if crate == nil {
		return fmt.Errorf("crate not found")
	}

	// Verify the new order contains the same tracks
	existingSet := make(map[string]bool)
	for _, id := range crate.TrackIDs {
		existingSet[id] = true
	}

	newSet := make(map[string]bool)
	for _, id := range req.TrackIDs {
		newSet[id] = true
	}

	// Check same tracks
	if len(existingSet) != len(newSet) {
		return fmt.Errorf("track list must contain the same tracks")
	}
	for id := range existingSet {
		if !newSet[id] {
			return fmt.Errorf("track %s missing from reordered list", id)
		}
	}
	for id := range newSet {
		if !existingSet[id] {
			return fmt.Errorf("unknown track %s in reordered list", id)
		}
	}

	crate.TrackIDs = req.TrackIDs
	crate.SortOrder = models.CrateSortCustom
	crate.UpdatedAt = time.Now()

	return s.repo.UpdateCrate(ctx, *crate)
}

// sortCrateTracks sorts tracks based on the crate's sort order
func (s *CrateService) sortCrateTracks(ctx context.Context, userID string, crate *models.Crate) {
	if len(crate.TrackIDs) == 0 || crate.SortOrder == models.CrateSortCustom {
		return
	}

	// Fetch all tracks
	tracks := make(map[string]*models.Track)
	for _, id := range crate.TrackIDs {
		track, err := s.repo.GetTrack(ctx, userID, id)
		if err == nil && track != nil {
			tracks[id] = track
		}
	}

	// Sort based on sort order
	sort.Slice(crate.TrackIDs, func(i, j int) bool {
		trackI := tracks[crate.TrackIDs[i]]
		trackJ := tracks[crate.TrackIDs[j]]

		if trackI == nil || trackJ == nil {
			return false
		}

		switch crate.SortOrder {
		case models.CrateSortBPM:
			return trackI.BPM < trackJ.BPM
		case models.CrateSortKey:
			return trackI.MusicalKey < trackJ.MusicalKey
		case models.CrateSortArtist:
			return trackI.Artist < trackJ.Artist
		case models.CrateSortTitle:
			return trackI.Title < trackJ.Title
		case models.CrateSortAdded:
			return trackI.CreatedAt.Before(trackJ.CreatedAt)
		default:
			return false
		}
	})
}
