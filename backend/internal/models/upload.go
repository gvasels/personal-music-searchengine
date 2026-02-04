package models

import (
	"fmt"
	"time"
)

// Upload represents a file upload and its processing status
type Upload struct {
	ID          string       `json:"id" dynamodbav:"id"`
	UserID      string       `json:"userId" dynamodbav:"userId"`
	FileName    string       `json:"fileName" dynamodbav:"fileName"`
	FileSize    int64        `json:"fileSize" dynamodbav:"fileSize"`
	ContentType string       `json:"contentType" dynamodbav:"contentType"`
	S3Key       string       `json:"s3Key" dynamodbav:"s3Key"`
	Status      UploadStatus `json:"status" dynamodbav:"status"`
	ErrorMsg    string       `json:"errorMsg,omitempty" dynamodbav:"errorMsg,omitempty"`
	TrackID     string       `json:"trackId,omitempty" dynamodbav:"trackId,omitempty"` // Set after successful processing
	Timestamps
	CompletedAt *time.Time `json:"completedAt,omitempty" dynamodbav:"completedAt,omitempty"`

	// Step tracking for partial success recovery
	MetadataExtracted bool `json:"metadataExtracted" dynamodbav:"metadataExtracted"`
	CoverArtExtracted bool `json:"coverArtExtracted" dynamodbav:"coverArtExtracted"`
	TrackCreated      bool `json:"trackCreated" dynamodbav:"trackCreated"`
	Indexed           bool `json:"indexed" dynamodbav:"indexed"`
	FileMoved         bool `json:"fileMoved" dynamodbav:"fileMoved"`

	// Multipart upload tracking
	IsMultipart   bool   `json:"isMultipart,omitempty" dynamodbav:"isMultipart,omitempty"`
	MultipartID   string `json:"multipartId,omitempty" dynamodbav:"multipartId,omitempty"`
	PartsUploaded int    `json:"partsUploaded,omitempty" dynamodbav:"partsUploaded,omitempty"`
	TotalParts    int    `json:"totalParts,omitempty" dynamodbav:"totalParts,omitempty"`
}

// UploadItem represents an Upload in DynamoDB single-table design
type UploadItem struct {
	DynamoDBItem
	Upload
}

// NewUploadItem creates a DynamoDB item for an upload
func NewUploadItem(upload Upload) UploadItem {
	return UploadItem{
		DynamoDBItem: DynamoDBItem{
			PK:     fmt.Sprintf("USER#%s", upload.UserID),
			SK:     fmt.Sprintf("UPLOAD#%s", upload.ID),
			GSI1PK: fmt.Sprintf("UPLOAD#STATUS#%s", upload.Status),
			GSI1SK: upload.CreatedAt.Format(time.RFC3339),
			Type:   string(EntityUpload),
		},
		Upload: upload,
	}
}

// PresignedUploadRequest represents a request to get a presigned URL for uploading
type PresignedUploadRequest struct {
	FileName    string `json:"fileName" validate:"required,min=1,max=500"`
	FileSize    int64  `json:"fileSize" validate:"required,min=1,max=1073741824"` // max 1GB
	ContentType string `json:"contentType" validate:"required,oneof=audio/mpeg audio/flac audio/wav audio/aac audio/ogg audio/x-flac"`
	IsMultipart bool   `json:"isMultipart,omitempty"` // Request multipart upload for large files
}

// PresignedUploadResponse represents a response with presigned URL for uploading
type PresignedUploadResponse struct {
	UploadID     string            `json:"uploadId"`
	UploadURL    string            `json:"uploadUrl"`
	Fields       map[string]string `json:"fields,omitempty"` // For POST uploads
	ExpiresAt    time.Time         `json:"expiresAt"`
	MaxFileSize  int64             `json:"maxFileSize"`

	// Multipart upload fields
	IsMultipart bool                     `json:"isMultipart,omitempty"`
	MultipartID string                   `json:"multipartId,omitempty"`
	PartURLs    []MultipartUploadPartURL `json:"partUrls,omitempty"` // Presigned URLs for each part
}

// MultipartUploadPartURL represents a presigned URL for a single multipart upload part
type MultipartUploadPartURL struct {
	PartNumber int       `json:"partNumber"`
	UploadURL  string    `json:"uploadUrl"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

// ConfirmUploadRequest represents a request to confirm an upload
type ConfirmUploadRequest struct {
	UploadID string `json:"uploadId" validate:"required,uuid"`
}

// ConfirmUploadResponse represents a response after confirming an upload
type ConfirmUploadResponse struct {
	UploadID string       `json:"uploadId"`
	Status   UploadStatus `json:"status"`
	Message  string       `json:"message"`
}

// UploadResponse represents an upload in API responses
type UploadResponse struct {
	ID          string       `json:"id"`
	FileName    string       `json:"fileName"`
	FileSize    int64        `json:"fileSize"`
	FileSizeStr string       `json:"fileSizeStr"`
	ContentType string       `json:"contentType"`
	Status      UploadStatus `json:"status"`
	ErrorMsg    string       `json:"errorMsg,omitempty"`
	TrackID     string       `json:"trackId,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	CompletedAt *time.Time   `json:"completedAt,omitempty"`

	// Step tracking for partial success
	Steps UploadSteps `json:"steps"`
}

