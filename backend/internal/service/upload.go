package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

const (
	// Upload URL expiration times
	uploadURLExpiry = 15 * time.Minute
	partURLExpiry   = 1 * time.Hour

	// Multipart upload thresholds
	multipartThreshold = 100 * 1024 * 1024 // 100 MB
	partSize           = 5 * 1024 * 1024   // 5 MB parts
)

// uploadService implements UploadService
type uploadService struct {
	repo             repository.Repository
	s3Repo           repository.S3Repository
	mediaBucket      string
	stepFunctionsARN string
}

// NewUploadService creates a new upload service
func NewUploadService(repo repository.Repository, s3Repo repository.S3Repository, mediaBucket string, stepFunctionsARN string) UploadService {
	return &uploadService{
		repo:             repo,
		s3Repo:           s3Repo,
		mediaBucket:      mediaBucket,
		stepFunctionsARN: stepFunctionsARN,
	}
}

func (s *uploadService) CreatePresignedUpload(ctx context.Context, userID string, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error) {
	// Check user storage limit
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}
	if user != nil && user.StorageUsed+req.FileSize > user.StorageLimit {
		return nil, models.ErrStorageLimitExceeded
	}

	// Generate upload ID and S3 key
	uploadID := uuid.New().String()
	s3Key := fmt.Sprintf("uploads/%s/%s/%s", userID, uploadID, req.FileName)

	// Create upload record
	now := time.Now()
	upload := models.Upload{
		ID:          uploadID,
		UserID:      userID,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		ContentType: req.ContentType,
		S3Key:       s3Key,
		Status:      models.UploadStatusPending,
		IsMultipart: req.IsMultipart || req.FileSize > multipartThreshold,
	}
	upload.CreatedAt = now
	upload.UpdatedAt = now

	if err := s.repo.CreateUpload(ctx, upload); err != nil {
		return nil, err
	}

	response := &models.PresignedUploadResponse{
		UploadID:    uploadID,
		ExpiresAt:   now.Add(uploadURLExpiry),
		MaxFileSize: req.FileSize,
		IsMultipart: upload.IsMultipart,
	}

	// Generate presigned URL(s)
	if upload.IsMultipart {
		// Initiate multipart upload
		multipartID, err := s.s3Repo.InitiateMultipartUpload(ctx, s3Key, req.ContentType)
		if err != nil {
			return nil, fmt.Errorf("failed to initiate multipart upload: %w", err)
		}

		// Update upload with multipart ID
		upload.MultipartID = multipartID
		numParts := int(req.FileSize/partSize) + 1
		upload.TotalParts = numParts
		if err := s.repo.UpdateUpload(ctx, upload); err != nil {
			return nil, err
		}

		// Generate presigned URLs for parts
		partURLs, err := s.s3Repo.GenerateMultipartUploadURLs(ctx, s3Key, multipartID, numParts, partURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate part URLs: %w", err)
		}

		response.MultipartID = multipartID
		response.PartURLs = partURLs
	} else {
		// Generate single presigned URL
		uploadURL, err := s.s3Repo.GeneratePresignedUploadURL(ctx, s3Key, req.ContentType, uploadURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
		}
		response.UploadURL = uploadURL
	}

	return response, nil
}

func (s *uploadService) ConfirmUpload(ctx context.Context, userID string, req models.ConfirmUploadRequest) (*models.ConfirmUploadResponse, error) {
	upload, err := s.repo.GetUpload(ctx, userID, req.UploadID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.ErrUploadNotFound
		}
		return nil, err
	}

	if upload.Status != models.UploadStatusPending {
		return nil, models.NewValidationError(map[string]string{
			"status": fmt.Sprintf("Upload is already %s", upload.Status),
		})
	}

	// Verify file exists in S3
	exists, err := s.s3Repo.ObjectExists(ctx, upload.S3Key)
	if err != nil {
		return nil, fmt.Errorf("failed to verify upload: %w", err)
	}
	if !exists {
		return nil, models.NewValidationError(map[string]string{
			"file": "File not found in upload location",
		})
	}

	// Update status to processing
	if err := s.repo.UpdateUploadStatus(ctx, userID, req.UploadID, models.UploadStatusProcessing, "", ""); err != nil {
		return nil, err
	}

	// TODO: Trigger Step Functions workflow for processing
	// For now, return processing status

	return &models.ConfirmUploadResponse{
		UploadID: req.UploadID,
		Status:   models.UploadStatusProcessing,
		Message:  "Upload confirmed, processing started",
	}, nil
}

