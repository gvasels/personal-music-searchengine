# CLAUDE.md - Personal Music Search Engine

## HARD REQUIREMENTS

**Non-negotiable rules. Violations require rollback and restart.**

| # | Rule | Requirement | If Violated |
|---|------|-------------|-------------|
| 1 | **SDLC Workflow** | Use `/sdlc` for ALL features: SPEC → TEST → CODE → BUILD → DOCS | Start over with `/sdlc` |
| 2 | **TDD Enforcement** | Spawn `test-engineer` BEFORE `implementation-agent`. Tests must fail first (Red). | Delete code, restart Phase 2 |
| 3 | **Spec-First** | All features need specs in `.spec-workflow/specs/` before coding | Create specs first |
| 4 | **Documentation** | Update CLAUDE.md + CHANGELOG.md in every directory with changes | Update before PR |
| 5 | **Quality Gates** | 80% test coverage, all tests pass, no critical security vulnerabilities | Fix before merge |
| 6 | **Wiring Checklist** | Follow `.claude/docs/wiring-checklist.md` for new services/handlers/routes | Complete checklist before PR |

---

## Project Overview

| Field | Value |
|-------|-------|
| **Name** | Personal Music Search Engine |
| **Description** | Multi-user music library platform with role-based access, artist profiles, public playlists, and follow system |
| **Status** | Development |
| **AWS Account** | 887395463840 |
| **AWS Profile** | `gvasels-muza` |
| **Region** | us-east-1 |

## Technology Stack

| Layer | Technologies |
|-------|--------------|
| **Backend** | Go 1.22+, Echo v4, AWS Lambda, dhowden/tag |
| **Frontend** | React 18, TanStack Router/Query, Tailwind + DaisyUI 5, Zustand, Vite |
| **Infrastructure** | AWS, OpenTofu 1.8+ (NOT Terraform), GitHub Actions |
| **Data** | DynamoDB (single-table), S3, Nixiesearch, Cognito (with Groups), CloudFront |
| **Auth** | Cognito User Pools with Groups for role-based access control |

---

## SDLC Workflow

```
/sdlc → SPEC → TEST → CODE → BUILD → DOCS → PR
         ↓      ↓      ↓       ↓       ↓
       spec-  test-  impl-   lint/  docs-
       workflow engineer agent  test  generator
```

| Phase | Purpose | Artifacts | Quality Gate |
|-------|---------|-----------|--------------|
| 1. SPEC | Requirements & design | `requirements.md`, `design.md`, `tasks.md` | Design approved |
| 2. TEST | TDD Red - failing tests | `*_test.go`, `*.test.ts` | Tests written & failing |
| 3. CODE | TDD Green - minimal code | Source files | All tests passing |
| 4. BUILD | Lint, type check, coverage | Build output | Coverage 80%+ |
| 5. DOCS | API docs, CLAUDE.md | Documentation | Files updated |

**Prerequisites**: Each phase requires previous phase completion.

---

## Repository Structure

```
├── backend/                 # Go Lambda services
│   ├── cmd/{api,indexer,processor}/  # Lambda entrypoints
│   └── internal/
│       ├── handlers/        # HTTP handlers + middleware/auth.go
│       ├── models/          # Domain models (role.go, artist_profile.go, follow.go)
│       ├── repository/      # DynamoDB access (artist_profile.go, follow.go)
│       ├── service/         # Business logic (role.go, artist_profile.go, follow.go)
│       └── {metadata,search}/
├── frontend/src/            # React app
│   ├── components/
│   │   ├── artist-profile/  # ArtistProfileCard, EditArtistProfileModal
│   │   ├── follow/          # FollowButton, FollowersList, FollowingList
│   │   ├── playlist/        # VisibilitySelector, CreatePlaylistModal
│   │   └── studio/          # ModuleCard, CrateList, HotCueBar
│   ├── hooks/               # useAuth, useFollow, useArtistProfiles, useFeatureFlags
│   ├── lib/api/             # API clients (artistProfiles.ts, follows.ts)
│   └── routes/artists/entity/  # Artist profile pages
├── infrastructure/          # OpenTofu IaC
│   └── {global,shared,backend,frontend}/
├── scripts/                 # Utility scripts
│   ├── bootstrap-admin.sh   # Admin user promotion
│   └── migrations/          # Data migration scripts
├── .claude/                 # Claude Code automation
├── .spec-workflow/specs/    # Feature specifications
└── .github/workflows/       # CI/CD
```

---

## Role-Based Access Control

### User Roles

| Role | Description | Cognito Group |
|------|-------------|---------------|
| `guest` | Unauthenticated, browse only | (none) |
| `subscriber` | Default authenticated user | `subscriber` |
| `artist` | Can create artist profile, upload tracks | `artist` |
| `admin` | Full system access | `admin` |

