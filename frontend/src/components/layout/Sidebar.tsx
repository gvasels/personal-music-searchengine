/**
 * Sidebar - Desktop Navigation
 * Shows navigation menu on desktop, hidden on mobile (MobileNav handles mobile)
 * Respects role simulation for admin testing
 */
import { Link } from '@tanstack/react-router';
import { useAuth } from '../../hooks/useAuth';
import { useFeatureFlags } from '../../hooks/useFeatureFlags';

const navItems = [
  { to: '/', label: 'Home', icon: '🏠', minRole: 'guest' as const },
  { to: '/tracks', label: 'Tracks', icon: '🎵', minRole: 'subscriber' as const },
  { to: '/albums', label: 'Albums', icon: '💿', minRole: 'subscriber' as const },
  { to: '/artists', label: 'Artists', icon: '🎤', minRole: 'subscriber' as const },
  { to: '/playlists', label: 'Playlists', icon: '📝', minRole: 'subscriber' as const },
  { to: '/tags', label: 'Tags', icon: '🏷️', minRole: 'subscriber' as const },
  { to: '/upload', label: 'Upload', icon: '⬆️', minRole: 'artist' as const },
  { to: '/settings', label: 'Settings', icon: '⚙️', minRole: 'subscriber' as const },
];

// Admin-only navigation items
const adminNavItems = [
  { to: '/admin/users', label: 'User Management', icon: '👥' },
];

export function Sidebar() {
  const { isAuthenticated, isLoading } = useAuth();
  const { role, hasRole, isSimulating } = useFeatureFlags();

  // Check if current (effective) role is admin
  const showAdminSection = hasRole('admin');

  // Don't show sidebar when not authenticated or still loading
  if (isLoading || !isAuthenticated) {
    return null;
  }

  // Filter nav items based on role requirements
  const visibleNavItems = navItems.filter((item) => {
    if ('minRole' in item && item.minRole) {
      return hasRole(item.minRole);
    }
    return true;
  });

  return (
    <nav
      role="navigation"
      className="hidden md:block w-64 bg-base-100 p-4 overflow-y-auto"
    >
      {/* Show simulation indicator */}
      {isSimulating && (
        <div className="mb-4 p-2 bg-warning/20 rounded-lg text-center text-sm">
          Viewing as: <span className="font-semibold capitalize">{role}</span>
        </div>
      )}

      <ul className="menu">
        {visibleNavItems.map((item) => (
          <li key={item.to}>
            <Link
              to={item.to}
              className="flex items-center gap-2 rounded-lg transition-colors"
              activeProps={{
                className: 'flex items-center gap-2 rounded-lg transition-colors bg-primary/20 border-l-3 border-primary font-semibold',
                'aria-current': 'page',
              }}
            >
              <span>{item.icon}</span>
              <span>{item.label}</span>
            </Link>
          </li>
        ))}

        {/* Admin section - only visible when effective role is admin */}
        {showAdminSection && (
          <>
            <li className="menu-title mt-4">
              <span className="text-xs uppercase tracking-wider text-base-content/50">Admin</span>
            </li>
            {adminNavItems.map((item) => (
              <li key={item.to}>
                <Link
                  to={item.to}
                  className="flex items-center gap-2 rounded-lg transition-colors"
                  activeProps={{
                    className: 'flex items-center gap-2 rounded-lg transition-colors bg-primary/20 border-l-3 border-primary font-semibold',
                    'aria-current': 'page',
                  }}
                >
                  <span>{item.icon}</span>
                  <span>{item.label}</span>
                </Link>
              </li>
            ))}
          </>
        )}
      </ul>
    </nav>
  );
}
