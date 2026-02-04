package service

import (
	"context"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// TrackVisibilityRepository defines the repository operations needed for track visibility.
type TrackVisibilityRepository interface {
	ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error)
	ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error)
	GetUserDisplayName(ctx context.Context, userID string) (string, error)
}

// RoleChecker defines the interface for checking user roles and permissions.
type RoleChecker interface {
	GetUserRole(ctx context.Context, userID string) (models.UserRole, error)
	HasPermission(ctx context.Context, userID string, permission models.Permission) (bool, error)
}

// TrackVisibilityService handles track listing with visibility filtering based on user role.
type TrackVisibilityService interface {
	// ListTracksWithVisibility returns tracks visible to the user based on their role.
	// - Admin/GlobalReaders: see all tracks
	// - Regular users: see own tracks + public tracks (if IncludePublic is true)
	ListTracksWithVisibility(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error)
}

// trackVisibilityService implements TrackVisibilityService.
type trackVisibilityService struct {
	repo TrackVisibilityRepository
	role RoleChecker
}

// NewTrackVisibilityService creates a new TrackVisibilityService.
func NewTrackVisibilityService(repo TrackVisibilityRepository, role RoleChecker) TrackVisibilityService {
	return &trackVisibilityService{
		repo: repo,
		role: role,
	}
}

// ListTracksWithVisibility returns tracks visible to the user based on their role.
func (s *trackVisibilityService) ListTracksWithVisibility(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error) {
	// Check if user has global view permission (admin or GlobalReaders)
	role, err := s.role.GetUserRole(ctx, userID)
	if err != nil {
		return nil, err
	}

	canViewGlobal := role == models.RoleAdmin
	if !canViewGlobal {
		hasGlobalPerm, err := s.role.HasPermission(ctx, userID, models.PermissionViewGlobal)
		if err == nil && hasGlobalPerm {
			canViewGlobal = true
		}
	}

	if canViewGlobal {
		return s.listAllTracksForGlobalViewer(ctx, userID, filter)
	}

	return s.listTracksForRegularUser(ctx, userID, filter)
}

// listAllTracksForGlobalViewer returns all tracks for admin/global viewers.
func (s *trackVisibilityService) listAllTracksForGlobalViewer(ctx context.Context, viewerID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error) {
	// Set global scope to get all tracks
	filter.GlobalScope = true

	result, err := s.repo.ListTracks(ctx, "", filter)
	if err != nil {
		return nil, err
	}

	// Convert to responses and add owner display names
	responses := make([]models.TrackResponse, 0, len(result.Items))
	displayNameCache := make(map[string]string)

	for i := range result.Items {
		track := &result.Items[i]

		// Get owner display name
		ownerDisplayName, ok := displayNameCache[track.UserID]
		if !ok {
			name, err := s.repo.GetUserDisplayName(ctx, track.UserID)
			if err != nil {
				name = "Unknown"
			}
			displayNameCache[track.UserID] = name
			ownerDisplayName = name
		}

		// Set "You" for own tracks
		if track.UserID == viewerID {
			ownerDisplayName = "You"
		}

		track.OwnerDisplayName = ownerDisplayName
		responses = append(responses, track.ToResponse(""))
	}

	return &repository.PaginatedResult[models.TrackResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

// listTracksForRegularUser returns own tracks + optionally public tracks for regular users.
func (s *trackVisibilityService) listTracksForRegularUser(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error) {
	// Get user's own tracks
	ownResult, err := s.repo.ListTracks(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Track IDs we've seen for deduplication
	seenIDs := make(map[string]bool)
	responses := make([]models.TrackResponse, 0, len(ownResult.Items))

	// Process own tracks
	for i := range ownResult.Items {
		track := &ownResult.Items[i]
		seenIDs[track.ID] = true
		track.OwnerDisplayName = "You"
		responses = append(responses, track.ToResponse(""))
	}

	// If IncludePublic is true, also fetch public tracks from other users
	if filter.IncludePublic {
		publicResult, err := s.repo.ListPublicTracks(ctx, filter.Limit, "")
		if err != nil {
			// Log error but don't fail - just return own tracks
			return &repository.PaginatedResult[models.TrackResponse]{
				Items:      responses,
				NextCursor: ownResult.NextCursor,
				HasMore:    ownResult.HasMore,
			}, nil
		}

		displayNameCache := make(map[string]string)

		for i := range publicResult.Items {
			track := &publicResult.Items[i]

			// Skip if already seen (user's own public tracks)
			if seenIDs[track.ID] {
				continue
			}
			seenIDs[track.ID] = true

			// Get owner display name
			ownerDisplayName, ok := displayNameCache[track.UserID]
			if !ok {
				name, err := s.repo.GetUserDisplayName(ctx, track.UserID)
				if err != nil {
					name = "Unknown"
				}
				displayNameCache[track.UserID] = name
				ownerDisplayName = name
			}

			track.OwnerDisplayName = ownerDisplayName
			responses = append(responses, track.ToResponse(""))
		}
	}

	return &repository.PaginatedResult[models.TrackResponse]{
		Items:      responses,
		NextCursor: ownResult.NextCursor,
		HasMore:    ownResult.HasMore,
	}, nil
}
