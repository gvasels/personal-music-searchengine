/**
 * Artist Profile Detail Page - Global User Type Feature
 *
 * Shows a single artist profile with their info, followers, and tracks.
 */
import { useState } from 'react';
import { useParams, useNavigate } from '@tanstack/react-router';
import { useArtistProfile } from '../../../hooks/useArtistProfiles';
import { FollowButton, FollowersList } from '../../../components/follow';
import { EditArtistProfileModal } from '../../../components/artist-profile';
import { useAuth } from '../../../hooks/useAuth';

export default function ArtistProfileDetailPage() {
  const params = useParams({ strict: false }) as { artistId?: string };
  const artistId = params.artistId || '';
  const navigate = useNavigate();
  const { user } = useAuth();
  const { data: profile, isLoading, error } = useArtistProfile(artistId);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<'tracks' | 'followers' | 'about'>('about');

  const isOwner = user?.userId === profile?.userId;

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error || !profile) {
    return (
      <div className="alert alert-error">
        <span>Failed to load artist profile</span>
        <button onClick={() => void navigate({ to: '/artists/entity' as string })} className="btn btn-sm">
          Back to Artists
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="relative">
        {/* Cover Image */}
        <div className="h-48 bg-gradient-to-r from-primary/20 to-secondary/20 rounded-lg overflow-hidden">
          {profile.headerImageUrl && (
            <img
              src={profile.headerImageUrl}
              alt=""
              className="w-full h-full object-cover"
            />
          )}
        </div>

        {/* Profile Info */}
        <div className="flex flex-col sm:flex-row gap-4 -mt-12 px-4">
          {/* Avatar */}
          <div className="avatar">
            <div className="w-24 h-24 sm:w-32 sm:h-32 rounded-full ring ring-base-100 ring-offset-2 bg-base-300">
              {profile.avatarUrl ? (
                <img src={profile.avatarUrl} alt={profile.displayName} />
              ) : (
                <div className="flex items-center justify-center w-full h-full text-4xl font-bold bg-primary text-primary-content">
                  {profile.displayName.charAt(0).toUpperCase()}
                </div>
              )}
            </div>
          </div>

          {/* Name and Actions */}
          <div className="flex-1 pt-14 sm:pt-16">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div>
                <h1 className="text-2xl sm:text-3xl font-bold flex items-center gap-2">
                  {profile.displayName}
                  {profile.isVerified && (
                    <span className="badge badge-primary">Verified</span>
                  )}
                </h1>
                {profile.location && (
                  <p className="text-base-content/60 flex items-center gap-1 mt-1">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    {profile.location}
                  </p>
                )}
              </div>

              <div className="flex gap-2">
                {isOwner ? (
                  <button
                    className="btn btn-outline"
                    onClick={() => setIsEditModalOpen(true)}
                  >
                    Edit Profile
                  </button>
                ) : (
                  <FollowButton artistId={artistId} />
                )}
              </div>
            </div>

            {/* Stats */}
            <div className="flex gap-6 mt-4 text-sm">
              <div>
                <span className="font-bold text-lg">{profile.followerCount}</span>
                <span className="text-base-content/60 ml-1">followers</span>
              </div>
              <div>
                <span className="font-bold text-lg">{profile.followingCount}</span>
                <span className="text-base-content/60 ml-1">following</span>
              </div>
              <div>
                <span className="font-bold text-lg">{profile.trackCount}</span>
                <span className="text-base-content/60 ml-1">tracks</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="tabs tabs-bordered">
        <button
          className={`tab ${activeTab === 'about' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('about')}
        >
          About
        </button>
        <button
          className={`tab ${activeTab === 'tracks' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('tracks')}
        >
          Tracks
        </button>
        <button
          className={`tab ${activeTab === 'followers' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('followers')}
        >
          Followers
        </button>
      </div>

      {/* Tab Content */}
      <div className="min-h-[200px]">
        {activeTab === 'about' && (
          <div className="space-y-4">
            {profile.bio ? (
              <p className="whitespace-pre-wrap">{profile.bio}</p>
            ) : (
              <p className="text-base-content/60 italic">No bio yet</p>
            )}

            {profile.website && (
              <div>
                <h3 className="font-semibold mb-2">Website</h3>
                <a
                  href={profile.website}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="link link-primary"
                >
                  {profile.website}
                </a>
              </div>
            )}

            {profile.socialLinks && Object.keys(profile.socialLinks).length > 0 && (
              <div>
                <h3 className="font-semibold mb-2">Social Links</h3>
                <div className="flex gap-2 flex-wrap">
                  {Object.entries(profile.socialLinks).map(([platform, url]) => (
                    <a
                      key={platform}
                      href={url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="btn btn-sm btn-outline"
                    >
                      {platform}
                    </a>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'tracks' && (
          <div className="text-center py-8 text-base-content/60">
            <p>Track listing coming soon</p>
          </div>
        )}

        {activeTab === 'followers' && <FollowersList artistId={artistId} />}
      </div>

      <EditArtistProfileModal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        profile={profile}
      />
    </div>
  );
}
