/**
 * Tracks List Page - Wave 2
 */
import { useNavigate } from '@tanstack/react-router';
import { useTracksQuery } from '../../hooks/useTracks';
import type { Track } from '../../types';

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export default function TracksPage() {
  const navigate = useNavigate();
  const { data, isLoading, isError, error, refetch } = useTracksQuery();

  const handleTrackClick = (track: Track) => {
    void navigate({
      to: '/tracks/$trackId',
      params: { trackId: track.id },
    });
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    void navigate({
      to: '/tracks',
      search: { sortBy: e.target.value },
    });
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="alert alert-error">
        <span>Failed to load tracks: {error?.message}</span>
        <button className="btn btn-sm" onClick={() => void refetch()}>
          Retry
        </button>
      </div>
    );
  }

  if (!data?.items.length) {
    return (
      <div className="text-center py-12">
        <p className="text-base-content/60">No tracks found</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{data.total} tracks</h1>
        <label className="form-control w-auto">
          <select
            className="select select-bordered select-sm"
            onChange={handleSortChange}
            aria-label="Sort by"
          >
            <option value="">Sort by...</option>
            <option value="title">Title</option>
            <option value="artist">Artist</option>
            <option value="album">Album</option>
            <option value="createdAt">Date Added</option>
          </select>
        </label>
      </div>

      <div className="overflow-x-auto">
        <table className="table table-zebra">
          <thead>
            <tr>
              <th>Title</th>
              <th>Artist</th>
              <th>Album</th>
              <th>Duration</th>
              <th>Format</th>
            </tr>
          </thead>
          <tbody>
            {data.items.map((track) => (
              <tr
                key={track.id}
                className="hover cursor-pointer"
                onClick={() => handleTrackClick(track)}
              >
                <td className="font-medium">{track.title}</td>
                <td>{track.artist}</td>
                <td>{track.album}</td>
                <td>{formatDuration(track.duration)}</td>
                <td className="uppercase text-xs">{track.format}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
