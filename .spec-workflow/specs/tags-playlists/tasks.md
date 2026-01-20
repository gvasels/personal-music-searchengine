# Tasks - Tags & Playlists (Epic 4)

## Epic: Tags & Playlists
**Status**: Not Started
**Wave**: 3

---

## Design Decisions

Based on clarification:
- **Tag duplicates in search**: Deduplicate silently
- **Tag name case**: Case-insensitive (normalize to lowercase throughout)
- **Invalid playlist position**: Append to end
- **Tag filter limit**: No limit

---

## Group 9: Tag Filtering in Search

### Task 9.1: Add filterByTags Method to Search Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/search.go`

**Description**: Implement tag filtering that validates tags exist and filters search results using AND logic.

**Functions**:
| Function | Description |
|----------|-------------|
| `filterByTags(ctx, userID, results, tags)` | Validates all tags exist, then filters results to only tracks with ALL specified tags. Returns NotFoundError if any tag doesn't exist. |
| `convertFilters` (update) | Pass tags through to SearchFilters for later filtering. |
| `Search` (update) | Call filterByTags after Nixiesearch query when tags are specified. |

**Implementation Details**:
```go
func (s *searchServiceImpl) filterByTags(ctx context.Context, userID string, results []search.SearchResult, tags []string) ([]search.SearchResult, error) {
    if len(tags) == 0 {
        return results, nil
    }

    // 1. Deduplicate and normalize tags (lowercase)
    seen := make(map[string]bool)
    uniqueTags := make([]string, 0, len(tags))
    for _, tag := range tags {
        normalized := strings.ToLower(strings.TrimSpace(tag))
        if normalized != "" && !seen[normalized] {
            seen[normalized] = true
            uniqueTags = append(uniqueTags, normalized)
        }
    }

    if len(uniqueTags) == 0 {
        return results, nil
    }

    // 2. Validate all tags exist
    for _, tagName := range uniqueTags {
        _, err := s.repo.GetTag(ctx, userID, tagName)
        if err != nil {
            if err == repository.ErrNotFound {
                return nil, models.NewNotFoundError("Tag", tagName)
            }
            return nil, fmt.Errorf("failed to check tag %s: %w", tagName, err)
        }
    }

    // 3. Build intersection of track IDs (AND logic)
    var validTrackIDs map[string]bool
    for i, tagName := range uniqueTags {
        taggedTracks, err := s.repo.GetTracksByTag(ctx, userID, tagName)
        if err != nil {
            return nil, err
        }

        tagTrackIDs := make(map[string]bool)
        for _, track := range taggedTracks {
            tagTrackIDs[track.ID] = true
        }

        if i == 0 {
            validTrackIDs = tagTrackIDs
        } else {
            // Intersect
            for id := range validTrackIDs {
                if !tagTrackIDs[id] {
                    delete(validTrackIDs, id)
                }
            }
        }

        if len(validTrackIDs) == 0 {
            return []search.SearchResult{}, nil
        }
    }

    // 4. Filter results
    filtered := make([]search.SearchResult, 0)
    for _, result := range results {
        if validTrackIDs[result.ID] {
            filtered = append(filtered, result)
        }
    }

    return filtered, nil
}
```

**Acceptance Criteria**:
- [ ] filterByTags returns results unchanged when no tags specified
- [ ] filterByTags returns NotFoundError when any tag doesn't exist
- [ ] filterByTags uses AND logic for multiple tags
- [ ] Search method integrates filterByTags correctly
- [ ] SearchResponse.TotalResults reflects filtered count when tags applied

---

### Task 9.2: Add filterByTags Unit Tests
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/search_test.go`

**Description**: Add comprehensive tests for the filterByTags function.

**Tests**:
| Test | Description |
|------|-------------|
| `TestFilterByTags_EmptyTags` | Empty tags array returns all results unchanged |
| `TestFilterByTags_TagNotFound` | Non-existent tag returns NotFoundError with tag name |
| `TestFilterByTags_SingleTag_Success` | Single tag filters results correctly |
| `TestFilterByTags_MultipleTags_ANDLogic` | Multiple tags returns only tracks with ALL tags |
| `TestFilterByTags_SecondTagNotFound` | Error on second tag still returns NotFoundError |
| `TestFilterByTags_NoMatchingTracks` | Returns empty array when no tracks match |

**Mock Requirements**:
- Create `MockTagFilterRepository` with mockable `GetTag` and `GetTracksByTag`
- Implement full `repository.Repository` interface (stubs for unused methods)

**Acceptance Criteria**:
- [ ] All 6 tests pass
- [ ] Tests cover success, error, and edge cases
- [ ] Mock repository properly set up

---

## Group 10: Tag Name Normalization

### Task 10.1: Normalize Tag Names to Lowercase
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/tag.go`

