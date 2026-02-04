# Follow Components - CLAUDE.md

## Overview

Components for the follow/unfollow system in the Global User Type feature. Allows subscribers and artists to follow other artists and view follow relationships.

## Files

| File | Description |
|------|-------------|
| `FollowButton.tsx` | Button to follow/unfollow an artist |
| `FollowersList.tsx` | List of users following an artist |
| `FollowingList.tsx` | List of artists the current user follows |
| `index.ts` | Barrel export for follow components |
| `__tests__/FollowButton.test.tsx` | Unit tests for FollowButton |

## Components

### FollowButton

Toggle button for following/unfollowing an artist. Requires subscriber role or higher.

```typescript
interface FollowButtonProps {
  artistId: string;
  className?: string;
  size?: 'sm' | 'md' | 'lg';  // Default: 'md'
}

export function FollowButton({ artistId, className, size }: FollowButtonProps): JSX.Element | null
```

**Behavior:**
- Returns `null` if user is not authenticated or not at least subscriber role
- Shows loading spinner while fetching follow status
- Displays "Follow" (primary button) or "Following" (outline button)
- Shows spinner during follow/unfollow toggle
- Uses optimistic updates via `useFollowToggle` hook

### FollowersList

Displays a paginated list of users who follow a specific artist.

```typescript
interface FollowersListProps {
  artistId: string;
  limit?: number;  // Default: 20
}

export function FollowersList({ artistId, limit }: FollowersListProps): JSX.Element
```

**Features:**
- Shows total follower count
- Displays avatar with initial letter fallback
- Shows display name, location, and verified badge
- Click to navigate to follower's profile
- Loading and error states
- Empty state message

### FollowingList

Displays artists that the current user is following.

```typescript
interface FollowingListProps {
  limit?: number;  // Default: 20
}

export function FollowingList({ limit }: FollowingListProps): JSX.Element
```

**Features:**
- Shows total following count
- Displays artist avatar, name, verified badge
- Shows follower and track counts for each artist
- Includes inline FollowButton for quick unfollow
- Click to navigate to artist profile
- Loading and error states
- Empty state with discovery prompt

## Dependencies

| Package | Usage |
|---------|-------|
| `@tanstack/react-router` | Navigation to artist profiles |
| `../../hooks/useFollow` | Follow queries and mutations |
| `../../hooks/useAuth` | Authentication and role checking |
| `../../types` | ArtistProfile type |

## Hooks Used

```typescript
// From useFollow.ts
const { isFollowing, isLoading, isToggling, toggle } = useFollowToggle(artistId);
const { data, isLoading, error } = useFollowers(artistId, { limit });
const { data, isLoading, error } = useFollowing({ limit });
```

## Permission Requirements

| Component | Required Role |
|-----------|---------------|
| `FollowButton` | `subscriber` or higher |
| `FollowersList` | Public (any authenticated user) |
| `FollowingList` | Authenticated (shows current user's following) |

## Usage

```typescript
import { FollowButton, FollowersList, FollowingList } from '@/components/follow';

// Follow button on artist card
<FollowButton artistId="artist-123" size="sm" />

// Display artist's followers
<FollowersList artistId="artist-123" limit={10} />

// Display current user's following
<FollowingList limit={20} />
```

## Related Components

- `ArtistProfileCard` - Uses FollowButton for inline follow action
- `VisibilitySelector` - Related visibility feature for playlists
