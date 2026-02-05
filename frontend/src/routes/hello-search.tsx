import { useState } from 'react';
import { createFileRoute } from '@tanstack/react-router';
import { useHelloSearch, useHelloFeatured } from '../hooks/useHelloSearch';
import { SearchInput } from '../components/hello/SearchInput';
import { TrackCard } from '../components/hello/TrackCard';
import { TrackCardSkeleton } from '../components/hello/TrackCardSkeleton';

function HelloSearchPage() {
  const [query, setQuery] = useState('');
  const featured = useHelloFeatured();
  const search = useHelloSearch(query);

  const isLoading = featured.isLoading || search.isLoading;
  const isError = featured.isError || search.isError;

  // Determine which tracks to display
  const displayData = query.length > 0 ? search.data : featured.data;

  return (
    <div className="container mx-auto p-4">
      <SearchInput
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search tracks..."
      />

      {isLoading && (
        <div className="grid gap-4 mt-4">
          <TrackCardSkeleton />
          <TrackCardSkeleton />
          <TrackCardSkeleton />
        </div>
      )}

      {isError && (
        <div role="alert" className="alert alert-error mt-4">
          An error occurred while loading tracks.
        </div>
      )}

      {search.data && search.data.items.length === 0 && (
        <p className="mt-4 text-center">No results found</p>
      )}

      {!isLoading && !isError && displayData && displayData.items.length > 0 && (
        <div className="grid gap-4 mt-4">
          {displayData.items.map((track) => (
            <TrackCard key={track.id} track={track} />
          ))}
        </div>
      )}
    </div>
  );
}

export const Route = createFileRoute('/hello-search')({
  component: HelloSearchPage,
});

export default HelloSearchPage;
