// Package main implements the Nixiesearch Lambda function.
// This is a simple search engine with S3-based index storage.
// For MVP, uses in-memory index with S3 persistence.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client    *s3.Client
	indexBucket string
	indexPath   string
	index       *SearchIndex
	indexMutex  sync.RWMutex
	initialized bool
)

// SearchIndex holds the in-memory search index
type SearchIndex struct {
	Documents map[string]Document `json:"documents"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

// Document represents a searchable track
type Document struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Title     string    `json:"title"`
	Artist    string    `json:"artist"`
	Album     string    `json:"album"`
	Genre     string    `json:"genre"`
	Year      int       `json:"year"`
	Duration  int       `json:"duration"`
	Filename  string    `json:"filename"`
	IndexedAt time.Time `json:"indexedAt"`
}

// Request represents the incoming Lambda request
type Request struct {
	Operation string      `json:"operation"` // search, index, delete, bulk_index
	Payload   interface{} `json:"payload"`
}

// Response represents the Lambda response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SearchQuery represents a search request
type SearchQuery struct {
	Query   string        `json:"query"`
	Filters SearchFilters `json:"filters"`
	Sort    *SortOption   `json:"sort,omitempty"`
	Limit   int           `json:"limit"`
	Cursor  string        `json:"cursor"`
}

// SearchFilters for narrowing results
type SearchFilters struct {
	UserID   string `json:"userId"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Genre    string `json:"genre"`
	YearFrom int    `json:"yearFrom"`
	YearTo   int    `json:"yearTo"`
}

// SortOption for result ordering
type SortOption struct {
	Field string `json:"field"`
	Order string `json:"order"` // asc or desc
}

// SearchResponse contains search results
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Cursor  string         `json:"cursor,omitempty"`
}

// SearchResult represents a single search hit
type SearchResult struct {
	Document Document `json:"document"`
	Score    float64  `json:"score"`
}

// IndexRequest for adding a document
type IndexRequest struct {
	Document Document `json:"document"`
}

// IndexResponse after indexing
type IndexResponse struct {
	ID      string `json:"id"`
	Indexed bool   `json:"indexed"`
}

// DeleteRequest for removing a document
type DeleteRequest struct {
	ID string `json:"id"`
}

// DeleteResponse after deletion
type DeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// BulkIndexRequest for batch indexing
type BulkIndexRequest struct {
	Documents []Document `json:"documents"`
}

