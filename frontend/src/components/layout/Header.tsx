import { useNavigate } from '@tanstack/react-router';
import { useThemeStore } from '@/lib/store/themeStore';
import { useAuth } from '@/hooks/useAuth';
import { SearchBar } from '@/components/search/SearchBar';

export function Header() {
  const navigate = useNavigate();
  const { theme, toggleTheme } = useThemeStore();
  const { isAuthenticated, signOut, isSigningOut } = useAuth();

  const handleSignOut = async () => {
    await signOut();
    navigate({ to: '/login' });
  };

  return (
    <header role="banner" className="navbar bg-base-100 shadow-md">
      <div className="flex-1">
        <span className="text-xl font-bold">Music Library</span>
      </div>
      <div className="flex-1 px-4 hidden md:block">
        {isAuthenticated && <SearchBar />}
      </div>
      <div className="flex-none gap-2">
        <button
          onClick={toggleTheme}
          className="btn btn-ghost btn-circle"
          aria-label="Toggle theme"
        >
          {theme === 'dark' ? '☀️' : '🌙'}
        </button>
        {isAuthenticated && (
          <button
            onClick={handleSignOut}
            className="btn btn-ghost btn-sm"
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
