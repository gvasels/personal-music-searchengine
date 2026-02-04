# Requirements Document: Rights Management Foundation

## Introduction

This spec establishes the foundation for comprehensive content rights management across geographic scopes - from local venues to global distribution. It includes the rights database schema, territory model, license tracking, and rights holder registry. This depends on WS1 (Artist Entity) and WS2 (Bedrock/Marengo for visual content fingerprinting).

## Alignment with Product Vision

This directly supports:
- **Rights Management System** - Full Local/Regional/National/Global rights handling
- **Phase 2: Global User Type** - Content ownership and licensing per track
- **Creator Studio** - Sync licensing for video creators, DJ set licensing
- **Royalty & Payment System** - Foundation for royalty calculations

## Prerequisites

- **WS1: Data Model Foundation** - Artist Entity required for rights holder relationships
- **WS2: Bedrock/Marengo** - Video embeddings enable visual content identification for copyright

## Requirements

### Requirement 1: Rights Database Schema

**User Story:** As a platform operator, I want to store rights information per track with geographic scope, so that I can enforce content availability and calculate royalties correctly.

#### Acceptance Criteria

1. WHEN a track is created THEN the system SHALL create a default TrackRights record (owner = uploader)
2. WHEN rights are stored THEN the system SHALL use DynamoDB with `PK: TRACK#{trackId}`, `SK: RIGHTS#{rightType}#{holderId}`
3. WHEN storing rights THEN the schema SHALL include:
   - trackId, holderId, rightType (mechanical, performance, sync, master, print)
   - sharePercent (0-100, must sum to 100 per rightType)
   - territories[] (array of ISO country/region codes)
   - startDate, endDate (null = perpetual)
   - restrictions (JSON: {explicit: bool, territories_blocked: []})
4. WHEN multiple holders exist THEN the system SHALL enforce sharePercent sums to 100
5. WHEN querying rights THEN the system SHALL support: by track, by holder, by territory

### Requirement 2: Territory Model

**User Story:** As a rights manager, I want to define geographic scopes (local, regional, national, global), so that I can manage rights at appropriate granularity.

#### Acceptance Criteria

1. WHEN territories are defined THEN the system SHALL support hierarchy:
   - Global: "WORLD"
   - National: ISO 3166-1 alpha-2 (US, GB, DE, JP, etc.)
   - Regional: ISO 3166-2 (US-CA, US-NY, GB-ENG, etc.)
   - Local: Custom venue/event codes prefixed with "LOC:" (LOC:VENUE123)
2. WHEN storing territory THEN the schema SHALL include:
   - code, scope (global/national/regional/local), parentCode
   - proId (associated PRO: ASCAP, BMI, PRS, GEMA, etc.)
   - royaltyRates (JSON map of rightType → rate)
3. WHEN resolving territory THEN the system SHALL walk hierarchy (local → regional → national → global)
4. WHEN rights specify "US" THEN the system SHALL apply to all US regional/local territories
5. WHEN rights specify "LOC:VENUE123" THEN the system SHALL apply only to that venue

### Requirement 3: License Tracking

**User Story:** As a content owner, I want to track active licenses with expiration dates and terms, so that I know what content is licensed where and when to renew.

#### Acceptance Criteria

1. WHEN a license is created THEN the system SHALL store:
   - licenseId, trackId, licenseeId, rightType, territories[]
   - startDate, endDate, autoRenew (bool)
   - terms (JSON: usage limits, attribution requirements, exclusivity)
   - fee, currency, paymentStatus
2. WHEN viewing license THEN the system SHALL show status: active, expired, pending, terminated
3. WHEN license nears expiration (30 days) THEN the system SHALL send notification
4. WHEN querying licenses THEN the system SHALL filter by: track, licensee, status, territory, date range
5. WHEN a license expires THEN the system SHALL update track availability automatically
6. WHEN auto-renew is enabled THEN the system SHALL process renewal 7 days before expiration

### Requirement 4: Rights Holder Registry

**User Story:** As a platform, I want a database of rights holders (labels, publishers, PROs, artists), so that I can attribute royalties and manage relationships.

#### Acceptance Criteria

