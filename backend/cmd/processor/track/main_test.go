package main

import (
	"context"
	"sync"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Repository ────────────────────────────────────────────────────────

// mockRepo implements repository.Repository with only the methods used by
// handleRequest. All other methods panic so we catch unexpected calls.
type mockRepo struct {
	repository.Repository // embed to satisfy interface; unused methods panic

	mu             sync.Mutex
	createdTracks  []models.Track
	tracksByArtist map[string][]models.Track   // key = "userId|artist"
	uploadSteps    []uploadStepCall
	albums         map[string]*models.Album    // key = "userId|albumName"
	uploads        map[string]*models.Upload   // key = "userId|uploadId"
	updatedUploads []models.Upload
	createdArtists []models.Artist
	artistsByName  map[string][]*models.Artist // key = "userId|artistName"
}

type uploadStepCall struct {
	UserID   string
	UploadID string
	Step     models.ProcessingStep
	Success  bool
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		tracksByArtist: make(map[string][]models.Track),
		albums:         make(map[string]*models.Album),
		uploads:        make(map[string]*models.Upload),
		artistsByName:  make(map[string][]*models.Artist),
	}
}

func (m *mockRepo) CreateTrack(ctx context.Context, track models.Track) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createdTracks = append(m.createdTracks, track)
	return nil
}

func (m *mockRepo) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := userID + "|" + artist
	return m.tracksByArtist[key], nil
}

func (m *mockRepo) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.uploadSteps = append(m.uploadSteps, uploadStepCall{userID, uploadID, step, success})
	return nil
}

func (m *mockRepo) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := userID + "|" + albumName
	if a, ok := m.albums[key]; ok {
		return a, nil
	}
	a := &models.Album{ID: "album-1", UserID: userID, Title: albumName, Artist: artist}
	m.albums[key] = a
	return a, nil
}

func (m *mockRepo) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := userID + "|" + uploadID
	if u, ok := m.uploads[key]; ok {
		return u, nil
	}
	// Return a minimal upload so the duplicate marking path doesn't error
	return &models.Upload{ID: uploadID, UserID: userID}, nil
}

func (m *mockRepo) UpdateUpload(ctx context.Context, upload models.Upload) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updatedUploads = append(m.updatedUploads, upload)
	return nil
}

func (m *mockRepo) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := userID + "|" + name
	if artists, ok := m.artistsByName[key]; ok {
		return artists, nil
	}
	return nil, nil
}

func (m *mockRepo) CreateArtist(ctx context.Context, artist models.Artist) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createdArtists = append(m.createdArtists, artist)
	// Store for future lookups
	key := artist.UserID + "|" + artist.Name
	a := artist
	m.artistsByName[key] = append(m.artistsByName[key], &a)
	return nil
}

// ─── Trigger Spy ────────────────────────────────────────────────────────────

// triggerCall records the arguments of an audio pipeline trigger invocation.
type triggerCall struct {
	TrackID string
	UserID  string
	S3Key   string
	Title   string
	Artist  string
}

// triggerSpy captures all calls to the audio pipeline trigger function.
type triggerSpy struct {
	mu    sync.Mutex
	calls []triggerCall
}

func (s *triggerSpy) record(ctx context.Context, trackID, userID, s3Key, title, artist string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, triggerCall{trackID, userID, s3Key, title, artist})
}

func (s *triggerSpy) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.calls)
}

// ─── Test Helpers ───────────────────────────────────────────────────────────

// setupTest wires the mock repository and trigger spy into the package-level
// variables so handleRequest uses them. It returns a cleanup function.
func setupTest(t *testing.T, mr *mockRepo, spy *triggerSpy) func() {
	t.Helper()

	origRepo := repo
	origTrigger := doTriggerAudioPipeline

	repo = mr
	doTriggerAudioPipeline = spy.record

	return func() {
		repo = origRepo
		doTriggerAudioPipeline = origTrigger
	}
}

func validEvent() Event {
	return Event{
		UploadID: "11111111-1111-1111-1111-111111111111",
		UserID:   "22222222-2222-2222-2222-222222222222",
		S3Key:    "uploads/22222222-2222-2222-2222-222222222222/song.mp3",
		FileName: "song.mp3",
		Metadata: &models.UploadMetadata{
			Title:    "Test Song",
			Artist:   "Test Artist",
			Duration: 180,
			Format:   "MP3",
		},
	}
}

