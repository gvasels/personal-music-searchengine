/**
 * Search Page - Wave 4
 */
import { useState, useEffect } from 'react';
import { useNavigate, useSearch as useRouterSearch } from '@tanstack/react-router';
import { useSearchQuery } from '../hooks/useSearch';
import type { Track, Playlist } from '../types';

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export default function SearchPage() {
  const navigate = useNavigate();
  const routerSearch = useRouterSearch({ from: '/search' });
  const query = (routerSearch as { q?: string }).q || '';
  const artistFilter = (routerSearch as { artist?: string }).artist || '';
  const albumFilter = (routerSearch as { album?: string }).album || '';

  const [artistInput, setArtistInput] = useState(artistFilter);
  const [albumInput, setAlbumInput] = useState(albumFilter);

  const { data, isLoading, isError, error } = useSearchQuery({
    query,
    artist: artistFilter || undefined,
    album: albumFilter || undefined,
  });

  // Update filters after debounce
  useEffect(() => {
    const timer = setTimeout(() => {
      if (artistInput !== artistFilter || albumInput !== albumFilter) {
        void navigate({
          to: '/search',
          search: { q: query, artist: artistInput, album: albumInput },
        });
      }
    }, 300);
    return () => clearTimeout(timer);
  }, [artistInput, albumInput, artistFilter, albumFilter, query, navigate]);

  const handleTrackClick = (track: Track) => {
    void navigate({
      to: '/tracks/$trackId',
      params: { trackId: track.id },
    });
  };

  const handlePlaylistClick = (playlist: Playlist) => {
    void navigate({
      to: '/playlists/$playlistId',
      params: { playlistId: playlist.id },
    });
  };

  const handleArtistClick = (e: React.MouseEvent, artist: string) => {
    e.stopPropagation();
    void navigate({
      to: '/artists/$artistName',
      params: { artistName: artist },
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
        <span>Search failed: {error?.message}</span>
      </div>
    );
  }

  const trackCount = data?.tracks?.length || 0;
  const playlistCount = data?.playlists?.length || 0;
  const totalCount = data?.totalResults || trackCount;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Search results for "{query}"</h1>
        <p className="text-base-content/60">{totalCount} {totalCount === 1 ? 'result' : 'results'}</p>
      </div>

      {/* Filters */}
      <div className="flex gap-4">
        <div className="form-control">
          <label className="label" htmlFor="artist-filter">
            <span className="label-text">Artist</span>
          </label>
          <input
            id="artist-filter"
            type="text"
            placeholder="Filter by artist"
            className="input input-bordered input-sm"
            value={artistInput}
            onChange={(e) => setArtistInput(e.target.value)}
          />
        </div>
        <div className="form-control">
          <label className="label" htmlFor="album-filter">
            <span className="label-text">Album</span>
          </label>
          <input
            id="album-filter"
            type="text"
            placeholder="Filter by album"
            className="input input-bordered input-sm"
            value={albumInput}
            onChange={(e) => setAlbumInput(e.target.value)}
          />
        </div>
      </div>

      {/* Results */}
      {trackCount === 0 && playlistCount === 0 ? (
        <div className="text-center py-12">
          <p className="text-base-content/60">No results found for "{query}"</p>
        </div>
      ) : (
        <>
          {/* Playlists Section */}
          {playlistCount > 0 && (
            <div className="space-y-3">
              <h2 className="text-lg font-semibold">Playlists</h2>
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                {data?.playlists?.map((playlist) => (
                  <div
                    key={playlist.id}
                    className="card bg-base-200 hover:bg-base-300 cursor-pointer transition-colors"
                    onClick={() => handlePlaylistClick(playlist)}
                  >
                    <div className="card-body p-4">
                      <div className="w-full aspect-square bg-base-300 rounded-lg flex items-center justify-center mb-2">
                        {playlist.coverArt ? (
                          <img
                            src={playlist.coverArt}
                            alt={playlist.name}
                            className="w-full h-full object-cover rounded-lg"
                          />
                        ) : (
                          <svg className="w-12 h-12 text-base-content/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
                          </svg>
                        )}
                      </div>
                      <h3 className="font-medium truncate">{playlist.name}</h3>
                      <p className="text-sm text-base-content/60">{playlist.trackCount} tracks</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Tracks Section */}
          {trackCount > 0 && (
            <div className="space-y-3">
              <h2 className="text-lg font-semibold">Tracks</h2>
              <div className="overflow-x-auto">
                <table className="table table-zebra">
                  <thead>
                    <tr>
                      <th>Title</th>
                      <th>Artist</th>
                      <th>Album</th>
                      <th>Duration</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data?.tracks?.map((track) => (
                      <tr
                        key={track.id}
                        className="hover cursor-pointer"
                        onClick={() => handleTrackClick(track)}
                      >
                        <td className="font-medium">{track.title}</td>
                        <td>
                          <button
                            className="text-base-content/70 hover:text-primary hover:underline"
                            onClick={(e) => handleArtistClick(e, track.artist)}
                          >
                            {track.artist}
                          </button>
                        </td>
                        <td>{track.album}</td>
                        <td>{formatDuration(track.duration)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
