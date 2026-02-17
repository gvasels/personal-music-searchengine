// Track Updater Lambda - Updates DynamoDB with audio analysis results
// Called as final step in audio understanding pipeline
package main

import (
	"context"
	"fmt"
	"os"

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

	// Extract embedding ID from nested result
	if event.EmbeddingResult != nil && event.EmbeddingResult.EmbeddingID != "" {
		analysis.EmbeddingID = event.EmbeddingResult.EmbeddingID
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

func main() {
	lambda.Start(handler)
}
