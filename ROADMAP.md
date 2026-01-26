# Personal Music Search Engine - Roadmap

## Overview

This document tracks planned features, improvements, and technical debt for the Personal Music Search Engine project.

---

## Phase 2: Global User Type (Streaming Service Model) ✅ COMPLETE

Foundation for a multi-user streaming platform with content rights management.

**Completed: 2026-01-26** | [Full Spec](.spec-workflow/specs/global-user-type/)

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Global User Role | User roles (admin, artist, subscriber, guest) with Cognito groups | ✅ Complete |
| High | Public Playlists | Visibility options (private/unlisted/public) with discovery | ✅ Complete |
| High | Artist Profiles | Artist profile management with catalog linking | ✅ Complete |
| High | Follow System | Follow/unfollow artists with follower counts | ✅ Complete |
| High | Admin Panel | User management with search, role changes, status toggle | ✅ Complete |
| Medium | Content Rights | Track ownership and licensing per track/album | Planned |
| Medium | User Discovery | Browse other users' public libraries | Planned |

### Implementation Details

- **Backend**: Go services for roles, artist profiles, follows, admin operations
- **Frontend**: React components for visibility selector, follow button, admin panel
- **Infrastructure**: Cognito groups (admin, artist, subscriber, GlobalReaders)
- **Data**: DynamoDB entities for ArtistProfile, Follow with GSI access patterns

---

## AI Chatbot & Agent System

Intelligent assistants for artists, actors, and creators to engage with fans.

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Artist Agent Framework | Infrastructure for artist-specific AI agents | Planned |
| High | Agent Knowledge Base | Per-artist knowledge store (bio, discography, FAQ) | Planned |
| High | Tour Information | Agent can provide tour dates, venues, ticket links | Planned |
| High | New Releases | Agent announces and discusses new music/content | Planned |
| Medium | Merchandise Store | Agent can browse/recommend merch, link to store | Planned |
| Medium | Fan Q&A | Answer common fan questions from knowledge base | Planned |
| Medium | Event Reminders | Notify users of upcoming shows in their area | Planned |
| Medium | Pre-save/Pre-order | Agent facilitates pre-saves for upcoming releases | Planned |
| Low | Voice Interface | Voice-based interaction with artist agents | Planned |
| Low | Multi-language | Agent responds in user's preferred language | Planned |
| Low | Personality Customization | Artists can customize agent tone/personality | Planned |

### Agent Architecture (Planned)

```
┌─────────────────────────────────────────────────────────────┐
│                      User Interface                          │
│                 (Chat widget, voice, API)                    │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Agent Router                              │
│         (Routes to appropriate artist agent)                 │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Artist Agent                               │
│    ┌─────────────┬─────────────┬─────────────┐              │
│    │ Knowledge   │ Tool        │ Response    │              │
│    │ Retrieval   │ Execution   │ Generation  │              │
│    └─────────────┴─────────────┴─────────────┘              │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   External Integrations                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │Ticketmaster│ │ Shopify  │ │ Bandcamp │ │ Songkick │       │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
└─────────────────────────────────────────────────────────────┘
```

### Use Cases

1. **Fan asks about tour**: "When is [Artist] playing near Chicago?"
   - Agent queries tour database, returns dates/venues/ticket links

2. **New release inquiry**: "What's [Artist]'s latest album?"
   - Agent returns album info, streaming links, track list

3. **Merch discovery**: "Does [Artist] have any hoodies?"
   - Agent searches merch catalog, shows options with purchase links

4. **Artist background**: "How did [Artist] get started?"
   - Agent retrieves bio from knowledge base, provides conversational response

---

## Creator Studio

Professional tools for DJs, podcasters, and content creators with role-based feature add-ons.

### Creator Types & Tiers

| Creator Type | Base Features | Premium Add-ons |
|--------------|---------------|-----------------|
| **DJ** | Library management, BPM/key display, basic playlists | Mix recorder, beat matching, live streaming |
| **Podcaster** | Episode hosting, RSS feed, basic analytics | Transcription, chapter markers, ad insertion |
| **Music Producer** | Track hosting, collaboration invites | Stem separation, version control, licensing tools |
| **Video Creator** | Music licensing search, sync rights | Batch licensing, project folders, usage tracking |
| **Radio/Station** | Scheduled playlists, automation | DMCA compliance, royalty reporting, ad scheduling |

