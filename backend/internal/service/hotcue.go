package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// HotCueRepository defines the repository interface for hot cue operations
type HotCueRepository interface {
	GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error)
	UpdateTrack(ctx context.Context, track models.Track) error
}

// HotCueService handles hot cue operations
type HotCueService struct {
	repo       HotCueRepository
	featureSvc *FeatureService
}

// NewHotCueService creates a new hot cue service
func NewHotCueService(repo HotCueRepository, featureSvc *FeatureService) *HotCueService {
	return &HotCueService{
		repo:       repo,
		featureSvc: featureSvc,
	}
}

// SetHotCue sets or updates a hot cue at a specific slot
func (s *HotCueService) SetHotCue(ctx context.Context, userID, trackID string, slot int, req models.SetHotCueRequest) (*models.HotCue, error) {
	// Validate slot
	if !models.IsValidSlot(slot) {
		return nil, fmt.Errorf("invalid slot: must be between 1 and %d", models.MaxHotCuesPerTrack)
	}

	// Check feature access
	enabled, err := s.featureSvc.IsEnabled(ctx, userID, models.FeatureHotCues)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, fmt.Errorf("hot cues feature is not enabled for your subscription tier")
	}

	// Get track
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}
	if track == nil {
		return nil, fmt.Errorf("track not found")
	}

	// Validate position is within track duration
	if req.Position > float64(track.Duration) {
		return nil, fmt.Errorf("position exceeds track duration")
	}

	// Initialize hot cues map if nil
	if track.HotCues == nil {
		track.HotCues = make(map[int]*models.HotCue)
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = models.GetDefaultColorForSlot(slot)
	}

	now := time.Now()
	hotCue := &models.HotCue{
		Slot:      slot,
		Position:  req.Position,
		Label:     req.Label,
		Color:     color,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Preserve original creation time if updating existing cue
	if existing, ok := track.HotCues[slot]; ok && existing != nil {
		hotCue.CreatedAt = existing.CreatedAt
	}

	track.HotCues[slot] = hotCue
	track.UpdatedAt = now

	if err := s.repo.UpdateTrack(ctx, *track); err != nil {
		return nil, fmt.Errorf("failed to save hot cue: %w", err)
	}

	return hotCue, nil
}

// DeleteHotCue removes a hot cue from a specific slot
func (s *HotCueService) DeleteHotCue(ctx context.Context, userID, trackID string, slot int) error {
	if !models.IsValidSlot(slot) {
		return fmt.Errorf("invalid slot: must be between 1 and %d", models.MaxHotCuesPerTrack)
	}

	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return fmt.Errorf("failed to get track: %w", err)
	}
	if track == nil {
		return fmt.Errorf("track not found")
	}

	if track.HotCues == nil {
		return nil // No hot cues to delete
	}

	delete(track.HotCues, slot)
	track.UpdatedAt = time.Now()

	if err := s.repo.UpdateTrack(ctx, *track); err != nil {
		return fmt.Errorf("failed to delete hot cue: %w", err)
	}

	return nil
}

// GetHotCues retrieves all hot cues for a track
func (s *HotCueService) GetHotCues(ctx context.Context, userID, trackID string) (*models.TrackHotCuesResponse, error) {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}
	if track == nil {
		return nil, fmt.Errorf("track not found")
	}

	hotCues := make([]models.HotCueResponse, 0)
	if track.HotCues != nil {
		for _, cue := range track.HotCues {
			if cue != nil {
				hotCues = append(hotCues, models.HotCueResponse{
					Slot:      cue.Slot,
					Position:  cue.Position,
					Label:     cue.Label,
					Color:     cue.Color,
					CreatedAt: cue.CreatedAt,
					UpdatedAt: cue.UpdatedAt,
				})
			}
		}
	}

	return &models.TrackHotCuesResponse{
		TrackID:  trackID,
		HotCues:  hotCues,
		MaxSlots: models.MaxHotCuesPerTrack,
	}, nil
}

// ClearAllHotCues removes all hot cues from a track
func (s *HotCueService) ClearAllHotCues(ctx context.Context, userID, trackID string) error {
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return fmt.Errorf("failed to get track: %w", err)
	}
	if track == nil {
		return fmt.Errorf("track not found")
	}

	track.HotCues = make(map[int]*models.HotCue)
	track.UpdatedAt = time.Now()

	if err := s.repo.UpdateTrack(ctx, *track); err != nil {
		return fmt.Errorf("failed to clear hot cues: %w", err)
	}

	return nil
}
