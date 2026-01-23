/**
 * Creator Studio Dashboard
 * Main hub for creator features with feature-gated modules
 */
import { createFileRoute } from '@tanstack/react-router';
import { useFeatureFlags } from '@/hooks/useFeatureFlags';
import { DashboardStats, ModuleCard, studioModules } from '@/components/studio';

// @ts-expect-error - Route path will be registered when route tree regenerates
export const Route = createFileRoute('/studio/')({
  component: StudioDashboard,
});

function StudioDashboard() {
  const { tier, isLoading } = useFeatureFlags();

  return (
    <div className="container mx-auto p-6 max-w-6xl">
      {/* Header */}
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold">Creator Studio</h1>
          <p className="text-base-content/70">
            {isLoading ? (
              <span className="loading loading-dots loading-sm" />
            ) : (
              <>Your tools for creating and managing music</>
            )}
          </p>
        </div>
        <div className="badge badge-lg badge-primary capitalize">
          {tier} Plan
        </div>
      </div>

      {/* Stats */}
      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-4">Overview</h2>
        <DashboardStats />
      </section>

      {/* Modules Grid */}
      <section>
        <h2 className="text-xl font-semibold mb-4">Studio Modules</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {studioModules.map((module) => (
            <ModuleCard key={module.feature} {...module} />
          ))}
        </div>
      </section>

      {/* Upgrade CTA for free users */}
      {tier === 'free' && (
        <section className="mt-8">
          <div className="alert alert-info">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-current shrink-0 w-6 h-6">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div>
              <h3 className="font-bold">Unlock Creator Features</h3>
              <div className="text-sm">
                Upgrade to Creator or Pro to access DJ crates, hot cues, BPM/key matching, and more.
              </div>
            </div>
            <a href="/subscription" className="btn btn-sm btn-primary">
              View Plans
            </a>
          </div>
        </section>
      )}
    </div>
  );
}
