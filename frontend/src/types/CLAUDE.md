# Types - CLAUDE.md

## Overview

Shared TypeScript type definitions for the frontend application. Contains domain models, API response types, and enums.

## File Descriptions

| File | Purpose |
|------|---------|
| `index.ts` | All shared type definitions |

## Key Types

### Music Library Types
- `Track`, `Album`, `Artist`, `Playlist`, `Tag`, `Upload`
- `PaginatedResponse<T>`, `SearchResponse`

### Artist Events Types
| Type | Description |
|------|-------------|
| `EventStatus` | `'scheduled' \| 'cancelled' \| 'postponed' \| 'rescheduled'` |
| `ArtistEvent` | Event with id, artistName, title, date, venue, city, region, country, ticketUrl, status, source |
| `ArtistSearchResult` | Artist search result with name, upcomingEvents count, source |
| `WatchedArtist` | Watched artist with artistName, watchedAt |
| `ArtistEventsResponse` | `{ artistName, events: ArtistEvent[], totalCount, source }` |
| `WatchResponse` | `{ artistName, watching, watchedAt }` |
| `WatchStatusResponse` | `{ watching, artistName }` |
| `WatchedArtistsResponse` | `{ items: WatchedArtist[], hasMore, nextCursor? }` |
| `ArtistSearchResponse` | `{ items: ArtistSearchResult[], total }` |

## Dependencies

No external dependencies - pure TypeScript types.
