# Layout Components - CLAUDE.md

## Overview

App shell components providing the main layout structure including header, sidebar navigation, and responsive container.

## Files

| File | Description |
|------|-------------|
| `Layout.tsx` | Main app shell wrapping content with Header, Sidebar, and PlayerBar |
| `Header.tsx` | Top navigation bar with theme toggle button |
| `Sidebar.tsx` | Navigation menu with links to main sections |
| `index.ts` | Barrel export for all layout components |
| `__tests__/Layout.test.tsx` | Unit tests for layout components (5 tests) |

## Key Functions

### Layout.tsx
```typescript
export function Layout({ children }: { children: React.ReactNode }): JSX.Element
```
Renders the main app shell with:
- Sidebar navigation (left)
- Header (top)
- Main content area (center)
- PlayerBar (bottom, fixed)

### Header.tsx
```typescript
export function Header(): JSX.Element
```
Renders the top navigation with:
- Theme toggle button (dark/light)
- Uses `useThemeStore` for theme state

### Sidebar.tsx
```typescript
export function Sidebar(): JSX.Element
```
Renders navigation menu with TanStack Router `Link` components:
- Home, Tracks, Albums, Artists, Playlists, Tags, Upload

## Dependencies

| Package | Usage |
|---------|-------|
| `@tanstack/react-router` | `Link` for navigation |
| `@/lib/store/themeStore` | Theme state management |
| `@/lib/store/playerStore` | Player state for PlayerBar |
| `@/components/player` | PlayerBar component |
