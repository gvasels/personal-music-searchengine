/**
 * TrackCardSkeleton Component
 *
 * Loading skeleton placeholder for TrackCard.
 * Uses DaisyUI skeleton classes for consistent loading states.
 */

export function TrackCardSkeleton() {
  return (
    <div className="card bg-base-200">
      <div className="card-body">
        <div className="skeleton h-6 w-3/4"></div>
        <div className="skeleton h-4 w-1/2"></div>
        <div className="skeleton h-4 w-2/3"></div>
        <div className="flex gap-2">
          <div className="skeleton h-4 w-16"></div>
          <div className="skeleton h-4 w-12"></div>
        </div>
      </div>
    </div>
  );
}
