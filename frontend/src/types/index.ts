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
  bpm?: number;           // Beats per minute (20-300)
  musicalKey?: string;    // e.g., "Am", "C", "F#m"
  keyMode?: string;       // "major" or "minor"
  keyCamelot?: string;    // Camelot notation, e.g., "8A", "11B"
  createdAt: string;
  updatedAt: string;
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
