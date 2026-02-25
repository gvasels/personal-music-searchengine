package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/google/uuid"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	UploadID   string                 `json:"uploadId"`
	UserID     string                 `json:"userId"`
	S3Key      string                 `json:"s3Key"`
	FileName   string                 `json:"fileName"`
	Metadata   *models.UploadMetadata `json:"metadata"`
	CoverArt   *CoverArtResult        `json:"coverArt"`
	Analysis   *AnalysisResult        `json:"analysis"`
	BucketName string                 `json:"bucketName"`
	TableName  string                 `json:"tableName"`
}

// CoverArtResult represents the cover art extraction result
type CoverArtResult struct {
	CoverArtKey string `json:"coverArtKey"`
}

// AnalysisResult represents the audio analysis result
type AnalysisResult struct {
	BPM        int    `json:"bpm,omitempty"`
	MusicalKey string `json:"musicalKey,omitempty"`
	KeyMode    string `json:"keyMode,omitempty"`
	KeyCamelot string `json:"keyCamelot,omitempty"`
	Analyzed   bool   `json:"analyzed"`
	Error      string `json:"error,omitempty"`
}

// Response represents the output to Step Functions
type Response struct {
	TrackID     string `json:"trackId"`
	AlbumID     string `json:"albumId,omitempty"`
	IsDuplicate bool   `json:"isDuplicate"`
}

var repo repository.Repository
var s3Client *s3.Client
var sfnClient *sfn.Client
var audioPipelineARN = os.Getenv("AUDIO_PIPELINE_ARN")
var mediaBucket = os.Getenv("MEDIA_BUCKET")

