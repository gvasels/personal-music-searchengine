// User roles for role-based access control
export type UserRole = 'guest' | 'subscriber' | 'artist' | 'admin';

// Permissions that can be checked
export type Permission =
  | 'browse'
  | 'listen'
  | 'create_playlist'
  | 'edit_playlist'
  | 'delete_playlist'
  | 'upload_tracks'
  | 'edit_tracks'
  | 'delete_tracks'
  | 'publish_tracks'
  | 'manage_users'
  | 'manage_content';

// Role permission mappings
export const ROLE_PERMISSIONS: Record<UserRole, Permission[]> = {
  guest: ['browse'],
  subscriber: ['browse', 'listen', 'create_playlist', 'edit_playlist', 'delete_playlist'],
  artist: [
    'browse',
    'listen',
    'create_playlist',
    'edit_playlist',
    'delete_playlist',
    'upload_tracks',
    'edit_tracks',
    'delete_tracks',
    'publish_tracks',
  ],
  admin: [
    'browse',
    'listen',
    'create_playlist',
    'edit_playlist',
    'delete_playlist',
    'upload_tracks',
    'edit_tracks',
    'delete_tracks',
    'publish_tracks',
    'manage_users',
    'manage_content',
  ],
};

// Check if a role has a permission
export function hasPermission(role: UserRole, permission: Permission): boolean {
  return ROLE_PERMISSIONS[role]?.includes(permission) ?? false;
}

// Playlist visibility levels
export type PlaylistVisibility = 'private' | 'unlisted' | 'public';

// Track visibility levels (same as playlist visibility)
export type TrackVisibility = 'private' | 'unlisted' | 'public';

// Artist role for track contributions
export type ArtistRole = 'main' | 'featuring' | 'remixer' | 'producer';

// Artist contribution on a track (multi-artist support)
export interface ArtistContribution {
  artistId: string;
  artistName?: string;
  role: ArtistRole;
}

export interface Track {
  id: string;
  title: string;
  artist: string;
  artistId?: string;                    // Reference to Artist entity
  artists?: ArtistContribution[];       // Multi-artist support
  album: string;
  albumId?: string;
  duration: number;
  trackNumber?: number;
  year?: number;
  genre?: string;
  format: string;
  fileSize: number;
  s3Key: string;
  coverArtUrl?: string;
  tags: string[];
  // Audio analysis fields
  bpm?: number;           // Beats per minute (20-300)
  bpmConfidence?: number; // Detection confidence (0-1)
  musicalKey?: string;    // e.g., "Am", "C", "F#m"
  keyMode?: string;       // "major" or "minor"
  keyCamelot?: string;    // Camelot notation, e.g., "8A", "11B"
  keyConfidence?: number; // Detection confidence (0-1)
  energy?: number;        // Energy level (0-1)
  loudness?: number;      // Loudness in LUFS
  // GenAI analysis fields (Epic 2.3)
  subGenre?: string;        // Sub-genre classification
  mood?: string;            // One-word mood
  toneDescription?: string; // 2-3 sentence description
  sections?: Section[];     // Section analysis
  instrumentation?: string; // Detected instruments
  vocalPresence?: string;   // none/male/female/mixed
  energyProfile?: string;   // Energy arc description
  embeddingId?: string;      // S3 Vectors embedding ID
  embeddingStatus?: string;  // PENDING/GENERATING/COMPLETED/FAILED
  analysisStatus?: string;   // PENDING/ANALYZING/COMPLETED/FAILED
  analyzedAt?: string;      // Analysis timestamp
  // Visibility fields
  visibility?: TrackVisibility;         // Track visibility level (default: private)
  ownerDisplayName?: string;            // Display name of track owner (for admin/global users)
  publishedAt?: string;                 // When track was made public
  createdAt: string;
  updatedAt: string;
}

// Section represents a detected section in audio analysis
export interface Section {
  name: string;        // intro, buildup, drop, breakdown, outro
  startSec: number;    // Start time in seconds
  endSec: number;      // End time in seconds
  description: string; // Section description
}

export interface Album {
  id: string;
  name: string;
  artist: string;
  year?: number;
  trackCount: number;
  coverArt?: string;
  createdAt: string;
}

// Artist entity (full entity model)
export interface Artist {
  id: string;
  name: string;
  sortName?: string;
  bio?: string;
  imageUrl?: string;
  externalLinks?: Record<string, string>;
  isActive: boolean;
  trackCount: number;
  albumCount: number;
  totalPlays?: number;
  createdAt: string;
  updatedAt: string;
}

// Lightweight artist summary (for lists from album-based aggregation)
export interface ArtistSummary {
  name: string;
  trackCount: number;
  albumCount: number;
  coverArtUrl?: string;
}

// Request types for artist CRUD
export interface CreateArtistRequest {
  name: string;
  sortName?: string;
  bio?: string;
  imageUrl?: string;
  externalLinks?: Record<string, string>;
}

export interface UpdateArtistRequest {
  name?: string;
  sortName?: string;
  bio?: string;
  imageUrl?: string;
  externalLinks?: Record<string, string>;
}

export interface Playlist {
  id: string;
  name: string;
  description?: string;
  trackIds: string[];
  trackCount: number;
  coverArt?: string;
  visibility: PlaylistVisibility;
  ownerUserId?: string;
  ownerDisplayName?: string;
  createdAt: string;
  updatedAt: string;
}

