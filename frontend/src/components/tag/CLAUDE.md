# Tag Components - CLAUDE.md

## Overview

Tag management components for adding and removing tags from tracks with inline editing.

## Files

| File | Description |
|------|-------------|
| `TagInput.tsx` | Inline tag editor with add/remove functionality |
| `index.ts` | Barrel export for tag components |
| `__tests__/TagInput.test.tsx` | Unit tests (6 tests, 82.75% coverage) |

## Key Functions

### TagInput.tsx
```typescript
interface TagInputProps {
  trackId: string;
  tags: string[];
}

export function TagInput({ trackId, tags }: TagInputProps): JSX.Element
```
Renders inline tag editor with:
- Badge display for existing tags (with remove button)
- "Add tag" button that expands to input field
- Auto-lowercase normalization on submission

**Tag Display:**
- Each tag shown as badge with emoji prefix
- Remove button (×) on each tag
- Add tag button at end of list

**Adding Tags:**
- Click "Add tag" to show input
- Type tag name and press Enter
- Tags automatically normalized to lowercase
- Input clears and hides after submission
- Escape key cancels input

**API Integration:**
- `addTagToTrack(trackId, tagName)` - Add new tag
- `removeTagFromTrack(trackId, tagName)` - Remove tag
- Invalidates `track` query on success

## Dependencies

| Package | Usage |
|---------|-------|
| `@tanstack/react-query` | `useMutation`, `useQueryClient` |
| `@/lib/api/client` | `addTagToTrack`, `removeTagFromTrack` |
| `react-hot-toast` | Toast notifications |

## Usage

```typescript
import { TagInput } from '@/components/tag';

<TagInput trackId={track.id} tags={track.tags} />
```

## Tag Normalization

All tags are normalized to lowercase before submission:
- "Rock" → "rock"
- "JAZZ" → "jazz"
- "Hip-Hop" → "hip-hop"