### DJ Studio Features

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Crate System | Organize tracks into DJ crates/folders | Planned |
| High | BPM/Key Match | Find compatible tracks by BPM range and harmonic key | Planned |
| High | Hot Cues | Set and recall cue points within tracks | Planned |
| Medium | Mix Recorder | Record DJ sets with tracklist generation | Planned |
| Medium | Beat Grid Editor | Manually adjust beat grids for better sync | Planned |
| Medium | Transition Suggestions | AI suggests next track based on energy/key/BPM | Planned |
| Medium | Live Streaming | Stream DJ sets to platforms (Twitch, YouTube) | Planned |
| Low | Deck Preview | Preview tracks in headphones before playing | Planned |
| Low | Effects Rack | Apply effects (echo, filter, reverb) | Planned |
| Low | MIDI Controller Support | Map physical DJ controllers | Planned |

### Podcaster Studio Features

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Episode Management | Upload, schedule, publish episodes | Planned |
| High | RSS Feed Generation | Auto-generate podcast RSS for distribution | Planned |
| High | Show Notes Editor | Rich text editor for episode descriptions | Planned |
| Medium | Auto-Transcription | AI transcription of episodes | Planned |
| Medium | Chapter Markers | Define chapters with timestamps | Planned |
| Medium | Analytics Dashboard | Downloads, listeners, geographic data | Planned |
| Low | Dynamic Ad Insertion | Insert/update ads in back catalog | Planned |
| Low | Guest Booking | Scheduling tools for remote guests | Planned |

### Music Producer Features

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Private Sharing | Share unreleased tracks with collaborators | Planned |
| High | Version History | Track revisions and compare versions | Planned |
| Medium | Stem Separation | AI-powered stem extraction (vocals, drums, etc.) | Planned |
| Medium | Collaboration Rooms | Real-time feedback and comments on tracks | Planned |
| Medium | License Generator | Create and attach licenses to tracks | Planned |
| Low | Sample Clearance | Tools to document and clear samples | Planned |

### Video Creator Features

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Sync License Search | Find tracks available for video sync | Planned |
| High | License Management | Track which songs used in which projects | Planned |
| Medium | Project Folders | Organize music by video project | Planned |
| Medium | Batch Licensing | License multiple tracks for a project | Planned |
| Medium | Usage Reporting | Generate reports for rights holders | Planned |
| Low | Direct Integration | Plugins for Premiere, Final Cut, DaVinci | Planned |

### Studio Architecture (Planned)

```
┌─────────────────────────────────────────────────────────────┐
│                     Creator Dashboard                        │
│              (Role-specific UI and features)                 │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Feature Flag System                        │
│         (Enable/disable features per subscription)           │
└─────────────────────────┬───────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│   DJ Module   │ │ Podcast Module│ │Producer Module│
│               │ │               │ │               │
│ • Crates      │ │ • Episodes    │ │ • Versions    │
│ • BPM Match   │ │ • RSS Feed    │ │ • Stems       │
│ • Mix Record  │ │ • Transcript  │ │ • Collab      │
└───────────────┘ └───────────────┘ └───────────────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    Core Music Library                        │
│            (Shared tracks, playlists, metadata)              │
└─────────────────────────────────────────────────────────────┘
```

### Subscription Model (Planned)

| Tier | Price | Includes |
|------|-------|----------|
| Free | $0/mo | Basic library, 5GB storage, standard streaming |
| Creator | $9.99/mo | 100GB storage, 1 creator module, basic analytics |
| Pro | $19.99/mo | 500GB storage, all modules, advanced analytics, priority support |
| Studio | $49.99/mo | Unlimited storage, white-label options, API access, dedicated support |

---

## Rights Management System

Comprehensive content rights management across geographic scopes - from local venues to global distribution.

### Geographic Scope Hierarchy

