package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// UploadService handles upload-related operations
type UploadService struct {
	repo        repository.Repository
	s3Client    *s3.Client
	mediaBucket string
}

// NewUploadService creates a new UploadService
func NewUploadService(repo repository.Repository, s3Client *s3.Client) *UploadService {
	mediaBucket := os.Getenv("MEDIA_BUCKET")
	if mediaBucket == "" {
		mediaBucket = "music-library-media"
	}

	return &UploadService{
		repo:        repo,
		s3Client:    s3Client,
		mediaBucket: mediaBucket,
	}
}

// CreatePresignedUpload creates a presigned URL for uploading
func (s *UploadService) CreatePresignedUpload(ctx context.Context, userID string, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error) {
	uploadID := uuid.New().String()
	s3Key := fmt.Sprintf("uploads/%s/%s/%s", userID, uploadID, req.FileName)

	// Create upload record
	now := time.Now()
	upload := &models.Upload{
		ID:          uploadID,
		UserID:      userID,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		ContentType: req.ContentType,
		S3Key:       s3Key,
		Status:      models.UploadStatusPending,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err := s.repo.CreateUpload(ctx, upload); err != nil {
		return nil, err
	}

	// Generate presigned URL
	presignClient := s3.NewPresignClient(s.s3Client)

	presignedReq, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.mediaBucket),
		Key:         aws.String(s3Key),
		ContentType: aws.String(req.ContentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create presigned URL: %w", err)
	}

	return &models.PresignedUploadResponse{
		UploadID:    uploadID,
		UploadURL:   presignedReq.URL,
		ExpiresAt:   now.Add(15 * time.Minute),
		MaxFileSize: 500 * 1024 * 1024, // 500MB
	}, nil
}

// ConfirmUpload confirms that a file has been uploaded
func (s *UploadService) ConfirmUpload(ctx context.Context, userID, uploadID string) (*models.ConfirmUploadResponse, error) {
	upload, err := s.repo.GetUpload(ctx, userID, uploadID)
	if err != nil {
		return nil, err
	}

	if upload.Status != models.UploadStatusPending {
		return nil, models.NewConflictError("Upload has already been confirmed")
	}

	// Verify file exists in S3
	_, err = s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.mediaBucket),
		Key:    aws.String(upload.S3Key),
	})
	if err != nil {
		return nil, models.NewNotFoundError("Uploaded file", upload.S3Key)
	}

	// Update status to processing
	upload.Status = models.UploadStatusProcessing
	upload.UpdatedAt = time.Now()

	if err := s.repo.UpdateUpload(ctx, upload); err != nil {
		return nil, err
	}

	return &models.ConfirmUploadResponse{
		UploadID: uploadID,
		Status:   models.UploadStatusProcessing,
		Message:  "Upload confirmed, processing metadata",
	}, nil
}

// ListUploads lists uploads with filtering
func (s *UploadService) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResponse[models.UploadResponse], error) {
	result, err := s.repo.ListUploads(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.UploadResponse, len(result.Items))
	for i, upload := range result.Items {
		responses[i] = upload.ToResponse()
	}

	return &models.PaginatedResponse[models.UploadResponse]{
		Items:      responses,
		Pagination: result.Pagination,
	}, nil
}
