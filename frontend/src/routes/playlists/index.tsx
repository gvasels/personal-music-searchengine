/**
 * Playlists Page - Wave 5
 */
import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { usePlaylistsQuery } from '../../hooks/usePlaylists';
import { CreatePlaylistModal } from '../../components/playlist';

export default function PlaylistsPage() {
  const navigate = useNavigate();
  const { data, isLoading, isError, error } = usePlaylistsQuery();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-64">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="alert alert-error">
        <span>{error?.message || 'Failed to load playlists'}</span>
      </div>
    );
  }

  const playlists = data?.items || [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Playlists</h1>
        <button
          className="btn btn-primary"
          onClick={() => setIsCreateModalOpen(true)}
        >
          Create Playlist
        </button>
      </div>

      {playlists.length === 0 ? (
        <div className="text-center py-12 text-base-content/60">
          <p>No playlists found</p>
          <p className="text-sm mt-2">Create your first playlist to get started</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {playlists.map((playlist) => (
            <div
              key={playlist.id}
              className="card bg-base-200 hover:bg-base-300 cursor-pointer transition-colors"
              onClick={() =>
                navigate({
                  to: '/playlists/$playlistId',
                  params: { playlistId: playlist.id },
                })
              }
            >
              <div className="card-body">
                <h2 className="card-title">{playlist.name}</h2>
                {playlist.description && (
                  <p className="text-base-content/70 line-clamp-2">{playlist.description}</p>
                )}
                <div className="flex items-center gap-2 text-sm text-base-content/60">
                  <span>{playlist.trackCount} tracks</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <CreatePlaylistModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
      />
    </div>
  );
}
