/**
 * Artist Detail Page
 * Shows artist info, albums grid, and top tracks
 */
import { useParams, Link } from '@tanstack/react-router';
import { useArtistQuery } from '@/hooks/useArtists';
import { usePlayerStore } from '@/lib/store/playerStore';
import { TagsCell } from '@/components/library/TagsCell';
import type { Track } from '@/types';

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export default function ArtistDetailPage() {
  const { artistName } = useParams({ strict: false });
  const { data: artist, isLoading, isError, error, refetch } = useArtistQuery(artistName);
  const { currentTrack, isPlaying, setQueue, pause } = usePlayerStore();

  const handlePlayAll = () => {
    if (artist?.recentTracks) {
      setQueue(artist.recentTracks, 0);
    }
  };

  const handlePlayTrack = (track: Track, index: number) => {
    if (currentTrack?.id === track.id) {
      if (isPlaying) {
        pause();
      } else {
        usePlayerStore.getState().play();
      }
    } else if (artist?.recentTracks) {
      setQueue(artist.recentTracks, index);
    }
  };

  if (isLoading) {
    return (
      <main className="p-4 md:p-8">
        <div className="mb-4">
          <Link to="/artists" className="link link-primary">
            ← Back to Artists
          </Link>
        </div>
        <div className="animate-pulse space-y-8">
          <div className="flex flex-col gap-4">
            <div className="h-10 bg-base-300 rounded w-1/3" />
            <div className="h-6 bg-base-300 rounded w-1/4" />
            <div className="h-10 bg-base-300 rounded w-32" />
          </div>
          <div>
            <div className="h-6 bg-base-300 rounded w-24 mb-4" />
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="aspect-square bg-base-300 rounded-lg" />
              ))}
            </div>
          </div>
          <div>
            <div className="h-6 bg-base-300 rounded w-24 mb-4" />
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="h-12 bg-base-300 rounded mb-2" />
            ))}
          </div>
        </div>
      </main>
    );
  }

  if (isError) {
    return (
      <main className="p-4 md:p-8">
        <div className="mb-4">
          <Link to="/artists" className="link link-primary">
            ← Back to Artists
          </Link>
        </div>
        <div className="alert alert-error">
          <span>Failed to load artist: {error?.message}</span>
          <button className="btn btn-sm" onClick={() => void refetch()}>
            Retry
          </button>
        </div>
      </main>
    );
  }

  if (!artist) {
    return (
      <main className="p-4 md:p-8">
        <div className="mb-4">
          <Link to="/artists" className="link link-primary">
            ← Back to Artists
          </Link>
        </div>
        <div className="text-center py-12">
          <p className="text-base-content/60">Artist not found</p>
        </div>
      </main>
    );
  }

  const albums = artist.albums || [];
  const tracks = artist.recentTracks || [];

  return (
    <main className="p-4 md:p-8">
      <div className="mb-4">
        <Link to="/artists" className="link link-primary">
          ← Back to Artists
        </Link>
      </div>

      {/* Artist Header */}
      <div className="mb-8">
        <p className="text-sm text-base-content/60 uppercase tracking-wide">Artist</p>
        <h1 className="text-3xl md:text-4xl font-bold mt-1">{artist.name}</h1>
        <div className="flex items-center gap-4 mt-2 text-base-content/80">
          <span>{artist.trackCount} tracks</span>
          <span>·</span>
          <span>{artist.albumCount} albums</span>
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

      {/* Albums Section */}
      {albums.length > 0 && (
        <section className="mb-8">
          <h2 className="text-xl font-bold mb-4">Albums</h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
            {albums.map((album) => (
              <Link
                key={album.id}
                to="/albums/$albumId"
                params={{ albumId: album.id }}
                className="group"
              >
                <div className="aspect-square bg-base-300 rounded-lg overflow-hidden shadow-md transition-transform group-hover:scale-105">
                  {album.coverArt ? (
                    <img
                      src={album.coverArt}
                      alt={`${album.name} cover`}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center">
                      <svg className="w-12 h-12 text-base-content/20" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z" />
                      </svg>
                    </div>
                  )}
                </div>
                <h3 className="mt-2 font-medium truncate group-hover:text-primary">
                  {album.name}
                </h3>
                {album.year && (
                  <p className="text-sm text-base-content/60">{album.year}</p>
                )}
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Popular Tracks Section */}
      {tracks.length > 0 && (
        <section>
          <h2 className="text-xl font-bold mb-4">Popular Tracks</h2>
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th className="w-12">#</th>
                  <th className="w-12"></th>
                  <th>Title</th>
                  <th>Album</th>
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
                      <td className="text-base-content/60">{index + 1}</td>
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
                      <td>
                        {track.albumId ? (
                          <Link
                            to="/albums/$albumId"
                            params={{ albumId: track.albumId }}
                            className="hover:text-primary"
                          >
                            {track.album}
                          </Link>
                        ) : (
                          <span className="text-base-content/60">{track.album}</span>
                        )}
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
        </section>
      )}

      {/* Empty State */}
      {albums.length === 0 && tracks.length === 0 && (
        <div className="text-center py-12">
          <p className="text-base-content/60">No music found for this artist</p>
        </div>
      )}
    </main>
  );
}