// ─── Tests ──────────────────────────────────────────────────────────────────

// TestNewTrack_DoesNotTriggerAudioPipeline verifies that when a brand-new
// track is created (no duplicate found), the audio pipeline is NOT triggered.
//
// EXPECTED TO FAIL (Red phase): current code triggers the pipeline for every
// new audio file (lines 180-182 in main.go). The desired behavior is to
// remove that trigger entirely — new tracks should only be triggered later
// when their HLS transcoding completes.
func TestNewTrack_DoesNotTriggerAudioPipeline(t *testing.T) {
	mr := newMockRepo()
	spy := &triggerSpy{}
	cleanup := setupTest(t, mr, spy)
	defer cleanup()

	// No existing tracks → no duplicate will be found
	// (tracksByArtist map is empty)

	event := validEvent()
	resp, err := handleRequest(context.Background(), event)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.TrackID, "track should be created")
	assert.False(t, resp.IsDuplicate, "should not be a duplicate")

	// Core assertion: audio pipeline must NOT be triggered for new tracks
	assert.Equal(t, 0, spy.callCount(),
		"audio pipeline should NOT be triggered for new tracks; "+
			"new tracks wait for HLS transcoding to complete before analysis")
}

// TestDuplicate_WithHLSReady_TriggersAudioPipeline verifies that when a
// duplicate is detected and the existing track has HLSStatus == READY (plus
// needs analysis), the audio pipeline IS triggered.
//
// EXPECTED TO PASS with current code (needsAnalysis is already checked),
// but serves as a regression guard after the HLS check is added.
func TestDuplicate_WithHLSReady_TriggersAudioPipeline(t *testing.T) {
	mr := newMockRepo()
	spy := &triggerSpy{}
	cleanup := setupTest(t, mr, spy)
	defer cleanup()

	existingTrack := models.Track{
		ID:             "existing-track-id",
		UserID:         "22222222-2222-2222-2222-222222222222",
		Title:          "Test Song",
		Artist:         "Test Artist",
		Duration:       180,
		S3Key:          "media/22222222-2222-2222-2222-222222222222/existing.mp3",
		HLSStatus:      models.HLSStatusReady,
		BPM:            0,  // needs analysis
		AnalysisStatus: "", // needs analysis
		EmbeddingID:    "", // needs analysis
	}
	mr.tracksByArtist["22222222-2222-2222-2222-222222222222|Test Artist"] = []models.Track{existingTrack}

	event := validEvent()
	resp, err := handleRequest(context.Background(), event)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.IsDuplicate, "should detect duplicate")
	assert.Equal(t, "existing-track-id", resp.TrackID)

	// Core assertion: pipeline should trigger for duplicate with HLS ready + needs analysis
	assert.Equal(t, 1, spy.callCount(),
		"audio pipeline should be triggered for duplicate with HLS ready and missing analysis")

	// Verify upload was marked as duplicate
	mr.mu.Lock()
	defer mr.mu.Unlock()
	require.Len(t, mr.updatedUploads, 1, "upload should be updated with isDuplicate flag")
	assert.True(t, mr.updatedUploads[0].IsDuplicate, "upload should be marked as duplicate")
}

