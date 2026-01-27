package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ListTracks returns a paginated list of tracks
// If user has GLOBAL permission, returns tracks from all users
func (h *Handlers) ListTracks(c echo.Context) error {
	auth := getAuthContext(c)
	if auth.UserID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var filter models.TrackFilter
	if err := c.Bind(&filter); err != nil {
		return handleError(c, models.ErrBadRequest)
	}

	// Set global scope if user has GLOBAL permission
	filter.GlobalScope = auth.HasGlobal

	tracks, err := h.services.Track.ListTracks(c.Request().Context(), auth.UserID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, tracks)
}

// GetTrack returns a single track by ID
// Respects visibility rules: owners and admins can access any track,
// others can only access public/unlisted tracks
func (h *Handlers) GetTrack(c echo.Context) error {
	auth := getAuthContext(c)
	if auth.UserID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("id")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	track, err := h.services.Track.GetTrack(c.Request().Context(), auth.UserID, trackID, auth.HasGlobal)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, track)
}

// UpdateTrack updates a track's metadata
func (h *Handlers) UpdateTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("id")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.UpdateTrackRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	track, err := h.services.Track.UpdateTrack(c.Request().Context(), userID, trackID, req)
	if err != nil {
		return handleError(c, err)
	}

	// Re-index in search to reflect metadata changes (best effort)
	if h.services.Search != nil {
		trackModel := models.Track{
			ID:       track.ID,
			UserID:   userID,
			Title:    track.Title,
			Artist:   track.Artist,
			Album:    track.Album,
			Genre:    track.Genre,
			Year:     track.Year,
			Duration: track.Duration,
		}
		_ = h.services.Search.IndexTrack(c.Request().Context(), trackModel)
	}

	return success(c, track)
}

// DeleteTrack deletes a track
func (h *Handlers) DeleteTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("id")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	if err := h.services.Track.DeleteTrack(c.Request().Context(), userID, trackID); err != nil {
		return handleError(c, err)
	}

	// Remove from search index (best effort - don't fail if search removal fails)
	if h.services.Search != nil {
		_ = h.services.Search.RemoveTrack(c.Request().Context(), trackID)
	}

	return noContent(c)
}

// AddTagsToTrack adds tags to a track
func (h *Handlers) AddTagsToTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("id")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.AddTagsToTrackRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	tags, err := h.services.Tag.AddTagsToTrack(c.Request().Context(), userID, trackID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, map[string][]string{"tags": tags})
}

// RemoveTagFromTrack removes a tag from a track
func (h *Handlers) RemoveTagFromTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("id")
	tagName := c.Param("tag")
	if trackID == "" || tagName == "" {
		return handleError(c, models.ErrBadRequest)
	}

	if err := h.services.Tag.RemoveTagFromTrack(c.Request().Context(), userID, trackID, tagName); err != nil {
		return handleError(c, err)
	}

	return noContent(c)
}

// UploadCoverArt generates a presigned URL for cover art upload
func (h *Handlers) UploadCoverArt(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("id")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.CoverArtUploadRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.services.Upload.UploadCoverArt(c.Request().Context(), userID, trackID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}

// UpdateTrackVisibilityRequest is the request body for updating track visibility.
type UpdateTrackVisibilityRequest struct {
	Visibility models.TrackVisibility `json:"visibility" validate:"required"`
}

// UpdateTrackVisibility updates the visibility of a track
func (h *Handlers) UpdateTrackVisibility(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("id")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req UpdateTrackVisibilityRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	err := h.services.Track.UpdateVisibility(c.Request().Context(), userID, trackID, req.Visibility)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, map[string]interface{}{
		"trackId":    trackID,
		"visibility": req.Visibility,
	})
}