1. WHEN storing rights holder THEN the system SHALL include:
   - holderId, name, type (label, publisher, PRO, distributor, artist)
   - territories[] (where they operate)
   - paymentInfo (encrypted: bank, paypal, etc.)
   - ipiNumber (for publishers), isni (for artists)
2. WHEN holder is a PRO THEN the system SHALL link to territory-specific PRO mapping
3. WHEN holder is an Artist THEN the system SHALL link to Artist entity (WS1)
4. WHEN searching holders THEN the system SHALL support: name, type, territory, IPI/ISNI
5. WHEN creating tracks THEN the system SHALL autocomplete holder selection
6. WHEN a holder is inactive THEN the system SHALL retain for historical records

### Requirement 5: Geo-Detection and Rights Enforcement

**User Story:** As a user, I want to only see content I'm allowed to access in my location, so that I don't encounter geo-blocked playback errors.

#### Acceptance Criteria

1. WHEN a user requests content THEN the system SHALL detect their location via IP geolocation
2. WHEN location is detected THEN the system SHALL check rights for that territory
3. IF content is blocked in territory THEN the system SHALL return 451 (Unavailable for Legal Reasons)
4. WHEN listing content THEN the system SHALL filter out unavailable items (not show then block)
5. WHEN using VPN detected THEN the system SHALL fall back to account country
6. WHEN location cannot be determined THEN the system SHALL use most restrictive interpretation
7. WHEN explicit content filtering applies THEN the system SHALL check territory regulations

### Requirement 6: Rights Inheritance and Conflict Resolution

**User Story:** As a rights manager, I want rights to cascade from global to local scopes with conflict resolution, so that I don't need to specify rights at every level.

#### Acceptance Criteria

1. WHEN rights are not specified at local level THEN the system SHALL inherit from regional
2. WHEN rights are not specified at regional level THEN the system SHALL inherit from national
3. WHEN rights are not specified at national level THEN the system SHALL inherit from global
4. WHEN conflict exists (different holders claim same right/territory) THEN the system SHALL:
   - Flag conflict for manual resolution
   - Apply most restrictive interpretation (block) until resolved
   - Notify affected parties
5. WHEN resolving conflicts THEN admin SHALL adjudicate with documentation
6. WHEN inheritance is overridden THEN the system SHALL store explicit override with reason

### Requirement 7: Integration with Content Identification

**User Story:** As a platform, I want to identify uploaded content that may infringe existing rights, so that I can prevent unauthorized distribution.

#### Acceptance Criteria

1. WHEN audio is uploaded THEN the system SHALL generate audio fingerprint (existing BPM/key detection enhanced)
2. WHEN video is uploaded THEN the system SHALL use Marengo embeddings (WS2) for visual fingerprinting
3. WHEN fingerprint matches known content THEN the system SHALL:
   - Block upload with "Potential copyright match detected"
   - Link to original content and rights holder
   - Offer "Claim this is mine" or "I have license" dispute flow
4. WHEN dispute is filed THEN the system SHALL hold content pending review
5. WHEN content is confirmed infringing THEN the system SHALL:
   - Remove content
   - Log strike against uploader
   - Notify rights holder
6. WHEN 3 strikes occur THEN the system SHALL suspend account pending review

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: RightsRepository, TerritoryService, LicenseService, GeoService, FingerprintService
- **Modular Design**: Rights checking should be middleware usable across all content endpoints
- **Dependency Management**: Rights system integrates with Track, Artist, and User but is independently testable
- **Clear Interfaces**: Define `TrackRights`, `Territory`, `License`, `RightsHolder` types

### Performance
- Rights check must complete in <50ms (cached per user+territory)
- Territory resolution must be precomputed (not walk hierarchy on each request)
- Fingerprint matching must complete in <5 seconds per upload
- Geo-detection must complete in <100ms

### Security
- Payment info must be encrypted at rest (AWS KMS)
- Rights modifications must be audit logged
- Geo-detection must not be spoofable via client headers
- Admin actions must require MFA

### Reliability
- Rights cache must invalidate on any rights change
- License expiration jobs must run with at-least-once semantics
- Fingerprint service unavailability must not block uploads (queue for later)

### Compliance
- DMCA takedown requests must be processed within 24 hours
- Rights holder data must be GDPR-compliant (deletable on request)
- Audit trail must be retained for 7 years (legal requirement)
- Cross-border data transfer must comply with local regulations
