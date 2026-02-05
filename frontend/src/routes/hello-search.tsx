import { createFileRoute } from '@tanstack/react-router';
import { useState } from 'react';
import { useHelloSearch, useFeaturedTracks } from '@/hooks/useHelloSearch';
import { SearchInput } from '@/components/hello/SearchInput';
import { TrackCard } from '@/components/hello/TrackCard';
import { TrackCardSkeleton } from '@/components/hello/TrackCardSkeleton';

export const Route = createFileRoute('/hello-search')({
  component: HelloSearchPage,
});

export default function HelloSearchPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const searchResults = useHelloSearch(searchQuery);
  const featured = useFeaturedTracks();

  const isSearching = searchQuery.length > 0;
  const activeQuery = isSearching ? searchResults : featured;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-4xl font-bold text-primary mb-6">Hello Music Search</h1>

      <SearchInput value={searchQuery} onChange={setSearchQuery} />

      {activeQuery.isLoading && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
          {[...Array(6)].map((_, i) => (
            <TrackCardSkeleton key={i} />
          ))}
        </div>
      )}

      {activeQuery.isError && (
        <div className="alert alert-error mt-4">Error loading tracks</div>
      )}

      {activeQuery.data?.tracks?.map((track) => (
        <TrackCard key={track.id} track={track} />
      ))}
    </div>
  );
}