### Permissions by Role

| Permission | Guest | Subscriber | Artist | Admin |
|------------|-------|------------|--------|-------|
| Browse public content | ✓ | ✓ | ✓ | ✓ |
| Listen to music | | ✓ | ✓ | ✓ |
| Create playlists | | ✓ | ✓ | ✓ |
| Follow artists | | ✓ | ✓ | ✓ |
| Upload tracks | | | ✓ | ✓ |
| Manage artist profile | | | ✓ | ✓ |
| View ALL tracks globally | | | | ✓ |
| Delete ANY track | | | | ✓ |
| Manage users/roles | | | | ✓ |

### Admin-Specific Capabilities

| Capability | Description |
|------------|-------------|
| **Global Track View** | Admins see ALL tracks from all users (not just own + public) |
| **Delete Any Track** | Admins can delete tracks owned by any user (including S3 + HLS files) |
| **User Management** | Search users by email, view details, change roles, enable/disable |
| **Role Simulation** | Test UI as different roles via Admin Panel without changing real role |
| **Cognito Sync** | Sync DynamoDB roles with Cognito groups for consistency |

### Track Visibility Levels

| Level | Who Can Access |
|-------|----------------|
| `private` | Owner only (and admins with global access) |
| `unlisted` | Anyone with direct link |
| `public` | Everyone, discoverable in search |

### Access Control Implementation

- **403 Forbidden**: Returned when user lacks permission to access a resource
- **404 Not Found**: Returned ONLY when the resource truly doesn't exist
- Service layer enforces visibility (not just handlers)
- `hasGlobal` parameter determines admin/global read permissions
- Real-time DB role checking overrides JWT claims for critical operations

### Playlist Visibility

| Level | Description |
|-------|-------------|
| `private` | Only owner can see |
| `unlisted` | Anyone with link can see |
| `public` | Discoverable by all users |

---

## Quick Reference

### Commands

| Command | Purpose |
|---------|---------|
| `/sdlc` | Start SDLC workflow (REQUIRED for features) |
| `/update-claudemd` | Update CLAUDE.md from changes |
| `/code-review` | Run code review |
| `/test-file <path>` | Generate tests for file |

### Build & Test

```bash
# Backend
cd backend && go test ./...
cd backend && go build ./cmd/api

# Frontend (437 tests)
cd frontend && npm test && npm run build
cd frontend && npm run lint

# Infrastructure
cd infrastructure/{global,shared,backend} && tofu plan
```

### Key Agents

| Agent | When to Spawn |
|-------|---------------|
| `test-engineer` | Phase 2 - ALWAYS FIRST |
| `implementation-agent` | Phase 3 - AFTER tests exist |
| `documentation-generator` | Phase 5 |
| `code-review` | Before PR |

### MCP Servers

| Server | Purpose |
|--------|---------|
| `spec-workflow` | Spec management & approvals |
| `context7` | Library documentation |
| `github` | GitHub API |
| `dynamodb-mcp-server` | DynamoDB design |
| `opentofu` | OpenTofu modules |

---

## Branching Strategy

```
main (protected)
  └── group-N/{name}           # Feature group
        ├── task-N.1/{desc}    # Task branches
        └── task-N.2/{desc}
              └── [PR to group → PR to main]
```

| Type | Pattern | Example |
|------|---------|---------|
| Group | `group-N/{name}` | `group-1/foundation-backend` |
| Task | `task-N.X/{desc}` | `task-1.1/domain-models` |
| Feature | `feature/{name}` | `feature/global-user-type` |
| Hotfix | `hotfix/{desc}` | `hotfix/lambda-timeout` |

---

## DynamoDB Schema (Single-Table)

| Entity | PK | SK | GSI1-PK | GSI1-SK |
|--------|----|----|---------|---------|
| User | `USER#{userId}` | `PROFILE` | | |
| Track | `USER#{userId}` | `TRACK#{trackId}` | | |
| Album | `USER#{userId}` | `ALBUM#{albumId}` | | |
| Playlist | `USER#{userId}` | `PLAYLIST#{playlistId}` | `PUBLIC_PLAYLIST` (if public) | `{createdAt}` |
| ArtistProfile | `ARTIST_PROFILE#{id}` | `PROFILE` | `USER#{userId}` | `ARTIST_PROFILE` |
| Follow | `FOLLOW#{followerId}` | `FOLLOWING#{followedId}` | `ARTIST_PROFILE#{followedId}` | `FOLLOWER#{followerId}` |
| Upload | `USER#{userId}` | `UPLOAD#{uploadId}` | | |
| Tag | `USER#{userId}` | `TAG#{tagName}` | | |

---

## API Endpoints

