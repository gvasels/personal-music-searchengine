import { createFileRoute } from '@tanstack/react-router';
import { useState } from 'react';
import { useHelloSearch, useFeaturedTracks } from '../hooks/useHelloSearch';
import { SearchInput } from '../components/hello/SearchInput';
import { TrackCard } from '../components/hello/TrackCard';
import { TrackCardSkeleton } from '../components/hello/TrackCardSkeleton';
import { HelloNav } from '../components/hello/HelloNav';

export const Route = createFileRoute('/hello-search')({
  component: HelloSearchPage,
});

function HelloSearchPage() {
  const [query, setQuery] = useState('');
  const searchResult = useHelloSearch(query);
  const featuredResult = useFeaturedTracks();

  const isSearching = query.length > 0;
  const activeResult = isSearching ? searchResult : featuredResult;
  const tracks = isSearching
    ? (searchResult.data?.tracks ?? [])
    : (featuredResult.data?.tracks ?? []);
  const isLoading = activeResult.isLoading;
  const isError = activeResult.isError;

  return (
    <div className="min-h-screen bg-base-100">
      <HelloNav searchValue={query} onSearchChange={setQuery} />

      {/* Hero section */}
      <div className="hero min-h-[40vh] bg-base-200">
        <div className="hero-content text-center">
          <div className="max-w-xl">
            <h1 className="text-5xl font-bold">Music Search</h1>
            <p className="py-6 text-base-content/70">
              Find your next favorite track
            </p>
            <div className="flex justify-center sm:hidden">
              <SearchInput value={query} onChange={setQuery} autoFocus />
            </div>
          </div>
        </div>
      </div>

      {/* Results section */}
      <div className="container mx-auto px-4 py-8">
        {/* Result heading */}
        {isSearching && tracks.length > 0 && (
          <p className="text-sm text-base-content/60 mb-4">
            Found {tracks.length} tracks for &apos;{query}&apos;
          </p>
        )}
        {!isSearching && !isLoading && (
          <h2 className="text-2xl font-bold mb-4">Featured Tracks</h2>
        )}

        {/* Error state */}
        {isError && (
          <div className="alert alert-error">
            <span>Failed to load tracks. Please try again.</span>
          </div>
        )}

        {/* Loading skeletons */}
        {isLoading && (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {Array.from({ length: 8 }).map((_, i) => (
              <TrackCardSkeleton key={i} />
            ))}
          </div>
        )}

        {/* Track cards grid */}
        {!isLoading && tracks.length > 0 && (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {tracks.map((track) => (
              <TrackCard key={track.id} track={track} />
            ))}
          </div>
        )}

        {/* Empty state */}
        {!isLoading && isSearching && tracks.length === 0 && !isError && (
          <div className="text-center py-12">
            <p className="text-base-content/60">
              No tracks found for &apos;{query}&apos;
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
