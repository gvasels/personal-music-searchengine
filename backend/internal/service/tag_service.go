package service

import (
	"context"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// TagService handles tag-related operations
type TagService struct {
	repo repository.Repository
}

// NewTagService creates a new TagService
func NewTagService(repo repository.Repository) *TagService {
	return &TagService{repo: repo}
}

// ListTags lists tags with filtering
func (s *TagService) ListTags(ctx context.Context, userID string, filter models.TagFilter) (*models.PaginatedResponse[models.TagResponse], error) {
	result, err := s.repo.ListTags(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.TagResponse, len(result.Items))
	for i, tag := range result.Items {
		responses[i] = tag.ToResponse()
	}

	return &models.PaginatedResponse[models.TagResponse]{
		Items:      responses,
		Pagination: result.Pagination,
	}, nil
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(ctx context.Context, userID string, req models.CreateTagRequest) (*models.Tag, error) {
	now := time.Now()
	tag := &models.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// UpdateTag updates a tag
func (s *TagService) UpdateTag(ctx context.Context, userID, tagName string, req models.UpdateTagRequest) (*models.Tag, error) {
	tag, err := s.repo.GetTag(ctx, userID, tagName)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Color != nil {
		tag.Color = *req.Color
	}

	tag.UpdatedAt = time.Now()

	if err := s.repo.UpdateTag(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(ctx context.Context, userID, tagName string) error {
	return s.repo.DeleteTag(ctx, userID, tagName)
}
