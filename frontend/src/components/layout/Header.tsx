import { useThemeStore } from '@/lib/store/themeStore';

export function Header() {
  const { theme, toggleTheme } = useThemeStore();

  return (
    <header role="banner" className="navbar bg-base-100 shadow-md">
      <div className="flex-1">
        <span className="text-xl font-bold">Music Library</span>
      </div>
      <div className="flex-none gap-2">
        <button
          onClick={toggleTheme}
          className="btn btn-ghost btn-circle"
          aria-label="Toggle theme"
        >
          {theme === 'dark' ? '☀️' : '🌙'}
        </button>
      </div>
    </header>
  );
}
