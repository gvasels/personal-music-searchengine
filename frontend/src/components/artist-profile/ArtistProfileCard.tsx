/**
 * ArtistProfileCard Component - Display an artist profile card
 */
import { useNavigate } from '@tanstack/react-router';
import { FollowButton } from '../follow/FollowButton';
import type { ArtistProfile } from '../../types';

interface ArtistProfileCardProps {
  profile: ArtistProfile;
  showFollowButton?: boolean;
}

export function ArtistProfileCard({ profile, showFollowButton = true }: ArtistProfileCardProps) {
  const navigate = useNavigate();

  const handleClick = () => {
    void navigate({ to: `/artists/entity/${profile.userId}` as string });
  };

  return (
    <div className="card bg-base-100 shadow-md hover:shadow-lg transition-shadow">
      {/* Header Image */}
      <figure className="relative h-24 bg-gradient-to-r from-primary/20 to-secondary/20">
        {profile.headerImageUrl && (
          <img
            src={profile.headerImageUrl}
            alt=""
            className="w-full h-full object-cover"
          />
        )}
        {/* Avatar overlapping header */}
        <div className="absolute -bottom-8 left-4">
          <div className="avatar">
            <div className="w-16 h-16 rounded-full ring ring-base-100 ring-offset-0 bg-base-300">
              {profile.avatarUrl ? (
                <img src={profile.avatarUrl} alt={profile.displayName} />
              ) : (
                <div className="flex items-center justify-center w-full h-full text-2xl font-bold bg-primary text-primary-content">
                  {profile.displayName.charAt(0).toUpperCase()}
                </div>
              )}
            </div>
          </div>
        </div>
      </figure>

      <div className="card-body pt-10">
        <div
          onClick={handleClick}
          className="hover:underline cursor-pointer"
        >
          <h3 className="card-title text-lg flex items-center gap-2">
            {profile.displayName}
            {profile.isVerified && (
              <span className="badge badge-primary badge-sm">Verified</span>
            )}
          </h3>
        </div>

        {profile.bio && (
          <p className="text-sm text-base-content/70 line-clamp-2">{profile.bio}</p>
        )}

        {profile.location && (
          <div className="text-xs text-base-content/50 flex items-center gap-1">
            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            {profile.location}
          </div>
        )}

        <div className="flex gap-4 text-sm text-base-content/60 mt-2">
          <span>
            <strong className="text-base-content">{profile.followerCount}</strong> followers
          </span>
          <span>
            <strong className="text-base-content">{profile.trackCount}</strong> tracks
          </span>
        </div>

        {showFollowButton && (
          <div className="card-actions justify-end mt-2">
            <FollowButton artistId={profile.userId} size="sm" />
          </div>
        )}
      </div>
    </div>
  );
}
