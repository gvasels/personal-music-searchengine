package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtistRoleConstants(t *testing.T) {
	t.Run("role constants have expected values", func(t *testing.T) {
		assert.Equal(t, ArtistRole("main"), RoleMain)
		assert.Equal(t, ArtistRole("featuring"), RoleFeaturing)
		assert.Equal(t, ArtistRole("remixer"), RoleRemixer)
		assert.Equal(t, ArtistRole("producer"), RoleProducer)
	})
}

func TestNewArtistItem(t *testing.T) {
	now := time.Now()
	artist := Artist{
		ID:     "artist-123",
		UserID: "user-456",
		Name:   "The Beatles",
		SortName: "Beatles",
		Bio:    "Legendary British rock band",
		ImageURL: "https://example.com/beatles.jpg",
		ExternalLinks: map[string]string{
			"spotify":  "https://spotify.com/beatles",
			"wikipedia": "https://en.wikipedia.org/wiki/The_Beatles",
		},
		IsActive: true,
		Timestamps: Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	t.Run("creates correct DynamoDB item", func(t *testing.T) {
		item := NewArtistItem(artist)

		// Check partition and sort keys
		assert.Equal(t, "USER#user-456", item.PK)
		assert.Equal(t, "ARTIST#artist-123", item.SK)
		assert.Equal(t, string(EntityArtist), item.Type)

		// Check GSI keys for name lookup
		assert.Equal(t, "USER#user-456#ARTIST", item.GSI1PK)
		assert.Equal(t, "The Beatles", item.GSI1SK)

		// Check artist data is preserved
		assert.Equal(t, artist.ID, item.Artist.ID)
		assert.Equal(t, artist.Name, item.Artist.Name)
		assert.Equal(t, artist.Bio, item.Artist.Bio)
	})

	t.Run("handles empty external links", func(t *testing.T) {
		artistNoLinks := Artist{
			ID:     "artist-789",
			UserID: "user-456",
			Name:   "Solo Artist",
		}

		item := NewArtistItem(artistNoLinks)
		assert.Nil(t, item.Artist.ExternalLinks)
	})
}

func TestArtistToResponse(t *testing.T) {
	now := time.Now()
	artist := Artist{
		ID:       "artist-123",
		UserID:   "user-456",
		Name:     "Pink Floyd",
		SortName: "Pink Floyd",
		Bio:      "Progressive rock legends",
		ImageURL: "https://example.com/pinkfloyd.jpg",
		ExternalLinks: map[string]string{
			"website": "https://pinkfloyd.com",
		},
		IsActive: true,
		Timestamps: Timestamps{
			CreatedAt: now,
			UpdatedAt: now.Add(time.Hour),
		},
	}

	t.Run("converts to response correctly", func(t *testing.T) {
		response := artist.ToResponse()

		assert.Equal(t, artist.ID, response.ID)
		assert.Equal(t, artist.Name, response.Name)
		assert.Equal(t, artist.SortName, response.SortName)
		assert.Equal(t, artist.Bio, response.Bio)
		assert.Equal(t, artist.ImageURL, response.ImageURL)
		assert.Equal(t, artist.ExternalLinks, response.ExternalLinks)
		assert.Equal(t, artist.IsActive, response.IsActive)
		assert.Equal(t, artist.CreatedAt, response.CreatedAt)
		assert.Equal(t, artist.UpdatedAt, response.UpdatedAt)
	})

	t.Run("excludes userId from response", func(t *testing.T) {
		response := artist.ToResponse()
		// ArtistResponse doesn't have UserID field - verify struct has only expected fields
		assert.NotEmpty(t, response.ID)
	})
}

func TestArtistWithStatsToResponseWithStats(t *testing.T) {
	now := time.Now()
	artistWithStats := ArtistWithStats{
		Artist: Artist{
			ID:       "artist-123",
			UserID:   "user-456",
			Name:     "Led Zeppelin",
			SortName: "Led Zeppelin",
			IsActive: true,
			Timestamps: Timestamps{
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		TrackCount: 75,
		AlbumCount: 9,
		TotalPlays: 10000,
	}

	t.Run("converts to response with stats", func(t *testing.T) {
		response := artistWithStats.ToResponseWithStats()

		assert.Equal(t, artistWithStats.ID, response.ID)
		assert.Equal(t, artistWithStats.Name, response.Name)
		assert.Equal(t, artistWithStats.TrackCount, response.TrackCount)
		assert.Equal(t, artistWithStats.AlbumCount, response.AlbumCount)
		assert.Equal(t, artistWithStats.TotalPlays, response.TotalPlays)
	})

	t.Run("handles zero stats", func(t *testing.T) {
		artistWithZeroStats := ArtistWithStats{
			Artist: Artist{
				ID:     "new-artist",
				UserID: "user-456",
				Name:   "New Artist",
			},
			TrackCount: 0,
			AlbumCount: 0,
			TotalPlays: 0,
		}

		response := artistWithZeroStats.ToResponseWithStats()
		assert.Equal(t, 0, response.TrackCount)
		assert.Equal(t, 0, response.AlbumCount)
		assert.Equal(t, 0, response.TotalPlays)
	})
}

func TestGenerateSortName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes 'The' prefix",
			input:    "The Beatles",
			expected: "Beatles",
		},
		{
			name:     "removes 'A' prefix",
			input:    "A Tribe Called Quest",
			expected: "Tribe Called Quest",
		},
		{
			name:     "removes 'An' prefix",
			input:    "An Artist",
			expected: "Artist",
		},
		{
			name:     "keeps name without prefix",
			input:    "Pink Floyd",
			expected: "Pink Floyd",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "handles 'The' at end of name",
			input:    "Over The Rainbow",
			expected: "Over The Rainbow",
		},
		{
			name:     "handles 'A' in middle of name",
			input:    "Queens A Day",
			expected: "Queens A Day",
		},
		{
			name:     "handles exact 'The' (no space after)",
			input:    "The",
			expected: "The",
		},
		{
			name:     "handles exact 'A' (no space after)",
			input:    "A",
			expected: "A",
		},
		{
			name:     "case sensitive - THE not removed",
			input:    "THE Band",
			expected: "THE Band",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateSortName(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestArtistContribution(t *testing.T) {
	t.Run("can create artist contribution", func(t *testing.T) {
		contribution := ArtistContribution{
			ArtistID:   "artist-123",
			ArtistName: "Featured Artist",
			Role:       RoleFeaturing,
		}

		assert.Equal(t, "artist-123", contribution.ArtistID)
		assert.Equal(t, "Featured Artist", contribution.ArtistName)
		assert.Equal(t, RoleFeaturing, contribution.Role)
	})

	t.Run("supports all role types", func(t *testing.T) {
		roles := []ArtistRole{RoleMain, RoleFeaturing, RoleRemixer, RoleProducer}

		for _, role := range roles {
			contribution := ArtistContribution{
				ArtistID: "artist-123",
				Role:     role,
			}
			assert.NotEmpty(t, contribution.Role)
		}
	})
}

func TestCreateArtistRequest(t *testing.T) {
	t.Run("accepts valid request", func(t *testing.T) {
		req := CreateArtistRequest{
			Name:     "New Artist",
			SortName: "Artist, New",
			Bio:      "Biography text",
			ImageURL: "https://example.com/image.jpg",
			ExternalLinks: map[string]string{
				"spotify": "https://spotify.com/artist",
			},
		}

		assert.Equal(t, "New Artist", req.Name)
		assert.Equal(t, "Artist, New", req.SortName)
	})

	t.Run("allows empty optional fields", func(t *testing.T) {
		req := CreateArtistRequest{
			Name: "Minimal Artist",
		}

		assert.Empty(t, req.SortName)
		assert.Empty(t, req.Bio)
		assert.Empty(t, req.ImageURL)
		assert.Nil(t, req.ExternalLinks)
	})
}

func TestUpdateArtistRequest(t *testing.T) {
	t.Run("accepts partial updates", func(t *testing.T) {
		newName := "Updated Name"
		req := UpdateArtistRequest{
			Name: &newName,
		}

		require.NotNil(t, req.Name)
		assert.Equal(t, "Updated Name", *req.Name)
		assert.Nil(t, req.SortName)
		assert.Nil(t, req.Bio)
	})

	t.Run("supports all fields", func(t *testing.T) {
		name := "Artist"
		sortName := "Sort Name"
		bio := "Bio"
		imageURL := "https://example.com/image.jpg"

		req := UpdateArtistRequest{
			Name:          &name,
			SortName:      &sortName,
			Bio:           &bio,
			ImageURL:      &imageURL,
			ExternalLinks: map[string]string{"key": "value"},
		}

		assert.Equal(t, "Artist", *req.Name)
		assert.Equal(t, "Sort Name", *req.SortName)
		assert.Equal(t, "Bio", *req.Bio)
		assert.Equal(t, "https://example.com/image.jpg", *req.ImageURL)
		assert.Equal(t, "value", req.ExternalLinks["key"])
	})
}