**Description**: Update tag service to normalize all tag names to lowercase for case-insensitive matching.

**Functions to Update**:
| Function | Change |
|----------|--------|
| `CreateTag` | Normalize `req.Name` to lowercase before storage |
| `GetTag` | Normalize `tagName` param to lowercase before lookup |
| `UpdateTag` | Normalize both old and new tag names |
| `DeleteTag` | Normalize `tagName` param to lowercase |
| `AddTagsToTrack` | Normalize all tag names in `req.Tags` |
| `RemoveTagFromTrack` | Normalize `tagName` param |
| `GetTracksByTag` | Normalize `tagName` param |

**Helper Function**:
```go
// normalizeTagName converts tag name to lowercase for consistent storage/lookup
func normalizeTagName(name string) string {
    return strings.ToLower(strings.TrimSpace(name))
}
```

**Tests to Add** (in Task 11.1):
| Test | Description |
|------|-------------|
| `TestCreateTag_NormalizesName` | "Rock" stored as "rock" |
| `TestGetTag_CaseInsensitive` | "ROCK" finds tag created as "rock" |
| `TestAddTagsToTrack_NormalizesNames` | Tags normalized before storage |

**Acceptance Criteria**:
- [ ] All tag operations use normalized (lowercase) names
- [ ] Existing case-sensitive tests updated to expect lowercase
- [ ] Helper function added and used consistently

---

## Group 11: Tag Service Tests

### Task 11.1: Create Tag Service Unit Tests
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/tag_test.go` (new)

**Description**: Comprehensive unit tests for all TagService methods.

**Tests**:
| Test | Description |
|------|-------------|
| `TestCreateTag_Success` | Creates tag and returns response |
| `TestCreateTag_AlreadyExists` | Returns ConflictError for duplicate name |
| `TestGetTag_Success` | Returns tag by name |
| `TestGetTag_NotFound` | Returns NotFoundError for unknown tag |
| `TestUpdateTag_Success` | Updates color successfully |
| `TestUpdateTag_Rename` | Renames tag when new name provided |
| `TestUpdateTag_RenameConflict` | Returns ConflictError when renaming to existing name |
| `TestUpdateTag_NotFound` | Returns NotFoundError for unknown tag |
| `TestDeleteTag_Success` | Deletes existing tag |
| `TestDeleteTag_NotFound` | Returns NotFoundError for unknown tag |
| `TestListTags_Success` | Returns all user tags |
| `TestListTags_Empty` | Returns empty array when no tags |
| `TestAddTagsToTrack_Success` | Adds tags to track, creates new tags |
| `TestAddTagsToTrack_TrackNotFound` | Returns NotFoundError for unknown track |
| `TestRemoveTagFromTrack_Success` | Removes tag from track |
| `TestGetTracksByTag_Success` | Returns tracks with specified tag |
| `TestGetTracksByTag_TagNotFound` | Returns NotFoundError for unknown tag |

**Mock Requirements**:
```go
type MockTagRepository struct {
    mock.Mock
}

// Implement mockable methods for:
// - CreateTag, GetTag, UpdateTag, DeleteTag, ListTags
// - GetTrack, UpdateTrack
// - AddTagsToTrack, RemoveTagFromTrack, GetTracksByTag
```

**Acceptance Criteria**:
- [ ] All 17 tests pass
- [ ] 80%+ code coverage for tag.go
- [ ] Tests cover all error paths

---

## Group 12: Playlist Service Tests

### Task 12.1: Create Playlist Service Unit Tests
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/playlist_test.go` (new)

**Description**: Comprehensive unit tests for all PlaylistService methods.

