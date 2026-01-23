# Tasks Document: Data Model Foundation

## Task Overview
| Task | Description | Estimated Files |
|------|-------------|-----------------|
| 1 | Artist Model | 1 |
| 2 | Artist Repository | 1 |
| 3 | Artist Service | 1 |
| 4 | Artist Handlers | 1 |
| 5 | Track Model Updates | 1 |
| 6 | Migration Service | 1 |
| 7 | Frontend Types | 2 |
| 8 | Tests | 4 |

---

- [x] 1. Create Artist model in backend/internal/models/artist.go
  - File: backend/internal/models/artist.go
  - Define Artist struct with all fields (ID, UserID, Name, SortName, Bio, ImageURL, ExternalLinks, IsActive, CreatedAt, UpdatedAt)
  - Define ArtistContribution struct with ArtistID and Role
  - Define role constants (main, featuring, remixer, producer)
  - Define ArtistWithStats for aggregated responses
  - Add DynamoDB marshaling tags
  - Purpose: Establish Artist entity data model
  - _Leverage: backend/internal/models/track.go for struct patterns_
  - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - _Prompt: Implement the task for spec data-model-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer specializing in DynamoDB data modeling | Task: Create Artist model struct with UUID ID, DynamoDB tags following single-table design (PK=USER#{userId}, SK=ARTIST#{artistId}), include ArtistContribution for multi-artist support with role enum | Restrictions: Follow existing model patterns in track.go, use google/uuid for IDs, do not add validation logic here (goes in service) | Success: Artist and ArtistContribution structs compile, proper DynamoDB tags, role constants defined | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 2. Create Artist repository in backend/internal/repository/artist.go
  - File: backend/internal/repository/artist.go
  - Implement ArtistRepository interface (Create, GetByID, GetByName, List, Update, Delete)
  - Use existing DynamoDB client patterns
  - Add GSI query for name lookup
  - Implement batch get for resolving multiple artists
  - Purpose: Data access layer for Artist entities
  - _Leverage: backend/internal/repository/dynamodb.go, backend/internal/repository/track.go_
  - _Requirements: 1.1, 1.5_
  - _Prompt: Implement the task for spec data-model-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with DynamoDB expertise | Task: Implement ArtistRepository with CRUD operations using single-table design, GSI for name lookups (GSI1: PK=USER#{userId}#ARTIST, SK=name), batch get for resolving multiple artist IDs | Restrictions: Follow repository patterns from track.go, use context for all operations, return proper errors | Success: All repository methods implemented, GSI queries working, batch operations efficient | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 3. Create Artist service in backend/internal/service/artist.go
  - File: backend/internal/service/artist.go
  - Implement ArtistService interface
  - Add business logic for create (generate UUID, set timestamps)
  - Add stats aggregation for GetArtist (track count, album count)
  - Add search functionality using name GSI
  - Add soft-delete logic (set IsActive=false)
  - Purpose: Business logic layer for Artist operations
  - _Leverage: backend/internal/service/track.go_
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_
  - _Prompt: Implement the task for spec data-model-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer specializing in service layer architecture | Task: Implement ArtistService with UUID generation, timestamp management, stats aggregation by querying tracks/albums, soft-delete pattern | Restrictions: Use repository interface (not concrete), validate inputs, aggregate stats efficiently (avoid N+1) | Success: Service handles all business rules, stats calculated correctly, soft-delete works | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 4. Create Artist handlers in backend/internal/handlers/artist.go
  - File: backend/internal/handlers/artist.go
  - Implement HTTP handlers: CreateArtist, GetArtist, ListArtists, UpdateArtist, DeleteArtist, GetArtistTracks
  - Register routes in handlers.go
  - Add input validation and error handling
  - Purpose: HTTP API layer for Artist endpoints
  - _Leverage: backend/internal/handlers/track.go, backend/internal/handlers/handlers.go_
  - _Requirements: 1.5_
  - _Prompt: Implement the task for spec data-model-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with Echo framework expertise | Task: Implement HTTP handlers for Artist CRUD following existing patterns, register routes under /api/v1/artists, validate inputs, use getUserIDFromContext | Restrictions: Follow handler patterns from track.go, use handleError for errors, return proper status codes | Success: All endpoints working, proper validation, consistent error responses | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 5. Update Track model to support Artist linking
  - File: backend/internal/models/track.go (modify)
  - Add ArtistID field (string, required for new tracks)
  - Add Artists field ([]ArtistContribution for multi-artist)
  - Add ArtistLegacy field (backup during migration)
  - Update TrackRepository to populate artist info
  - Purpose: Link tracks to Artist entities
  - _Leverage: backend/internal/models/artist.go_
  - _Requirements: 2.1, 2.2, 2.3, 2.4_
  - _Prompt: Implement the task for spec data-model-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Extend Track model with artistId (UUID reference), artists (array of ArtistContribution), artistLegacy (backup field). Update track queries to resolve artistName from artistId | Restrictions: Maintain backward compatibility, existing tracks still work, don't break existing tests | Success: Track model extended, queries resolve artist names, backward compatible | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 6. Create Migration service in backend/internal/service/migration.go
  - File: backend/internal/service/migration.go
  - Implement MigrateArtists function that:
    - Scans all tracks for user
    - Extracts unique artist names
    - Creates Artist entities with deterministic UUIDs (hash of userId+artistName)
    - Updates tracks with artistId references
    - Tracks progress and handles resumption
  - Purpose: Migrate existing string-based artists to entity model
  - _Leverage: backend/internal/repository/track.go, backend/internal/repository/artist.go_
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_
  - _Prompt: Implement the task for spec data-model-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer specializing in data migrations | Task: Create idempotent migration service that extracts unique artists from tracks, creates Artist entities with deterministic UUIDs (UUID v5 from namespace+userId+artistName), updates track references, supports resumption | Restrictions: Must be idempotent (safe to re-run), use batched writes, track progress, don't lose data | Success: Migration completes successfully, all tracks linked, idempotent | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 7. Update Frontend types for Artist entity
  - Files: frontend/src/types/index.ts, frontend/src/lib/api/artists.ts
  - Add Artist interface with all fields
  - Add ArtistContribution interface
  - Update Track interface with artistId and artists fields
  - Create API functions: getArtists, getArtist, createArtist, updateArtist
  - Purpose: Frontend type definitions and API client for Artists
  - _Leverage: frontend/src/types/index.ts, frontend/src/lib/api/tracks.ts_
  - _Requirements: 1.5, 2.2_
  - _Prompt: Implement the task for spec data-model-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: TypeScript/React Developer | Task: Add Artist and ArtistContribution interfaces to types, extend Track with artistId/artists, create API client functions following existing patterns | Restrictions: Follow existing type patterns, maintain backward compatibility with optional fields | Success: Types compile, API functions work, no type errors in existing code | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 8. Write tests for Artist functionality
  - Files: backend/internal/models/artist_test.go, backend/internal/repository/artist_test.go, backend/internal/service/artist_test.go, backend/internal/handlers/artist_test.go
  - Unit tests for model validation
  - Repository tests with DynamoDB Local
  - Service tests with mocked repository
  - Handler tests with mocked service
  - Purpose: Ensure Artist functionality is reliable
  - _Leverage: Existing test patterns in backend/internal/*_test.go_
  - _Requirements: All_
  - _Prompt: Implement the task for spec data-model-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer | Task: Write comprehensive tests for Artist model, repository, service, and handlers following TDD patterns. Test success cases, error cases, edge cases (duplicate names, soft-delete, migration) | Restrictions: Use testify for assertions, mock interfaces properly, test in isolation | Success: 80%+ coverage, all tests pass, edge cases covered | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_