```
┌─────────────────────────────────────────────────────────────┐
│                         GLOBAL                               │
│                  (Worldwide distribution)                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                      NATIONAL                          │  │
│  │               (Country-specific rights)                │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │                   REGIONAL                       │  │  │
│  │  │            (State/Province/Territory)            │  │  │
│  │  │  ┌───────────────────────────────────────────┐  │  │  │
│  │  │  │                 LOCAL                      │  │  │  │
│  │  │  │         (City/Venue/Event-specific)        │  │  │  │
│  │  │  └───────────────────────────────────────────┘  │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Rights Types

| Right Type | Description | Scope Examples |
|------------|-------------|----------------|
| **Mechanical** | Reproduction/distribution of recordings | Per-country rates (Harry Fox, MCPS, JASRAC) |
| **Performance** | Public performance (streaming, radio, live) | PRO territories (ASCAP, BMI, PRS, GEMA) |
| **Sync** | Music paired with visual media | Per-project, often global |
| **Master** | Use of specific recording | Label/distributor territories |
| **Print** | Sheet music, lyrics display | Publisher territories |
| **Neighboring** | Rights for performers, producers | Varies by country (some don't recognize) |

### Core Rights Management Features

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Rights Database | Store rights info per track with geographic scope | Planned |
| High | Territory Restrictions | Block/allow content by country/region | Planned |
| High | License Tracking | Track active licenses, expiration dates, terms | Planned |
| High | Rights Holder Registry | Database of labels, publishers, PROs, distributors | Planned |
| Medium | Geo-Detection | Detect user location for rights enforcement | Planned |
| Medium | Rights Inheritance | Cascade rights from global → national → regional → local | Planned |
| Medium | Conflict Resolution | Handle overlapping or conflicting rights claims | Planned |
| Low | Rights Marketplace | Buy/sell/license rights within platform | Planned |

### Local Rights (Venue/Event Level)

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Venue Licensing | Track which venues have performance licenses | Planned |
| High | Event Permits | Temporary rights for specific events/dates | Planned |
| Medium | DJ Set Licensing | Bundle licenses for live DJ performances | Planned |
| Medium | Venue Reporting | Generate play reports for venue PRO compliance | Planned |
| Low | Proximity Licensing | Auto-detect venue and apply appropriate rights | Planned |

### Regional Rights (State/Province Level)

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Regional Blackouts | Restrict content in specific regions | Planned |
| High | State-Specific Rules | Handle varying state/province regulations | Planned |
| Medium | Regional PRO Mapping | Map regions to appropriate collection societies | Planned |
| Medium | Sub-National Licensing | Support regional licensing deals | Planned |
| Low | Regional Pricing | Different pricing tiers by region | Planned |

### National Rights (Country Level)

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Country Availability | Per-country content availability flags | Planned |
| High | National PRO Integration | Connect with ASCAP, BMI, PRS, GEMA, JASRAC, etc. | Planned |
| High | Royalty Rate Tables | Country-specific royalty calculations | Planned |
| High | DMCA/Copyright Compliance | Country-specific takedown procedures | Planned |
| Medium | Tax Withholding | Handle withholding requirements by country | Planned |
| Medium | Local Format Requirements | Support country-specific metadata standards | Planned |
| Medium | Censorship Compliance | Handle content restrictions (explicit, political) | Planned |
| Low | National Holidays | Adjust licensing for special broadcast days | Planned |

### Global Rights

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Multi-Territory Deals | Support rights spanning multiple countries | Planned |
| High | Worldwide Licensing | Flag content cleared for global distribution | Planned |
| High | Currency Handling | Multi-currency royalty calculations | Planned |
| Medium | International PRO Network | Reciprocal agreements between societies | Planned |
| Medium | Treaty Compliance | Berne Convention, WIPO, Rome Convention | Planned |
| Medium | Cross-Border Reporting | Consolidated reporting across territories | Planned |
| Low | Global Release Coordination | Timezone-aware release scheduling | Planned |

### Royalty & Payment System

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Play Tracking | Accurate per-stream counting by territory | Planned |
| High | Royalty Calculation | Apply correct rates per territory/right type | Planned |
| High | Split Management | Handle multiple rights holders per track | Planned |
| High | Payment Processing | Distribute payments to rights holders | Planned |
| Medium | Minimum Thresholds | Country-specific payout minimums | Planned |
| Medium | Advance Recoupment | Track and recoup advances against royalties | Planned |
| Medium | Audit Trail | Complete history for royalty audits | Planned |
| Low | Blockchain Verification | Immutable record of plays and payments | Planned |

### Compliance & Reporting

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Automated Reporting | Generate PRO reports per territory | Planned |
| High | Takedown Management | Process and track DMCA/copyright claims | Planned |
| High | License Expiry Alerts | Notify before licenses expire | Planned |
| Medium | Compliance Dashboard | Overview of rights status across territories | Planned |
| Medium | Audit Support | Export data for rights holder audits | Planned |
| Medium | Dispute Resolution | Track and manage rights disputes | Planned |
| Low | Regulatory Updates | Alert on changing regulations by country | Planned |

### Rights Data Model (Planned)

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│     Track       │────▶│   TrackRights   │◀────│  RightsHolder   │
│                 │     │                 │     │                 │
│ • id            │     │ • trackId       │     │ • id            │
│ • title         │     │ • holderId      │     │ • name          │
│ • artist        │     │ • rightType     │     │ • type (label,  │
└─────────────────┘     │ • sharePercent  │     │   publisher,    │
                        │ • territories[] │     │   PRO, artist)  │
                        │ • startDate     │     │ • territories[] │
                        │ • endDate       │     │ • paymentInfo   │
                        │ • restrictions  │     └─────────────────┘
                        └─────────────────┘
                                │
                                ▼
                        ┌─────────────────┐
                        │   Territory     │
                        │                 │
                        │ • code (ISO)    │
                        │ • scope (local, │
                        │   regional,     │
                        │   national,     │
                        │   global)       │
                        │ • parentCode    │
                        │ • proId         │
                        │ • royaltyRates  │
                        └─────────────────┘
```

