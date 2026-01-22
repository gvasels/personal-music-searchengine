import { Link } from '@tanstack/react-router';
import { useAuth } from '../../hooks/useAuth';

const navItems = [
  { to: '/', label: 'Home', icon: '🏠' },
  { to: '/tracks', label: 'Music', icon: '🎵' },
  { to: '/playlists', label: 'Playlists', icon: '📝' },
  { to: '/upload', label: 'Upload', icon: '⬆️' },
];

interface SidebarProps {
  isOpen?: boolean;
  onClose?: () => void;
}

export function Sidebar({ isOpen = false, onClose }: SidebarProps) {
  const { isAuthenticated, isLoading } = useAuth();

  // Don't show sidebar when not authenticated or still loading
  if (isLoading || !isAuthenticated) {
    return null;
  }

  const handleNavClick = () => {
    // Close sidebar on mobile after navigation
    if (onClose) {
      onClose();
    }
  };

  return (
    <>
      {/* Mobile overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 md:hidden"
          onClick={onClose}
        />
      )}

      {/* Sidebar */}
      <nav
        role="navigation"
        className={`
          fixed md:static inset-y-0 left-0 z-50
          w-64 bg-base-100 p-4
          transform transition-transform duration-200 ease-in-out
          md:transform-none md:block
          ${isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'}
        `}
      >
        {/* Close button for mobile */}
        <div className="flex justify-end md:hidden mb-4">
          <button
            onClick={onClose}
            className="btn btn-ghost btn-circle btn-sm"
            aria-label="Close menu"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <ul className="menu">
          {navItems.map((item) => (
            <li key={item.to}>
              <Link to={item.to} className="flex items-center gap-2" onClick={handleNavClick}>
                <span>{item.icon}</span>
                <span>{item.label}</span>
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </>
  );
}
