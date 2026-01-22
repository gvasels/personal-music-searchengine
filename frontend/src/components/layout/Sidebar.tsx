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
];

export function Sidebar() {
  const { isAuthenticated, isLoading } = useAuth();

  // Don't show sidebar when not authenticated or still loading
  if (isLoading || !isAuthenticated) {
    return null;
  }

  return (
    <nav role="navigation" className="w-64 bg-base-100 p-4 hidden md:block">
      <ul className="menu">
        {navItems.map((item) => (
          <li key={item.to}>
            <Link to={item.to} className="flex items-center gap-2">
              <span>{item.icon}</span>
              <span>{item.label}</span>
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}