// UploadSteps represents the completion status of each processing step
type UploadSteps struct {
	MetadataExtracted bool `json:"metadataExtracted"`
	CoverArtExtracted bool `json:"coverArtExtracted"`
	TrackCreated      bool `json:"trackCreated"`
	Indexed           bool `json:"indexed"`
	FileMoved         bool `json:"fileMoved"`
}

// ToResponse converts an Upload to an UploadResponse
func (u *Upload) ToResponse() UploadResponse {
	return UploadResponse{
		ID:          u.ID,
		FileName:    u.FileName,
		FileSize:    u.FileSize,
		FileSizeStr: formatFileSize(u.FileSize),
		ContentType: u.ContentType,
		Status:      u.Status,
		ErrorMsg:    u.ErrorMsg,
		TrackID:     u.TrackID,
		CreatedAt:   u.CreatedAt,
		CompletedAt: u.CompletedAt,
		Steps: UploadSteps{
			MetadataExtracted: u.MetadataExtracted,
			CoverArtExtracted: u.CoverArtExtracted,
			TrackCreated:      u.TrackCreated,
			Indexed:           u.Indexed,
			FileMoved:         u.FileMoved,
		},
	}
}

// UploadFilter represents filter options for listing uploads
type UploadFilter struct {
	Status    UploadStatus `query:"status"`
	SortBy    string       `query:"sortBy"`    // createdAt, fileName, fileSize
	SortOrder string       `query:"sortOrder"` // asc, desc
	Limit     int          `query:"limit"`
	LastKey   string       `query:"lastKey"`
}

// UploadMetadata represents metadata extracted from uploaded audio files
type UploadMetadata struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	AlbumArtist string `json:"albumArtist,omitempty"`
	Album       string `json:"album,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Year        int    `json:"year,omitempty"`
	TrackNumber int    `json:"trackNumber,omitempty"`
	DiscNumber  int    `json:"discNumber,omitempty"`
	Duration    int    `json:"duration"` // seconds
	Bitrate     int    `json:"bitrate,omitempty"`
	SampleRate  int    `json:"sampleRate,omitempty"`
	Channels    int    `json:"channels,omitempty"`
	Format      string `json:"format"`
	HasCoverArt bool   `json:"hasCoverArt"`
	Composer    string `json:"composer,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Lyrics      string `json:"lyrics,omitempty"`
}

// ProcessingStep represents a step in the upload processing pipeline
type ProcessingStep string

const (
	StepExtractMetadata ProcessingStep = "extract_metadata"
	StepExtractCover    ProcessingStep = "extract_cover"
	StepCreateTrack     ProcessingStep = "create_track"
	StepIndex           ProcessingStep = "index"
	StepMoveFile        ProcessingStep = "move_file"
)

// ReprocessUploadRequest represents a request to reprocess a failed upload
type ReprocessUploadRequest struct {
	FromStep ProcessingStep `json:"fromStep,omitempty" validate:"omitempty,oneof=extract_metadata extract_cover create_track index move_file"`
}

// CoverArtUploadRequest represents a request to upload cover art for a track
type CoverArtUploadRequest struct {
	FileName    string `json:"fileName" validate:"required,min=1,max=500"`
	FileSize    int64  `json:"fileSize" validate:"required,min=1,max=10485760"` // max 10MB
	ContentType string `json:"contentType" validate:"required,oneof=image/jpeg image/png image/webp"`
}

// CoverArtUploadResponse represents a response with presigned URL for cover art upload
type CoverArtUploadResponse struct {
	TrackID     string    `json:"trackId"`
	UploadURL   string    `json:"uploadUrl"`
	ExpiresAt   time.Time `json:"expiresAt"`
	MaxFileSize int64     `json:"maxFileSize"`
}

// CompleteMultipartUploadRequest represents a request to complete a multipart upload
type CompleteMultipartUploadRequest struct {
	UploadID string               `json:"uploadId" validate:"required,uuid"`
	Parts    []CompletedPartInfo `json:"parts" validate:"required,min=1"`
}

// CompletedPartInfo represents information about a completed multipart upload part
type CompletedPartInfo struct {
	PartNumber int    `json:"partNumber" validate:"required,min=1"`
	ETag       string `json:"etag" validate:"required"`
}