**Tests**:
| Test | Description |
|------|-------------|
| `TestCreatePlaylist_Success` | Creates playlist and returns response |
| `TestGetPlaylist_Success` | Returns playlist with tracks |
| `TestGetPlaylist_NotFound` | Returns NotFoundError for unknown playlist |
| `TestGetPlaylist_WithCoverArt` | Includes cover art URL when present |
| `TestUpdatePlaylist_Success` | Updates name, description, visibility |
| `TestUpdatePlaylist_NotFound` | Returns NotFoundError for unknown playlist |
| `TestDeletePlaylist_Success` | Deletes existing playlist |
| `TestDeletePlaylist_NotFound` | Returns NotFoundError for unknown playlist |
| `TestListPlaylists_Success` | Returns paginated playlists |
| `TestListPlaylists_Empty` | Returns empty array when no playlists |
| `TestAddTracks_Success` | Adds tracks at default position (append) |
| `TestAddTracks_AtPosition` | Adds tracks at specific position |
| `TestAddTracks_TrackNotFound` | Returns NotFoundError for unknown track |
| `TestAddTracks_PlaylistNotFound` | Returns NotFoundError for unknown playlist |
| `TestRemoveTracks_Success` | Removes tracks and updates stats |
| `TestRemoveTracks_PlaylistNotFound` | Returns NotFoundError for unknown playlist |

**Mock Requirements**:
```go
type MockPlaylistRepository struct {
    mock.Mock
}

type MockPlaylistS3Repository struct {
    mock.Mock
}

// Implement mockable methods for:
// - CreatePlaylist, GetPlaylist, UpdatePlaylist, DeletePlaylist, ListPlaylists
// - AddTracksToPlaylist, RemoveTracksFromPlaylist, GetPlaylistTracks
// - GetTrack (for validation)
// - GeneratePresignedDownloadURL (S3)
```

**Acceptance Criteria**:
- [ ] All 16 tests pass
- [ ] 80%+ code coverage for playlist.go
- [ ] Tests cover all error paths

---

## Group 13: Documentation & Cleanup

### Task 13.1: Update CHANGELOG
**Status**: [ ] Pending
**Files**:
- `CHANGELOG.md`

**Description**: Document Epic 4 changes following Keep a Changelog format.

**Content**:
```markdown
#### Epic 4: Tags & Playlists
- Tag filtering in search
  - Filter search results by multiple tags (AND logic)
  - Tags stored in DynamoDB, filtered post-search
  - Returns NotFoundError if tag doesn't exist
- Unit tests for tag service (17 tests)
- Unit tests for playlist service (16 tests)
- Unit tests for filterByTags (6 tests)
```

**Acceptance Criteria**:
- [ ] CHANGELOG updated under [Unreleased]
- [ ] All Epic 4 features documented

### Task 13.2: Update Service CLAUDE.md
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/CLAUDE.md`

**Description**: Update service documentation with filterByTags details.

**Acceptance Criteria**:
- [ ] filterByTags function documented
- [ ] Test files listed

---

## Summary

| Group | Tasks | Description |
|-------|-------|-------------|
| Group 9 | 2 | Tag filtering in search + tests |
| Group 10 | 1 | Tag name normalization (lowercase) |
| Group 11 | 1 | Tag service unit tests |
| Group 12 | 1 | Playlist service unit tests |
| Group 13 | 2 | Documentation updates |
| **Total** | **7** | |

---

## Test Plan Summary

### Unit Tests to Create

| File | Test Count | Coverage Target |
|------|------------|-----------------|
| `service/search_test.go` | +6 tests | filterByTags 100% |
| `service/tag_test.go` | 17 tests | tag.go 80%+ |
| `service/playlist_test.go` | 16 tests | playlist.go 80%+ |
| **Total** | **39 tests** | |

### Test Execution
```bash
# Run all service tests
go test ./internal/service/... -v

# Run with coverage
go test ./internal/service/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

---

## Dependencies

### Between Tasks
- Task 9.2 depends on Task 9.1 (tests need implementation)
- Tasks 10.1 and 11.1 can run in parallel (independent)
- Task 12.1 should run last (documents completed work)

### External
- `github.com/stretchr/testify` - already in go.mod

---

## PR Checklist

After completing all tasks:
- [ ] All tests pass (`go test ./internal/service/... -v`)
- [ ] Code builds (`go build ./...`)
- [ ] Coverage meets 80% target
- [ ] CHANGELOG.md updated
- [ ] CLAUDE.md files updated
- [ ] Create PR to main with all changes
