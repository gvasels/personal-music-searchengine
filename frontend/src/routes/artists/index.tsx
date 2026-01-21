/**
 * Artists List Page - Wave 2
 */
import { useNavigate } from '@tanstack/react-router';
import { useArtistsQuery } from '../../hooks/useArtists';
import type { Artist } from '../../types';

export default function ArtistsPage() {
  const navigate = useNavigate();
  const { data, isLoading, isError, error, refetch } = useArtistsQuery();

  const handleArtistClick = (artist: Artist) => {
    void navigate({
      to: '/artists/$artistName',
      params: { artistName: artist.name },
    });
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    void navigate({
      to: '/artists',
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
        <span>Failed to load artists: {error?.message}</span>
        <button className="btn btn-sm" onClick={() => void refetch()}>
          Retry
        </button>
      </div>
    );
  }

  if (!data?.items.length) {
    return (
      <div className="text-center py-12">
        <p className="text-base-content/60">No artists found</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{data.total} artists</h1>
        <label className="form-control w-auto">
          <select
            className="select select-bordered select-sm"
            onChange={handleSortChange}
            aria-label="Sort by"
          >
            <option value="">Sort by...</option>
            <option value="name">Name</option>
            <option value="trackCount">Track Count</option>
            <option value="albumCount">Album Count</option>
          </select>
        </label>
      </div>

      <div className="space-y-2">
        {data.items.map((artist) => (
          <div
            key={artist.name}
            className="flex items-center justify-between p-4 bg-base-200 rounded-lg cursor-pointer hover:bg-base-300 transition-colors"
            onClick={() => handleArtistClick(artist)}
          >
            <div>
              <h3 className="font-semibold">{artist.name}</h3>
              <p className="text-sm text-base-content/60">
                {artist.trackCount} tracks · {artist.albumCount} {artist.albumCount === 1 ? 'album' : 'albums'}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