// BulkIndexResponse after bulk indexing
type BulkIndexResponse struct {
	Indexed int      `json:"indexed"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

func init() {
	indexBucket = os.Getenv("SEARCH_INDEX_BUCKET")
	indexPath = os.Getenv("INDEX_PATH")
	if indexPath == "" {
		indexPath = "/tmp/nixiesearch"
	}
}

func initializeAWS(ctx context.Context) error {
	if s3Client != nil {
		return nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	return nil
}

func loadIndex(ctx context.Context) error {
	if initialized {
		return nil
	}

	indexMutex.Lock()
	defer indexMutex.Unlock()

	if initialized {
		return nil
	}

	// Try to load from S3
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &indexBucket,
		Key:    stringPtr("index.json"),
	})
	if err != nil {
		// Index doesn't exist yet, create empty
		index = &SearchIndex{
			Documents: make(map[string]Document),
			UpdatedAt: time.Now(),
		}
		initialized = true
		return nil
	}
	defer result.Body.Close()

	var loadedIndex SearchIndex
	if err := json.NewDecoder(result.Body).Decode(&loadedIndex); err != nil {
		return fmt.Errorf("failed to decode index: %w", err)
	}

	index = &loadedIndex
	initialized = true
	return nil
}

func saveIndex(ctx context.Context) error {
	indexMutex.RLock()
	data, err := json.Marshal(index)
	indexMutex.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &indexBucket,
		Key:         stringPtr("index.json"),
		Body:        strings.NewReader(string(data)),
		ContentType: stringPtr("application/json"),
	})
	if err != nil {
		return fmt.Errorf("failed to save index to S3: %w", err)
	}

	return nil
}

func handleRequest(ctx context.Context, req Request) (Response, error) {
	if err := initializeAWS(ctx); err != nil {
		return Response{Success: false, Error: err.Error()}, nil
	}

	if err := loadIndex(ctx); err != nil {
		return Response{Success: false, Error: err.Error()}, nil
	}

	switch req.Operation {
	case "search":
		return handleSearch(ctx, req.Payload)
	case "index":
		return handleIndex(ctx, req.Payload)
	case "delete":
		return handleDelete(ctx, req.Payload)
	case "bulk_index":
		return handleBulkIndex(ctx, req.Payload)
	default:
		return Response{Success: false, Error: fmt.Sprintf("unknown operation: %s", req.Operation)}, nil
	}
}

func handleSearch(ctx context.Context, payload interface{}) (Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Response{Success: false, Error: "invalid search payload"}, nil
	}

	var query SearchQuery
	if err := json.Unmarshal(data, &query); err != nil {
		return Response{Success: false, Error: "invalid search query"}, nil
	}

	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	indexMutex.RLock()
	defer indexMutex.RUnlock()

	var results []SearchResult
	queryLower := strings.ToLower(query.Query)

	for _, doc := range index.Documents {
		// Filter by user
		if query.Filters.UserID != "" && doc.UserID != query.Filters.UserID {
			continue
		}

		// Apply filters
		if query.Filters.Artist != "" && !strings.Contains(strings.ToLower(doc.Artist), strings.ToLower(query.Filters.Artist)) {
			continue
		}
		if query.Filters.Album != "" && !strings.Contains(strings.ToLower(doc.Album), strings.ToLower(query.Filters.Album)) {
			continue
		}
		if query.Filters.Genre != "" && doc.Genre != query.Filters.Genre {
			continue
		}
		if query.Filters.YearFrom > 0 && doc.Year < query.Filters.YearFrom {
			continue
		}
		if query.Filters.YearTo > 0 && doc.Year > query.Filters.YearTo {
			continue
		}

		// Calculate relevance score
		score := calculateScore(doc, queryLower)
		if queryLower == "" || score > 0 {
			results = append(results, SearchResult{
				Document: doc,
				Score:    score,
			})
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Apply limit
	total := len(results)
	if len(results) > query.Limit {
		results = results[:query.Limit]
	}

	return Response{
		Success: true,
		Data: SearchResponse{
			Results: results,
			Total:   total,
		},
	}, nil
}

func calculateScore(doc Document, query string) float64 {
	if query == "" {
		return 1.0
	}

	var score float64

	// Title match (highest weight)
	if strings.Contains(strings.ToLower(doc.Title), query) {
		score += 3.0
	}
	// Artist match
	if strings.Contains(strings.ToLower(doc.Artist), query) {
		score += 2.0
	}
	// Album match
	if strings.Contains(strings.ToLower(doc.Album), query) {
		score += 1.5
	}
	// Filename match
	if strings.Contains(strings.ToLower(doc.Filename), query) {
		score += 1.0
	}

	return score
}

func handleIndex(ctx context.Context, payload interface{}) (Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Response{Success: false, Error: "invalid index payload"}, nil
	}

	var req IndexRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return Response{Success: false, Error: "invalid index request"}, nil
	}

	req.Document.IndexedAt = time.Now()

	indexMutex.Lock()
	index.Documents[req.Document.ID] = req.Document
	index.UpdatedAt = time.Now()
	indexMutex.Unlock()

	if err := saveIndex(ctx); err != nil {
		return Response{Success: false, Error: err.Error()}, nil
	}

	return Response{
		Success: true,
		Data: IndexResponse{
			ID:      req.Document.ID,
			Indexed: true,
		},
	}, nil
}

func handleDelete(ctx context.Context, payload interface{}) (Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Response{Success: false, Error: "invalid delete payload"}, nil
	}

	var req DeleteRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return Response{Success: false, Error: "invalid delete request"}, nil
	}

	indexMutex.Lock()
	_, exists := index.Documents[req.ID]
	if exists {
		delete(index.Documents, req.ID)
		index.UpdatedAt = time.Now()
	}
	indexMutex.Unlock()

	if exists {
		if err := saveIndex(ctx); err != nil {
			return Response{Success: false, Error: err.Error()}, nil
		}
	}

	return Response{
		Success: true,
		Data: DeleteResponse{
			ID:      req.ID,
			Deleted: exists,
		},
	}, nil
}

func handleBulkIndex(ctx context.Context, payload interface{}) (Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Response{Success: false, Error: "invalid bulk index payload"}, nil
	}

	var req BulkIndexRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return Response{Success: false, Error: "invalid bulk index request"}, nil
	}

	indexMutex.Lock()
	indexed := 0
	for _, doc := range req.Documents {
		doc.IndexedAt = time.Now()
		index.Documents[doc.ID] = doc
		indexed++
	}
	index.UpdatedAt = time.Now()
	indexMutex.Unlock()

	if err := saveIndex(ctx); err != nil {
		return Response{Success: false, Error: err.Error()}, nil
	}

	return Response{
		Success: true,
		Data: BulkIndexResponse{
			Indexed: indexed,
			Failed:  0,
		},
	}, nil
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	lambda.Start(handleRequest)
}
