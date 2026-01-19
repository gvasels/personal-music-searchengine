package service

import (
	"context"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// tagService implements TagService
type tagService struct {
	repo repository.Repository
}

// NewTagService creates a new tag service
func NewTagService(repo repository.Repository) TagService {
	return &tagService{
		repo: repo,
	}
}

func (s *tagService) CreateTag(ctx context.Context, userID string, req models.CreateTagRequest) (*models.TagResponse, error) {
	// Check if tag already exists
	existing, err := s.repo.GetTag(ctx, userID, req.Name)
	if err == nil {
		return nil, models.NewConflictError("Tag with this name already exists")
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	now := time.Now()
	tag := models.Tag{
		UserID:     userID,
		Name:       req.Name,
		Color:      req.Color,
		TrackCount: 0,
	}
	tag.CreatedAt = now
	tag.UpdatedAt = now

	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return nil, err
	}

	_ = existing // silence unused variable warning
	response := tag.ToResponse()
	return &response, nil
}

func (s *tagService) GetTag(ctx context.Context, userID, tagName string) (*models.TagResponse, error) {
	tag, err := s.repo.GetTag(ctx, userID, tagName)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Tag", tagName)
		}
		return nil, err
	}

	response := tag.ToResponse()
	return &response, nil
}

func (s *tagService) UpdateTag(ctx context.Context, userID, tagName string, req models.UpdateTagRequest) (*models.TagResponse, error) {
	tag, err := s.repo.GetTag(ctx, userID, tagName)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Tag", tagName)
		}
		return nil, err
	}

	// Apply updates
	if req.Name != nil && *req.Name != tag.Name {
		// Renaming a tag - need to check if new name exists
		_, err := s.repo.GetTag(ctx, userID, *req.Name)
		if err == nil {
			return nil, models.NewConflictError("Tag with this name already exists")
		}
		if err != repository.ErrNotFound {
			return nil, err
		}
		tag.Name = *req.Name
	}
	if req.Color != nil {
		tag.Color = *req.Color
	}

	if err := s.repo.UpdateTag(ctx, *tag); err != nil {
		return nil, err
	}

	response := tag.ToResponse()
	return &response, nil
}

func (s *tagService) DeleteTag(ctx context.Context, userID, tagName string) error {
	_, err := s.repo.GetTag(ctx, userID, tagName)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("Tag", tagName)
		}
		return err
	}

	return s.repo.DeleteTag(ctx, userID, tagName)
}

func (s *tagService) ListTags(ctx context.Context, userID string) ([]models.TagResponse, error) {
	tags, err := s.repo.ListTags(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.TagResponse, 0, len(tags))
	for _, tag := range tags {
		responses = append(responses, tag.ToResponse())
	}

	return responses, nil
}

func (s *tagService) AddTagsToTrack(ctx context.Context, userID, trackID string, req models.AddTagsToTrackRequest) ([]string, error) {
	// Verify track exists
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Track", trackID)
		}
		return nil, err
	}

	// Create tags that don't exist
	for _, tagName := range req.Tags {
		_, err := s.repo.GetTag(ctx, userID, tagName)
		if err == repository.ErrNotFound {
			now := time.Now()
			tag := models.Tag{
				UserID:     userID,
				Name:       tagName,
				TrackCount: 0,
			}
			tag.CreatedAt = now
			tag.UpdatedAt = now
			_ = s.repo.CreateTag(ctx, tag) // Ignore errors, tag might have been created by another request
		}
	}

	// Add tags to track
	if err := s.repo.AddTagsToTrack(ctx, userID, trackID, req.Tags); err != nil {
		return nil, err
	}

	// Update track's tags list
	existingTags := make(map[string]bool)
	for _, t := range track.Tags {
		existingTags[t] = true
	}
	for _, t := range req.Tags {
		existingTags[t] = true
	}

	newTags := make([]string, 0, len(existingTags))
	for t := range existingTags {
		newTags = append(newTags, t)
	}

	track.Tags = newTags
	if err := s.repo.UpdateTrack(ctx, *track); err != nil {
		return nil, err
	}

	return newTags, nil
}

func (s *tagService) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	// Verify track exists
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("Track", trackID)
		}
		return err
	}

	// Remove tag association
	if err := s.repo.RemoveTagFromTrack(ctx, userID, trackID, tagName); err != nil {
		return err
	}

	// Update track's tags list
	newTags := make([]string, 0, len(track.Tags))
	for _, t := range track.Tags {
		if t != tagName {
			newTags = append(newTags, t)
		}
	}

	track.Tags = newTags
	return s.repo.UpdateTrack(ctx, *track)
}

func (s *tagService) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.TrackResponse, error) {
	// Verify tag exists
	_, err := s.repo.GetTag(ctx, userID, tagName)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Tag", tagName)
		}
		return nil, err
	}

	tracks, err := s.repo.GetTracksByTag(ctx, userID, tagName)
	if err != nil {
		return nil, err
	}

	responses := make([]models.TrackResponse, 0, len(tracks))
	for _, track := range tracks {
		responses = append(responses, track.ToResponse(""))
	}

	return responses, nil
}
