/**
 * Creator Studio Dashboard
 * Main hub for creator features with feature-gated modules
 */
import { createFileRoute } from '@tanstack/react-router';
import { useFeatureFlags } from '@/hooks/useFeatureFlags';
import { DashboardStats, ModuleCard, studioModules } from '@/components/studio';

export const Route = createFileRoute('/studio/')({
  component: StudioDashboard,
});

function StudioDashboard() {
  const { role, isLoading } = useFeatureFlags();

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
          {role}
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
    </div>
  );
}
