import { useNavigate } from '@tanstack/react-router';
import { useThemeStore } from '@/lib/store/themeStore';
import { useAuth } from '@/hooks/useAuth';
import { SearchBar } from '@/components/search/SearchBar';

interface HeaderProps {
  onMenuClick?: () => void;
}

export function Header({ onMenuClick }: HeaderProps) {
  const navigate = useNavigate();
  const { theme, toggleTheme } = useThemeStore();
  const { isAuthenticated, signOut, isSigningOut } = useAuth();

  const handleSignOut = async () => {
    await signOut();
    navigate({ to: '/login' });
  };

  return (
    <header role="banner" className="navbar bg-secondary text-secondary-content shadow-md">
      <div className="flex-none md:hidden">
        <button
          onClick={onMenuClick}
          className="btn btn-ghost btn-circle text-secondary-content"
          aria-label="Open menu"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
      </div>
      <div className="flex-none">
        <img src="/logo.png" alt="Music Library" className="h-10" />
      </div>
      <div className="flex-1"></div>
      {isAuthenticated && (
        <div className="absolute left-1/2 -translate-x-1/2 hidden md:block header-search-wrapper">
          <SearchBar />
        </div>
      )}
      <div className="flex-none gap-2 z-10">
        <button
          onClick={toggleTheme}
          className="btn btn-ghost btn-circle text-secondary-content hover:bg-secondary-content/20"
          aria-label="Toggle theme"
        >
          {theme === 'dark' ? '☀️' : '🌙'}
        </button>
        {isAuthenticated && (
          <button
            onClick={handleSignOut}
            className="btn btn-primary btn-sm"
            disabled={isSigningOut}
            aria-label="Sign out"
          >
            {isSigningOut ? 'Signing out...' : 'Logout'}
          </button>
        )}
      </div>
    </header>
  );
}
