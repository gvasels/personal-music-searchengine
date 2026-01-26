/**
 * FollowingList Component - Display artists the current user is following
 */
import { useNavigate } from '@tanstack/react-router';
import { useFollowing } from '../../hooks/useFollow';
import { FollowButton } from './FollowButton';
import type { ArtistProfile } from '../../types';

interface FollowingListProps {
  limit?: number;
}

export function FollowingList({ limit = 20 }: FollowingListProps) {
  const navigate = useNavigate();
  const { data, isLoading, error } = useFollowing({ limit });

  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <span>Failed to load following</span>
      </div>
    );
  }

  if (!data?.items?.length) {
    return (
      <div className="text-center py-8 text-base-content/60">
        <p>You're not following anyone yet</p>
        <p className="text-sm mt-2">Discover artists to follow</p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <div className="text-sm text-base-content/60 mb-4">
        Following {data.totalCount} {data.totalCount === 1 ? 'artist' : 'artists'}
      </div>
      <div className="grid gap-2">
        {data.items.map((profile: ArtistProfile) => (
          <div
            key={profile.userId}
            className="flex items-center gap-3 p-3 rounded-lg bg-base-100 border border-base-200"
          >
            <div
              onClick={() => void navigate({ to: `/artists/entity/${profile.userId}` as string })}
              className="flex items-center gap-3 flex-1 min-w-0 cursor-pointer"
            >
              <div className="avatar">
                <div className="w-12 h-12 rounded-full bg-base-300">
                  {profile.avatarUrl ? (
                    <img src={profile.avatarUrl} alt={profile.displayName} />
                  ) : (
                    <div className="flex items-center justify-center w-full h-full text-xl font-bold">
                      {profile.displayName.charAt(0).toUpperCase()}
                    </div>
                  )}
                </div>
              </div>
              <div className="flex-1 min-w-0">
                <div className="font-medium truncate flex items-center gap-2">
                  {profile.displayName}
                  {profile.isVerified && (
                    <span className="badge badge-primary badge-xs">Verified</span>
                  )}
                </div>
                <div className="text-sm text-base-content/60">
                  {profile.followerCount} {profile.followerCount === 1 ? 'follower' : 'followers'}
                  {' · '}
                  {profile.trackCount} {profile.trackCount === 1 ? 'track' : 'tracks'}
                </div>
              </div>
            </div>
            <FollowButton artistId={profile.userId} size="sm" />
          </div>
        ))}
      </div>
    </div>
  );
}
