/**
 * Tag Detail Page - Wave 5
 */
import { Link, useParams } from '@tanstack/react-router';
import { useTracksByTagQuery } from '../../hooks/useTags';
import { usePlayerStore } from '../../lib/store/playerStore';
import { TrackList } from '../../components/library';

export default function TagDetailPage() {
  const { tagName } = useParams({ strict: false });
  const { data, isLoading, isError, error } = useTracksByTagQuery(tagName || '');
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
        <span>{error?.message || 'Failed to load tracks'}</span>
      </div>
    );
  }

  const tracks = data?.items || [];

  const handlePlayAll = () => {
    if (tracks.length > 0) {
      setQueue(tracks, 0);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/tags" className="btn btn-ghost btn-sm">
          ← Back
        </Link>
      </div>

      <div className="flex flex-col md:flex-row md:items-end gap-6">
        <div className="w-24 h-24 bg-primary/20 rounded-lg flex items-center justify-center">
          <span className="text-4xl">🏷️</span>
        </div>

        <div className="space-y-2">
          <p className="text-sm uppercase tracking-wide text-base-content/60">Tag</p>
          <h1 className="text-4xl font-bold capitalize">{tagName}</h1>
          <p className="text-sm text-base-content/60">{data?.total || 0} tracks</p>
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
      </div>

      {tracks.length === 0 ? (
        <div className="text-center py-12 text-base-content/60">
          <p>No tracks with this tag</p>
        </div>
      ) : (
        <TrackList tracks={tracks} />
      )}
    </div>
  );
}
