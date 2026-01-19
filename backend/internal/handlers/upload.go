package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetPresignedUploadURL returns a presigned URL for uploading a file
func (h *Handlers) GetPresignedUploadURL(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.PresignedUploadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	// Validate content type
	validContentTypes := map[string]bool{
		"audio/mpeg":    true,
		"audio/flac":    true,
		"audio/wav":     true,
		"audio/aac":     true,
		"audio/ogg":     true,
		"audio/x-flac":  true,
	}
	if !validContentTypes[req.ContentType] {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrUnsupportedMediaType))
	}

	// Max file size: 500MB
	maxFileSize := int64(500 * 1024 * 1024)
	if req.FileSize > maxFileSize {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrPayloadTooLarge))
	}

	result, err := h.uploadService.CreatePresignedUpload(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// ConfirmUpload confirms that a file has been uploaded and triggers processing
func (h *Handlers) ConfirmUpload(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.ConfirmUploadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	result, err := h.uploadService.ConfirmUpload(c.Request().Context(), userID, req.UploadID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// ListUploads returns a list of uploads for the current user
func (h *Handlers) ListUploads(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var filter models.UploadFilter
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	// Set defaults
	if filter.Limit == 0 || filter.Limit > 100 {
		filter.Limit = 50
	}
	if filter.SortBy == "" {
		filter.SortBy = "createdAt"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	result, err := h.uploadService.ListUploads(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}
