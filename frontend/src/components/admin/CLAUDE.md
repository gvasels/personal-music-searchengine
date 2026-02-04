# Admin Components - CLAUDE.md

## Overview

Admin panel components for user management. Only visible and accessible to admin users. Provides search, viewing, and management of user accounts including role changes and account status.

## Files

| File | Description |
|------|-------------|
| `UserSearchForm.tsx` | Debounced search input for finding users by email or display name |
| `UserCard.tsx` | Card component displaying user summary in search results |
| `UserDetailModal.tsx` | Modal for viewing full user details and editing role/status |
| `RoleSwitcher.tsx` | Dropdown for admin role simulation |
| `SimulationBanner.tsx` | Alert banner shown when role simulation is active |
| `index.ts` | Barrel export for admin components |

## Components

### UserSearchForm

Search input with clear button and loading state.

```typescript
interface UserSearchFormProps {
  onSearch: (query: string) => void;
  isLoading?: boolean;  // Shows spinner when true
  placeholder?: string;  // Default: 'Search by email or name...'
}

export function UserSearchForm(props: UserSearchFormProps): JSX.Element
```

**Features:**
- Search icon with input field
- Loading spinner or clear button (contextual)
- Helper text with usage instructions
- Form submit handler

### UserCard

Clickable card showing user summary in search results.

```typescript
interface UserSummary {
  id: string;
  email: string;
  displayName: string;
  role: UserRole;
  disabled: boolean;
  createdAt: string;
}

interface UserCardProps {
  user: UserSummary;
  onSelect: (userId: string) => void;
  isSelected?: boolean;  // Adds ring highlight when true
}

export function UserCard(props: UserCardProps): JSX.Element
```

**Features:**
- Display name (or "No display name" fallback)
- Email address
- Role badge with color coding:
  - Guest: `badge-ghost`
  - Subscriber: `badge-info`
  - Artist: `badge-secondary`
  - Admin: `badge-primary`
- Disabled badge (red outline) when account is disabled
- Join date formatted
- Chevron indicator for click affordance
- Visual feedback for selected state

### UserDetailModal

Modal dialog for viewing and editing user details.

```typescript
interface UserDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  userId: string | null;  // Fetches user details when provided
}

export function UserDetailModal(props: UserDetailModalProps): JSX.Element | null
```

**Features:**
- User avatar with initials
- Full user information display
- Stats grid: tracks, playlists, albums, followers, following, storage
- Join date and last login time
- Role selector dropdown with confirmation dialog
- Account status toggle with confirmation dialog
- Loading states for data fetching and mutations
- Error handling with dismissible alerts

**Role Change Flow:**
1. User selects new role from dropdown
2. Warning dialog appears explaining the change
3. User confirms or cancels
4. On confirm: updates DynamoDB and Cognito user groups atomically

**Status Change Flow:**
1. User toggles disable/enable checkbox
2. Confirmation dialog appears
3. User confirms or cancels
4. On confirm: updates user status in DynamoDB and Cognito

## Dependencies

| Package | Usage |
|---------|-------|
| `../../hooks/useAdmin` | `useSearchUsers`, `useUserDetails`, mutations |
| `../../types` | `UserRole` type |
| `react-hot-toast` | Toast notifications (in parent page) |

## Types

```typescript
type UserRole = 'guest' | 'subscriber' | 'artist' | 'admin';

interface UserSummary {
  id: string;
  email: string;
  displayName: string;
  role: UserRole;
  disabled: boolean;
  createdAt: string;
}

interface UserDetails extends UserSummary {
  lastLoginAt?: string;
  trackCount: number;
  playlistCount: number;
  albumCount: number;
  storageUsed: number;
  followerCount: number;
  followingCount: number;
}
```

## Usage

```typescript
import { UserSearchForm, UserCard, UserDetailModal } from '@/components/admin';

// In admin users page
const [searchQuery, setSearchQuery] = useState('');
const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
const [isModalOpen, setIsModalOpen] = useState(false);

<UserSearchForm onSearch={setSearchQuery} isLoading={isLoading} />

{users.map((user) => (
  <UserCard
    key={user.id}
    user={user}
    onSelect={(id) => {
      setSelectedUserId(id);
      setIsModalOpen(true);
    }}
  />
))}

<UserDetailModal
  isOpen={isModalOpen}
  onClose={() => setIsModalOpen(false)}
  userId={selectedUserId}
/>
```

## Access Control

These components are designed for admin-only access. The parent route (`/admin/users`) handles:
- Checking `isAdmin` from `useAuth()`
- Redirecting non-admin users with toast notification
- Showing loading state during auth check

### RoleSwitcher

Admin-only dropdown component for simulating different user roles.

```typescript
export function RoleSwitcher(): JSX.Element | null
```

**Features:**
- Only renders for admin users (returns null otherwise)
- Dropdown with all available roles: Guest, Subscriber, Artist, Admin
- Visual indicator (checkmark) for currently simulated role
- Warning style (btn-warning) when simulation is active
- Normal style (btn-ghost) when not simulating

**Role Simulation Behavior:**
- Selecting a role activates simulation mode
- Selecting "Admin" stops simulation (returns to actual admin role)
- Persists to localStorage via roleSimulationStore

### SimulationBanner

Alert banner displayed when role simulation is active.

```typescript
export function SimulationBanner(): JSX.Element | null
```

**Features:**
- Only renders when simulation is active (returns null otherwise)
- Shows the simulated role name (Guest/Subscriber/Artist)
- "Exit Simulation" button to stop simulation
- Fixed position at top of viewport
- Warning color scheme (alert-warning)
- ARIA role="alert" for accessibility

## Related Components

- `useAdmin` hook - Data fetching and mutations for admin operations
- `useRoleSimulation` hook - Role simulation state and actions
- `roleSimulationStore` - Zustand store for simulation persistence
- Admin API (`/lib/api/admin.ts`) - API client functions
- Sidebar admin section - Navigation to admin pages
