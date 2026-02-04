/**
 * Artist Profiles Discovery Page - Global User Type Feature
 *
 * Lists artist profiles for discovery. Uses the /artists/entity endpoint.
 */
import { useState } from 'react';
import { useArtistProfiles, useSearchArtists } from '../../../hooks/useArtistProfiles';
import { ArtistProfileCard } from '../../../components/artist-profile';
import { useAuth } from '../../../hooks/useAuth';
import { EditArtistProfileModal } from '../../../components/artist-profile';

export default function ArtistProfilesPage() {
  const { isArtist, user } = useAuth();
  const [searchQuery, setSearchQuery] = useState('');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const { data: profilesData, isLoading: isLoadingProfiles } = useArtistProfiles({ limit: 20 });
  const { data: searchData, isLoading: isSearching } = useSearchArtists({
    query: searchQuery,
    limit: 10,
  });

  const isLoading = searchQuery ? isSearching : isLoadingProfiles;
  const profiles = searchQuery ? searchData?.items : profilesData?.items;

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <h1 className="text-2xl font-bold">Discover Artists</h1>
        {isArtist && (
          <button
            className="btn btn-primary"
            onClick={() => setIsCreateModalOpen(true)}
          >
            Create Artist Profile
          </button>
        )}
      </div>

      {/* Search */}
      <div className="form-control">
        <input
          type="text"
          placeholder="Search artists..."
          className="input input-bordered w-full max-w-md"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      {isLoading ? (
        <div className="flex justify-center py-12">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : !profiles?.length ? (
        <div className="text-center py-12 text-base-content/60">
          {searchQuery ? (
            <p>No artists found for "{searchQuery}"</p>
          ) : (
            <div>
              <p>No artist profiles yet</p>
              {isArtist && (
                <button
                  className="btn btn-primary mt-4"
                  onClick={() => setIsCreateModalOpen(true)}
                >
                  Be the first to create one
                </button>
              )}
            </div>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {profiles.map((profile) => (
            <ArtistProfileCard
              key={profile.userId}
              profile={profile}
              showFollowButton={profile.userId !== user?.userId}
            />
          ))}
        </div>
      )}

      <EditArtistProfileModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
      />
    </div>
  );
}
