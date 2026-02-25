package events

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// venueInfo holds a predefined venue with its location details.
type venueInfo struct {
	Venue   string
	City    string
	Region  string
	Country string
}

// venues is a pool of real venue/city/region/country combinations.
var venues = []venueInfo{
	{"Madison Square Garden", "New York", "NY", "US"},
	{"The O2", "London", "", "UK"},
	{"Ziggo Dome", "Amsterdam", "", "NL"},
	{"Olympiastadion", "Berlin", "", "DE"},
	{"Forum Melbourne", "Melbourne", "VIC", "AU"},
	{"Budokan", "Tokyo", "", "JP"},
}

// mockArtists is the built-in list of fake artist names for SearchArtists.
var mockArtists = []string{
	"Kylie Minogue",
	"MÖWE",
	"Xell Project",
	"Riser",
	"Ed Sheeran",
	"Daft Punk",
	"Deadmau5",
	"Armin van Buuren",
	"Above & Beyond",
	"Tiësto",
}

// MockProvider generates deterministic fake event data for development and testing.
type MockProvider struct{}

// NewMockProvider creates a new MockProvider.
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// GetArtistEvents generates 2-4 deterministic fake events for the given artist.
// The same artist name always produces the same events (same IDs, venues, offsets).
// Dates are relative to time.Now() so they shift forward as real time passes.
func (p *MockProvider) GetArtistEvents(_ context.Context, artistName string) ([]models.Event, error) {
	h := hashString(artistName)

	// Determine count: 2-4 events based on hash
	count := int(h%3) + 2 // 0,1,2 -> 2,3,4

	now := time.Now().UTC().Truncate(24 * time.Hour)

	// Normalize artist name for ID generation (lowercase, spaces to hyphens)
	normalizedName := normalizeForID(artistName)

	events := make([]models.Event, count)
	for i := 0; i < count; i++ {
		// Use a sub-hash for each event index to distribute venues and dates
		subHash := hashString(fmt.Sprintf("%s-%d", artistName, i))

		// Pick venue from pool
		venueIdx := int(subHash % uint32(len(venues)))
		v := venues[venueIdx]

		// Compute month offset: 1-6 months from now
		monthOffset := int(subHash%6) + 1
		eventDate := now.AddDate(0, monthOffset, 0).Add(20 * time.Hour) // 8 PM

		events[i] = models.Event{
			ID:         fmt.Sprintf("mock-%s-%d", normalizedName, i),
			ArtistName: artistName,
			Title:      fmt.Sprintf("%s Live", artistName),
			Date:       eventDate,
			Venue:      v.Venue,
			City:       v.City,
			Region:     v.Region,
			Country:    v.Country,
			TicketURL:  "",
			Status:     "scheduled",
			Source:     "mock",
		}
	}

	return events, nil
}

// SearchArtists searches the built-in list of mock artists using case-insensitive
// substring matching. Results respect the limit parameter.
func (p *MockProvider) SearchArtists(_ context.Context, query string, limit int) ([]models.ArtistSearchResult, error) {
	queryLower := strings.ToLower(query)

	var results []models.ArtistSearchResult
	for _, name := range mockArtists {
		if query != "" && !strings.Contains(strings.ToLower(name), queryLower) {
			continue
		}

		h := hashString(name)
		upcomingEvents := int(h%5) + 1 // 1-5

		results = append(results, models.ArtistSearchResult{
			Name:           name,
			UpcomingEvents: upcomingEvents,
			Source:         "mock",
		})

		if limit > 0 && len(results) >= limit {
			break
		}
	}

	// Ensure we never return nil -- always return an empty slice
	if results == nil {
		results = []models.ArtistSearchResult{}
	}

	return results, nil
}

// hashString returns a deterministic 32-bit hash of the input string using FNV-1a.
func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// normalizeForID converts an artist name to a lowercase, hyphen-separated form
// suitable for use in event IDs.
func normalizeForID(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	return strings.ReplaceAll(lower, " ", "-")
}
