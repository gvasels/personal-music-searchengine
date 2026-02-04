# Requirements Document: Creator Studio Foundation

## Introduction

This spec establishes the foundation for the Creator Studio - a modular platform for DJs, podcasters, music producers, and video creators. It includes the feature flag system, subscription tier management, creator dashboard, and the first creator module (DJ Studio). This depends on the Data Model Foundation (WS1) for proper Artist entity linking.

## Alignment with Product Vision

This directly supports:
- **Creator Studio** - Full DJ, Podcaster, Producer, Video Creator toolsets
- **Phase 2: Global User Type** - Creator roles and subscription tiers
- **Subscription Model** - Free, Creator, Pro, Studio tiers

## Prerequisites

- **WS1: Data Model Foundation** - Artist Entity must exist for creator profiles
- Artist linking required for track attribution in DJ sets and productions

## Requirements

### Requirement 1: Feature Flag System

**User Story:** As a platform operator, I want to enable/disable features per subscription tier, so that I can monetize premium features and A/B test new functionality.

#### Acceptance Criteria

1. WHEN a user requests a feature THEN the system SHALL check if feature is enabled for their tier
2. WHEN feature flags are defined THEN the system SHALL support flags at: global, tier, and user levels
3. WHEN flags are checked THEN the system SHALL apply hierarchy: user > tier > global (most specific wins)
4. WHEN a flag is toggled THEN the system SHALL take effect within 1 minute (cached with TTL)
5. WHEN defining flags THEN the format SHALL be: `{featureKey: string, enabled: boolean, tiers: string[], users: string[]}`
6. WHEN a feature is disabled THEN the UI SHALL hide the feature, API SHALL return 403 Forbidden
7. WHEN new flags are added THEN the system SHALL default to disabled (opt-in)

### Requirement 2: Subscription Tier Management

**User Story:** As a user, I want to subscribe to different tiers (Free, Creator, Pro, Studio), so that I can access features appropriate to my needs.

#### Acceptance Criteria

1. WHEN a user signs up THEN the system SHALL assign "Free" tier by default
2. WHEN tiers are defined THEN the system SHALL support: Free, Creator ($9.99/mo), Pro ($19.99/mo), Studio ($49.99/mo)
3. WHEN upgrading/downgrading THEN the system SHALL prorate billing appropriately
4. WHEN tier changes THEN the system SHALL update feature flags immediately
5. WHEN a tier includes storage THEN the system SHALL enforce: Free=5GB, Creator=100GB, Pro=500GB, Studio=Unlimited
6. WHEN storage limit is exceeded THEN the system SHALL block new uploads with clear error message
7. WHEN subscription lapses THEN the system SHALL downgrade to Free (data retained for 30 days)

### Requirement 3: Creator Dashboard

**User Story:** As a creator, I want a unified dashboard showing my content, analytics, and tools, so that I can manage my creator business efficiently.

#### Acceptance Criteria

1. WHEN a creator navigates to /studio THEN the system SHALL show personalized dashboard
2. WHEN dashboard loads THEN the system SHALL display:
   - Active modules (DJ, Podcast, Producer, Video)
   - Recent content (tracks, mixes, episodes)
   - Quick stats (plays, followers, storage used)
   - Pending tasks (uploads processing, scheduled releases)
3. WHEN a module is not enabled THEN the system SHALL show upgrade prompt with feature preview
4. WHEN viewing stats THEN the system SHALL show last 7 days by default (with date range selector)
5. WHEN on mobile THEN the dashboard SHALL show simplified card layout
6. WHEN creator has no content THEN the system SHALL show onboarding wizard

### Requirement 4: Creator Profile (Artist Integration)

**User Story:** As a creator, I want my profile linked to my Artist entity, so that my tracks, mixes, and productions are properly attributed.

#### Acceptance Criteria