### Integration Points

| System | Purpose | Priority |
|--------|---------|----------|
| DDEX | Music metadata & rights exchange standard | High |
| CIS-Net | International rights database | Medium |
| ISRC Registry | Recording identification | High |
| ISWC Database | Song/composition identification | High |
| PRO APIs | ASCAP, BMI, PRS, GEMA, etc. | High |
| MLC (US) | Mechanical licensing collective | High |
| SoundExchange | Digital performance royalties (US) | Medium |
| PPL/PRS (UK) | UK collection societies | Medium |

---

## Data Model Improvements

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Artist Entity | Create dedicated Artist table with UUIDs (names not unique) | Planned |
| High | Artist Linking | Link tracks/albums to artistId instead of artist name string | Planned |
| Medium | Multiple Artists | Support multiple artists per track (featuring, collaborations) | Planned |
| Medium | Artist Metadata | Bio, images, external links (Spotify, Discogs, etc.) | Planned |
| Low | Genre Taxonomy | Hierarchical genre system instead of free-form strings | Planned |

---

## Search & Discovery

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Playlist Search | Search includes playlists | Completed |
| Medium | Faceted Search | Filter by multiple criteria simultaneously | Planned |
| Medium | Similar Tracks | "More like this" recommendations based on audio features | Planned |
| Medium | Smart Playlists | Auto-generated playlists based on rules (BPM range, key, genre) | Planned |
| Low | Full-Text Lyrics | Search by lyrics (requires lyrics ingestion) | Planned |

---

## Audio Analysis

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | BPM Detection | Detect beats per minute | Completed |
| High | Key Detection | Detect musical key and mode | Completed |
| Medium | Waveform Display | Visual waveform in player | Planned |
| Medium | Beat Grid | Visual beat markers for DJ-style mixing | Planned |
| Low | Mood Classification | ML-based mood tagging (energetic, chill, etc.) | Planned |
| Low | Audio Fingerprinting | Identify tracks via acoustic fingerprint | Planned |

---

## Player & Playback

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Gapless Playback | Seamless transition between tracks | Planned |
| Medium | Crossfade | Configurable crossfade between tracks | Planned |
| Medium | Equalizer | Audio EQ controls | Planned |
| Medium | Playback Speed | Variable playback speed (0.5x - 2x) | Planned |
| Low | Chromecast Support | Cast audio to Chromecast devices | Planned |
| Low | AirPlay Support | Cast audio to AirPlay devices | Planned |