// doTriggerAudioPipeline is the function called to trigger the audio analysis
// pipeline. It defaults to triggerAudioPipeline but can be replaced in tests.
var doTriggerAudioPipeline = triggerAudioPipeline

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "MusicLibrary"
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo = repository.NewDynamoDBRepository(dynamoClient, tableName)
	s3Client = s3.NewFromConfig(cfg)
	sfnClient = sfn.NewFromConfig(cfg)
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Add timeout to context (5 seconds less than Lambda timeout)
	ctx, cancel := context.WithTimeout(ctx, validation.ProcessorTimeoutSeconds*time.Second)
	defer cancel()

	// Validate input UUIDs to prevent injection attacks
	if err := validation.ValidateUUID(event.UserID, "userId"); err != nil {
		return nil, err
	}
	if err := validation.ValidateUUID(event.UploadID, "uploadId"); err != nil {
		return nil, err
	}

	title := getOrDefault(event.Metadata, "title", event.FileName)
	artist := getOrDefault(event.Metadata, "artist", "Unknown Artist")
	duration := getIntOrDefault(event.Metadata, "duration", 0)

	// Check for duplicate track (same title + artist for this user)
	if existing := findDuplicate(ctx, event.UserID, title, artist, duration); existing != nil {
		fmt.Printf("Duplicate detected: existing track %s matches upload %s (%s - %s)\n",
			existing.ID, event.UploadID, title, artist)

		// Clean up the uploaded file
		cleanupUploadedFile(ctx, event.S3Key)

		// Trigger delta analysis if existing track has HLS ready but is missing AI analysis
		if needsAnalysis(existing) && existing.HLSStatus == models.HLSStatusReady {
			fmt.Printf("Existing track %s needs analysis, triggering audio pipeline\n", existing.ID)
			doTriggerAudioPipeline(ctx, existing.ID, event.UserID, existing.S3Key, existing.Title, existing.Artist)
		}

		// Update step progress
		if err := repo.UpdateUploadStep(ctx, event.UserID, event.UploadID, models.StepCreateTrack, true); err != nil {
			fmt.Printf("Warning: failed to update step progress: %v\n", err)
		}

		// Mark upload as duplicate so frontend can show appropriate message
		if upload, err := repo.GetUpload(ctx, event.UserID, event.UploadID); err == nil {
			upload.IsDuplicate = true
			if err := repo.UpdateUpload(ctx, *upload); err != nil {
				fmt.Printf("Warning: failed to mark upload as duplicate: %v\n", err)
			}
		}

		return &Response{TrackID: existing.ID, IsDuplicate: true}, nil
	}

	trackID := uuid.New().String()
	now := time.Now()

	// Determine format from metadata
	format := models.AudioFormatMP3
	if event.Metadata != nil && event.Metadata.Format != "" {
		format = models.AudioFormat(event.Metadata.Format)
	}

	// Create track record
	track := models.Track{
		ID:         trackID,
		UserID:     event.UserID,
		Title:      title,
		Artist:     artist,
		Album:      getOrDefault(event.Metadata, "album", ""),
		Genre:      getOrDefault(event.Metadata, "genre", ""),
		Year:       getIntOrDefault(event.Metadata, "year", 0),
		Duration:   duration,
		Format:     format,
		S3Key:      event.S3Key, // Will be updated after file is moved
		Visibility: models.VisibilityPrivate,
		PlayCount:  0,
	}
	track.CreatedAt = now
	track.UpdatedAt = now

	// Set cover art key if available
	if event.CoverArt != nil && event.CoverArt.CoverArtKey != "" {
		track.CoverArtKey = event.CoverArt.CoverArtKey
	}

	// Set audio analysis results if available
	if event.Analysis != nil && event.Analysis.Analyzed {
		track.BPM = event.Analysis.BPM
		track.MusicalKey = event.Analysis.MusicalKey
		track.KeyMode = event.Analysis.KeyMode
		track.KeyCamelot = event.Analysis.KeyCamelot
	}

	// Set additional metadata fields if available
	if event.Metadata != nil {
		track.Bitrate = event.Metadata.Bitrate
	}

	// Parse multi-artist metadata and create Artist entities
	contributions := models.ParseArtists(artist)
	if len(contributions) > 0 {
		for i, c := range contributions {
			artistEntity, err := getOrCreateArtist(ctx, event.UserID, c.ArtistName)
			if err != nil {
				fmt.Printf("Warning: failed to get/create artist %q: %v\n", c.ArtistName, err)
				continue
			}
			contributions[i].ArtistID = artistEntity.ID
			// Set primary ArtistID on track to first main artist
			if i == 0 {
				track.ArtistID = artistEntity.ID
			}
		}
		track.Artists = contributions
	}

	// Create the track
	if err := repo.CreateTrack(ctx, track); err != nil {
		return nil, fmt.Errorf("failed to create track: %w", err)
	}

	// Update step progress
	if err := repo.UpdateUploadStep(ctx, event.UserID, event.UploadID, models.StepCreateTrack, true); err != nil {
		fmt.Printf("Warning: failed to update step progress: %v\n", err)
	}

	response := &Response{TrackID: trackID}

	// Create or update album if album name is present
	if track.Album != "" {
		album, err := repo.GetOrCreateAlbum(ctx, event.UserID, track.Album, track.Artist)
		if err != nil {
			// Log error but don't fail - track is already created
			fmt.Printf("Warning: failed to create/update album: %v\n", err)
		} else {
			response.AlbumID = album.ID
		}
	}

	return response, nil
}

// findDuplicate checks if a track with the same title+artist already exists for this user.
// Returns the existing track if found, nil otherwise.
func findDuplicate(ctx context.Context, userID, title, artist string, duration int) *models.Track {
	tracks, err := repo.ListTracksByArtist(ctx, userID, artist)
	if err != nil {
		fmt.Printf("Warning: failed to check for duplicates: %v\n", err)
		return nil
	}

	for i := range tracks {
		if strings.EqualFold(tracks[i].Title, title) {
			// Skip tracks whose file was never moved to permanent storage —
			// they came from a failed/incomplete previous upload and aren't playable.
			if !strings.HasPrefix(tracks[i].S3Key, "media/") {
				continue
			}
			// If both have duration, check within 5-second tolerance
			if duration > 0 && tracks[i].Duration > 0 {
				if math.Abs(float64(tracks[i].Duration-duration)) > 5 {
					continue
				}
			}
			return &tracks[i]
		}
	}
	return nil
}

