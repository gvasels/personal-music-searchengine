# Tasks Document: Rights Management Foundation

## Task Overview
| Task | Description | Estimated Files |
|------|-------------|-----------------|
| 1 | Rights Models | 1 |
| 2 | Territory Service | 2 |
| 3 | Rights Repository | 1 |
| 4 | Rights Service | 1 |
| 5 | License Service | 2 |
| 6 | Geo-Detection | 2 |
| 7 | Rights Handlers | 1 |
| 8 | Content ID Integration | 2 |
| 9 | Rights Middleware | 1 |
| 10 | Tests | 4 |

**Prerequisites**:
- WS1 (Data Model Foundation) for Artist entity (rights holders link to artists)
- WS2 (Bedrock/Marengo) for content fingerprinting

---

- [ ] 1. Create Rights data models
  - File: backend/internal/models/rights.go
  - Define TrackRights struct with all fields
  - Define RightType enum (mechanical, performance, sync, master, print, neighboring)
  - Define Territory struct with scope hierarchy
  - Define RightsHolder struct linked to Artist (WS1)
  - Define License struct with terms
  - Purpose: Rights entity data models
  - _Leverage: backend/internal/models/artist.go (WS1)_
  - _Requirements: 1.1, 1.2, 1.3_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create rights models: TrackRights (trackId, holderId, rightType, sharePercent, territories, dates), Territory (code, scope, parent, proId, rates), RightsHolder (with ArtistID link to WS1), License (with terms and status) | Restrictions: DynamoDB single-table design, proper validation tags, share percent must sum to 100 | Success: All models defined, relationships clear, validation rules specified | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 2. Create Territory service with hierarchy
  - Files: backend/internal/repository/territory.go, backend/internal/service/territory.go
  - Seed territory data (ISO countries, major regions)
  - Resolve hierarchy (local → regional → national → global)
  - Get PRO for territory
  - Get royalty rates for territory/right type
  - Purpose: Geographic territory management
  - _Leverage: ISO 3166 country codes_
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create territory repository and service. Seed ISO 3166-1 (countries) and 3166-2 (regions) data. Implement ResolveHierarchy that walks from specific to global. Map territories to PROs (ASCAP, BMI for US, PRS for UK, etc.) | Restrictions: Efficient caching, handle LOC: custom codes for venues, pre-compute hierarchy | Success: Territory lookup fast, hierarchy resolution correct, PRO mapping accurate | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 3. Create Rights repository
  - File: backend/internal/repository/rights.go
  - CRUD for TrackRights (by trackId, by holderId)
  - CRUD for RightsHolder
  - Query rights by territory
  - Validate share percentages sum to 100
  - Purpose: Rights data access layer
  - _Leverage: Existing repository patterns_
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create RightsRepository with CRUD for TrackRights and RightsHolder. Query by trackId (PK=TRACK#{trackId}, SK=RIGHTS#*), by holderId (GSI). Validate share percentages in write operations | Restrictions: Atomic updates for rights changes, proper indexes for queries, validate sums | Success: CRUD operations work, queries efficient, validation enforced | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 4. Create Rights service
  - File: backend/internal/service/rights.go
  - GetTrackRights with resolved holder info
  - SetTrackRights with validation
  - CheckAccess for territory-based access control
  - Rights inheritance from hierarchy
  - Conflict detection and flagging
  - Purpose: Rights business logic
  - _Leverage: backend/internal/service/territory.go_
  - _Requirements: 1.1, 1.2, 1.3, 5.1, 5.2, 5.3, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create RightsService with business logic. CheckAccess checks rights for user's territory using inheritance. Detect conflicts (overlapping claims). Validate total share = 100%. Return AccessResult with allowed/reason | Restrictions: Cache access checks, handle missing rights gracefully, log conflicts | Success: Access checks accurate, inheritance works, conflicts detected | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 5. Create License service
  - Files: backend/internal/repository/license.go, backend/internal/service/license.go
  - License CRUD with validation
  - Expiration tracking and alerts
  - Auto-renewal logic
  - License status management
  - Purpose: License lifecycle management
  - _Leverage: Existing service patterns_
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create LicenseService with CRUD, expiration tracking (cron job for alerts), auto-renewal processing, status transitions (pending→active→expired). GetExpiringLicenses for notifications | Restrictions: Handle timezone-aware expiration, idempotent renewal, audit trail | Success: Licenses managed correctly, expirations tracked, auto-renew works | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 6. Create Geo-detection service
  - Files: backend/internal/service/geo.go, backend/internal/middleware/geo.go
  - MaxMind GeoIP integration
  - IP to territory resolution
  - VPN detection (optional)
  - CloudFront header support
  - Purpose: User location detection
  - _Leverage: MaxMind GeoLite2 database_
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create GeoService using MaxMind GeoLite2. Create middleware that extracts IP (CloudFront-Viewer-Country header or X-Forwarded-For), resolves territory, adds to context. Handle unknown IPs conservatively | Restrictions: Update GeoIP database monthly, handle IPv4/IPv6, cache lookups | Success: Location detected accurately, middleware sets territory, handles edge cases | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 7. Create Rights API handlers
  - File: backend/internal/handlers/rights.go
  - GET/PUT /api/v1/tracks/:id/rights
  - CRUD /api/v1/rights-holders
  - CRUD /api/v1/licenses
  - GET /api/v1/territories with hierarchy
  - Purpose: Rights management API
  - _Leverage: Existing handler patterns_
  - _Requirements: All API requirements_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create handlers for rights CRUD, rights holders CRUD, licenses CRUD, territories query. Validate user permissions (only owner can modify). Return proper error codes (403 for unauthorized, 451 for geo-blocked) | Restrictions: Proper authorization, validation, consistent error responses | Success: All endpoints work, authorization enforced, errors clear | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 8. Integrate Content ID with Marengo
  - Files: backend/internal/service/contentid.go, backend/internal/handlers/contentid.go
  - Use Marengo (WS2) for audio/video fingerprinting
  - Check uploads against rights database
  - Block/flag potential infringements
  - Dispute workflow (claim ownership, provide license)
  - Purpose: Copyright protection
  - _Leverage: backend/internal/clients/marengo.go (WS2)_
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create ContentIDService using Marengo embeddings (WS2). On upload, generate fingerprint and search for matches. If match found with >0.9 similarity, block and show rights info. Implement dispute flow | Restrictions: False positive handling (threshold tuning), async processing, fair use considerations | Success: Infringing content detected, disputes can be filed, false positives manageable | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 9. Create Rights-check middleware
  - File: backend/internal/middleware/rights.go
  - Middleware for content endpoints (stream, download)
  - Check rights for user's territory before serving
  - Return 451 for blocked content
  - Cache access decisions
  - Purpose: Enforce rights on content access
  - _Leverage: backend/internal/service/rights.go, backend/internal/middleware/geo.go_
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create RightsMiddleware that intercepts content requests (/stream/:id, /download/:id). Get territory from context (geo middleware), check access via RightsService. Return 451 with reason if blocked. Cache decisions | Restrictions: Efficient (cache), handle missing rights (deny), proper error response | Success: Blocked content returns 451, allowed content passes through, cached | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 10. Write tests for Rights Management
  - Files: backend/internal/service/rights_test.go, backend/internal/service/territory_test.go, backend/internal/service/license_test.go, backend/internal/middleware/rights_test.go
  - Test rights checking with various scenarios
  - Test territory hierarchy resolution
  - Test license expiration logic
  - Test geo-detection and middleware
  - Test content ID matching
  - Purpose: Ensure rights enforcement is reliable
  - _Leverage: Existing test patterns_
  - _Requirements: All_
  - _Prompt: Implement the task for spec rights-management-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer | Task: Write comprehensive tests for rights service (inheritance, conflicts), territory (hierarchy), licenses (expiration, auto-renew), geo-detection (IP resolution), content ID (fingerprint matching). Test middleware access control | Restrictions: Mock GeoIP, test edge cases (expired, conflicting, missing), compliance scenarios | Success: 80%+ coverage, all tests pass, edge cases covered | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_