---

## UI/UX Improvements

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Mobile Responsive | Improve mobile experience | In Progress |
| Medium | Keyboard Shortcuts | Play/pause, next, previous, search focus | Planned |
| Medium | Drag & Drop Reorder | Reorder playlist tracks via drag and drop | Planned |
| Medium | Batch Operations | Select multiple tracks for tagging, playlist add, delete | Planned |
| Low | Themes | Additional color themes beyond light/dark | Planned |
| Low | Customizable Layout | Configurable sidebar, column visibility | Planned |

---

## Upload & Processing

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| Medium | Bulk Upload | Upload entire folders/albums at once | Planned |
| Medium | Duplicate Detection | Warn on duplicate uploads (by fingerprint or metadata) | Planned |
| Medium | Metadata Editor | Edit track metadata after upload | Partial |
| Low | Cover Art Search | Auto-fetch cover art from external sources | Planned |
| Low | Lyrics Fetch | Auto-fetch lyrics from external sources | Planned |

---

## Infrastructure & Performance

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | CDN Optimization | Optimize CloudFront caching strategies | Planned |
| Medium | Offline Support | PWA with offline playback capability | Planned |
| Medium | Background Sync | Queue uploads when offline, sync when online | Planned |
| Low | Multi-Region | Deploy to multiple AWS regions | Planned |

---

## API & Integrations

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| Medium | Public API | REST API for third-party integrations | Planned |
| Medium | Webhooks | Notify external services on events | Planned |
| Low | Last.fm Scrobbling | Scrobble plays to Last.fm | Planned |
| Low | Spotify Import | Import playlists from Spotify | Planned |
| Low | Discord Rich Presence | Show currently playing in Discord | Planned |

---

## Bugs & Issues

| Priority | Issue | Description | Status |
|----------|-------|-------------|--------|
| ~~High~~ | ~~Missing Thumbnails~~ | ~~Cover art not displaying in UI despite being embedded in MP3 files~~ | Fixed |

### Resolved Issues

| Issue | Root Cause | Fix | Date |
|-------|------------|-----|------|
| Missing Thumbnails | Frontend type mismatch: expected `artworkS3Key`, backend sends `coverArtUrl` | Updated `Track` type and `PlayerBar` component to use `coverArtUrl` | 2026-01-23 |

---

## Technical Debt

| Priority | Item | Description | Status |
|----------|------|-------------|--------|
| High | Test Coverage | Increase backend test coverage to 80%+ | In Progress |
| High | GitHub Actions Fix | Fix deployment workflow - currently failing on push to main | Planned |
| Medium | Error Handling | Standardize error responses across all endpoints | Planned |
| Medium | Logging | Structured logging with correlation IDs | Planned |
| Medium | Route Warnings | Fix TanStack Router export warnings | Planned |
| Low | Bundle Size | Code splitting to reduce initial JS bundle | Planned |

---

## Future Enhancements (Roadmap)

### Admin Panel & User Management

| Priority | Item | Description | Status |
|----------|------|-------------|--------|
| Medium | DynamoDB User Sync | Power user management from DynamoDB with Cognito sync | Planned |
| Medium | User Creation Trigger | Lambda trigger on Cognito signup to create DynamoDB user record | Planned |
| Medium | Bulk User Operations | Admin batch operations for user role updates | Planned |
| Low | User Activity Logs | Track admin actions on user accounts | Planned |

---

## AI/ML Infrastructure - Bedrock Access Gateway + Marengo

