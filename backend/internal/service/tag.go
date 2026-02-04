package service

import (
	"context"
	"strings"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// normalizeTagName converts tag name to lowercase for consistent storage/lookup
func normalizeTagName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

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
	// Normalize tag name to lowercase
	normalizedName := normalizeTagName(req.Name)

	// Check if tag already exists
	existing, err := s.repo.GetTag(ctx, userID, normalizedName)
	if err == nil {
		return nil, models.NewConflictError("Tag with this name already exists")
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	now := time.Now()
	tag := models.Tag{
		UserID:     userID,
		Name:       normalizedName,
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
	normalizedName := normalizeTagName(tagName)

	tag, err := s.repo.GetTag(ctx, userID, normalizedName)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Tag", normalizedName)
		}
		return nil, err
	}

	response := tag.ToResponse()
	return &response, nil
}

func (s *tagService) UpdateTag(ctx context.Context, userID, tagName string, req models.UpdateTagRequest) (*models.TagResponse, error) {
	normalizedName := normalizeTagName(tagName)

	tag, err := s.repo.GetTag(ctx, userID, normalizedName)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Tag", normalizedName)
		}
		return nil, err
	}

	// Apply updates
	if req.Name != nil {
		normalizedNewName := normalizeTagName(*req.Name)
		if normalizedNewName != tag.Name {
			// Renaming a tag - need to check if new name exists
			_, err := s.repo.GetTag(ctx, userID, normalizedNewName)
			if err == nil {
				return nil, models.NewConflictError("Tag with this name already exists")
			}
			if err != repository.ErrNotFound {
				return nil, err
			}
			tag.Name = normalizedNewName
		}
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
	normalizedName := normalizeTagName(tagName)

	_, err := s.repo.GetTag(ctx, userID, normalizedName)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("Tag", normalizedName)
		}
		return err
	}

	return s.repo.DeleteTag(ctx, userID, normalizedName)
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

	// Normalize all tag names
	normalizedTags := make([]string, 0, len(req.Tags))
	for _, tagName := range req.Tags {
		normalizedTags = append(normalizedTags, normalizeTagName(tagName))
	}

	// Build set of existing tags on this track
	existingTagSet := make(map[string]bool)
	for _, t := range track.Tags {
		existingTagSet[t] = true
	}

	// Determine which tags are actually new (not already on track)
	newTagsToAdd := make([]string, 0)
	for _, tagName := range normalizedTags {
		if !existingTagSet[tagName] {
			newTagsToAdd = append(newTagsToAdd, tagName)
		}
	}

	// Create tags that don't exist and increment TrackCount for new associations
	for _, tagName := range newTagsToAdd {
		tag, err := s.repo.GetTag(ctx, userID, tagName)
		if err == repository.ErrNotFound {
			// Create new tag with TrackCount = 1
			now := time.Now()
			newTag := models.Tag{
				UserID:     userID,
				Name:       tagName,
				TrackCount: 1,
			}
			newTag.CreatedAt = now
			newTag.UpdatedAt = now
			_ = s.repo.CreateTag(ctx, newTag) // Ignore errors, tag might have been created by another request
		} else if err == nil {
			// Tag exists, increment TrackCount
			tag.TrackCount++
			tag.UpdatedAt = time.Now()
			_ = s.repo.UpdateTag(ctx, *tag)
		}
	}

	// Add tags to track (only if there are new tags to add)
	if len(newTagsToAdd) > 0 {
		if err := s.repo.AddTagsToTrack(ctx, userID, trackID, newTagsToAdd); err != nil {
			return nil, err
		}
	}

	// Update track's tags list
	for _, t := range normalizedTags {
		existingTagSet[t] = true
	}

	allTags := make([]string, 0, len(existingTagSet))
	for t := range existingTagSet {
		allTags = append(allTags, t)
	}

	track.Tags = allTags
	if err := s.repo.UpdateTrack(ctx, *track); err != nil {
		return nil, err
	}

	return allTags, nil
}

func (s *tagService) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	normalizedName := normalizeTagName(tagName)

	// Verify track exists
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("Track", trackID)
		}
		return err
	}

	// Check if track actually has this tag
	hasTag := false
	for _, t := range track.Tags {
		if t == normalizedName {
			hasTag = true
			break
		}
	}

	if !hasTag {
		// Tag not on track, nothing to do
		return nil
	}

	// Remove tag association
	if err := s.repo.RemoveTagFromTrack(ctx, userID, trackID, normalizedName); err != nil {
		return err
	}

	// Decrement TrackCount on the tag
	tag, err := s.repo.GetTag(ctx, userID, normalizedName)
	if err == nil && tag.TrackCount > 0 {
		tag.TrackCount--
		tag.UpdatedAt = time.Now()
		_ = s.repo.UpdateTag(ctx, *tag)
	}

	// Update track's tags list
	newTags := make([]string, 0, len(track.Tags))
	for _, t := range track.Tags {
		if t != normalizedName {
			newTags = append(newTags, t)
		}
	}

	track.Tags = newTags
	return s.repo.UpdateTrack(ctx, *track)
}

func (s *tagService) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.TrackResponse, error) {
	normalizedName := normalizeTagName(tagName)

	// Verify tag exists
	_, err := s.repo.GetTag(ctx, userID, normalizedName)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Tag", normalizedName)
		}
		return nil, err
	}

	tracks, err := s.repo.GetTracksByTag(ctx, userID, normalizedName)
	if err != nil {
		return nil, err
	}

	responses := make([]models.TrackResponse, 0, len(tracks))
	for _, track := range tracks {
		responses = append(responses, track.ToResponse(""))
	}

	return responses, nil
}
