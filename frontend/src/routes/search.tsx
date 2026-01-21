/**
 * Search Page - Wave 4
 */
import { useState, useEffect } from 'react';
import { useNavigate, useSearch as useRouterSearch } from '@tanstack/react-router';
import { useSearchQuery } from '../hooks/useSearch';
import type { Track } from '../types';

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

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Search results for "{query}"</h1>
        <p className="text-base-content/60">{data?.total || 0} {data?.total === 1 ? 'result' : 'results'}</p>
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
      {!data?.items.length ? (
        <div className="text-center py-12">
          <p className="text-base-content/60">No results found for "{query}"</p>
        </div>
      ) : (
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
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
