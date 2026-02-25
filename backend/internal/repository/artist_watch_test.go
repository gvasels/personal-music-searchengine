package repository_test

import (
	"context"
	"fmt"
	"sort"
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

func (m *mockArtistWatchRepo) CreateArtistWatch(_ context.Context, watch models.ArtistWatch) error {
	key := watchKey(watch.UserID, watch.ArtistName)
	m.watches[key] = watch
	return nil
}

func (m *mockArtistWatchRepo) DeleteArtistWatch(_ context.Context, userID, artistName string) error {
	key := watchKey(userID, artistName)
	if _, ok := m.watches[key]; !ok {
		return repository.ErrNotFound
	}
	delete(m.watches, key)
	return nil
}

func (m *mockArtistWatchRepo) GetArtistWatch(_ context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	key := watchKey(userID, artistName)
	w, ok := m.watches[key]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return &w, nil
}

func (m *mockArtistWatchRepo) ListWatchedArtists(_ context.Context, userID string, limit int, _ string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	prefix := userID + "#"
	var items []models.ArtistWatch
	for k, v := range m.watches {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			items = append(items, v)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].ArtistName < items[j].ArtistName
	})
	hasMore := false
	if limit > 0 && len(items) > limit {
		items = items[:limit]
		hasMore = true
	}
	return &repository.PaginatedResult[models.ArtistWatch]{
		Items:   items,
		HasMore: hasMore,
	}, nil
}

// compile-time check: the Repository interface must include ArtistWatch methods.
var _ interface {
	CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error
	DeleteArtistWatch(ctx context.Context, userID, artistName string) error
	GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error)
	ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error)
} = (repository.Repository)(nil)

// ---------------------------------------------------------------------------
// TestArtistWatch_CreateAndGet
// ---------------------------------------------------------------------------
func TestArtistWatch_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	repo := newMockArtistWatchRepo()

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
// ---------------------------------------------------------------------------
func TestArtistWatch_Delete(t *testing.T) {
	ctx := context.Background()
	repo := newMockArtistWatchRepo()

	watch := models.NewArtistWatch("user-456", "Deadmau5")
	err := repo.CreateArtistWatch(ctx, *watch)
	require.NoError(t, err)

	err = repo.DeleteArtistWatch(ctx, "user-456", "Deadmau5")
	require.NoError(t, err)

	got, err := repo.GetArtistWatch(ctx, "user-456", "Deadmau5")
	assert.ErrorIs(t, err, repository.ErrNotFound)
	assert.Nil(t, got, "deleted watch should return nil")
}

// ---------------------------------------------------------------------------
// TestArtistWatch_ListWatched
// ---------------------------------------------------------------------------
func TestArtistWatch_ListWatched(t *testing.T) {
	ctx := context.Background()
	repo := newMockArtistWatchRepo()

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
		result, err := repo.ListWatchedArtists(ctx, "user-789", 2, "")
		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
		assert.True(t, result.HasMore)
	})
}

// ---------------------------------------------------------------------------
// TestArtistWatch_CreateDuplicate
// ---------------------------------------------------------------------------
func TestArtistWatch_CreateDuplicate(t *testing.T) {
	ctx := context.Background()
	repo := newMockArtistWatchRepo()

	watch := models.NewArtistWatch("user-dup", "Moby")
	err := repo.CreateArtistWatch(ctx, *watch)
	require.NoError(t, err)

	err = repo.CreateArtistWatch(ctx, *watch)
	assert.NoError(t, err, "duplicate create should be idempotent (no error)")

	got, err := repo.GetArtistWatch(ctx, "user-dup", "Moby")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Moby", got.ArtistName)
}

// ---------------------------------------------------------------------------
// TestArtistWatch_GetNonExistent
// ---------------------------------------------------------------------------
func TestArtistWatch_GetNonExistent(t *testing.T) {
	ctx := context.Background()
	repo := newMockArtistWatchRepo()

	got, err := repo.GetArtistWatch(ctx, "no-such-user", "No Such Artist")
	assert.ErrorIs(t, err, repository.ErrNotFound)
	assert.Nil(t, got, "non-existent watch should return nil")
}
