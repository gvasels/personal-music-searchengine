package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// CreatePresignedUpload generates a presigned URL for file upload
func (h *Handlers) CreatePresignedUpload(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var req models.PresignedUploadRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.services.Upload.CreatePresignedUpload(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}

// ConfirmUpload confirms an upload and triggers processing
func (h *Handlers) ConfirmUpload(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var req models.ConfirmUploadRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.services.Upload.ConfirmUpload(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}

// CompleteMultipartUpload completes a multipart upload
func (h *Handlers) CompleteMultipartUpload(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var req models.CompleteMultipartUploadRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.services.Upload.CompleteMultipartUpload(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}

// ListUploads returns a paginated list of uploads
func (h *Handlers) ListUploads(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var filter models.UploadFilter
	if err := c.Bind(&filter); err != nil {
		return handleError(c, models.ErrBadRequest)
	}

	uploads, err := h.services.Upload.ListUploads(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, uploads)
}

// GetUploadStatus returns the status of an upload
func (h *Handlers) GetUploadStatus(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	uploadID := c.Param("id")
	if uploadID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	upload, err := h.services.Upload.GetUploadStatus(c.Request().Context(), userID, uploadID)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, upload)
}

// ReprocessUpload retries a failed upload from a specific step
func (h *Handlers) ReprocessUpload(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	uploadID := c.Param("id")
	if uploadID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.ReprocessUploadRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	upload, err := h.services.Upload.ReprocessUpload(c.Request().Context(), userID, uploadID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, upload)
}
