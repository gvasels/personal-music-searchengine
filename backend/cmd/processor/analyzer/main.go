package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gvasels/personal-music-searchengine/internal/analysis"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	UploadID   string `json:"uploadId"`
	UserID     string `json:"userId"`
	S3Key      string `json:"s3Key"`
	FileName   string `json:"fileName"`
	BucketName string `json:"bucketName"`
}

// Response represents the output to Step Functions
type Response struct {
	BPM        int    `json:"bpm,omitempty"`
	MusicalKey string `json:"musicalKey,omitempty"`
	KeyMode    string `json:"keyMode,omitempty"`
	KeyCamelot string `json:"keyCamelot,omitempty"`
	Analyzed   bool   `json:"analyzed"`
	Error      string `json:"error,omitempty"`
}

var s3Client *s3.Client
var analyzer *analysis.Analyzer

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}
	s3Client = s3.NewFromConfig(cfg)
	analyzer = analysis.NewAnalyzer()
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Add timeout to context (allow up to 25 seconds for analysis)
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	// Validate file size before download
	if err := validation.ValidateFileSize(ctx, s3Client, event.BucketName, event.S3Key); err != nil {
		// Return success with error message - don't fail the workflow
		return &Response{
			Analyzed: false,
			Error:    fmt.Sprintf("file validation failed: %v", err),
		}, nil
	}

	// Download file from S3
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &event.BucketName,
		Key:    &event.S3Key,
	})
	if err != nil {
		return &Response{
			Analyzed: false,
			Error:    fmt.Sprintf("failed to download from S3: %v", err),
		}, nil
	}
	defer result.Body.Close()

	// Analyze audio
	analysisResult, err := analyzer.Analyze(ctx, result.Body, event.FileName)
	if err != nil {
		// Return success with error - analysis failure shouldn't block upload
		return &Response{
			Analyzed: false,
			Error:    fmt.Sprintf("analysis failed: %v", err),
		}, nil
	}

	return &Response{
		BPM:        analysisResult.BPM,
		MusicalKey: analysisResult.MusicalKey,
		KeyMode:    analysisResult.KeyMode,
		KeyCamelot: analysisResult.KeyCamelot,
		Analyzed:   true,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
