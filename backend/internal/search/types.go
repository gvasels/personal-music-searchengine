// Package search provides Nixiesearch client and search functionality.
package search

import "time"

// Document represents a searchable document in the index.
type Document struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Title     string    `json:"title"`
	Artist    string    `json:"artist"`
	Album     string    `json:"album"`
	Genre     string    `json:"genre"`
	Year      int       `json:"year,omitempty"`
	Duration  int       `json:"duration,omitempty"`
	Filename  string    `json:"filename"`
	IndexedAt time.Time `json:"indexedAt"`
}

// SearchQuery represents a search request.
type SearchQuery struct {
	Query   string        `json:"query"`
	Filters SearchFilters `json:"filters,omitempty"`
	Sort    *SortOption   `json:"sort,omitempty"`
	Limit   int           `json:"limit,omitempty"`
	Cursor  string        `json:"cursor,omitempty"`
}

// SearchFilters represents optional filters for search.
type SearchFilters struct {
	UserID   string   `json:"userId,omitempty"`   // Required - scopes search to user
	Artist   string   `json:"artist,omitempty"`
	Album    string   `json:"album,omitempty"`
	Genre    string   `json:"genre,omitempty"`
	YearFrom int      `json:"yearFrom,omitempty"`
	YearTo   int      `json:"yearTo,omitempty"`
	Tags     []string `json:"tags,omitempty"` // Filter by tags (AND logic)
}

// SortOption represents sorting configuration.
type SortOption struct {
	Field string `json:"field"` // title, artist, album, year, duration
	Order string `json:"order"` // asc, desc
}

// SearchResult represents a single search result.
type SearchResult struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Artist      string  `json:"artist"`
	Album       string  `json:"album"`
	Genre       string  `json:"genre"`
	Year        int     `json:"year,omitempty"`
	Duration    int     `json:"duration,omitempty"`
	CoverArtURL string  `json:"coverArtUrl,omitempty"`
	Score       float64 `json:"score"`
}

// SearchResponse represents the response from a search query.
type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	Total      int            `json:"total"`
	NextCursor string         `json:"cursor,omitempty"`
}

// IndexRequest represents a request to index a document.
type IndexRequest struct {
	Document Document `json:"document"`
}

// IndexResponse represents the response from an index operation.
type IndexResponse struct {
	ID      string `json:"id"`
	Indexed bool   `json:"indexed"`
}

// DeleteRequest represents a request to delete a document.
type DeleteRequest struct {
	ID string `json:"id"`
}

// DeleteResponse represents the response from a delete operation.
type DeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// BulkIndexRequest represents a request to index multiple documents.
type BulkIndexRequest struct {
	Documents []Document `json:"documents"`
}

// BulkIndexResponse represents the response from a bulk index operation.
type BulkIndexResponse struct {
	Indexed int      `json:"indexed"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

// NixiesearchRequest represents a request to the Nixiesearch Lambda.
type NixiesearchRequest struct {
	Operation string      `json:"operation"` // search, index, delete, bulk_index
	Payload   interface{} `json:"payload"`
}

// NixiesearchResponse represents a response from the Nixiesearch Lambda.
type NixiesearchResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
