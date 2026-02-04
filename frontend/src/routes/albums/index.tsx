/**
 * Albums Grid Page - Wave 2
 */
import { useNavigate } from '@tanstack/react-router';
import { useAlbumsQuery } from '../../hooks/useAlbums';
import type { Album } from '../../types';

export default function AlbumsPage() {
  const navigate = useNavigate();
  const { data, isLoading, isError, error, refetch } = useAlbumsQuery();

  const handleAlbumClick = (album: Album) => {
    void navigate({
      to: '/albums/$albumId',
      params: { albumId: album.id },
    });
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    void navigate({
      to: '/albums',
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
        <span>Failed to load albums: {error?.message}</span>
        <button className="btn btn-sm" onClick={() => void refetch()}>
          Retry
        </button>
      </div>
    );
  }

  if (!data?.items.length) {
    return (
      <div className="text-center py-12">
        <p className="text-base-content/60">No albums found</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{data.total} albums</h1>
        <label className="form-control w-auto">
          <select
            className="select select-bordered select-sm"
            onChange={handleSortChange}
            aria-label="Sort by"
          >
            <option value="">Sort by...</option>
            <option value="title">Title</option>
            <option value="artist">Artist</option>
            <option value="year">Year</option>
            <option value="createdAt">Date Added</option>
          </select>
        </label>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
        {data.items.map((album) => (
          <div
            key={album.id}
            className="card bg-base-200 cursor-pointer hover:shadow-lg transition-shadow"
            onClick={() => handleAlbumClick(album)}
          >
            <figure className="aspect-square bg-base-300">
              {album.coverArt ? (
                <img src={album.coverArt} alt={album.name} className="w-full h-full object-cover" />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-4xl">🎵</div>
              )}
            </figure>
            <div className="card-body p-3">
              <h2 className="card-title text-sm truncate">{album.name}</h2>
              <p className="text-sm text-base-content/60 truncate">{album.artist}</p>
              <div className="flex justify-between text-xs text-base-content/40">
                <span>{album.year}</span>
                <span>{album.trackCount} tracks</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
