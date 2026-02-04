export function TrackCardSkeleton() {
  return (
    <div className="card bg-base-200 shadow-xl">
      <figure>
        <div className="skeleton h-48 w-full" />
      </figure>
      <div className="card-body p-4">
        <div className="skeleton h-4 w-3/4" />
        <div className="skeleton h-3 w-1/2" />
        <div className="skeleton h-3 w-1/3" />
        <div className="flex justify-between mt-2">
          <div className="skeleton h-5 w-20" />
          <div className="skeleton h-3 w-10" />
        </div>
      </div>
    </div>
  );
}
