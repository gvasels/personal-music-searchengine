package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewArtistWatch(t *testing.T) {
	userID := "user-123"
	artistName := "Kylie Minogue"

	watch := NewArtistWatch(userID, artistName)

	if watch.UserID != userID {
		t.Errorf("NewArtistWatch().UserID = %v, want %v", watch.UserID, userID)
	}

	if watch.ArtistName != artistName {
		t.Errorf("NewArtistWatch().ArtistName = %v, want %v", watch.ArtistName, artistName)
	}

	if watch.WatchedAt.IsZero() {
		t.Error("NewArtistWatch().WatchedAt should be set")
	}
}

func TestNewArtistWatchItem(t *testing.T) {
	watch := NewArtistWatch("user-123", "Kylie Minogue")
	item := NewArtistWatchItem(*watch)

	// Primary key: "what artists does this user watch?"
	expectedPK := "USER#user-123"
	expectedSK := "ARTIST_WATCH#kylie minogue"

	if item.PK != expectedPK {
		t.Errorf("ArtistWatchItem.PK = %v, want %v", item.PK, expectedPK)
	}

	if item.SK != expectedSK {
		t.Errorf("ArtistWatchItem.SK = %v, want %v", item.SK, expectedSK)
	}

	// GSI1: "which users watch this artist?"
	expectedGSI1PK := "ARTIST_WATCH#kylie minogue"
	expectedGSI1SK := "USER#user-123"

	if item.GSI1PK != expectedGSI1PK {
		t.Errorf("ArtistWatchItem.GSI1PK = %v, want %v", item.GSI1PK, expectedGSI1PK)
	}

	if item.GSI1SK != expectedGSI1SK {
		t.Errorf("ArtistWatchItem.GSI1SK = %v, want %v", item.GSI1SK, expectedGSI1SK)
	}

	if item.Type != string(EntityArtistWatch) {
		t.Errorf("ArtistWatchItem.Type = %v, want %v", item.Type, EntityArtistWatch)
	}
}

func TestNewArtistWatchItem_PreservesOriginalArtistName(t *testing.T) {
	watch := NewArtistWatch("user-123", "Kylie Minogue")
	item := NewArtistWatchItem(*watch)

	// The SK uses the normalized name, but the ArtistWatch data preserves original casing
	if item.ArtistName != "Kylie Minogue" {
		t.Errorf("ArtistWatchItem.ArtistName = %v, want %v (original casing)", item.ArtistName, "Kylie Minogue")
	}
}

func TestNormalizeArtistName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Mixed case to lowercase",
			input:    "Kylie Minogue",
			expected: "kylie minogue",
		},
		{
			name:     "All uppercase to lowercase",
			input:    "DEADMAU5",
			expected: "deadmau5",
		},
		{
			name:     "Trims whitespace and lowercases",
			input:    "  Tiesto  ",
			expected: "tiesto",
		},
		{
			name:     "Already lowercase",
			input:    "moby",
			expected: "moby",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Whitespace only",
			input:    "   ",
			expected: "",
		},
		{
			name:     "Unicode characters preserved",
			input:    " MOWE ",
			expected: "mowe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeArtistName(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeArtistName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestArtistWatch_Validate(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		artistName string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "Valid watch",
			userID:     "user-123",
			artistName: "Kylie Minogue",
			wantErr:    false,
		},
		{
			name:       "Empty user ID",
			userID:     "",
			artistName: "Kylie Minogue",
			wantErr:    true,
			errMsg:     "user ID cannot be empty",
		},
		{
			name:       "Empty artist name",
			userID:     "user-123",
			artistName: "",
			wantErr:    true,
			errMsg:     "artist name cannot be empty",
		},
		{
			name:       "Whitespace-only artist name",
			userID:     "user-123",
			artistName: "   ",
			wantErr:    true,
			errMsg:     "artist name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watch := NewArtistWatch(tt.userID, tt.artistName)
			err := watch.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() should return error for %s", tt.name)
				} else if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetArtistWatchPK(t *testing.T) {
	pk := GetArtistWatchPK("user-123")
	expected := "USER#user-123"

	if pk != expected {
		t.Errorf("GetArtistWatchPK() = %v, want %v", pk, expected)
	}
}