// Backend returns this structure for single playlist with tracks
export interface PlaylistWithTracks {
  playlist: Playlist;
  tracks: Track[];
}

export interface Tag {
  name: string;
  trackCount: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
}

// Library statistics
export type StatsScope = 'own' | 'public' | 'all';

export interface LibraryStats {
  totalTracks: number;
  totalAlbums: number;
  totalArtists: number;
  totalDuration: number; // in seconds
}

export type RepeatMode = 'off' | 'all' | 'one';
export type Theme = 'light' | 'dark';

// Subscription types
export type SubscriptionTier = 'free' | 'creator' | 'pro';
export type SubscriptionStatus = 'active' | 'canceled' | 'past_due' | 'trialing';
export type SubscriptionInterval = 'monthly' | 'yearly';

// Feature flags
export type FeatureKey =
  | 'DJ_MODULE'
  | 'CRATES'
  | 'HOT_CUES'
  | 'BPM_MATCHING'
  | 'KEY_MATCHING'
  | 'MIX_RECORDING'
  | 'BULK_EDIT'
  | 'ADVANCED_STATS'
  | 'API_ACCESS'
  | 'UNLIMITED_STORAGE'
  | 'HQ_STREAMING';

export interface UserFeaturesResponse {
  tier: SubscriptionTier;
  role: UserRole; // Current role from database (real-time, not from JWT)
  features: Record<FeatureKey, boolean>;
}

export interface TierConfig {
  tier: SubscriptionTier;
  name: string;
  description: string;
  monthlyPriceCents: number;
  yearlyPriceCents: number;
  storageLimitBytes: number; // -1 for unlimited
  features: FeatureKey[];
}

export interface SubscriptionResponse {
  userId: string;
  tier: SubscriptionTier;
  tierName: string;
  status: SubscriptionStatus;
  interval: SubscriptionInterval;
  currentPeriodStart: string;
  currentPeriodEnd: string;
  cancelAtPeriodEnd: boolean;
  trialEnd?: string;
  storageLimit: number;
  storageUsed: number;
  features: FeatureKey[];
}

export interface StorageUsageResponse {
  storageUsedBytes: number;
  storageLimitBytes: number;
  usagePercent: number;
}

// Hot cues
export type HotCueColor = '#FF0000' | '#FF8C00' | '#FFFF00' | '#00FF00' | '#00FFFF' | '#0000FF' | '#800080' | '#FF69B4';

export interface HotCue {
  slot: number;
  position: number;
  label?: string;
  color: HotCueColor;
  createdAt: string;
  updatedAt: string;
}

export interface TrackHotCuesResponse {
  trackId: string;
  hotCues: HotCue[];
  maxSlots: number;
}

// Crates
export type CrateSortOrder = 'custom' | 'bpm' | 'key' | 'artist' | 'title' | 'added';

export interface Crate {
  id: string;
  name: string;
  description?: string;
  color?: string;
  trackCount: number;
  sortOrder: CrateSortOrder;
  isSmartCrate: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CrateWithTracks extends Crate {
  tracks: Track[];
}

// BPM/Key Matching
export interface MatchResult {
  track: Track;
  bpmCompatibility: number;
  keyCompatibility: number;
  overallScore: number;
  bpmDiff: number;
  keyRelation: string;
}

// Artist profile for artists with accounts
export interface ArtistProfile {
  userId: string;
  displayName: string;
  bio?: string;
  avatarUrl?: string;
  headerImageUrl?: string;
  location?: string;
  website?: string;
  socialLinks?: Record<string, string>;
  isVerified: boolean;
  followerCount: number;
  followingCount: number;
  trackCount: number;
  createdAt: string;
  updatedAt: string;
}

// Follow relationship
export interface Follow {
  followerId: string;
  followedId: string;
  createdAt: string;
}

// User profile with role
export interface UserProfile {
  userId: string;
  email: string;
  displayName?: string;
  role: UserRole;
  avatarUrl?: string;
  createdAt: string;
  updatedAt: string;
}

// Admin Track Reprocess Result
export interface ReprocessResult {
  trackId: string;
  status: 'processing' | 'complete' | 'failed';
  bpm?: number;
  bpmConfidence?: number;
  musicalKey?: string;
  keyCamelot?: string;
  embeddingStatus?: string;
  error?: string;
  processedAt?: string;
}

// Artist Events & Shows types
export type EventStatus = 'scheduled' | 'cancelled' | 'postponed';

export interface ArtistEvent {
  id: string;
  artistName: string;
  title: string;
  date: string;
  venue: string;
  city: string;
  region: string;
  country: string;
  ticketUrl?: string;
  status: EventStatus;
  source: string;
}

export interface ArtistSearchResult {
  name: string;
  imageUrl?: string;
  upcomingEvents: number;
  source: string;
}

export interface WatchedArtist {
  artistName: string;
  watchedAt: string;
}

export interface ArtistEventsResponse {
  artistName: string;
  events: ArtistEvent[];
  totalCount: number;
  source: string;
}

export interface WatchResponse {
  artistName: string;
  watching: boolean;
  watchedAt: string;
}

export interface WatchStatusResponse {
  watching: boolean;
  artistName: string;
}

export interface WatchedArtistsResponse {
  items: WatchedArtist[];
  nextCursor?: string;
  hasMore: boolean;
}

export interface ArtistSearchResponse {
  items: ArtistSearchResult[];
}
