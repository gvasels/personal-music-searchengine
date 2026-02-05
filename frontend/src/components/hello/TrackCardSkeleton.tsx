export function TrackCardSkeleton() {
  return (
    <div className="card bg-base-200">
      <div className="card-body">
        <div className="skeleton h-6 w-3/4"></div>
        <div className="skeleton h-4 w-1/2"></div>
        <div className="skeleton h-4 w-1/3"></div>
      </div>
    </div>
  );
}
