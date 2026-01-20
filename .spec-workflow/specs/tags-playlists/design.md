# Design - Tags & Playlists (Epic 4)

## Overview

This epic adds tag filtering to search and comprehensive unit tests for existing tag and playlist services. The core functionality (models, repositories, services, handlers) already exists; this epic ensures proper testing and completes the search integration.

## Steering Document Alignment

### Technical Standards (tech.md)
- Go 1.22+ with standard library patterns
- testify for assertions and mocking
- Table-driven tests where applicable
- Repository pattern with interface-based mocking

### Project Structure (structure.md)
- Tests in same package as implementation (`*_test.go`)
- Mock implementations follow existing patterns in `search_test.go`
- Service tests use mock repositories

## Code Reuse Analysis

### Existing Components to Leverage
- **`service/search.go`**: Existing search service to extend with tag filtering
- **`service/search_test.go`**: Existing mock patterns for search tests
- **`repository/repository.go`**: Repository interface for mocking
- **`models/errors.go`**: NewNotFoundError, NewValidationError helpers

### Integration Points
- **Search Service**: Add `filterByTags` method, integrate with `Search` method
- **Repository Interface**: Use existing `GetTag`, `GetTracksByTag` methods
- **API Models**: Use existing `SearchFilters.Tags` field

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Search Request                           │
│                  (with filters.tags: ["tag1", "tag2"])          │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Search Service                               │
│                                                                  │
│  1. Validate query (existing)                                   │
│  2. Execute Nixiesearch query (existing)                        │
│  3. Apply tag filtering (NEW - filterByTags)                    │
│  4. Enrich with cover art (existing)                            │
│  5. Return filtered results                                      │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                    filterByTags Logic                            │
│                                                                  │
│  1. If no tags, return results unchanged                        │
│  2. Validate all tags exist (GetTag for each)                   │
│     - If any tag not found → return NotFoundError               │
│  3. For each tag, get tagged tracks (GetTracksByTag)            │
│  4. Intersect track IDs (AND logic)                             │
│  5. Filter search results to only matching IDs                  │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Component 1: Tag Filtering in Search
- **Purpose**: Filter search results to only include tracks with ALL specified tags
- **Location**: `backend/internal/service/search.go`
- **New Method**: `filterByTags(ctx, userID, results, tags) ([]SearchResult, error)`
- **Dependencies**: `repository.Repository` (GetTag, GetTracksByTag)
- **Reuses**: Existing `searchServiceImpl` struct

### Component 2: Tag Service Tests
- **Purpose**: Comprehensive unit tests for all TagService methods
- **Location**: `backend/internal/service/tag_test.go`
- **Dependencies**: Mock repository implementing `repository.Repository`
- **Reuses**: Mock patterns from `search_test.go`

### Component 3: Playlist Service Tests
- **Purpose**: Comprehensive unit tests for all PlaylistService methods
- **Location**: `backend/internal/service/playlist_test.go`
- **Dependencies**: Mock repository, mock S3 repository
- **Reuses**: Mock patterns from `search_test.go`

## Data Flow

### Tag Filtering in Search

```
Input: SearchRequest with filters.tags = ["favorites", "rock"]

1. Execute Nixiesearch query → results = [track-1, track-2, track-3]

2. filterByTags:
   a. GetTag("favorites") → exists ✓
   b. GetTag("rock") → exists ✓
   c. GetTracksByTag("favorites") → [track-1, track-2]
   d. GetTracksByTag("rock") → [track-1, track-3]
   e. Intersect → [track-1] (only track with BOTH tags)
   f. Filter results → [track-1]

3. Return: filtered results with only track-1
```

### Tag Not Found Flow

```
Input: SearchRequest with filters.tags = ["nonexistent"]

1. Execute Nixiesearch query → results = [track-1, track-2]

2. filterByTags:
   a. GetTag("nonexistent") → ErrNotFound
   b. Return NewNotFoundError("Tag", "nonexistent")

3. Return: error (not empty results)
```

## Error Handling

### Error Scenarios

1. **Tag Not Found**
   - **Handling**: Return `models.NewNotFoundError("Tag", tagName)`
   - **User Impact**: Clear error message: "Tag 'tagname' not found"
   - **HTTP Status**: 404

2. **Repository Error**
   - **Handling**: Wrap and return error
   - **User Impact**: "Internal server error" with logged details
   - **HTTP Status**: 500

3. **Empty Results After Filtering**
   - **Handling**: Return empty results (not an error)
   - **User Impact**: "No results match your filters"
   - **HTTP Status**: 200 with empty tracks array

## Testing Strategy

### Unit Testing Approach

**Tag Service Tests** (`tag_test.go`):
- Mock `repository.Repository` interface
- Test each service method independently
- Cover success, not found, and error scenarios

**Playlist Service Tests** (`playlist_test.go`):
- Mock both `repository.Repository` and `repository.S3Repository`
- Test CRUD operations
- Test track management (add/remove with positions)

**filterByTags Tests** (in `search_test.go`):
- Test tag validation (not found error)
- Test single tag filtering
- Test multiple tags with AND logic
- Test empty tags (no filtering)

### Test Coverage Targets

| Component | Target | Focus |
|-----------|--------|-------|
| Tag Service | 80%+ | All public methods |
| Playlist Service | 80%+ | All public methods |
| filterByTags | 100% | All branches |

### Mock Strategy

Use testify/mock with custom mock structs implementing repository interfaces:

```go
// MockTagRepository for tag service tests
type MockTagRepository struct {
    mock.Mock
}

func (m *MockTagRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
    args := m.Called(ctx, userID, tagName)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Tag), args.Error(1)
}
// ... other methods
```

## Dependencies

### Internal
- `internal/models` - Tag, Playlist, Track, errors
- `internal/repository` - Repository interface, ErrNotFound
- `internal/search` - SearchResult type

### External
- `github.com/stretchr/testify` - Testing assertions and mocks

## Files to Modify

| File | Change |
|------|--------|
| `service/search.go` | Add `filterByTags` method, integrate in `Search` |
| `service/search_test.go` | Add filterByTags tests |
| `service/tag_test.go` | New file - tag service unit tests |
| `service/playlist_test.go` | New file - playlist service unit tests |
| `CHANGELOG.md` | Document Epic 4 changes |

## Acceptance Criteria Summary

- [ ] filterByTags validates all tags exist before filtering
- [ ] filterByTags returns NotFoundError for missing tags
- [ ] filterByTags uses AND logic for multiple tags
- [ ] Search service integrates filterByTags correctly
- [ ] Tag service has 80%+ test coverage
- [ ] Playlist service has 80%+ test coverage
- [ ] All tests pass
