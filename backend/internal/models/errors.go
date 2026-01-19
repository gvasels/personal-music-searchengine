package models

import (
	"fmt"
	"net/http"
)

// APIError represents a structured API error response
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common API errors
var (
	ErrNotFound = &APIError{
		Code:       "NOT_FOUND",
		Message:    "The requested resource was not found",
		StatusCode: http.StatusNotFound,
	}

	ErrUnauthorized = &APIError{
		Code:       "UNAUTHORIZED",
		Message:    "Authentication is required",
		StatusCode: http.StatusUnauthorized,
	}

	ErrForbidden = &APIError{
		Code:       "FORBIDDEN",
		Message:    "You do not have permission to access this resource",
		StatusCode: http.StatusForbidden,
	}

	ErrBadRequest = &APIError{
		Code:       "BAD_REQUEST",
		Message:    "The request was invalid",
		StatusCode: http.StatusBadRequest,
	}

	ErrInternalServer = &APIError{
		Code:       "INTERNAL_ERROR",
		Message:    "An internal server error occurred",
		StatusCode: http.StatusInternalServerError,
	}

	ErrConflict = &APIError{
		Code:       "CONFLICT",
		Message:    "The request conflicts with existing data",
		StatusCode: http.StatusConflict,
	}

	ErrPayloadTooLarge = &APIError{
		Code:       "PAYLOAD_TOO_LARGE",
		Message:    "The file size exceeds the maximum allowed",
		StatusCode: http.StatusRequestEntityTooLarge,
	}

	ErrUnsupportedMediaType = &APIError{
		Code:       "UNSUPPORTED_MEDIA_TYPE",
		Message:    "The file format is not supported",
		StatusCode: http.StatusUnsupportedMediaType,
	}

	ErrStorageLimitExceeded = &APIError{
		Code:       "STORAGE_LIMIT_EXCEEDED",
		Message:    "Your storage limit has been exceeded",
		StatusCode: http.StatusPaymentRequired,
	}

	// Upload-specific errors
	ErrUploadExpired = &APIError{
		Code:       "UPLOAD_EXPIRED",
		Message:    "The upload session has expired",
		StatusCode: http.StatusGone,
	}

	ErrUploadNotFound = &APIError{
		Code:       "UPLOAD_NOT_FOUND",
		Message:    "The upload session was not found",
		StatusCode: http.StatusNotFound,
	}

	ErrMultipartUploadFailed = &APIError{
		Code:       "MULTIPART_UPLOAD_FAILED",
		Message:    "One or more parts of the multipart upload failed",
		StatusCode: http.StatusBadRequest,
	}

	ErrUploadProcessingFailed = &APIError{
		Code:       "UPLOAD_PROCESSING_FAILED",
		Message:    "Upload processing failed",
		StatusCode: http.StatusInternalServerError,
	}

	ErrInvalidCursor = &APIError{
		Code:       "INVALID_CURSOR",
		Message:    "The pagination cursor is invalid or expired",
		StatusCode: http.StatusBadRequest,
	}
)

// NewAPIError creates a new API error
func NewAPIError(code, message string, statusCode int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewValidationError creates a validation error with details
func NewValidationError(details any) *APIError {
	return &APIError{
		Code:       "VALIDATION_ERROR",
		Message:    "The request failed validation",
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a not found error for a specific resource
func NewNotFoundError(resource, id string) *APIError {
	return &APIError{
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s with ID '%s' was not found", resource, id),
		StatusCode: http.StatusNotFound,
	}
}

// NewConflictError creates a conflict error with details
func NewConflictError(message string) *APIError {
	return &APIError{
		Code:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error *APIError `json:"error"`
}

// NewErrorResponse creates an error response
func NewErrorResponse(err *APIError) ErrorResponse {
	return ErrorResponse{Error: err}
}
