/**
 * Sidebar - Desktop Navigation
 * Shows navigation menu on desktop, hidden on mobile (MobileNav handles mobile)
 */
import { Link } from '@tanstack/react-router';
import { useAuth } from '../../hooks/useAuth';

const navItems = [
  { to: '/', label: 'Home', icon: '🏠' },
  { to: '/tracks', label: 'Tracks', icon: '🎵' },
  { to: '/albums', label: 'Albums', icon: '💿' },
  { to: '/artists', label: 'Artists', icon: '🎤' },
  { to: '/playlists', label: 'Playlists', icon: '📝' },
  { to: '/tags', label: 'Tags', icon: '🏷️' },
  { to: '/upload', label: 'Upload', icon: '⬆️' },
  { to: '/settings', label: 'Settings', icon: '⚙️' },
];

// Admin-only navigation items
const adminNavItems = [
  { to: '/admin/users', label: 'User Management', icon: '👥' },
];

export function Sidebar() {
  const { isAuthenticated, isLoading, isAdmin } = useAuth();

  // Don't show sidebar when not authenticated or still loading
  if (isLoading || !isAuthenticated) {
    return null;
  }

  return (
    <nav
      role="navigation"
      className="hidden md:block w-64 bg-base-100 p-4 overflow-y-auto"
    >
      <ul className="menu">
        {navItems.map((item) => (
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

        {/* Admin section - only visible to admins */}
        {isAdmin && (
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