func TestGetArtistWatchSK(t *testing.T) {
	sk := GetArtistWatchSK("kylie minogue")
	expected := "ARTIST_WATCH#kylie minogue"

	if sk != expected {
		t.Errorf("GetArtistWatchSK() = %v, want %v", sk, expected)
	}
}

func TestGetArtistWatchGSI1PK(t *testing.T) {
	gsi1pk := GetArtistWatchGSI1PK("kylie minogue")
	expected := "ARTIST_WATCH#kylie minogue"

	if gsi1pk != expected {
		t.Errorf("GetArtistWatchGSI1PK() = %v, want %v", gsi1pk, expected)
	}
}

func TestEvent_JSON(t *testing.T) {
	eventDate := time.Date(2026, 6, 15, 20, 0, 0, 0, time.UTC)
	event := Event{
		ID:         "event-001",
		ArtistName: "Kylie Minogue",
		Title:      "Tension Tour 2026",
		Date:       eventDate,
		Venue:      "Madison Square Garden",
		City:       "New York",
		Region:     "NY",
		Country:    "US",
		TicketURL:  "https://tickets.example.com/kylie-msg",
		Status:     "on_sale",
		Source:     "ticketmaster",
	}

	t.Run("serializes all fields to JSON", func(t *testing.T) {
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("json.Marshal() error: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("json.Unmarshal() error: %v", err)
		}

		if result["id"] != "event-001" {
			t.Errorf("JSON id = %v, want %v", result["id"], "event-001")
		}
		if result["artistName"] != "Kylie Minogue" {
			t.Errorf("JSON artistName = %v, want %v", result["artistName"], "Kylie Minogue")
		}
		if result["title"] != "Tension Tour 2026" {
			t.Errorf("JSON title = %v, want %v", result["title"], "Tension Tour 2026")
		}
		if result["venue"] != "Madison Square Garden" {
			t.Errorf("JSON venue = %v, want %v", result["venue"], "Madison Square Garden")
		}
		if result["city"] != "New York" {
			t.Errorf("JSON city = %v, want %v", result["city"], "New York")
		}
		if result["region"] != "NY" {
			t.Errorf("JSON region = %v, want %v", result["region"], "NY")
		}
		if result["country"] != "US" {
			t.Errorf("JSON country = %v, want %v", result["country"], "US")
		}
		if result["ticketUrl"] != "https://tickets.example.com/kylie-msg" {
			t.Errorf("JSON ticketUrl = %v, want %v", result["ticketUrl"], "https://tickets.example.com/kylie-msg")
		}
		if result["status"] != "on_sale" {
			t.Errorf("JSON status = %v, want %v", result["status"], "on_sale")
		}
		if result["source"] != "ticketmaster" {
			t.Errorf("JSON source = %v, want %v", result["source"], "ticketmaster")
		}
	})

	t.Run("omits ticketUrl when empty", func(t *testing.T) {
		eventNoTicket := Event{
			ID:         "event-002",
			ArtistName: "Deadmau5",
			Title:      "Cube V4",
			Date:       eventDate,
			Venue:      "Avant Gardner",
			City:       "Brooklyn",
			Region:     "NY",
			Country:    "US",
			Status:     "announced",
			Source:     "songkick",
		}

		data, err := json.Marshal(eventNoTicket)
		if err != nil {
			t.Fatalf("json.Marshal() error: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("json.Unmarshal() error: %v", err)
		}

		if _, exists := result["ticketUrl"]; exists {
			t.Error("JSON should omit ticketUrl when empty")
		}
	})

	t.Run("roundtrips through JSON correctly", func(t *testing.T) {
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("json.Marshal() error: %v", err)
		}

		var decoded Event
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("json.Unmarshal() error: %v", err)
		}

		if decoded.ID != event.ID {
			t.Errorf("roundtrip ID = %v, want %v", decoded.ID, event.ID)
		}
		if decoded.ArtistName != event.ArtistName {
			t.Errorf("roundtrip ArtistName = %v, want %v", decoded.ArtistName, event.ArtistName)
		}
		if !decoded.Date.Equal(event.Date) {
			t.Errorf("roundtrip Date = %v, want %v", decoded.Date, event.Date)
		}
		if decoded.TicketURL != event.TicketURL {
			t.Errorf("roundtrip TicketURL = %v, want %v", decoded.TicketURL, event.TicketURL)
		}
	})
}