func (s *uploadService) CompleteMultipartUpload(ctx context.Context, userID string, req models.CompleteMultipartUploadRequest) (*models.ConfirmUploadResponse, error) {
	upload, err := s.repo.GetUpload(ctx, userID, req.UploadID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.ErrUploadNotFound
		}
		return nil, err
	}

	if !upload.IsMultipart {
		return nil, models.NewValidationError(map[string]string{
			"multipart": "This is not a multipart upload",
		})
	}

	if upload.MultipartID == "" {
		return nil, models.NewValidationError(map[string]string{
			"multipart": "Multipart upload ID not found",
		})
	}

	// Complete the multipart upload
	if err := s.s3Repo.CompleteMultipartUpload(ctx, upload.S3Key, upload.MultipartID, req.Parts); err != nil {
		return nil, models.ErrMultipartUploadFailed
	}

	// Confirm the upload (triggers processing)
	return s.ConfirmUpload(ctx, userID, models.ConfirmUploadRequest{UploadID: req.UploadID})
}

func (s *uploadService) GetUploadStatus(ctx context.Context, userID, uploadID string) (*models.UploadResponse, error) {
	upload, err := s.repo.GetUpload(ctx, userID, uploadID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.ErrUploadNotFound
		}
		return nil, err
	}

	response := upload.ToResponse()
	return &response, nil
}

func (s *uploadService) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.UploadResponse], error) {
	result, err := s.repo.ListUploads(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.UploadResponse, 0, len(result.Items))
	for _, upload := range result.Items {
		responses = append(responses, upload.ToResponse())
	}

	return &repository.PaginatedResult[models.UploadResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (s *uploadService) ReprocessUpload(ctx context.Context, userID, uploadID string, req models.ReprocessUploadRequest) (*models.UploadResponse, error) {
	upload, err := s.repo.GetUpload(ctx, userID, uploadID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.ErrUploadNotFound
		}
		return nil, err
	}

	if upload.Status != models.UploadStatusFailed {
		return nil, models.NewValidationError(map[string]string{
			"status": "Only failed uploads can be reprocessed",
		})
	}

	// Reset status to processing
	if err := s.repo.UpdateUploadStatus(ctx, userID, uploadID, models.UploadStatusProcessing, "", ""); err != nil {
		return nil, err
	}

	// TODO: Trigger Step Functions workflow from the specified step
	// For now, return processing status

	response := upload.ToResponse()
	response.Status = models.UploadStatusProcessing
	return &response, nil
}

func (s *uploadService) UploadCoverArt(ctx context.Context, userID, trackID string, req models.CoverArtUploadRequest) (*models.CoverArtUploadResponse, error) {
	// Verify track exists
	track, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Track", trackID)
		}
		return nil, err
	}

	// Generate S3 key for cover art
	s3Key := fmt.Sprintf("media/%s/%s/cover%s", userID, trackID, getFileExtension(req.FileName))

	// Generate presigned URL
	uploadURL, err := s.s3Repo.GeneratePresignedUploadURL(ctx, s3Key, req.ContentType, uploadURLExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Update track with cover art key (will be applied after upload)
	track.CoverArtKey = s3Key
	if err := s.repo.UpdateTrack(ctx, *track); err != nil {
		return nil, err
	}

	return &models.CoverArtUploadResponse{
		TrackID:     trackID,
		UploadURL:   uploadURL,
		ExpiresAt:   time.Now().Add(uploadURLExpiry),
		MaxFileSize: req.FileSize,
	}, nil
}

// getFileExtension extracts the file extension from a filename
func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}
