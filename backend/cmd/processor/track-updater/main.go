// Track Updater Lambda - Updates DynamoDB with audio analysis results
// Called as final step in audio understanding pipeline
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// Event combines results from audio analyzer and embedding generator
type Event struct {
	TrackID string `json:"trackId"`
	UserID  string `json:"userId"`
	S3Key   string `json:"s3Key"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Album   string `json:"album"`
	// Results from previous steps
	FeaturesResult  *FeaturesResult  `json:"featuresResult,omitempty"`
	AnalyzerResult  *AnalyzerResult  `json:"analyzerResult,omitempty"`
	EmbeddingResult *EmbeddingResult `json:"embeddingResult,omitempty"`
	// Error from catch blocks
	FeaturesError  map[string]interface{} `json:"featuresError,omitempty"`
	AnalyzerError  map[string]interface{} `json:"analyzerError,omitempty"`
	EmbeddingError map[string]interface{} `json:"embeddingError,omitempty"`
}

type FeaturesResult struct {
	TrackID     string `json:"trackId"`
	UserID      string `json:"userId"`
	BPM         int    `json:"bpm,omitempty"`
	Key         string `json:"key,omitempty"`
	CamelotCode string `json:"camelotCode,omitempty"`
	Error       string `json:"error,omitempty"`
}

type AnalyzerResult struct {
	TrackID  string   `json:"trackId"`
	UserID   string   `json:"userId"`
	Analysis Analysis `json:"analysis"`
	Error    string   `json:"error,omitempty"`
}

type EmbeddingResult struct {
	TrackID     string `json:"trackId"`
	UserID      string `json:"userId"`
	EmbeddingID string `json:"embeddingId,omitempty"`
	Error       string `json:"error,omitempty"`
}

type Analysis struct {
	Genre           string    `json:"genre"`
	SubGenre        string    `json:"subGenre"`
	Mood            string    `json:"mood"`
	ToneDescription string    `json:"toneDescription"`
	Sections        []Section `json:"sections"`
	Instrumentation string    `json:"instrumentation"`
	VocalPresence   string    `json:"vocalPresence"`
	EnergyProfile   string    `json:"energyProfile"`
}

type Section struct {
	Name        string `json:"name"`
	StartSec    int    `json:"startSec"`
	EndSec      int    `json:"endSec"`
	Description string `json:"description"`
}

type Result struct {
	TrackID string `json:"trackId"`
	UserID  string `json:"userId"`
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
}

var repo repository.Repository

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	repo = repository.NewDynamoDBRepository(dynamoClient, tableName)
}

func handler(ctx context.Context, event Event) (Result, error) {
	result := Result{
		TrackID: event.TrackID,
		UserID:  event.UserID,
	}

	// Extract analysis from nested result
	var analysis models.AudioAnalysis
	if event.AnalyzerResult != nil && event.AnalyzerResult.Analysis.Genre != "" {
		a := event.AnalyzerResult.Analysis
		sections := make([]models.Section, len(a.Sections))
		for i, s := range a.Sections {
			sections[i] = models.Section{
				Name:        s.Name,
				StartSec:    s.StartSec,
				EndSec:      s.EndSec,
				Description: s.Description,
			}
		}
		analysis = models.AudioAnalysis{
			Genre:           a.Genre,
			SubGenre:        a.SubGenre,
			Mood:            a.Mood,
			ToneDescription: a.ToneDescription,
			Sections:        sections,
			Instrumentation: a.Instrumentation,
			VocalPresence:   a.VocalPresence,
			EnergyProfile:   a.EnergyProfile,
		}
	}

	// Extract embedding ID and status from nested result
	if event.EmbeddingResult != nil && event.EmbeddingResult.EmbeddingID != "" {
		analysis.EmbeddingID = event.EmbeddingResult.EmbeddingID
		analysis.EmbeddingStatus = "COMPLETED"
	} else if event.EmbeddingError != nil || (event.EmbeddingResult != nil && event.EmbeddingResult.Error != "") {
		analysis.EmbeddingStatus = "FAILED"
	}

	// Extract BPM/Key from features result
	var bpm int
	var key, camelotCode string
	if event.FeaturesResult != nil {
		bpm = event.FeaturesResult.BPM
		key = event.FeaturesResult.Key
		camelotCode = event.FeaturesResult.CamelotCode
	}

	// Update track in DynamoDB
	if err := repo.UpdateTrackAnalysisWithFeatures(ctx, event.UserID, event.TrackID, analysis, bpm, key, camelotCode); err != nil {
		result.Status = "FAILED"
		result.Error = fmt.Sprintf("failed to update track: %v", err)
		return result, nil
	}

	// Auto-tag from analysis results (non-blocking — track update already succeeded)
	if analysis.Genre != "" || analysis.SubGenre != "" || analysis.Mood != "" {
		autoTagFromAnalysis(ctx, event.UserID, event.TrackID, analysis)
	}

	// Determine final status based on errors
	hasAnalysisError := event.AnalyzerError != nil || (event.AnalyzerResult != nil && event.AnalyzerResult.Error != "")
	hasEmbeddingError := event.EmbeddingError != nil || (event.EmbeddingResult != nil && event.EmbeddingResult.Error != "")

	if hasAnalysisError && hasEmbeddingError {
		result.Status = "FAILED"
		result.Error = "both analysis and embedding failed"
	} else if hasAnalysisError || hasEmbeddingError {
		result.Status = "PARTIAL"
		if hasEmbeddingError && event.EmbeddingResult != nil {
			result.Error = event.EmbeddingResult.Error
		}
	} else {
		result.Status = "COMPLETED"
	}

	return result, nil
}

// autoTagFromAnalysis creates tags from genre, sub-genre, and mood analysis
// fields and associates them with the track. All errors are logged as warnings
// since the core track update has already succeeded.
func autoTagFromAnalysis(ctx context.Context, userID, trackID string, analysis models.AudioAnalysis) {
	// Collect non-empty tag names, normalized to lowercase
	var tagNames []string
	for _, raw := range []string{analysis.Genre, analysis.SubGenre, analysis.Mood} {
		name := strings.TrimSpace(strings.ToLower(raw))
		if name != "" {
			tagNames = append(tagNames, name)
		}
	}
	if len(tagNames) == 0 {
		return
	}

	// Ensure each tag entity exists
	for _, name := range tagNames {
		if _, err := repo.GetTag(ctx, userID, name); err != nil {
			// Tag doesn't exist — create it
			tag := models.Tag{
				UserID: userID,
				Name:   name,
			}
			if createErr := repo.CreateTag(ctx, tag); createErr != nil {
				fmt.Printf("Warning: failed to create auto-tag %q: %v\n", name, createErr)
			}
		}
	}

	// Associate tags with the track
	if err := repo.AddTagsToTrack(ctx, userID, trackID, tagNames); err != nil {
		fmt.Printf("Warning: failed to add auto-tags to track %s: %v\n", trackID, err)
		return
	}

	// Update the track's Tags field to include new auto-tags
	track, err := repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		fmt.Printf("Warning: failed to get track %s for tag update: %v\n", trackID, err)
		return
	}

	// Merge new tags with existing, avoiding duplicates
	existingSet := make(map[string]bool, len(track.Tags))
	for _, t := range track.Tags {
		existingSet[t] = true
	}
	for _, name := range tagNames {
		if !existingSet[name] {
			track.Tags = append(track.Tags, name)
		}
	}
	if err := repo.UpdateTrack(ctx, *track); err != nil {
		fmt.Printf("Warning: failed to update track tags for %s: %v\n", trackID, err)
	}

	fmt.Printf("Auto-tagged track %s with: %v\n", trackID, tagNames)
}

func main() {
	lambda.Start(handler)
}