### Artist Profiles
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/artists/entity` | Artist+ | Create artist profile |
| GET | `/artists/entity` | Public | List artist profiles |
| GET | `/artists/entity/:id` | Public | Get artist profile |
| PUT | `/artists/entity/:id` | Owner | Update artist profile |
| GET | `/artists/entity/search` | Public | Search artists |

### Follow System
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/artists/entity/:id/follow` | Subscriber+ | Follow artist |
| DELETE | `/artists/entity/:id/follow` | Subscriber+ | Unfollow artist |
| GET | `/artists/entity/:id/followers` | Public | Get followers |
| GET | `/artists/entity/:id/following` | Public | Check if following |
| GET | `/users/me/following` | Subscriber+ | Get user's following list |

### Playlists
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/playlists/public` | Public | List public playlists |
| PUT | `/playlists/:id/visibility` | Owner | Update visibility |

### Roles (Admin only)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/users/:id/role` | Admin | Get user role |
| PUT | `/users/:id/role` | Admin | Update user role |

---

## Documentation Standards

### Subdirectory CLAUDE.md Requirements

Every major subdirectory needs a CLAUDE.md with:
1. Overview - Brief description
2. File Descriptions - Purpose of each file
3. Key Functions - Signatures and behavior
4. Dependencies - Internal and external
5. Usage Examples

**Exception**: `.claude/` uses root files: `AGENTS.md`, `COMMANDS.md`, `PLUGINS.md`, `SKILLS.md`

### CHANGELOG.md Format

```markdown
## [Unreleased]
### Added
### Changed
### Fixed
```

---

## Reference Docs

Located in `.claude/docs/`:
- `tdd-workflow.md` - TDD practices (includes wiring verification tests)
- `wiring-checklist.md` - **CRITICAL** checklist for new services/handlers/routes
- `wave-planning.md` - Parallel @claude execution
- `epic-completion-checklist.md` - Epic completion
- `task-granularity.md` - Task breakdown rules

Current specs: `.spec-workflow/specs/localstack-test-migration/`, `.spec-workflow/specs/global-user-type/`

**Important**: When adding new features, ALWAYS follow the wiring checklist to ensure:
- Services are added to Services struct
- Handlers are created and routes registered
- Frontend routes are added to router (code-based routing)
- Environment variables are configured in Lambda

---

## Recent Updates

### 2026-01-31: LocalStack Integration Test Migration (Epic 9)

**44 integration tests** covering the full backend stack against real AWS services via LocalStack:
- Repository tests: DynamoDB CRUD + GSI queries, S3 operations + presigned URLs
- Service tests: Track visibility, playlist discovery, tag normalization, follow system, Cognito auth
- API tests: Full HTTP endpoint testing with auth middleware, role-based access, admin routes
- CI: GitHub Actions workflow (`.github/workflows/integration.yml`) with LocalStack service container
- Test infrastructure: `testutil/server.go` (full Echo server), `testutil/http_helpers.go` (request helpers)
- Run with: `cd backend && go test -tags=integration ./internal/repository/ ./internal/service/ ./test/`

### 2025-01-27: Admin Access Control Bug Fixes

**Bug Fixes:**
- **Admin Track Listing**: Fixed DynamoDB Scan pagination - was only showing 5 of 9 tracks. Now scans in batches of 100 until enough filtered tracks are collected.
- **Admin Track Deletion**: Fixed 404 error when admin deletes other users' tracks. Now uses `GetTrackByID` for global lookup, then deletes with actual owner's ID.
- **Clean Track Deletion**: Added `DeleteByPrefix` to S3Repository for batch deletion of HLS transcoded files at `hls/{userId}/{trackId}/`.
- **UI Layout Overlap**: Fixed table overlapping player bar with `h-screen overflow-hidden` and `min-h-0` flex constraints.

**API Changes:**
- `DeleteTrack(ctx, userID, trackID, hasGlobal)` - Added `hasGlobal` parameter for admin access

### 2025-01-26: Global User Type Feature (PR #18)
Replaced subscription tier system with role-based access control:

**Backend Changes:**
- New models: `UserRole`, `Permission`, `ArtistProfile`, `Follow`, `PlaylistVisibility`
- New services: `RoleService`, `ArtistProfileService`, `FollowService`
- New handlers with authorization middleware (`RequireRole`, `RequireAuth`)
- Repository methods for artist profiles and follow relationships

**Frontend Changes:**
- Role-based feature gating (replaced `hasTier()` with `hasRole()`)
- New components: `ArtistProfileCard`, `FollowButton`, `VisibilitySelector`
- New routes: `/artists/entity`, `/artists/entity/$artistId`, `/playlists/public`
- 437 tests passing (39 new component tests)

**Infrastructure:**
- Cognito user groups: admin, artist, subscriber
- Bootstrap script: `scripts/bootstrap-admin.sh`
- Migration scripts in `scripts/migrations/`
