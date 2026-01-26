# Requirements: Global User Type

## Introduction

This specification defines the foundation for transforming the Personal Music Search Engine from a single-user library application into a multi-user streaming platform with role-based access control. This enables the platform to support different user types (artists, subscribers, guests) with appropriate permissions and capabilities.

## Alignment with Product Vision

This feature implements **Phase 2: Global User Type (Streaming Service Model)** from the roadmap, establishing the core infrastructure needed for:
- Multi-user content sharing
- Artist profile management and linking to catalog
- Social features foundation (follows, public playlists)

---

## Requirements

### REQ-1: User Roles System

**User Story:** As a platform administrator, I want to assign different roles to users so that I can control what features and content each user type can access.

#### Acceptance Criteria

1. WHEN a new user registers THEN the system SHALL assign the default role `subscriber`
2. WHEN an admin promotes a user to `artist` role THEN the system SHALL enable artist-specific features
3. IF a user has `admin` role THEN the system SHALL grant access to all administrative functions
4. WHEN an unauthenticated user accesses the platform THEN the system SHALL treat them as `guest` with read-only access to public content
5. IF a user's role changes THEN the system SHALL force a token refresh to apply new permissions

#### Role Definitions

| Role | Description | Permissions |
|------|-------------|-------------|
| `admin` | Platform administrator | Full system access, user management, role assignment |
| `artist` | Verified content creator | Manage artist profile, link to catalog, view analytics |
| `subscriber` | Registered user | Stream content, create playlists, upload to personal library |
| `guest` | Unauthenticated visitor | View public content, browse public playlists, follow artists (read-only) |

#### Implementation Notes

- Roles managed via **Cognito Groups** (admin, artist, subscriber)
- Guest role is implicit (unauthenticated users)
- Initial admin: `gvasels90@gmail.com` - added via bootstrap script
- **Replaces existing SubscriptionTier system** - subscriptions to be rebuilt in future phase

---

### REQ-2: Public Playlists

**User Story:** As a subscriber, I want to make my playlists public so that other users can discover and enjoy my curated music collections.

#### Acceptance Criteria

1. WHEN creating a playlist THEN the system SHALL default visibility to `private`
2. WHEN a user sets a playlist to `public` THEN the system SHALL make it discoverable to all users
3. IF a playlist is `public` THEN `guest` users SHALL be able to view (but not modify) it
4. WHEN viewing a public playlist THEN the system SHALL display the creator's username and avatar
5. IF a playlist contains tracks that become unavailable THEN the system SHALL hide those tracks from public view
6. WHEN a user searches playlists THEN the system SHALL include public playlists in results

#### Visibility Options

| Visibility | Discoverable | Accessible By |
|------------|--------------|---------------|
| `private` | No | Owner only |
| `unlisted` | No | Anyone with link |
| `public` | Yes | All users including guests |

#### Migration

- Existing playlists with `IsPublic: true` → `Visibility: public`
- Existing playlists with `IsPublic: false` → `Visibility: private`

---

### REQ-3: Artist Profiles

**User Story:** As an artist, I want a profile page so that fans can find my music, bio, and official content in one place.

#### Acceptance Criteria

1. WHEN a user is granted `artist` role THEN the system SHALL allow them to create an artist profile
2. WHEN an artist profile exists THEN it SHALL be linked to existing Artist catalog entries (claim mechanism)
3. WHEN viewing an artist profile THEN the system SHALL show bio, discography, and social links
4. IF an artist has followers THEN the system SHALL display follower count on their profile
5. WHEN an artist uploads content THEN it SHALL automatically appear on their artist profile
6. IF searching for an artist THEN the system SHALL include artist profiles in results

#### Artist Profile Fields

| Field | Required | Description |
|-------|----------|-------------|
| `displayName` | Yes | Artist stage name |
| `bio` | No | Artist biography (max 2000 chars) |
| `avatarUrl` | No | Profile image |
| `headerImageUrl` | No | Profile banner image |
| `socialLinks` | No | Links to social media profiles |
| `linkedArtistId` | No | Link to existing Artist catalog entity |
| `followerCount` | Auto | Number of followers |

#### Artist Linking

Artists can "claim" existing Artist entities in the catalog:
- ArtistProfile.linkedArtistId → Artist.id
- When linked, profile shows all tracks/albums by that Artist
- Multiple ArtistProfiles cannot claim the same Artist entity

---

### REQ-4: Follow System

**User Story:** As a user, I want to follow artists so that I can stay updated on their activity.

#### Acceptance Criteria

1. WHEN a user follows an artist THEN the system SHALL add the artist to their followed list
2. IF a user follows an artist THEN the artist's follower count SHALL increase by one
3. IF a user unfollows an artist THEN the system SHALL remove them from the followed list
4. WHEN viewing a user's profile THEN the system SHALL display their followed artists count
5. IF a guest user follows an artist THEN they SHALL be prompted to sign in/register

#### Notes

- Both authenticated users AND guests can initiate follow (guests redirected to auth)
- Notifications for new content from followed artists → Future phase (roadmap)

---

## Non-Functional Requirements

### Code Architecture and Modularity

- **Single Responsibility**: Role checking and profile management in separate service modules
- **Modular Design**: Authorization middleware separate from business logic
- **Dependency Management**: Role/permission checks isolated to authorization layer
- **Clear Interfaces**: Define permission checking interfaces that can be mocked for testing

### Performance

- Role/permission checks SHALL complete within 10ms (claims from Cognito groups)
- Public playlist queries SHALL use GSI for efficient discovery
- Follower counts SHALL be denormalized for fast reads (eventual consistency acceptable)

### Security

- Role elevation SHALL require admin authorization via Cognito group management
- JWT tokens SHALL include Cognito group claims for role
- Token refresh forced on role change to apply new permissions immediately

### Reliability

- Role changes SHALL be reflected on next token refresh (forced)
- Public content visibility changes SHALL propagate within 1 minute
- Follower count updates are eventually consistent

### Usability

- Role SHALL be displayed in user profile
- Public/private toggle SHALL be obvious in playlist UI
- Artist profile link to catalog SHALL be clear

---

## Out of Scope (Roadmap Items)

The following features are explicitly NOT part of this specification:

- **Content Rights Management** - Dedicated spec needed for ownership, licensing, rights types
- **Artist Verification Flow** - Formal application/verification process
- **Notification System** - For follows, new releases, etc.
- **Likes/Reactions** - Generic interaction entity for likes, comments
- **User Discovery** - Browse other users' public libraries
- **Discovery Algorithms** - Recommendations based on follows, plays
- **Role Downgrade Handling** - What happens to content when artist becomes subscriber
- **Denormalized Data Sync Strategy** - Detailed eventual consistency approach
- **Rate Limiting** - For public content endpoints

---

## Dependencies

- Existing User model and Cognito authentication
- Existing Playlist model (will be extended with Visibility)
- Existing Artist model (will be linked from ArtistProfile)
- DynamoDB single-table design patterns

---

## Glossary

| Term | Definition |
|------|------------|
| Role | A user's platform-wide permission level (admin, artist, subscriber, guest) |
| Visibility | Who can discover and access a piece of content |
| Artist Profile | A user-facing profile for artists, linkable to catalog Artist entities |
| Follower | A user who has subscribed to updates from an artist profile |
| Artist (catalog) | Existing metadata entity representing a music artist in the library |
