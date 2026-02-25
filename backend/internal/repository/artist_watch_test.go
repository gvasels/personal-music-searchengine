package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockArtistWatchRepo is a minimal mock that implements the ArtistWatch subset
// of the Repository interface. It stores watches in memory keyed by "userID#artistName".
type mockArtistWatchRepo struct {
	watches map[string]models.ArtistWatch
}

func newMockArtistWatchRepo() *mockArtistWatchRepo {
	return &mockArtistWatchRepo{
		watches: make(map[string]models.ArtistWatch),
	}
}

func watchKey(userID, artistName string) string {
	return fmt.Sprintf("%s#%s", userID, models.NormalizeArtistName(artistName))
}

// compile-time check: the Repository interface must include ArtistWatch methods.
// This line will fail to compile until the interface is updated.
var _ interface {
	CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error
	DeleteArtistWatch(ctx context.Context, userID, artistName string) error
	GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error)
	ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error)
} = (repository.Repository)(nil)

// ---------------------------------------------------------------------------
// TestArtistWatch_CreateAndGet
// Creates a watch, retrieves it, and verifies all fields match.
// ---------------------------------------------------------------------------
func TestArtistWatch_CreateAndGet(t *testing.T) {
	ctx := context.Background()

	// Build a concrete repo (DynamoDBRepository) — the method calls below will
	// fail to compile because CreateArtistWatch / GetArtistWatch do not exist
	// on DynamoDBRepository (or the Repository interface) yet.
	var repo repository.Repository

	watch := models.NewArtistWatch("user-123", "Kylie Minogue")

	err := repo.CreateArtistWatch(ctx, *watch)
	require.NoError(t, err)

	got, err := repo.GetArtistWatch(ctx, "user-123", "Kylie Minogue")
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "user-123", got.UserID)
	assert.Equal(t, "Kylie Minogue", got.ArtistName)
	assert.False(t, got.WatchedAt.IsZero(), "WatchedAt should be set")
}

// ---------------------------------------------------------------------------
// TestArtistWatch_Delete
// Creates a watch, deletes it, then verifies GetArtistWatch returns ErrNotFound.
// ---------------------------------------------------------------------------
func TestArtistWatch_Delete(t *testing.T) {
	ctx := context.Background()
	var repo repository.Repository

	watch := models.NewArtistWatch("user-456", "Deadmau5")
	err := repo.CreateArtistWatch(ctx, *watch)
	require.NoError(t, err)

	// Delete the watch
	err = repo.DeleteArtistWatch(ctx, "user-456", "Deadmau5")
	require.NoError(t, err)

	// Verify it's gone
	got, err := repo.GetArtistWatch(ctx, "user-456", "Deadmau5")
	assert.ErrorIs(t, err, repository.ErrNotFound)
	assert.Nil(t, got, "deleted watch should return nil")
}

// ---------------------------------------------------------------------------
// TestArtistWatch_ListWatched
// Creates 3 watches for a single user, lists them, and verifies count and
// pagination behavior.
// ---------------------------------------------------------------------------
func TestArtistWatch_ListWatched(t *testing.T) {
	ctx := context.Background()
	var repo repository.Repository

	artists := []string{"Kylie Minogue", "Deadmau5", "Tiesto"}
	for _, name := range artists {
		watch := models.NewArtistWatch("user-789", name)
		err := repo.CreateArtistWatch(ctx, *watch)
		require.NoError(t, err, "creating watch for %s", name)
	}

	t.Run("list all watched artists", func(t *testing.T) {
		result, err := repo.ListWatchedArtists(ctx, "user-789", 10, "")
		require.NoError(t, err)
		assert.Len(t, result.Items, 3)
		assert.False(t, result.HasMore, "all 3 should fit in limit=10")
	})

	t.Run("paginate with limit=2", func(t *testing.T) {
		allArtists := make(map[string]bool)
		cursor := ""

		for {
			result, err := repo.ListWatchedArtists(ctx, "user-789", 2, cursor)
			require.NoError(t, err)

			for _, w := range result.Items {
				allArtists[w.ArtistName] = true
			}

			if !result.HasMore {
				break
			}
			cursor = result.NextCursor
		}

		assert.Len(t, allArtists, 3, "pagination should yield all 3 watched artists")
	})
}

// ---------------------------------------------------------------------------
// TestArtistWatch_CreateDuplicate
// Creates the same watch twice. The second call should be idempotent (no
// error) or return a well-known duplicate error — either is acceptable.
// ---------------------------------------------------------------------------
func TestArtistWatch_CreateDuplicate(t *testing.T) {
	ctx := context.Background()
	var repo repository.Repository

	watch := models.NewArtistWatch("user-dup", "Moby")
	err := repo.CreateArtistWatch(ctx, *watch)
	require.NoError(t, err)

	// Second create for same user + artist should not error
	err = repo.CreateArtistWatch(ctx, *watch)
	assert.NoError(t, err, "duplicate create should be idempotent (no error)")

	// Verify we still get the watch back
	got, err := repo.GetArtistWatch(ctx, "user-dup", "Moby")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Moby", got.ArtistName)
}

// ---------------------------------------------------------------------------
// TestArtistWatch_GetNonExistent
// Attempts to get a watch that was never created and verifies it returns
// ErrNotFound.
// ---------------------------------------------------------------------------
func TestArtistWatch_GetNonExistent(t *testing.T) {
	ctx := context.Background()
	var repo repository.Repository

	got, err := repo.GetArtistWatch(ctx, "no-such-user", "No Such Artist")
	assert.ErrorIs(t, err, repository.ErrNotFound)
	assert.Nil(t, got, "non-existent watch should return nil")
}