// needsAnalysis returns true if the track is missing AI analysis data.
func needsAnalysis(track *models.Track) bool {
	if track.AnalysisStatus == "COMPLETED" {
		return false
	}
	// Missing BPM, or missing embedding, or no analysis status at all
	return track.BPM == 0 || track.EmbeddingID == "" || track.AnalysisStatus == "" || track.AnalysisStatus == "FAILED"
}

// triggerAudioPipeline starts the audio analysis Step Functions pipeline for a track.
func triggerAudioPipeline(ctx context.Context, trackID, userID, s3Key, title, artist string) {
	if audioPipelineARN == "" || sfnClient == nil {
		return
	}
	audioInput, _ := json.Marshal(map[string]string{
		"trackId": trackID,
		"userId":  userID,
		"s3Key":   s3Key,
		"title":   title,
		"artist":  artist,
	})
	_, err := sfnClient.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(audioPipelineARN),
		Name:            aws.String(fmt.Sprintf("audio-%s-%d", trackID, time.Now().Unix())),
		Input:           aws.String(string(audioInput)),
	})
	if err != nil {
		fmt.Printf("Warning: failed to start audio pipeline for track %s: %v\n", trackID, err)
	}
}

// cleanupUploadedFile deletes the uploaded file from S3 since it's a duplicate.
func cleanupUploadedFile(ctx context.Context, s3Key string) {
	if s3Client == nil || mediaBucket == "" {
		return
	}
	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(mediaBucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		fmt.Printf("Warning: failed to cleanup duplicate upload %s: %v\n", s3Key, err)
	}
}

func getOrDefault(meta *models.UploadMetadata, field, defaultVal string) string {
	if meta == nil {
		return defaultVal
	}
	switch field {
	case "title":
		if meta.Title != "" {
			return meta.Title
		}
	case "artist":
		if meta.Artist != "" {
			return meta.Artist
		}
	case "album":
		if meta.Album != "" {
			return meta.Album
		}
	case "genre":
		if meta.Genre != "" {
			return meta.Genre
		}
	}
	return defaultVal
}

func getIntOrDefault(meta *models.UploadMetadata, field string, defaultVal int) int {
	if meta == nil {
		return defaultVal
	}
	switch field {
	case "year":
		if meta.Year != 0 {
			return meta.Year
		}
	case "duration":
		if meta.Duration != 0 {
			return meta.Duration
		}
	}
	return defaultVal
}

func isAudioFile(key string) bool {
	lower := strings.ToLower(key)
	return strings.HasSuffix(lower, ".mp3") || strings.HasSuffix(lower, ".flac") ||
		strings.HasSuffix(lower, ".wav") || strings.HasSuffix(lower, ".m4a") ||
		strings.HasSuffix(lower, ".aac") || strings.HasSuffix(lower, ".ogg")
}

// getOrCreateArtist looks up an Artist entity by name for the user. If none
// exists, it creates one and returns it.
func getOrCreateArtist(ctx context.Context, userID, name string) (*models.Artist, error) {
	existing, err := repo.GetArtistByName(ctx, userID, name)
	if err != nil {
		return nil, fmt.Errorf("lookup failed: %w", err)
	}
	if len(existing) > 0 {
		return existing[0], nil
	}

	artist := models.Artist{
		ID:       uuid.New().String(),
		UserID:   userID,
		Name:     name,
		SortName: models.GenerateSortName(name),
		IsActive: true,
	}
	if err := repo.CreateArtist(ctx, artist); err != nil {
		// Handle race condition — another Lambda may have just created it
		// (condition expression failure means the PK already exists)
		existing, lookupErr := repo.GetArtistByName(ctx, userID, name)
		if lookupErr == nil && len(existing) > 0 {
			return existing[0], nil
		}
		return nil, fmt.Errorf("create failed: %w", err)
	}
	return &artist, nil
}

func main() {
	lambda.Start(handleRequest)
}
