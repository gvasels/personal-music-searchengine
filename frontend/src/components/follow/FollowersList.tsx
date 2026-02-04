/**
 * FollowersList Component - Display followers of an artist
 */
import { useNavigate } from '@tanstack/react-router';
import { useFollowers } from '../../hooks/useFollow';
import type { ArtistProfile } from '../../types';

interface FollowersListProps {
  artistId: string;
  limit?: number;
}

export function FollowersList({ artistId, limit = 20 }: FollowersListProps) {
  const navigate = useNavigate();
  const { data, isLoading, error } = useFollowers(artistId, { limit });

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
        <span>Failed to load followers</span>
      </div>
    );
  }

  if (!data?.items?.length) {
    return (
      <div className="text-center py-8 text-base-content/60">
        <p>No followers yet</p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <div className="text-sm text-base-content/60 mb-4">
        {data.totalCount} {data.totalCount === 1 ? 'follower' : 'followers'}
      </div>
      <div className="grid gap-2">
        {data.items.map((profile: ArtistProfile) => (
          <div
            key={profile.userId}
            onClick={() => void navigate({ to: `/artists/entity/${profile.userId}` as string })}
            className="flex items-center gap-3 p-3 rounded-lg hover:bg-base-200 transition-colors cursor-pointer"
          >
            <div className="avatar">
              <div className="w-10 h-10 rounded-full bg-base-300">
                {profile.avatarUrl ? (
                  <img src={profile.avatarUrl} alt={profile.displayName} />
                ) : (
                  <div className="flex items-center justify-center w-full h-full text-lg font-bold">
                    {profile.displayName.charAt(0).toUpperCase()}
                  </div>
                )}
              </div>
            </div>
            <div className="flex-1 min-w-0">
              <div className="font-medium truncate">{profile.displayName}</div>
              {profile.location && (
                <div className="text-sm text-base-content/60 truncate">{profile.location}</div>
              )}
            </div>
            {profile.isVerified && (
              <div className="badge badge-primary badge-sm">Verified</div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
