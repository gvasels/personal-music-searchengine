# Artist Profile Components - CLAUDE.md

## Overview

Components for displaying and editing artist profiles as part of the Global User Type feature. Enables artists to create profiles linked to the music catalog.

## Files

| File | Description |
|------|-------------|
| `ArtistProfileCard.tsx` | Card component displaying artist profile with avatar, stats, and follow button |
| `EditArtistProfileModal.tsx` | Modal form for creating or editing artist profiles |
| `index.ts` | Barrel export for artist profile components |
| `__tests__/ArtistProfileCard.test.tsx` | Unit tests for ArtistProfileCard |

## Components

### ArtistProfileCard

Displays an artist profile card with header image, avatar, bio, location, and follower/track counts.

```typescript
interface ArtistProfileCardProps {
  profile: ArtistProfile;
  showFollowButton?: boolean;  // Default: true
}

export function ArtistProfileCard({ profile, showFollowButton }: ArtistProfileCardProps): JSX.Element
```

**Features:**
- Header image with gradient fallback
- Avatar with initial letter fallback
- Verified badge indicator
- Bio with line clamping (2 lines)
- Location with map pin icon
- Follower and track count stats
- Integrated FollowButton component
- Click to navigate to artist profile page

### EditArtistProfileModal

Modal dialog for creating new or editing existing artist profiles.

```typescript
interface EditArtistProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
  profile?: ArtistProfile | null;  // If provided, edit mode; otherwise, create mode
  onSuccess?: (profile: ArtistProfile) => void;
}

export function EditArtistProfileModal(props: EditArtistProfileModalProps): JSX.Element | null
```

**Form Fields:**
- Display Name (required, max 100 chars)
- Bio (optional, max 500 chars with counter)
- Location (optional, max 100 chars)
- Website (optional, URL validation)

**Behavior:**
- Detects create vs edit mode from `profile` prop
- Pre-populates form in edit mode
- Validates display name is required
- Shows loading state during submission
- Calls `onSuccess` callback with saved profile
- Closes modal on successful save

## Dependencies

| Package | Usage |
|---------|-------|
| `@tanstack/react-router` | Navigation to artist pages |
| `../../hooks/useArtistProfiles` | Create/update mutations |
| `../../hooks/useFollow` | Follow functionality |
| `../../types` | ArtistProfile type |

## Types

```typescript
interface ArtistProfile {
  userId: string;
  displayName: string;
  bio?: string;
  avatarUrl?: string;
  headerImageUrl?: string;
  location?: string;
  website?: string;
  socialLinks?: Record<string, string>;
  isVerified?: boolean;
  followerCount: number;
  trackCount: number;
  linkedArtistId?: string;
  createdAt: string;
  updatedAt: string;
}
```

## Usage

```typescript
import { ArtistProfileCard, EditArtistProfileModal } from '@/components/artist-profile';

// Display an artist card
<ArtistProfileCard profile={artistProfile} />

// Edit/create modal
const [isOpen, setIsOpen] = useState(false);
const [editingProfile, setEditingProfile] = useState<ArtistProfile | null>(null);

<EditArtistProfileModal
  isOpen={isOpen}
  onClose={() => setIsOpen(false)}
  profile={editingProfile}
  onSuccess={(profile) => console.log('Saved:', profile)}
/>
```

## Related Components

- `FollowButton` - Used within ArtistProfileCard for follow functionality
- `FollowersList` / `FollowingList` - Display follow relationships
