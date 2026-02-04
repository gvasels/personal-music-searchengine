# Tasks Document: Creator Studio Foundation

## Task Overview
| Task | Description | Estimated Files |
|------|-------------|-----------------|
| 1 | Feature Flag Backend | 3 |
| 2 | Subscription Backend | 3 |
| 3 | Feature Flag Frontend | 2 |
| 4 | Creator Dashboard | 3 |
| 5 | Crate Backend | 3 |
| 6 | Crate Frontend | 3 |
| 7 | BPM/Key Matching | 2 |
| 8 | Hot Cues Backend | 2 |
| 9 | Hot Cues Frontend | 2 |
| 10 | Tests | 5 |

**Prerequisites**: WS1 (Data Model Foundation) must complete Task 1-4 first for Artist entity.

---

- [ ] 1. Create Feature Flag backend service
  - Files: backend/internal/models/feature.go, backend/internal/repository/feature.go, backend/internal/service/feature.go
  - Define FeatureFlag model with tier-based enablement
  - Repository for DynamoDB CRUD
  - Service with IsEnabled(userID, featureKey) logic
  - Tier hierarchy check (global → tier → user override)
  - Purpose: Feature access control backend
  - _Leverage: Existing repository patterns_
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create FeatureFlag model, repository (DynamoDB: PK=FEATURE, SK=FLAG#{key}), and service. Implement IsEnabled that checks global → tier → user hierarchy. Cache flags with 1-min TTL. Return enabled features list for user | Restrictions: Efficient queries, cache invalidation on tier change, default to disabled | Success: Feature checks work, tier hierarchy respected, caching efficient | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 2. Create Subscription backend service
  - Files: backend/internal/models/subscription.go, backend/internal/service/subscription.go, backend/internal/handlers/subscription.go
  - Define Tier enum and TierConfig
  - Subscription service with Stripe integration
  - Webhook handler for Stripe events
  - Storage limit enforcement
  - Purpose: Subscription management backend
  - _Leverage: Stripe Go SDK_
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with Stripe expertise | Task: Create subscription service with Tier enum (free/creator/pro/studio). Integrate Stripe SDK for checkout sessions and customer portal. Handle webhooks for subscription changes. Track storage usage and enforce limits | Restrictions: Secure webhook validation, handle Stripe errors, idempotent webhook processing | Success: Users can upgrade/downgrade, webhooks update tier, storage enforced | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 3. Create Feature Flag frontend hook
  - Files: frontend/src/hooks/useFeatureFlags.ts, frontend/src/lib/api/features.ts
  - Fetch enabled features on auth
  - Cache in React Query
  - isEnabled(featureKey) helper
  - Expose current tier
  - Purpose: Feature access control frontend
  - _Leverage: TanStack Query_
  - _Requirements: 1.1, 1.6_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/TypeScript Developer | Task: Create useFeatureFlags hook that fetches /api/v1/features on auth. Cache with React Query. Expose isEnabled(key), tier, isLoading. Refetch on subscription change | Restrictions: Handle loading state, efficient caching, type-safe | Success: Features accessible, shows loading state, updates on tier change | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 4. Create Creator Dashboard
  - Files: frontend/src/routes/studio/index.tsx, frontend/src/components/studio/DashboardStats.tsx, frontend/src/components/studio/ModuleCard.tsx
  - Main /studio route with dashboard layout
  - Show enabled modules, stats, recent content
  - Module cards with feature-gate (show upgrade for disabled)
  - Onboarding wizard for new creators
  - Purpose: Central creator hub
  - _Leverage: frontend/src/hooks/useFeatureFlags.ts_
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create /studio route with dashboard layout. Show stats (plays, followers, storage). Display module cards (DJ, Podcast, Producer). Gate disabled modules with upgrade CTA. Add onboarding for new users | Restrictions: Responsive design, clear upgrade paths, accessible | Success: Dashboard shows user's modules, stats accurate, upgrade flows work | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 5. Create Crate backend service
  - Files: backend/internal/models/crate.go, backend/internal/repository/crate.go, backend/internal/service/crate.go
  - Crate model with tracks list and custom sort
  - CRUD repository operations
  - Service with add/remove track methods
  - Feature gate for DJ module
  - Purpose: DJ crate management backend
  - _Leverage: Existing repository patterns_
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create Crate model (DynamoDB: PK=USER#{userId}, SK=CRATE#{crateId}). Implement CRUD + AddTrack/RemoveTrack. Support custom sort order. Feature-gate under DJ module | Restrictions: Limit 1000 tracks per crate, validate track ownership, efficient queries | Success: Crates CRUD works, tracks can be added/removed, limit enforced | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 6. Create Crate frontend components
  - Files: frontend/src/components/studio/dj/CrateSidebar.tsx, frontend/src/components/studio/dj/CrateView.tsx, frontend/src/routes/studio/dj/index.tsx
  - Sidebar showing user's crates
  - Crate view with track list (BPM, Key columns)
  - Drag to add tracks to crate
  - Create/edit/delete crate modals
  - Purpose: DJ crate management UI
  - _Leverage: DaisyUI components, @dnd-kit_
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create DJ Studio page at /studio/dj. Crate sidebar with list and create button. Crate view shows tracks with BPM/Key columns. Drag tracks from library to crate. Edit/delete crate modals | Restrictions: Feature-gated (show upgrade if not DJ tier), responsive, fast filtering | Success: Crate management works, drag-drop adds tracks, DJ-friendly UI | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 7. Implement BPM/Key matching
  - Files: backend/internal/service/matching.go, frontend/src/components/studio/dj/CompatibleTracks.tsx
  - Backend: Find tracks within BPM range and harmonic keys
  - Implement Camelot wheel neighbor logic
  - Frontend: "Find Compatible" button that shows matching tracks
  - Sort by compatibility score
  - Purpose: DJ track matching feature
  - _Leverage: Existing track query patterns_
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Full-stack Developer with music theory knowledge | Task: Create matching service that finds tracks by BPM range (±3%) and harmonic keys (Camelot neighbors). Score compatibility. Create CompatibleTracks component showing results sorted by score | Restrictions: Efficient queries (use indexes), support Camelot and standard key notation | Success: Find compatible returns relevant tracks, scored correctly, UI shows clearly | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 8. Create Hot Cues backend
  - Files: backend/internal/service/hotcue.go, backend/internal/handlers/hotcue.go
  - Hot cue storage in Track metadata (8 slots)
  - Set/get/delete hot cue endpoints
  - Validate slot (1-8) and timestamp
  - Purpose: Per-track cue point storage
  - _Leverage: Track model extension_
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Add hotCues []HotCue to Track model. Create endpoints PUT/DELETE /api/v1/tracks/:id/hotcues/:slot. Validate slot 1-8, timestamp within duration. Store color/label | Restrictions: Feature-gated to DJ module, validate ownership, 8 slots max | Success: Hot cues can be set/deleted, persist with track, validated | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 9. Create Hot Cues frontend component
  - Files: frontend/src/components/studio/dj/HotCueBar.tsx, frontend/src/hooks/useHotCues.ts
  - Display 8 hot cue buttons
  - Click to set cue at current position
  - Click existing cue to jump
  - Long-press/right-click to delete
  - Show on waveform
  - Purpose: DJ hot cue interaction UI
  - _Leverage: frontend/src/components/player/Waveform.tsx (from WS3)_
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create HotCueBar with 8 buttons. Click empty slot to set cue. Click filled to jump. Long-press/right-click to delete. Show cue markers on waveform. Color-coded buttons | Restrictions: Feature-gated, responsive, touch-friendly for mobile DJ use | Success: Hot cues work in real-time, markers on waveform, mobile-friendly | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 10. Write tests for Creator Studio
  - Files: backend/internal/service/feature_test.go, backend/internal/service/subscription_test.go, backend/internal/service/crate_test.go, frontend/src/hooks/__tests__/useFeatureFlags.test.ts, frontend/src/components/studio/__tests__/CrateView.test.tsx
  - Test feature flag evaluation
  - Test subscription tier changes
  - Test crate CRUD and track limits
  - Test frontend feature gating
  - Purpose: Ensure creator features are reliable
  - _Leverage: Existing test patterns_
  - _Requirements: All_
  - _Prompt: Implement the task for spec creator-studio-foundation, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Test Engineer | Task: Write tests for feature flags (tier hierarchy), subscription (Stripe mocks), crates (CRUD, limits), frontend feature gating (show/hide based on tier). Test BPM/key matching accuracy | Restrictions: Mock Stripe properly, test edge cases, integration tests for critical paths | Success: 80%+ coverage, all tests pass, critical paths covered | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_
