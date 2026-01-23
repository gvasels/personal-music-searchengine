/**
 * Mobile Navigation Overlay
 * Full-screen navigation for mobile devices with swipe-to-close gesture
 */
import { useEffect, useRef, useCallback, useState } from 'react';
import { Link, useLocation } from '@tanstack/react-router';
import { usePlayerStore } from '@/lib/store/playerStore';

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

interface MobileNavProps {
  isOpen: boolean;
  onClose: () => void;
}

export function MobileNav({ isOpen, onClose }: MobileNavProps) {
  const location = useLocation();
  const navRef = useRef<HTMLDivElement>(null);
  const { currentTrack, isPlaying, play, pause } = usePlayerStore();

  // Touch handling for swipe-to-close
  const [touchStart, setTouchStart] = useState<number | null>(null);
  const [touchCurrent, setTouchCurrent] = useState<number | null>(null);
  const [translateX, setTranslateX] = useState(0);

  // Prevent body scroll when nav is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  // Close on escape key
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
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

    // Only allow swiping left (negative direction)
    const diff = currentX - touchStart;
    if (diff < 0) {
      setTranslateX(diff);
    }
  }, [touchStart]);

  const handleTouchEnd = useCallback(() => {
    if (touchStart === null || touchCurrent === null) return;

    const diff = touchCurrent - touchStart;
    const threshold = -100; // Close if swiped left more than 100px

    if (diff < threshold) {
      onClose();
    }

    setTouchStart(null);
    setTouchCurrent(null);
    setTranslateX(0);
  }, [touchStart, touchCurrent, onClose]);

  const formatDuration = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 md:hidden">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200"
        onClick={onClose}
      />

      {/* Navigation panel */}
      <div
        ref={navRef}
        className="absolute inset-y-0 left-0 w-full max-w-sm bg-base-100 shadow-2xl animate-in slide-in-from-left duration-200"
        style={{
          transform: `translateX(${translateX}px)`,
          transition: touchStart === null ? 'transform 0.2s ease-out' : 'none',
        }}
        onTouchStart={handleTouchStart}
        onTouchMove={handleTouchMove}
        onTouchEnd={handleTouchEnd}
      >
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-base-300">
            <h2 className="text-lg font-bold">Menu</h2>
            <button
              onClick={onClose}
              className="btn btn-ghost btn-circle btn-sm"
              aria-label="Close menu"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {/* Swipe hint */}
          <div className="px-4 py-2 text-xs text-base-content/50 text-center">
            Swipe left to close
          </div>

          {/* Navigation items */}
          <nav className="flex-1 overflow-y-auto p-4">
            <ul className="space-y-1">
              {navItems.map((item) => {
                const isActive = location.pathname === item.to ||
                  (item.to !== '/' && location.pathname.startsWith(item.to));

                return (
                  <li key={item.to}>
                    <Link
                      to={item.to}
                      onClick={onClose}
                      className={`
                        flex items-center gap-4 p-4 rounded-xl
                        text-lg font-medium transition-all
                        ${isActive
                          ? 'bg-primary/20 text-primary border-l-4 border-primary'
                          : 'hover:bg-base-200 active:bg-base-300'
                        }
                      `}
                      aria-current={isActive ? 'page' : undefined}
                    >
                      <span className="text-2xl">{item.icon}</span>
                      <span>{item.label}</span>
                    </Link>
                  </li>
                );
              })}
            </ul>
          </nav>

          {/* Mini player at bottom */}
          {currentTrack && (
            <div className="border-t border-base-300 p-4 bg-base-200">
              <div className="flex items-center gap-3">
                {/* Cover art */}
                <div className="w-12 h-12 bg-base-300 rounded-lg overflow-hidden flex-shrink-0">
                  {currentTrack.coverArtUrl ? (
                    <img
                      src={currentTrack.coverArtUrl}
                      alt=""
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-2xl">
                      🎵
                    </div>
                  )}
                </div>

                {/* Track info */}
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{currentTrack.title}</p>
                  <p className="text-sm text-base-content/60 truncate">
                    {currentTrack.artist}
                  </p>
                </div>

                {/* Play/pause button */}
                <button
                  className="btn btn-circle btn-primary"
                  onClick={() => (isPlaying ? pause() : play())}
                  aria-label={isPlaying ? 'Pause' : 'Play'}
                >
                  {isPlaying ? (
                    <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                    </svg>
                  ) : (
                    <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M8 5v14l11-7z" />
                    </svg>
                  )}
                </button>
              </div>

              {/* Now playing indicator */}
              {isPlaying && (
                <div className="mt-2 flex items-center gap-2 text-xs text-primary">
                  <span className="animate-pulse">▶</span>
                  <span>Now playing</span>
                  <span className="text-base-content/50">
                    {formatDuration(currentTrack.duration)}
                  </span>
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