1. WHEN a user enables creator mode THEN the system SHALL create/link Artist entity
2. WHEN creator uploads content THEN the system SHALL automatically attribute to their Artist entity
3. WHEN viewing creator profile THEN the system SHALL show same info as Artist profile plus:
   - Creator tier badge
   - Enabled modules
   - Public contact/booking info (if public)
4. WHEN creator has multiple personas THEN the system SHALL support multiple linked Artist entities
5. WHEN switching personas THEN the system SHALL allow selecting which Artist to attribute content to

### Requirement 5: DJ Studio - Crate System

**User Story:** As a DJ, I want to organize tracks into crates (folders), so that I can quickly access tracks organized by set type, venue, or mood.

#### Acceptance Criteria

1. WHEN viewing DJ Studio THEN the system SHALL show crate sidebar
2. WHEN creating a crate THEN the system SHALL accept: name, color, icon (optional)
3. WHEN adding tracks to crate THEN the system SHALL support drag-drop or "Add to crate" menu
4. WHEN a track is in multiple crates THEN the system SHALL track all crate memberships
5. WHEN viewing crate THEN the system SHALL show tracks with BPM, Key, Duration columns
6. WHEN sorting crate THEN the system SHALL support: name, BPM, key, date added, custom order
7. WHEN searching in DJ mode THEN the system SHALL include crate filter

### Requirement 6: DJ Studio - BPM/Key Matching

**User Story:** As a DJ, I want to find tracks compatible with my current track's BPM and key, so that I can mix harmonically and rhythmically.

#### Acceptance Criteria

1. WHEN a track is playing THEN the system SHALL show "Find Compatible" button
2. WHEN finding compatible tracks THEN the system SHALL filter by:
   - BPM range: ±3% of current (configurable)
   - Key: same key or harmonic neighbors (Camelot wheel)
3. WHEN displaying compatible tracks THEN the system SHALL show energy level indicator
4. WHEN key matching THEN the system SHALL support Camelot notation (e.g., 8A, 8B) alongside standard keys
5. WHEN no current track THEN the system SHALL allow manual BPM/Key input for matching
6. WHEN results are shown THEN the system SHALL sort by compatibility score (combined BPM + key match)

### Requirement 7: DJ Studio - Hot Cues

**User Story:** As a DJ, I want to set and recall cue points within tracks, so that I can jump to specific moments (drops, breakdowns) during performance.

#### Acceptance Criteria

1. WHEN viewing track in DJ mode THEN the system SHALL show 8 hot cue slots
2. WHEN setting a cue THEN the user SHALL click/tap cue slot at current position
3. WHEN a cue is set THEN the system SHALL store: timestamp, color, label (optional)
4. WHEN recalling a cue THEN playback SHALL jump to that position instantly
5. WHEN cues exist THEN the system SHALL show markers on the waveform
6. WHEN deleting a cue THEN the user SHALL long-press/right-click the cue slot
7. WHEN cues are saved THEN the system SHALL persist per-track to DynamoDB

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: FeatureFlagService, SubscriptionService, CrateService, HotCueService
- **Modular Design**: Each creator module should be loadable independently
- **Dependency Management**: Creator modules depend on core library, not each other
- **Clear Interfaces**: Define `CreatorModule`, `FeatureFlag`, `Crate`, `HotCue` types

### Performance
- Feature flag checks must complete in <10ms (cached)
- Dashboard must load in <1 second
- Crate switching must be instant (preloaded)
- Hot cue recall must happen in <50ms

### Security
- Subscription tier must be verified server-side (not client-controlled)
- Creator profile edits must verify ownership
- Feature flag bypass attempts must be logged and blocked

### Reliability
- Subscription lapse must not lose data immediately (grace period)
- Crate/hot cue data must be backed up with tracks
- Feature flag cache must invalidate on tier changes

### Usability
- Clear upgrade prompts when hitting tier limits
- Feature previews for locked features
- DJ-specific keyboard shortcuts (cue triggers)
- Touch-friendly cue buttons for mobile DJ use
