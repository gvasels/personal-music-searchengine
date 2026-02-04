# Playlist Components - CLAUDE.md

## Overview

Playlist management components including modal for creating new playlists and visibility controls for the Global User Type feature.

## Files

| File | Description |
|------|-------------|
| `CreatePlaylistModal.tsx` | Modal dialog for creating new playlists |
| `VisibilitySelector.tsx` | Visibility controls: dropdown, badge, and radio group |
| `index.ts` | Barrel export for playlist components |
| `__tests__/CreatePlaylistModal.test.tsx` | Unit tests (7 tests, 88.79% coverage) |
| `__tests__/VisibilitySelector.test.tsx` | Unit tests for visibility components |

## Components

### CreatePlaylistModal

```typescript
interface CreatePlaylistModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CreatePlaylistModal({ isOpen, onClose }: CreatePlaylistModalProps): JSX.Element | null
```

Renders a modal dialog with:
- Name input (required, min 3 characters)
- Description textarea (optional)
- Cancel and Create buttons
- Loading state during submission

**Form Validation:**
- Name is required
- Name must be at least 3 characters
- Trims whitespace before submission

**Keyboard Support:**
- Escape key closes modal
- Enter submits form

**API Integration:**
- Uses `createPlaylist` mutation from API client
- Invalidates `playlists` query on success
- Shows toast notifications for success/error

### VisibilitySelector

Dropdown select for choosing playlist visibility.

```typescript
interface VisibilitySelectorProps {
  value: PlaylistVisibility;
  onChange: (visibility: PlaylistVisibility) => void;
  disabled?: boolean;
  size?: 'sm' | 'md' | 'lg';  // Default: 'md'
}

export function VisibilitySelector(props: VisibilitySelectorProps): JSX.Element
```

**Visibility Options:**
| Value | Label | Icon | Description |
|-------|-------|------|-------------|
| `private` | Private | Lock | Only you can see this playlist |
| `unlisted` | Unlisted | Link | Anyone with the link can see |
| `public` | Public | Globe | Visible to everyone |

### VisibilityBadge

Badge displaying current visibility status.

```typescript
interface VisibilityBadgeProps {
  visibility: PlaylistVisibility;
  size?: 'sm' | 'md' | 'lg';  // Default: 'md'
}

export function VisibilityBadge({ visibility, size }: VisibilityBadgeProps): JSX.Element | null
```

**Badge Styles:**
| Visibility | Badge Class |
|------------|-------------|
| `private` | `badge-ghost` |
| `unlisted` | `badge-warning` |
| `public` | `badge-success` |

### VisibilityRadioGroup

Radio button group for visibility selection with descriptions.

```typescript
interface VisibilityRadioGroupProps {
  value: PlaylistVisibility;
  onChange: (visibility: PlaylistVisibility) => void;
  disabled?: boolean;
}

export function VisibilityRadioGroup(props: VisibilityRadioGroupProps): JSX.Element
```

**Features:**
- Larger touch targets with descriptions
- Visual highlight for selected option
- Disabled state styling

## Types

```typescript
type PlaylistVisibility = 'private' | 'unlisted' | 'public';
```

## Dependencies

| Package | Usage |
|---------|-------|
| `@tanstack/react-query` | `useMutation`, `useQueryClient` |
| `@/lib/api/client` | `createPlaylist` function |
| `react-hot-toast` | Toast notifications |
| `../../types` | `PlaylistVisibility` type |

## Usage

```typescript
import {
  CreatePlaylistModal,
  VisibilitySelector,
  VisibilityBadge,
  VisibilityRadioGroup
} from '@/components/playlist';

// Create playlist modal
const [isOpen, setIsOpen] = useState(false);
<CreatePlaylistModal isOpen={isOpen} onClose={() => setIsOpen(false)} />

// Visibility dropdown
const [visibility, setVisibility] = useState<PlaylistVisibility>('private');
<VisibilitySelector value={visibility} onChange={setVisibility} />

// Visibility badge (read-only display)
<VisibilityBadge visibility="public" size="sm" />

// Visibility radio group (for forms)
<VisibilityRadioGroup value={visibility} onChange={setVisibility} />
```

## Related Components

- `FollowButton` - Related social feature
- `ArtistProfileCard` - Artist profiles can have public/private content