// TestDuplicate_WithoutHLSReady_DoesNotTrigger verifies that when a
// duplicate is detected but the existing track does NOT have
// HLSStatus == READY, the audio pipeline is NOT triggered, even if
// needsAnalysis would otherwise return true.
//
// EXPECTED TO FAIL (Red phase): current code only checks needsAnalysis()
// without gating on HLSStatus. After the change, duplicates must also have
// HLSStatus == READY before triggering analysis.
func TestDuplicate_WithoutHLSReady_DoesNotTrigger(t *testing.T) {
	hlsStatuses := []struct {
		name   string
		status models.HLSStatus
	}{
		{"empty HLS status", ""},
		{"HLS pending", models.HLSStatusPending},
		{"HLS processing", models.HLSStatusProcessing},
		{"HLS failed", models.HLSStatusFailed},
	}

	for _, tc := range hlsStatuses {
		t.Run(tc.name, func(t *testing.T) {
			mr := newMockRepo()
			spy := &triggerSpy{}
			cleanup := setupTest(t, mr, spy)
			defer cleanup()

			existingTrack := models.Track{
				ID:             "existing-track-id",
				UserID:         "22222222-2222-2222-2222-222222222222",
				Title:          "Test Song",
				Artist:         "Test Artist",
				Duration:       180,
				S3Key:          "media/22222222-2222-2222-2222-222222222222/existing.mp3",
				HLSStatus:      tc.status,
				BPM:            0,  // needs analysis
				AnalysisStatus: "", // needs analysis
				EmbeddingID:    "", // needs analysis
			}
			mr.tracksByArtist["22222222-2222-2222-2222-222222222222|Test Artist"] = []models.Track{existingTrack}

			event := validEvent()
			resp, err := handleRequest(context.Background(), event)

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.True(t, resp.IsDuplicate, "should detect duplicate")
			assert.Equal(t, "existing-track-id", resp.TrackID)

			// Core assertion: pipeline must NOT trigger when HLS is not ready
			assert.Equal(t, 0, spy.callCount(),
				"audio pipeline should NOT be triggered when HLS status is %q; "+
					"analysis requires HLS transcoding to be complete", tc.status)
		})
	}
}

// TestDuplicate_WithHLSReady_AnalysisCompleted_DoesNotTrigger verifies that
// even with HLS ready, if analysis is already completed, no trigger fires.
// This is a regression test — current needsAnalysis() already handles this.
func TestDuplicate_WithHLSReady_AnalysisCompleted_DoesNotTrigger(t *testing.T) {
	mr := newMockRepo()
	spy := &triggerSpy{}
	cleanup := setupTest(t, mr, spy)
	defer cleanup()

	existingTrack := models.Track{
		ID:             "existing-track-id",
		UserID:         "22222222-2222-2222-2222-222222222222",
		Title:          "Test Song",
		Artist:         "Test Artist",
		Duration:       180,
		S3Key:          "media/22222222-2222-2222-2222-222222222222/existing.mp3",
		HLSStatus:      models.HLSStatusReady,
		BPM:            128,
		AnalysisStatus: "COMPLETED",
		EmbeddingID:    "emb-123",
	}
	mr.tracksByArtist["22222222-2222-2222-2222-222222222222|Test Artist"] = []models.Track{existingTrack}

	event := validEvent()
	resp, err := handleRequest(context.Background(), event)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.IsDuplicate)

	// Should NOT trigger — analysis already complete
	assert.Equal(t, 0, spy.callCount(),
		"audio pipeline should NOT trigger when analysis is already completed")
}

// TestDuplicate_IncompleteTrack_NotTreatedAsDuplicate verifies that a track
// whose file was never moved to permanent storage (S3Key still in uploads/)
// is NOT treated as a duplicate. This prevents the edge case where a previous
// failed upload leaves a broken track record, and re-uploading the same file
// would clean up the new upload's file while pointing at an inaccessible track.
func TestDuplicate_IncompleteTrack_NotTreatedAsDuplicate(t *testing.T) {
	mr := newMockRepo()
	spy := &triggerSpy{}
	cleanup := setupTest(t, mr, spy)
	defer cleanup()

	// Simulate a track from a failed previous upload — file never moved from uploads/
	incompleteTrack := models.Track{
		ID:       "incomplete-track-id",
		UserID:   "22222222-2222-2222-2222-222222222222",
		Title:    "Test Song",
		Artist:   "Test Artist",
		Duration: 180,
		S3Key:    "uploads/22222222-2222-2222-2222-222222222222/song.mp3", // never moved
	}
	mr.tracksByArtist["22222222-2222-2222-2222-222222222222|Test Artist"] = []models.Track{incompleteTrack}

	event := validEvent()
	resp, err := handleRequest(context.Background(), event)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.IsDuplicate, "incomplete track should NOT be treated as a duplicate")
	assert.NotEqual(t, "incomplete-track-id", resp.TrackID, "should create a new track, not reuse the incomplete one")

	// A new track should have been created
	mr.mu.Lock()
	defer mr.mu.Unlock()
	require.Len(t, mr.createdTracks, 1, "a new track should be created")
}
