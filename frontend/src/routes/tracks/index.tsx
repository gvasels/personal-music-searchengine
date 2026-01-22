/**
 * Tracks List Page - Wave 2
 * Updated with play button and tags column
 */
import { useNavigate } from '@tanstack/react-router';
import { useTracksQuery } from '../../hooks/useTracks';
import { usePlayerStore } from '@/lib/store/playerStore';
import { TagsCell } from '@/components/library/TagsCell';
import type { Track } from '../../types';

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export default function TracksPage() {
  const navigate = useNavigate();
  const { data, isLoading, isError, error, refetch } = useTracksQuery();
  const { currentTrack, isPlaying, setQueue, pause } = usePlayerStore();

  const handleTrackClick = (track: Track) => {
    void navigate({
      to: '/tracks/$trackId',
      params: { trackId: track.id },
    });
  };

  const handlePlayClick = (e: React.MouseEvent, track: Track, index: number) => {
    e.stopPropagation(); // Don't navigate when clicking play

    if (currentTrack?.id === track.id) {
      // Toggle play/pause for current track
      if (isPlaying) {
        pause();
      } else {
        usePlayerStore.getState().play();
      }
    } else {
      // Play new track, queue all tracks starting from this one
      setQueue(data?.items || [], index);
    }
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
              <th className="w-12"></th>
              <th>Title</th>
              <th>Artist</th>
              <th>Album</th>
              <th>Tags</th>
              <th>Duration</th>
            </tr>
          </thead>
          <tbody>
            {data.items.map((track, index) => {
              const isCurrentTrack = currentTrack?.id === track.id;
              const isCurrentlyPlaying = isCurrentTrack && isPlaying;

              return (
                <tr
                  key={track.id}
                  className={`hover cursor-pointer ${isCurrentTrack ? 'bg-primary/10' : ''}`}
                  onClick={() => handleTrackClick(track)}
                >
                  <td className="w-12" onClick={(e) => e.stopPropagation()}>
                    <button
                      className={`btn btn-ghost btn-sm btn-circle ${isCurrentlyPlaying ? 'text-primary' : ''}`}
                      onClick={(e) => handlePlayClick(e, track, index)}
                      aria-label={isCurrentlyPlaying ? 'Pause' : 'Play'}
                    >
                      {isCurrentlyPlaying ? (
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                          <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                        </svg>
                      ) : isCurrentTrack ? (
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                          <path d="M8 5v14l11-7z" />
                        </svg>
                      ) : (
                        <svg className="w-5 h-5 opacity-50 group-hover:opacity-100" fill="currentColor" viewBox="0 0 24 24">
                          <path d="M8 5v14l11-7z" />
                        </svg>
                      )}
                    </button>
                  </td>
                  <td className="font-medium">{track.title}</td>
                  <td>{track.artist}</td>
                  <td>{track.album}</td>
                  <td onClick={(e) => e.stopPropagation()}>
                    <TagsCell trackId={track.id} tags={track.tags || []} maxVisible={2} />
                  </td>
                  <td>{formatDuration(track.duration)}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
