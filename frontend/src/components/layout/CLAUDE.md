# Layout Components - CLAUDE.md

## Overview

App shell components providing the main layout structure including header, sidebar navigation, and responsive container. Supports role-based access control and admin role simulation.

## Files

| File | Description |
|------|-------------|
| `Layout.tsx` | Main app shell with role simulation support |
| `Header.tsx` | Top navigation bar with theme toggle button |
| `Sidebar.tsx` | Role-aware navigation menu |
| `MobileNav.tsx` | Mobile navigation overlay |
| `index.ts` | Barrel export for all layout components |
| `__tests__/Layout.test.tsx` | Unit tests for layout components |

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
- SimulationBanner (when admin is simulating another role)

**Role Simulation Behavior:**
- When admin simulates **Guest** role → automatically redirects to `/permission-denied`
- Uses `useFeatureFlags()` to check `isSimulating` and `role`

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
Role-aware navigation menu with TanStack Router `Link` components.

**Navigation Items by Role:**
| Route | Required Role | Description |
|-------|---------------|-------------|
| `/` (Home) | guest | Dashboard - accessible to all |
| `/tracks` | subscriber | Track listing |
| `/albums` | subscriber | Album grid |
| `/artists` | subscriber | Artist listing |
| `/playlists` | subscriber | Playlist management |
| `/tags` | subscriber | Tag cloud |
| `/upload` | artist | File upload |
| `/settings` | subscriber | User settings |
| `/admin/users` | admin | User management (admin section) |

**Simulation Indicator:** Shows "Viewing as: [role]" badge when simulating.

## Role Simulation

Admins can simulate other roles to test the UI experience:

| Simulated Role | Behavior |
|----------------|----------|
| **Guest** | Immediately redirected to permission-denied page |
| **Subscriber** | Sees subscriber nav items, no upload/admin |
| **Artist** | Sees upload option, no admin section |

## Dependencies

| Package | Usage |
|---------|-------|
| `@tanstack/react-router` | `Link`, `useNavigate`, `useLocation` |
| `@/lib/store/themeStore` | Theme state management |
| `@/lib/store/playerStore` | Player state for PlayerBar |
| `@/components/player` | PlayerBar component |
| `@/hooks/useFeatureFlags` | Role and simulation state |
