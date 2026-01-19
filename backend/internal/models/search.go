package models

// SearchRequest represents a search query
type SearchRequest struct {
	Query   string        `json:"query" validate:"required,min=1,max=500"`
	Filters SearchFilters `json:"filters,omitempty"`
	Sort    SearchSort    `json:"sort,omitempty"`
	Limit   int           `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Cursor  string        `json:"cursor,omitempty"` // Opaque base64-encoded pagination cursor
}

// SearchFilters represents filters for search
type SearchFilters struct {
	Artists []string `json:"artists,omitempty"`
	Albums  []string `json:"albums,omitempty"`
	Genres  []string `json:"genres,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Years   []int    `json:"years,omitempty"`
	Formats []string `json:"formats,omitempty"`
}

// SearchSort represents sort options for search
type SearchSort struct {
	Field string `json:"field,omitempty"` // relevance, title, artist, album, year, playCount, createdAt
	Order string `json:"order,omitempty"` // asc, desc
}

// SearchResponse represents search results
type SearchResponse struct {
	Query        string          `json:"query"`
	TotalResults int             `json:"totalResults"`
	Tracks       []TrackResponse `json:"tracks"`
	Albums       []AlbumResponse `json:"albums,omitempty"`
	Artists      []ArtistSummary `json:"artists,omitempty"`
	Facets       SearchFacets    `json:"facets,omitempty"`
	Limit        int             `json:"limit"`
	NextCursor   string          `json:"nextCursor,omitempty"` // Next page cursor (empty if no more results)
	HasMore      bool            `json:"hasMore"`
}

// SearchFacets represents aggregated facets for filtering
type SearchFacets struct {
	Artists []FacetItem `json:"artists,omitempty"`
	Albums  []FacetItem `json:"albums,omitempty"`
	Genres  []FacetItem `json:"genres,omitempty"`
	Years   []FacetItem `json:"years,omitempty"`
	Tags    []FacetItem `json:"tags,omitempty"`
	Formats []FacetItem `json:"formats,omitempty"`
}

// FacetItem represents a single facet value with count
type FacetItem struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// SearchSuggestion represents an autocomplete suggestion
type SearchSuggestion struct {
	Text     string `json:"text"`
	Type     string `json:"type"` // track, artist, album, tag
	ID       string `json:"id,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
}

// AutocompleteResponse represents autocomplete suggestions
type AutocompleteResponse struct {
	Query       string             `json:"query"`
	Suggestions []SearchSuggestion `json:"suggestions"`
}

// NixieIndexDocument represents a document in the Nixiesearch index
type NixieIndexDocument struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	Type        string   `json:"type"` // track, album, artist
	Title       string   `json:"title"`
	Artist      string   `json:"artist"`
	AlbumArtist string   `json:"album_artist,omitempty"`
	Album       string   `json:"album,omitempty"`
	Genre       string   `json:"genre,omitempty"`
	Year        int      `json:"year,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Duration    int      `json:"duration,omitempty"`
	Format      string   `json:"format,omitempty"`
	PlayCount   int      `json:"play_count,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// NixieSearchQuery represents a query to Nixiesearch
type NixieSearchQuery struct {
	Query       string           `json:"query"`
	Fields      []string         `json:"fields,omitempty"`
	Filter      map[string]any   `json:"filter,omitempty"`
	Sort        []NixieSortField `json:"sort,omitempty"`
	Size        int              `json:"size,omitempty"`
	SearchAfter []any            `json:"search_after,omitempty"` // Cursor-based pagination for Nixiesearch
}

// NixieSortField represents a sort field in Nixiesearch
type NixieSortField struct {
	Field string `json:"field"`
	Order string `json:"order"` // asc, desc
}

// NixieSearchResult represents a result from Nixiesearch
type NixieSearchResult struct {
	Hits       []NixieHit `json:"hits"`
	TotalHits  int        `json:"total_hits"`
	MaxScore   float64    `json:"max_score,omitempty"`
}

// NixieHit represents a single hit from Nixiesearch
type NixieHit struct {
	ID     string             `json:"_id"`
	Score  float64            `json:"_score"`
	Source NixieIndexDocument `json:"_source"`
}
