/**
 * FollowButton Component - Follow/Unfollow an artist
 */
import { useFollowToggle } from '../../hooks/useFollow';
import { useAuth } from '../../hooks/useAuth';

interface FollowButtonProps {
  artistId: string;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

export function FollowButton({ artistId, className = '', size = 'md' }: FollowButtonProps) {
  const { isAuthenticated, isSubscriber } = useAuth();
  const { isFollowing, isLoading, isToggling, toggle } = useFollowToggle(artistId);

  // Must be authenticated and at least subscriber to follow
  if (!isAuthenticated || !isSubscriber) {
    return null;
  }

  const sizeClasses = {
    sm: 'btn-sm',
    md: '',
    lg: 'btn-lg',
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
      className={`btn ${sizeClasses[size]} ${isFollowing ? 'btn-outline' : 'btn-primary'} ${className}`}
      onClick={toggle}
      disabled={isToggling}
    >
      {isToggling ? (
        <span className="loading loading-spinner loading-xs"></span>
      ) : isFollowing ? (
        'Following'
      ) : (
        'Follow'
      )}
    </button>
  );
}
