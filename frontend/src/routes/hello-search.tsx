/**
 * Hello Search Page
 *
 * Full-stack demonstration page for the Hello World local dev feature.
 * Displays featured tracks on load and supports searching.
 */
import { createFileRoute } from '@tanstack/react-router';
import { useState } from 'react';
import { useHelloSearch, useHelloFeatured } from '../hooks/useHelloSearch';
import { HelloNav } from '../components/hello/HelloNav';
import { SearchInput } from '../components/hello/SearchInput';
import { TrackCard } from '../components/hello/TrackCard';
import { TrackCardSkeleton } from '../components/hello/TrackCardSkeleton';

function HelloSearchPage() {
  const [query, setQuery] = useState('');
  const searchResult = useHelloSearch(query);
  const featuredResult = useHelloFeatured();

  // Check search state first - if searching, use search results
  // Otherwise fall back to featured tracks
  const isSearching = query.length > 0 || searchResult.isLoading || searchResult.isError || searchResult.data !== undefined;
  const isLoading = isSearching ? searchResult.isLoading : featuredResult.isLoading;
  const isError = isSearching ? searchResult.isError : featuredResult.isError;
  const data = isSearching ? searchResult.data : featuredResult.data;

  // Show "No results found" when search returns empty items
  const showNoResults = !isLoading && !isError && data?.items.length === 0;

  return (
    <div className="min-h-screen bg-base-100">
      <HelloNav />
      <div className="container mx-auto p-4">
        <SearchInput value={query} onChange={setQuery} placeholder="Search tracks..." />

        {isLoading && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
            {[...Array(6)].map((_, i) => (
              <TrackCardSkeleton key={i} />
            ))}
          </div>
        )}

        {isError && (
          <div role="alert" className="alert alert-error mt-4">
            <span>Error loading tracks</span>
          </div>
        )}

        {showNoResults && (
          <div className="text-center mt-4">No results found</div>
        )}

        {!isLoading && !isError && data?.items && data.items.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
            {data.items.map((track) => (
              <TrackCard key={track.id} track={track} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export const Route = createFileRoute('/hello-search')({
  component: HelloSearchPage,
});

export default HelloSearchPage;
