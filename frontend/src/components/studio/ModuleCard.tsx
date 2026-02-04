/**
 * ModuleCard Component
 * Displays a studio module with feature gating
 */
import { useFeatureGate } from '@/hooks/useFeatureFlags';
import type { FeatureKey, UserRole } from '@/types';

interface ModuleCardProps {
  title: string;
  description: string;
  icon: string;
  feature: FeatureKey;
  href: string;
  requiredRole: UserRole;
}

export function ModuleCard({
  title,
  description,
  icon,
  feature,
  href,
  requiredRole,
}: ModuleCardProps) {
  const { isLoading, isLocked } = useFeatureGate(feature);

  if (isLoading) {
    return (
      <div className="card bg-base-200 animate-pulse">
        <div className="card-body">
          <div className="h-12 w-12 bg-base-300 rounded-full mb-4" />
          <div className="h-6 bg-base-300 rounded w-3/4 mb-2" />
          <div className="h-4 bg-base-300 rounded w-full" />
        </div>
      </div>
    );
  }

  if (isLocked) {
    return (
      <div className="card bg-base-200 border border-base-300 opacity-75">
        <div className="card-body">
          <div className="text-4xl mb-2 grayscale">{icon}</div>
          <h3 className="card-title text-base-content/70">{title}</h3>
          <p className="text-sm text-base-content/50">{description}</p>
          <div className="card-actions justify-end mt-4">
            <span className="badge badge-outline capitalize">
              {requiredRole} role required
            </span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <a href={href} className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer">
      <div className="card-body">
        <div className="text-4xl mb-2">{icon}</div>
        <h3 className="card-title">{title}</h3>
        <p className="text-sm text-base-content/70">{description}</p>
        <div className="card-actions justify-end mt-4">
          <span className="badge badge-primary">Available</span>
        </div>
      </div>
    </a>
  );
}

// Preset module configurations - features gated by role
// eslint-disable-next-line react-refresh/only-export-components
export const studioModules = [
  {
    title: 'DJ Crates',
    description: 'Organize tracks into crates for your DJ sets',
    icon: '📦',
    feature: 'CRATES' as FeatureKey,
    href: '/studio/crates',
    requiredRole: 'artist' as UserRole,
  },
  {
    title: 'BPM Matching',
    description: 'Find tracks with compatible tempos',
    icon: '🎚️',
    feature: 'BPM_MATCHING' as FeatureKey,
    href: '/studio/matching',
    requiredRole: 'artist' as UserRole,
  },
  {
    title: 'Key Matching',
    description: 'Find harmonically compatible tracks',
    icon: '🎹',
    feature: 'KEY_MATCHING' as FeatureKey,
    href: '/studio/matching',
    requiredRole: 'artist' as UserRole,
  },
  {
    title: 'Hot Cues',
    description: 'Set cue points on your tracks',
    icon: '🎯',
    feature: 'HOT_CUES' as FeatureKey,
    href: '/studio/hotcues',
    requiredRole: 'artist' as UserRole,
  },
  {
    title: 'Bulk Edit',
    description: 'Edit multiple tracks at once',
    icon: '✏️',
    feature: 'BULK_EDIT' as FeatureKey,
    href: '/studio/bulk-edit',
    requiredRole: 'artist' as UserRole,
  },
  {
    title: 'Mix Recording',
    description: 'Record your DJ mixes',
    icon: '🎙️',
    feature: 'MIX_RECORDING' as FeatureKey,
    href: '/studio/recording',
    requiredRole: 'artist' as UserRole,
  },
  {
    title: 'Advanced Stats',
    description: 'Detailed listening analytics',
    icon: '📊',
    feature: 'ADVANCED_STATS' as FeatureKey,
    href: '/studio/stats',
    requiredRole: 'admin' as UserRole,
  },
];
