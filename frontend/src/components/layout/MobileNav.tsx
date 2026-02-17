/**
 * Mobile Navigation Overlay
 * Section-based navigation matching Sidebar structure
 */
import { useEffect, useRef, useCallback, useState } from 'react';
import { Link, useLocation } from '@tanstack/react-router';
import { usePlayerStore } from '@/lib/store/playerStore';
import { useAuth } from '@/hooks/useAuth';

interface NavItem { to: string; label: string; icon: string; }

const sections: (NavItem & { children?: NavItem[] })[] = [
  { to: '/', label: 'Home', icon: '🏠' },
  {
    to: '/music', label: 'Music', icon: '🎵',
    children: [
      { to: '/music/tracks', label: 'Tracks', icon: '🎵' },
      { to: '/music/albums', label: 'Albums', icon: '💿' },
      { to: '/music/artists', label: 'Artists', icon: '🎤' },
      { to: '/music/playlists', label: 'Playlists', icon: '📝' },
      { to: '/music/tags', label: 'Tags', icon: '🏷️' },
    ],
  },
  { to: '/videos', label: 'Videos', icon: '🎬' },
  { to: '/gaming', label: 'Gaming', icon: '🎮' },
  { to: '/upload', label: 'Upload', icon: '⬆️' },
  { to: '/settings', label: 'Settings', icon: '⚙️' },
];

interface MobileNavProps { isOpen: boolean; onClose: () => void; }

export function MobileNav({ isOpen, onClose }: MobileNavProps) {
  const location = useLocation();
  const navRef = useRef<HTMLDivElement>(null);
  const { currentTrack, isPlaying, play, pause } = usePlayerStore();
  const { isAdmin } = useAuth();
  const [touchStart, setTouchStart] = useState<number | null>(null);
  const [touchCurrent, setTouchCurrent] = useState<number | null>(null);
  const [translateX, setTranslateX] = useState(0);
  const [musicExpanded, setMusicExpanded] = useState(location.pathname.startsWith('/music'));

  useEffect(() => {
    document.body.style.overflow = isOpen ? 'hidden' : '';
    return () => { document.body.style.overflow = ''; };
  }, [isOpen]);

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape' && isOpen) onClose();
    }
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    setTouchStart(e.touches[0].clientX);
    setTouchCurrent(e.touches[0].clientX);
  }, []);

  const handleTouchMove = useCallback((e: React.TouchEvent) => {
    if (touchStart === null) return;
    const currentX = e.touches[0].clientX;
    setTouchCurrent(currentX);
    const diff = currentX - touchStart;
    if (diff < 0) setTranslateX(diff);
  }, [touchStart]);

  const handleTouchEnd = useCallback(() => {
    if (touchStart !== null && touchCurrent !== null && touchCurrent - touchStart < -100) {
      onClose();
    }
    setTouchStart(null);
    setTouchCurrent(null);
    setTranslateX(0);
  }, [touchStart, touchCurrent, onClose]);

  if (!isOpen) return null;

  const isActive = (to: string) => location.pathname === to || (to !== '/' && location.pathname.startsWith(to));
  const itemClass = (to: string) =>
    `flex items-center gap-4 p-4 rounded-xl text-lg font-medium transition-all ${
      isActive(to) ? 'bg-primary/20 text-primary border-l-4 border-primary' : 'hover:bg-base-200 active:bg-base-300'
    }`;

  return (
    <div className="fixed inset-0 z-50 md:hidden">
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200" onClick={onClose} />
      <div
        ref={navRef}
        className="absolute inset-y-0 left-0 w-full max-w-sm bg-base-100 shadow-2xl animate-in slide-in-from-left duration-200"
        style={{ transform: `translateX(${translateX}px)`, transition: touchStart === null ? 'transform 0.2s ease-out' : 'none' }}
        onTouchStart={handleTouchStart}
        onTouchMove={handleTouchMove}
        onTouchEnd={handleTouchEnd}
      >
        <div className="flex flex-col h-full">
          <div className="flex items-center justify-between p-4 border-b border-base-300">
            <h2 className="text-lg font-bold">Menu</h2>
            <button onClick={onClose} className="btn btn-ghost btn-circle btn-sm" aria-label="Close menu">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div className="px-4 py-2 text-xs text-base-content/50 text-center">Swipe left to close</div>
          <nav className="flex-1 overflow-y-auto p-4">
            <ul className="space-y-1">
              {sections.map((item) => (
                <li key={item.to}>
                  {item.children ? (
                    <>
                      <button
                        onClick={() => setMusicExpanded(!musicExpanded)}
                        className={`${itemClass(item.to)} w-full justify-between`}
                      >
                        <span className="flex items-center gap-4">
                          <span className="text-2xl">{item.icon}</span>
                          <span>{item.label}</span>
                        </span>
                        <span className="text-sm">{musicExpanded ? '▾' : '▸'}</span>
                      </button>
                      {musicExpanded && (
                        <ul className="ml-8 mt-1 space-y-1">
                          {item.children.map((child) => (
                            <li key={child.to}>
                              <Link to={child.to} onClick={onClose} className={itemClass(child.to)} aria-current={isActive(child.to) ? 'page' : undefined}>
                                <span className="text-xl">{child.icon}</span>
                                <span>{child.label}</span>
                              </Link>
                            </li>
                          ))}
                        </ul>
                      )}
                    </>
                  ) : (
                    <Link to={item.to} onClick={onClose} className={itemClass(item.to)} aria-current={isActive(item.to) ? 'page' : undefined}>
                      <span className="text-2xl">{item.icon}</span>
                      <span>{item.label}</span>
                    </Link>
                  )}
                </li>
              ))}
            </ul>
            {isAdmin && (
              <>
                <div className="divider text-xs text-base-content/50 my-2">Admin</div>
                <ul className="space-y-1">
                  <li>
                    <Link to="/admin/users" onClick={onClose} className={itemClass('/admin/users')} aria-current={isActive('/admin/users') ? 'page' : undefined}>
                      <span className="text-2xl">👥</span>
                      <span>User Management</span>
                    </Link>
                  </li>
                </ul>
              </>
            )}
          </nav>

          {currentTrack && (
            <div className="border-t border-base-300 p-4 bg-base-200">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 bg-base-300 rounded-lg overflow-hidden flex-shrink-0 flex items-center justify-center text-2xl">
                  {currentTrack.coverArtUrl ? <img src={currentTrack.coverArtUrl} alt="" className="w-full h-full object-cover" /> : '🎵'}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{currentTrack.title}</p>
                  <p className="text-sm text-base-content/60 truncate">{currentTrack.artist}</p>
                </div>
                <button className="btn btn-circle btn-primary" onClick={() => (isPlaying ? pause() : play())} aria-label={isPlaying ? 'Pause' : 'Play'}>
                  {isPlaying ? (
                    <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"><path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" /></svg>
                  ) : (
                    <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
                  )}
                </button>
              </div>
              {isPlaying && (
                <div className="mt-2 flex items-center gap-2 text-xs text-primary">
                  <span className="animate-pulse">▶</span>
                  <span>Now playing</span>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default MobileNav;