func TestArtistSearchResult_JSON(t *testing.T) {
	t.Run("serializes all fields to JSON", func(t *testing.T) {
		result := ArtistSearchResult{
			Name:           "Kylie Minogue",
			ImageURL:       "https://example.com/kylie.jpg",
			UpcomingEvents: 5,
			Source:         "ticketmaster",
		}

		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("json.Marshal() error: %v", err)
		}

		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("json.Unmarshal() error: %v", err)
		}

		if decoded["name"] != "Kylie Minogue" {
			t.Errorf("JSON name = %v, want %v", decoded["name"], "Kylie Minogue")
		}
		if decoded["imageUrl"] != "https://example.com/kylie.jpg" {
			t.Errorf("JSON imageUrl = %v, want %v", decoded["imageUrl"], "https://example.com/kylie.jpg")
		}
		// JSON numbers decode as float64
		if decoded["upcomingEvents"] != float64(5) {
			t.Errorf("JSON upcomingEvents = %v, want %v", decoded["upcomingEvents"], 5)
		}
		if decoded["source"] != "ticketmaster" {
			t.Errorf("JSON source = %v, want %v", decoded["source"], "ticketmaster")
		}
	})

	t.Run("omits imageUrl when empty", func(t *testing.T) {
		result := ArtistSearchResult{
			Name:           "Unknown Artist",
			UpcomingEvents: 0,
			Source:         "songkick",
		}

		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("json.Marshal() error: %v", err)
		}

		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("json.Unmarshal() error: %v", err)
		}

		if _, exists := decoded["imageUrl"]; exists {
			t.Error("JSON should omit imageUrl when empty")
		}
	})

	t.Run("roundtrips through JSON correctly", func(t *testing.T) {
		original := ArtistSearchResult{
			Name:           "Deadmau5",
			ImageURL:       "https://example.com/deadmau5.jpg",
			UpcomingEvents: 12,
			Source:         "bandsintown",
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("json.Marshal() error: %v", err)
		}

		var decoded ArtistSearchResult
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("json.Unmarshal() error: %v", err)
		}

		if decoded.Name != original.Name {
			t.Errorf("roundtrip Name = %v, want %v", decoded.Name, original.Name)
		}
		if decoded.ImageURL != original.ImageURL {
			t.Errorf("roundtrip ImageURL = %v, want %v", decoded.ImageURL, original.ImageURL)
		}
		if decoded.UpcomingEvents != original.UpcomingEvents {
			t.Errorf("roundtrip UpcomingEvents = %v, want %v", decoded.UpcomingEvents, original.UpcomingEvents)
		}
		if decoded.Source != original.Source {
			t.Errorf("roundtrip Source = %v, want %v", decoded.Source, original.Source)
		}
	})
}

func TestArtistWatch_ToResponse(t *testing.T) {
	watch := NewArtistWatch("user-123", "Kylie Minogue")

	response := watch.ToResponse()

	if response.ArtistName != watch.ArtistName {
		t.Errorf("Response.ArtistName = %v, want %v", response.ArtistName, watch.ArtistName)
	}

	if response.WatchedAt != watch.WatchedAt {
		t.Errorf("Response.WatchedAt = %v, want %v", response.WatchedAt, watch.WatchedAt)
	}
}
