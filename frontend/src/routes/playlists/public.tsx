/**
 * Public Playlists Discovery Page - Global User Type Feature
 *
 * Lists public playlists for discovery. No authentication required.
 */
import { Link } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import { getPublicPlaylists } from '../../lib/api/playlists';
import { VisibilityBadge } from '../../components/playlist';
import type { Playlist } from '../../types';

export default function PublicPlaylistsPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['playlists', 'public'],
    queryFn: () => getPublicPlaylists({ limit: 50 }),
  });

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <span>Failed to load public playlists</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Discover Playlists</h1>
      </div>

      {!data?.items?.length ? (
        <div className="text-center py-12 text-base-content/60">
          <p>No public playlists yet</p>
          <p className="text-sm mt-2">Be the first to share a playlist!</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {data.items.map((playlist: Playlist) => (
            <Link
              key={playlist.id}
              to="/playlists/$playlistId"
              params={{ playlistId: playlist.id }}
              className="card bg-base-100 shadow hover:shadow-lg transition-shadow"
            >
              <figure className="h-32 bg-gradient-to-br from-primary/20 to-secondary/20">
                {playlist.coverArt && (
                  <img
                    src={playlist.coverArt}
                    alt={playlist.name}
                    className="w-full h-full object-cover"
                  />
                )}
              </figure>
              <div className="card-body p-4">
                <h3 className="card-title text-base">{playlist.name}</h3>
                {playlist.description && (
                  <p className="text-sm text-base-content/60 line-clamp-2">
                    {playlist.description}
                  </p>
                )}
                <div className="flex items-center justify-between mt-2">
                  <span className="text-sm text-base-content/60">
                    {playlist.trackCount} {playlist.trackCount === 1 ? 'track' : 'tracks'}
                  </span>
                  <VisibilityBadge visibility={playlist.visibility} size="sm" />
                </div>
                {playlist.ownerDisplayName && (
                  <div className="text-xs text-base-content/50 mt-1">
                    by {playlist.ownerDisplayName}
                  </div>
                )}
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
