// Audio Analyzer Lambda - Combines signal analysis and GenAI analysis
// Uses Bedrock for genre/mood classification and Marengo for embeddings
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Event struct {
	TrackID string `json:"trackId"`
	UserID  string `json:"userId"`
	S3Key   string `json:"s3Key"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Album   string `json:"album"`
}

type AnalysisResult struct {
	TrackID  string   `json:"trackId"`
	UserID   string   `json:"userId"`
	Analysis Analysis `json:"analysis"`
	Error    string   `json:"error,omitempty"`
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

var bedrockClient *bedrockruntime.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}
	bedrockClient = bedrockruntime.NewFromConfig(cfg)
}

func handler(ctx context.Context, event Event) (AnalysisResult, error) {
	result := AnalysisResult{
		TrackID: event.TrackID,
		UserID:  event.UserID,
	}

	// Build prompt for genre/mood analysis
	prompt := buildAnalysisPrompt(event)

	// Call Bedrock Claude for analysis
	analysis, err := analyzeWithBedrock(ctx, prompt)
	if err != nil {
		result.Error = fmt.Sprintf("bedrock analysis failed: %v", err)
		return result, nil // Return partial result, don't fail pipeline
	}

	result.Analysis = analysis
	return result, nil
}

func buildAnalysisPrompt(event Event) string {
	return fmt.Sprintf(`Analyze this music track and provide structured metadata.

Track info:
- Title: %s
- Artist: %s
- Album: %s

Based on the artist, title, and album information, infer the likely genre, mood, and characteristics.

Provide JSON output with these fields:
{
  "genre": "primary genre (e.g., Electronic, Rock, Hip Hop, Jazz, Classical)",
  "subGenre": "sub-genre (for Electronic: Deep House, Tech House, Techno, Trance, etc.)",
  "mood": "one-word mood (e.g., energetic, melancholic, uplifting, dark, chill)",
  "toneDescription": "2-3 sentence description of the likely tone and feel",
  "sections": [
    {"name": "intro", "startSec": 0, "endSec": 30, "description": "typical intro description"},
    {"name": "main", "startSec": 30, "endSec": 180, "description": "main section description"}
  ],
  "instrumentation": "likely instruments/sounds (e.g., synthesizers, guitar, drums)",
  "vocalPresence": "none|male|female|mixed",
  "energyProfile": "description of energy arc"
}

For electronic music, use this taxonomy:
- House: Deep House, Tech House, Progressive House, Electro House, Future House
- Techno: Minimal Techno, Industrial Techno, Melodic Techno, Acid Techno
- Trance: Progressive Trance, Uplifting Trance, Psytrance, Vocal Trance
- Drum and Bass: Liquid DnB, Neurofunk, Jump-Up, Jungle
- Dubstep: Brostep, Riddim, Deep Dubstep
- Ambient: Dark Ambient, Ambient Dub, Drone

Return ONLY valid JSON, no other text.`, event.Title, event.Artist, event.Album)
}

func analyzeWithBedrock(ctx context.Context, prompt string) (Analysis, error) {
	modelID := os.Getenv("BEDROCK_MODEL_ID")
	if modelID == "" {
		modelID = "anthropic.claude-3-haiku-20240307-v1:0"
	}

	// Build request body for Claude
	requestBody := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        1024,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return Analysis{}, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := bedrockClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     &modelID,
		Body:        bodyBytes,
		ContentType: stringPtr("application/json"),
	})
	if err != nil {
		return Analysis{}, fmt.Errorf("invoke model: %w", err)
	}

	// Parse response
	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return Analysis{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(response.Content) == 0 {
		return Analysis{}, fmt.Errorf("empty response from model")
	}

	// Parse the JSON from the response
	var analysis Analysis
	if err := json.Unmarshal([]byte(response.Content[0].Text), &analysis); err != nil {
		return Analysis{}, fmt.Errorf("parse analysis JSON: %w", err)
	}

	return analysis, nil
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	lambda.Start(handler)
}
