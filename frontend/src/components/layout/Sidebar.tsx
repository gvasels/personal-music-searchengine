/**
 * Sidebar - Desktop Navigation
 * Section-based navigation: Home, Music, Videos, Gaming
 * Music sub-nav shown when on /music/* routes
 */
import { Link, useLocation } from '@tanstack/react-router';
import { useAuth } from '../../hooks/useAuth';
import { useFeatureFlags } from '../../hooks/useFeatureFlags';

type MinRole = 'guest' | 'subscriber' | 'artist' | 'admin';

interface NavItem {
  to: string;
  label: string;
  icon: string;
}

interface NavSection {
  id: string;
  label: string;
  icon: string;
  to: string;
  minRole: MinRole;
  children?: NavItem[];
}

const navSections: NavSection[] = [
  { id: 'home', label: 'Home', icon: '🏠', to: '/', minRole: 'guest' },
  {
    id: 'music', label: 'Music', icon: '🎵', to: '/music', minRole: 'subscriber',
    children: [
      { to: '/music/tracks', label: 'Tracks', icon: '🎵' },
      { to: '/music/albums', label: 'Albums', icon: '💿' },
      { to: '/music/artists', label: 'Artists', icon: '🎤' },
      { to: '/music/playlists', label: 'Playlists', icon: '📝' },
      { to: '/music/tags', label: 'Tags', icon: '🏷️' },
    ],
  },
  { id: 'videos', label: 'Videos', icon: '🎬', to: '/videos', minRole: 'subscriber' },
  { id: 'gaming', label: 'Gaming', icon: '🎮', to: '/gaming', minRole: 'subscriber' },
];

const utilityItems: (NavItem & { minRole: MinRole })[] = [
  { to: '/upload', label: 'Upload', icon: '⬆️', minRole: 'artist' },
  { to: '/search', label: 'Search', icon: '🔍', minRole: 'subscriber' },
  { to: '/settings', label: 'Settings', icon: '⚙️', minRole: 'subscriber' },
];

const adminNavItems: NavItem[] = [
  { to: '/admin/users', label: 'User Management', icon: '👥' },
];

const linkClass = 'flex items-center gap-2 rounded-lg transition-colors';
const activeLinkClass = `${linkClass} bg-primary/20 border-l-3 border-primary font-semibold`;

export function Sidebar() {
  const { isAuthenticated, isLoading } = useAuth();
  const { role, hasRole, isSimulating } = useFeatureFlags();
  const location = useLocation();

  if (isLoading || !isAuthenticated) return null;

  const inMusicSection = location.pathname.startsWith('/music');

  return (
    <nav role="navigation" className="hidden md:block w-64 bg-base-100 p-4 overflow-y-auto">
      {isSimulating && (
        <div className="mb-4 p-2 bg-warning/20 rounded-lg text-center text-sm">
          Viewing as: <span className="font-semibold capitalize">{role}</span>
        </div>
      )}

      <ul className="menu">
        {navSections.filter((s) => hasRole(s.minRole)).map((section) => (
          <li key={section.id}>
            <Link
              to={section.to}
              className={linkClass}
              activeProps={{ className: activeLinkClass, 'aria-current': 'page' }}
            >
              <span>{section.icon}</span>
              <span>{section.label}</span>
            </Link>
            {section.children && inMusicSection && section.id === 'music' && (
              <ul className="ml-4 mt-1">
                {section.children.map((child) => (
                  <li key={child.to}>
                    <Link
                      to={child.to}
                      className={linkClass}
                      activeProps={{ className: activeLinkClass, 'aria-current': 'page' }}
                    >
                      <span>{child.icon}</span>
                      <span>{child.label}</span>
                    </Link>
                  </li>
                ))}
              </ul>
            )}
          </li>
        ))}

        {/* Utility items */}
        <li className="menu-title mt-4">
          <span className="text-xs uppercase tracking-wider text-base-content/50">Tools</span>
        </li>
        {utilityItems.filter((i) => hasRole(i.minRole)).map((item) => (
          <li key={item.to}>
            <Link to={item.to} className={linkClass} activeProps={{ className: activeLinkClass, 'aria-current': 'page' }}>
              <span>{item.icon}</span>
              <span>{item.label}</span>
            </Link>
          </li>
        ))}

        {/* Admin section */}
        {hasRole('admin') && (
          <>
            <li className="menu-title mt-4">
              <span className="text-xs uppercase tracking-wider text-base-content/50">Admin</span>
            </li>
            {adminNavItems.map((item) => (
              <li key={item.to}>
                <Link to={item.to} className={linkClass} activeProps={{ className: activeLinkClass, 'aria-current': 'page' }}>
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
