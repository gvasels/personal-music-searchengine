/**
 * WatchButton Component - Watch/Unwatch an artist for event notifications
 */
import { useWatchToggle } from '../../hooks/useArtistWatch';
import { useAuth } from '../../hooks/useAuth';

interface WatchButtonProps {
  artistName: string;
  className?: string;
  size?: 'sm' | 'md';
}

export function WatchButton({ artistName, className = '', size = 'md' }: WatchButtonProps) {
  const { isAuthenticated, isSubscriber } = useAuth();
  const { isWatching, isLoading, toggle } = useWatchToggle(artistName);

  if (!isAuthenticated || !isSubscriber) {
    return null;
  }

  const sizeClasses = {
    sm: 'btn-sm',
    md: '',
  };

  if (isLoading) {
    return (
      <button className={`btn ${sizeClasses[size]} ${className}`} disabled>
        <span className="loading loading-spinner loading-xs"></span>
      </button>
    );
  }

  return (
    <button
      className={`btn ${sizeClasses[size]} ${isWatching ? 'btn-outline' : 'btn-primary'} ${className}`}
      onClick={toggle}
    >
      {isWatching ? 'Watching' : 'Watch'}
    </button>
  );
}
