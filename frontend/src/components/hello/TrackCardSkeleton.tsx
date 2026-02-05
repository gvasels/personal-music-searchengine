export function TrackCardSkeleton() {
  return (
    <div className="card bg-base-200">
      <div className="card-body">
        <div data-testid="skeleton-line" className="skeleton h-6 w-3/4" />
        <div data-testid="skeleton-line" className="skeleton h-4 w-1/2" />
        <div data-testid="skeleton-line" className="skeleton h-4 w-1/2" />
        <div data-testid="skeleton-line" className="skeleton h-4 w-1/4" />
        <div data-testid="skeleton-line" className="skeleton h-4 w-1/6" />
        <div data-testid="skeleton-line" className="skeleton h-4 w-1/6" />
      </div>
    </div>
  );
}
