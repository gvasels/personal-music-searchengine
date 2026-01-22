/**
 * Playlist Detail Page - Wave 5
 */
import { Link, useParams } from '@tanstack/react-router';
import { usePlaylistQuery } from '../../hooks/usePlaylists';
import { usePlayerStore } from '../../lib/store/playerStore';
import { TrackList } from '../../components/library';

export default function PlaylistDetailPage() {
  const { playlistId } = useParams({ strict: false });
  const { data: playlist, isLoading, isError, error } = usePlaylistQuery(playlistId);
  const { setQueue } = usePlayerStore();

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
        <span>{error?.message || 'Playlist not found'}</span>
      </div>
    );
  }

  if (!playlist) {
    return (
      <div className="alert alert-error">
        <span>Playlist not found</span>
      </div>
    );
  }

  const tracks = (playlist as any).tracks || [];

  const handlePlayAll = () => {
    if (tracks.length > 0) {
      setQueue(tracks, 0);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/playlists" className="btn btn-ghost btn-sm">
          ← Back
        </Link>
      </div>

      <div className="flex flex-col md:flex-row md:items-end gap-6">
        <div className="w-48 h-48 bg-base-300 rounded-lg flex items-center justify-center">
          <span className="text-6xl">🎵</span>
        </div>

        <div className="space-y-2">
          <p className="text-sm uppercase tracking-wide text-base-content/60">Playlist</p>
          <h1 className="text-4xl font-bold">{playlist.name}</h1>
          {playlist.description && (
            <p className="text-base-content/70">{playlist.description}</p>
          )}
          <p className="text-sm text-base-content/60">{playlist.trackCount} tracks</p>
        </div>
      </div>

      <div className="flex gap-2">
        <button
          className="btn btn-primary"
          onClick={handlePlayAll}
          disabled={tracks.length === 0}
        >
          Play All
        </button>
        <button className="btn btn-outline">Edit</button>
        <button className="btn btn-outline btn-error">Delete</button>
      </div>

      {tracks.length === 0 ? (
        <div className="text-center py-12 text-base-content/60">
          <p>No tracks in this playlist</p>
          <p className="text-sm mt-2">Add tracks from your library</p>
        </div>
      ) : (
        <TrackList tracks={tracks} showDownload showAddedDate />
      )}
    </div>
  );
}
