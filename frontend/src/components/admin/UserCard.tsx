/**
 * UserCard - Admin Panel Component
 * Displays user summary in search results with click to select
 */
import type { UserRole } from '../../types';

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
  isSelected?: boolean;
}

// Role badge colors
const ROLE_BADGE_COLORS: Record<UserRole, string> = {
  guest: 'badge-ghost',
  subscriber: 'badge-info',
  artist: 'badge-secondary',
  admin: 'badge-primary',
};

// Role display names
const ROLE_DISPLAY_NAMES: Record<UserRole, string> = {
  guest: 'Guest',
  subscriber: 'Subscriber',
  artist: 'Artist',
  admin: 'Admin',
};

export function UserCard({ user, onSelect, isSelected = false }: UserCardProps) {
  const formattedDate = new Date(user.createdAt).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });

  return (
    <button
      type="button"
      onClick={() => onSelect(user.id)}
      className={`card bg-base-100 shadow-sm hover:shadow-md transition-shadow w-full text-left ${
        isSelected ? 'ring-2 ring-primary' : ''
      } ${user.disabled ? 'opacity-60' : ''}`}
    >
      <div className="card-body p-4">
        <div className="flex items-start justify-between gap-4">
          {/* User info */}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <h3 className="font-semibold text-base-content truncate">
                {user.displayName || 'No display name'}
              </h3>

              {/* Role badge */}
              <span className={`badge badge-sm ${ROLE_BADGE_COLORS[user.role]}`}>
                {ROLE_DISPLAY_NAMES[user.role]}
              </span>

              {/* Disabled badge */}
              {user.disabled && (
                <span className="badge badge-sm badge-error badge-outline">
                  Disabled
                </span>
              )}
            </div>

            <p className="text-sm text-base-content/70 truncate mt-1">
              {user.email}
            </p>

            <p className="text-xs text-base-content/50 mt-1">
              Joined {formattedDate}
            </p>
          </div>

          {/* Chevron indicator */}
          <svg
            className="w-5 h-5 text-base-content/40 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5l7 7-7 7"
            />
          </svg>
        </div>
      </div>
    </button>
  );
}
