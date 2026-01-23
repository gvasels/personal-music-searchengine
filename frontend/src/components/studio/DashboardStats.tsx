/**
 * DashboardStats Component
 * Displays user statistics in the creator dashboard
 */
import { useQuery } from '@tanstack/react-query';
import { getStorageUsage } from '@/lib/api/features';
import { getTracks } from '@/lib/api/client';

function formatBytes(bytes: number): string {
  if (bytes === -1) return 'Unlimited';
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

export function DashboardStats() {
  const { data: storage, isLoading: storageLoading } = useQuery({
    queryKey: ['storage'],
    queryFn: getStorageUsage,
    staleTime: 1000 * 60 * 5,
  });

  const { data: tracks, isLoading: tracksLoading } = useQuery({
    queryKey: ['tracks', { limit: 1 }],
    queryFn: () => getTracks({ limit: 1 }),
    staleTime: 1000 * 60 * 5,
  });

  const isLoading = storageLoading || tracksLoading;

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {[1, 2, 3].map((i) => (
          <div key={i} className="stat bg-base-200 rounded-lg animate-pulse">
            <div className="h-4 bg-base-300 rounded w-20 mb-2" />
            <div className="h-8 bg-base-300 rounded w-24" />
          </div>
        ))}
      </div>
    );
  }

  const usagePercent = storage?.usagePercent ?? 0;
  const isUnlimited = storage?.storageLimitBytes === -1;

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      {/* Track Count */}
      <div className="stat bg-base-200 rounded-lg">
        <div className="stat-title">Total Tracks</div>
        <div className="stat-value text-primary">{tracks?.total ?? 0}</div>
        <div className="stat-desc">In your library</div>
      </div>

      {/* Storage Used */}
      <div className="stat bg-base-200 rounded-lg">
        <div className="stat-title">Storage Used</div>
        <div className="stat-value text-secondary">
          {formatBytes(storage?.storageUsedBytes ?? 0)}
        </div>
        <div className="stat-desc">
          {isUnlimited
            ? 'Unlimited storage'
            : `of ${formatBytes(storage?.storageLimitBytes ?? 0)}`}
        </div>
      </div>

      {/* Storage Progress */}
      <div className="stat bg-base-200 rounded-lg">
        <div className="stat-title">Storage Usage</div>
        {isUnlimited ? (
          <div className="stat-value text-accent">∞</div>
        ) : (
          <>
            <div className="stat-value">{usagePercent.toFixed(1)}%</div>
            <progress
              className={`progress w-full ${
                usagePercent > 90 ? 'progress-error' : usagePercent > 70 ? 'progress-warning' : 'progress-primary'
              }`}
              value={usagePercent}
              max="100"
            />
          </>
        )}
        <div className="stat-desc">
          {usagePercent > 90 && !isUnlimited
            ? 'Consider upgrading'
            : isUnlimited
            ? 'Pro storage'
            : 'Available space'}
        </div>
      </div>
    </div>
  );
}