Unified AI infrastructure using [Bedrock Access Gateway](https://github.com/aws-samples/bedrock-access-gateway) with extended support for TwelveLabs Marengo video understanding models.

### Overview

The Bedrock Access Gateway provides OpenAI-compatible APIs for Amazon Bedrock, enabling:
- Standard OpenAI SDK/API patterns for all Bedrock models
- Streaming responses via SSE
- Chat completions, embeddings, function calling
- Cross-region model inference
- Prompt caching (up to 90% cost reduction)

**Extended Requirement**: Add support for [TwelveLabs Marengo](https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-twelvelabs.html) video embedding models, enabling multimodal video understanding for music videos, concert footage, DJ sets, and visual content.

### Marengo Integration Features

| Priority | Feature | Description | Status |
|----------|---------|-------------|--------|
| High | Gateway Deployment | Deploy Bedrock Access Gateway (Lambda or Fargate) | Planned |
| High | Marengo Embedding API | Extend gateway to support Marengo video embeddings | Planned |
| High | Video Upload Pipeline | Process uploaded videos through Marengo for embedding | Planned |
| High | Vector Storage | Store 1024-dim Marengo embeddings in vector DB | Planned |
| Medium | Text-to-Video Search | Search videos using natural language queries | Planned |
| Medium | Audio-to-Video Search | Find videos matching audio content | Planned |
| Medium | Image-to-Video Search | Find videos containing similar visual content | Planned |
| Medium | Video Similarity | Find similar videos based on embedding distance | Planned |
| Low | Music Video Analysis | Extract scenes, moods, visual themes from music videos | Planned |
| Low | Concert Footage Index | Index and search concert recordings | Planned |
| Low | DJ Set Visual Sync | Match audio tracks to video moments in DJ sets | Planned |

### Marengo Capabilities (v3.0)

| Capability | Description | Use Case |
|------------|-------------|----------|
| Multi-Vector Embeddings | Separate vectors for visual, audio, speech | Fine-grained search |
| 4-Hour Video Support | Process long-form content | Full concerts, DJ sets |
| Sports Intelligence | Player/jersey/action tracking | Performance analysis |
| Composed Queries | Image + text combined search | "Find videos like this but with drums" |
| Cross-Modal Search | Any-to-any (text, image, audio, video) | Universal content discovery |

### Architecture (Planned)

```
┌─────────────────────────────────────────────────────────────┐
│                   Client Applications                        │
│            (Frontend, Mobile, External APIs)                 │
└─────────────────────────┬───────────────────────────────────┘
                          │ OpenAI-compatible API
┌─────────────────────────▼───────────────────────────────────┐
│              Bedrock Access Gateway                          │
│         (API Gateway + Lambda / ALB + Fargate)              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              Extended Embedding Handler               │    │
│  │    • Standard embeddings (text)                      │    │
│  │    • Marengo video embeddings (extended)             │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│ Amazon Bedrock │ │ TwelveLabs    │ │ Vector Store  │
│ (Claude, Nova) │ │ Marengo 3.0   │ │ (OpenSearch/  │
│               │ │ (via Bedrock) │ │  Pinecone)    │
└───────────────┘ └───────────────┘ └───────────────┘
```

### Video Processing Pipeline

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  Video   │───▶│   S3     │───▶│ Marengo  │───▶│  Vector  │
│  Upload  │    │  Bucket  │    │ Embedding│    │   DB     │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
                     │                               │
                     ▼                               ▼
              ┌──────────┐                    ┌──────────┐
              │ Metadata │                    │  Search  │
              │ Extract  │                    │  Index   │
              └──────────┘                    └──────────┘
```

### Integration Points with Existing Features

| Feature | Integration |
|---------|-------------|
| Creator Studio | Video creators can search/analyze their content |
| Rights Management | Visual content identification for copyright |
| AI Agents | Agents can search and reference video content |
| Search & Discovery | Unified multimodal search across audio + video |
| DJ Studio | Match visuals to audio for sync licensing |

---

## Completed Features

| Feature | Description | Completed |
|---------|-------------|-----------|
| Core Upload Flow | Upload, process, index audio files | Yes |
| Search Engine | Full-text search with Nixiesearch | Yes |
| Playlist Management | Create, edit, delete playlists | Yes |
| Tag System | Add/remove tags on tracks | Yes |
| Audio Streaming | CloudFront signed URL streaming | Yes |
| BPM/Key Detection | Audio analysis for BPM and musical key | Yes |
| Playlist Search | Search returns matching playlists | Yes |
| Clickable Artists | Artist names link to artist page | Yes |

---

## Notes

- Priorities: High = next sprint, Medium = next quarter, Low = backlog
- Status: Planned, In Progress, Completed, Blocked
- Update this document as features are completed or priorities change
