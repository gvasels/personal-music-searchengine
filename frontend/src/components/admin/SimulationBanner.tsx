/**
 * SimulationBanner Component - Admin Role Switching Feature
 *
 * Persistent banner showing when viewing as a different role.
 */
import { useRoleSimulation } from '../../hooks/useRoleSimulation';

const ROLE_LABELS: Record<string, string> = {
  guest: 'Guest',
  subscriber: 'Subscriber',
  artist: 'Artist',
  admin: 'Admin',
};

export function SimulationBanner() {
  const { isSimulating, effectiveRole, stopSimulation } = useRoleSimulation();

  if (!isSimulating) {
    return null;
  }

  const roleLabel = ROLE_LABELS[effectiveRole] ?? effectiveRole;

  return (
    <div className="bg-warning text-warning-content px-4 py-2 flex items-center justify-between sticky top-0 z-50">
      <div className="flex items-center gap-2">
        <svg
          className="h-5 w-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
          />
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
          />
        </svg>
        <span className="font-medium">
          Viewing as {roleLabel}
        </span>
        <span className="text-sm opacity-75">
          (Write actions blocked)
        </span>
      </div>
      <button
        onClick={stopSimulation}
        className="btn btn-sm btn-ghost"
      >
        Exit Simulation
      </button>
    </div>
  );
}
