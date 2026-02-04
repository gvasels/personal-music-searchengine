/**
 * Album Detail Page
 * Shows album cover, metadata, and track list with play functionality
 */
import { useParams, Link } from '@tanstack/react-router';
import { useAlbumQuery } from '@/hooks/useAlbums';
import { usePlayerStore } from '@/lib/store/playerStore';
import { TagsCell } from '@/components/library/TagsCell';
import type { Track } from '@/types';

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

function formatTotalDuration(tracks: Track[]): string {
  const totalSeconds = tracks.reduce((sum, track) => sum + track.duration, 0);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);

  if (hours > 0) {
    return `${hours}h ${minutes}m`;
  }
  return `${minutes} min`;
}

export default function AlbumDetailPage() {
  const { albumId } = useParams({ strict: false });
  const { data: album, isLoading, isError, error, refetch } = useAlbumQuery(albumId);
  const { currentTrack, isPlaying, setQueue, pause } = usePlayerStore();

  const handlePlayAll = () => {
    if (album?.tracks) {
      setQueue(album.tracks, 0);
    }
  };

  const handlePlayTrack = (track: Track, index: number) => {
    if (currentTrack?.id === track.id) {
      if (isPlaying) {
        pause();
      } else {
        usePlayerStore.getState().play();
      }
    } else if (album?.tracks) {
      setQueue(album.tracks, index);
    }
  };

  if (isLoading) {
    return (
      <main className="p-4 md:p-8">
        <div className="mb-4">
          <Link to="/albums" className="link link-primary">
            ← Back to Albums
          </Link>
        </div>
        <div className="flex flex-col md:flex-row gap-8 animate-pulse">
          <div className="w-48 h-48 md:w-64 md:h-64 bg-base-300 rounded-lg flex-shrink-0" />
          <div className="flex-1 space-y-4">
            <div className="h-8 bg-base-300 rounded w-3/4" />
            <div className="h-6 bg-base-300 rounded w-1/2" />
            <div className="h-6 bg-base-300 rounded w-1/4" />
            <div className="h-10 bg-base-300 rounded w-32" />
          </div>
        </div>
        <div className="mt-8 space-y-2">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="h-12 bg-base-300 rounded animate-pulse" />
          ))}
        </div>
      </main>
    );
  }

  if (isError) {
    return (
      <main className="p-4 md:p-8">
        <div className="mb-4">
          <Link to="/albums" className="link link-primary">
            ← Back to Albums
          </Link>
        </div>
        <div className="alert alert-error">
          <span>Failed to load album: {error?.message}</span>
          <button className="btn btn-sm" onClick={() => void refetch()}>
            Retry
          </button>
        </div>
      </main>
    );
  }

  if (!album) {
    return (
      <main className="p-4 md:p-8">
        <div className="mb-4">
          <Link to="/albums" className="link link-primary">
            ← Back to Albums
          </Link>
        </div>
        <div className="text-center py-12">
          <p className="text-base-content/60">Album not found</p>
        </div>
      </main>
    );
  }

  const tracks = album.tracks || [];

  return (
    <main className="p-4 md:p-8">
      <div className="mb-4">
        <Link to="/albums" className="link link-primary">
          ← Back to Albums
        </Link>
      </div>

      {/* Album Header */}
      <div className="flex flex-col md:flex-row gap-6 md:gap-8 mb-8">
        {/* Cover Art */}
        <div className="flex-shrink-0">
          {album.coverArt ? (
            <img
              src={album.coverArt}
              alt={`${album.name} cover`}
              className="w-48 h-48 md:w-64 md:h-64 object-cover rounded-lg shadow-lg"
            />
          ) : (
            <div className="w-48 h-48 md:w-64 md:h-64 bg-base-300 rounded-lg shadow-lg flex items-center justify-center">
              <svg className="w-24 h-24 text-base-content/20" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z" />
              </svg>
            </div>
          )}
        </div>

        {/* Album Info */}
        <div className="flex flex-col justify-end">
          <p className="text-sm text-base-content/60 uppercase tracking-wide">Album</p>
          <h1 className="text-3xl md:text-4xl font-bold mt-1">{album.name}</h1>
          <div className="flex items-center gap-2 mt-2 text-base-content/80">
            <Link
              to="/artists/$artistName"
              params={{ artistName: album.artist }}
              className="font-medium hover:underline"
            >
              {album.artist}
            </Link>
            {album.year && (
              <>
                <span>·</span>
                <span>{album.year}</span>
              </>
            )}
            <span>·</span>
            <span>{tracks.length} tracks</span>
            {tracks.length > 0 && (
              <>
                <span>·</span>
                <span>{formatTotalDuration(tracks)}</span>
              </>
            )}
          </div>
          <div className="mt-4">
            <button
              className="btn btn-primary gap-2"
              onClick={handlePlayAll}
              disabled={tracks.length === 0}
            >
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z" />
              </svg>
              Play All
            </button>
          </div>
        </div>
      </div>

      {/* Track List */}
      {tracks.length === 0 ? (
        <div className="text-center py-8">
          <p className="text-base-content/60">No tracks in this album</p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="table table-zebra">
            <thead>
              <tr>
                <th className="w-12">#</th>
                <th className="w-12"></th>
                <th>Title</th>
                <th>Tags</th>
                <th className="text-right">Duration</th>
              </tr>
            </thead>
            <tbody>
              {tracks.map((track, index) => {
                const isCurrentTrack = currentTrack?.id === track.id;
                const isCurrentlyPlaying = isCurrentTrack && isPlaying;

                return (
                  <tr
                    key={track.id}
                    className={`hover ${isCurrentTrack ? 'bg-primary/10' : ''}`}
                  >
                    <td className="text-base-content/60">{track.trackNumber || index + 1}</td>
                    <td>
                      <button
                        className={`btn btn-ghost btn-sm btn-circle ${isCurrentlyPlaying ? 'text-primary' : ''}`}
                        onClick={() => handlePlayTrack(track, index)}
                        aria-label={isCurrentlyPlaying ? 'Pause' : 'Play'}
                      >
                        {isCurrentlyPlaying ? (
                          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                          </svg>
                        ) : (
                          <svg className="w-5 h-5 opacity-50 group-hover:opacity-100" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M8 5v14l11-7z" />
                          </svg>
                        )}
                      </button>
                    </td>
                    <td>
                      <Link
                        to="/tracks/$trackId"
                        params={{ trackId: track.id }}
                        className="font-medium hover:text-primary"
                      >
                        {track.title}
                      </Link>
                    </td>
                    <td onClick={(e) => e.stopPropagation()}>
                      <TagsCell trackId={track.id} tags={track.tags || []} maxVisible={2} />
                    </td>
                    <td className="text-right text-base-content/60">
                      {formatDuration(track.duration)}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </main>
  );
}
