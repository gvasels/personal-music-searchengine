# Playlist Components - CLAUDE.md

## Overview

Playlist management components including modal for creating new playlists.

## Files

| File | Description |
|------|-------------|
| `CreatePlaylistModal.tsx` | Modal dialog for creating new playlists |
| `index.ts` | Barrel export for playlist components |
| `__tests__/CreatePlaylistModal.test.tsx` | Unit tests (7 tests, 88.79% coverage) |

## Key Functions

### CreatePlaylistModal.tsx
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

## Dependencies

| Package | Usage |
|---------|-------|
| `@tanstack/react-query` | `useMutation`, `useQueryClient` |
| `@/lib/api/client` | `createPlaylist` function |
| `react-hot-toast` | Toast notifications |

## Usage

```typescript
import { CreatePlaylistModal } from '@/components/playlist';

const [isOpen, setIsOpen] = useState(false);

<CreatePlaylistModal isOpen={isOpen} onClose={() => setIsOpen(false)} />
```
