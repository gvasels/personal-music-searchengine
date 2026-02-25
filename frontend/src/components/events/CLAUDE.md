# Events Components - CLAUDE.md

## Overview

React components for the Artist Events & Shows feature. Displays event information, artist watch toggles, and event lists.

## File Descriptions

| File | Purpose |
|------|---------|
| `WatchButton.tsx` | Toggle button to watch/unwatch an artist for event notifications |
| `EventCard.tsx` | Display a single artist event with venue, date, status badge, ticket link |
| `ArtistEventsSection.tsx` | Fetch and display upcoming events for an artist with loading/error/empty states |
| `index.ts` | Barrel export for all components |

## Components

### WatchButton
```typescript
interface WatchButtonProps {
  artistName: string;
  className?: string;
  size?: 'sm' | 'md';
}
```
- Requires authentication (subscriber+) - renders `null` for guests
- Uses `useWatchToggle` hook for state and toggle action
- Shows loading spinner when fetching status
- Toggles between "Watch" (btn-primary) and "Watching" (btn-outline)

### EventCard
```typescript
interface EventCardProps {
  event: ArtistEvent;
}
```
- Displays: title, venue, city/region/country, formatted date
- Status badge: "Cancelled" (badge-error) or "Postponed" (badge-warning) - hidden for "scheduled"
- Ticket link: opens in new tab with `rel="noopener noreferrer"`

### ArtistEventsSection
```typescript
interface ArtistEventsSectionProps {
  artistName: string;
}
```
- Uses `useArtistEvents` hook to fetch events
- States: loading ("Loading events..."), error ("Failed to load events"), empty ("No upcoming events"), data (EventCard list)

## Dependencies

### Internal
- `../../hooks/useArtistWatch` - `useWatchToggle`
- `../../hooks/useArtistEvents` - `useArtistEvents`
- `../../hooks/useAuth` - `useAuth`
- `../../types` - `ArtistEvent`

## Testing

Tests in `__tests__/`:
- `WatchButton.test.tsx` - 11 tests (render states, toggle, auth, styling)
- `EventCard.test.tsx` - 9 tests (content rendering, ticket link, status badges)
- `ArtistEventsSection.test.tsx` - 4 tests (loading, data, empty, error states)
