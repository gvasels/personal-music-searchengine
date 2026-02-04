/**
 * Tags Page - Wave 5
 */
import { useNavigate } from '@tanstack/react-router';
import { useTagsQuery } from '../../hooks/useTags';

export default function TagsPage() {
  const navigate = useNavigate();
  const { data, isLoading, isError, error } = useTagsQuery();

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-64">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="alert alert-error">
        <span>{error?.message || 'Failed to load tags'}</span>
      </div>
    );
  }

  const tags = data?.items || [];

  // Calculate tag sizes for tag cloud effect
  const maxCount = Math.max(...tags.map((t) => t.trackCount), 1);
  const minCount = Math.min(...tags.map((t) => t.trackCount), 0);

  const getTagSize = (count: number) => {
    if (maxCount === minCount) return 'text-base';
    const ratio = (count - minCount) / (maxCount - minCount);
    if (ratio < 0.25) return 'text-sm';
    if (ratio < 0.5) return 'text-base';
    if (ratio < 0.75) return 'text-lg';
    return 'text-xl font-semibold';
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Tags</h1>

      {tags.length === 0 ? (
        <div className="text-center py-12 text-base-content/60">
          <p>No tags found</p>
          <p className="text-sm mt-2">Add tags to your tracks to organize your library</p>
        </div>
      ) : (
        <div className="flex flex-wrap gap-3">
          {tags.map((tag) => (
            <button
              key={tag.name}
              className={`btn btn-outline ${getTagSize(tag.trackCount)}`}
              onClick={() =>
                navigate({
                  to: '/tags/$tagName',
                  params: { tagName: tag.name },
                })
              }
            >
              <span>{tag.name}</span>
              <span className="badge badge-sm ml-1">{tag.trackCount}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
