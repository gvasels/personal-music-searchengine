/**
 * My Shows Page - Watched artist events and event search
 */
import { useState, useEffect } from 'react';
import { createFileRoute } from '@tanstack/react-router';
import { useWatchedArtists } from '../hooks/useArtistWatch';
import { useSearchArtistEvents } from '../hooks/useArtistEvents';

function useDebouncedValue<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);
  useEffect(() => {
    const handler = setTimeout(() => setDebouncedValue(value), delay);
    return () => clearTimeout(handler);
  }, [value, delay]);
  return debouncedValue;
}

export const Route = createFileRoute('/shows')({
  component: ShowsPage,
});

export default function ShowsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const debouncedQuery = useDebouncedValue(searchQuery, 300);
  const { data, isLoading, isError } = useWatchedArtists();
  const searchResult = useSearchArtistEvents(debouncedQuery);

  if (isLoading) {
    return (
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-6">My Shows</h1>
        <div className="flex justify-center py-8">
          <span className="loading loading-spinner loading-lg" role="status"></span>
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-6">My Shows</h1>
        <div className="text-center py-8 text-error">Failed to load watched artists</div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">My Shows</h1>

      <div className="mb-6">
        <input
          type="search"
          role="searchbox"
          placeholder="Search for artists..."
          className="input input-bordered w-full max-w-md"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      {debouncedQuery && searchResult.data?.items && (
        <div className="mb-6">
          <h2 className="text-lg font-semibold mb-3">Search Results</h2>
          <div className="space-y-2">
            {searchResult.data.items.map((result) => (
              <div key={result.name} className="card bg-base-200 p-3">
                <div className="flex justify-between items-center">
                  <span className="font-medium">{result.name}</span>
                  <span className="text-sm text-base-content/60">
                    {result.upcomingEvents} upcoming events
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {(!data?.items || data.items.length === 0) ? (
        <div className="text-center py-8 text-base-content/60">
          No watched artists yet. Watch artists to see their upcoming events here.
        </div>
      ) : (
        <div className="space-y-3">
          <h2 className="text-lg font-semibold">Watched Artists</h2>
          {data.items.map((artist) => (
            <div key={artist.artistName} className="card bg-base-200 p-4">
              <span className="font-medium">{artist.artistName}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
