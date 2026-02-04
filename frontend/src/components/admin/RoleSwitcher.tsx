/**
 * RoleSwitcher Component - Admin Role Switching Feature
 *
 * Dropdown for admins to select which role to view the app as.
 */
import { useRoleSimulation } from '../../hooks/useRoleSimulation';
import type { UserRole } from '../../types';

interface RoleOption {
  value: UserRole;
  label: string;
  description: string;
  icon: string;
}

const ROLE_OPTIONS: RoleOption[] = [
  { value: 'admin', label: 'Admin', description: 'Full access (actual role)', icon: '👑' },
  { value: 'artist', label: 'Artist', description: 'Artist features enabled', icon: '🎨' },
  { value: 'subscriber', label: 'Subscriber', description: 'Standard user access', icon: '👤' },
  { value: 'guest', label: 'Guest', description: 'Unauthenticated view', icon: '👻' },
];

export function RoleSwitcher() {
  const { canSimulate, effectiveRole, startSimulation, isSimulating } = useRoleSimulation();

  if (!canSimulate) {
    return null;
  }

  const currentOption = ROLE_OPTIONS.find((opt) => opt.value === effectiveRole) ?? ROLE_OPTIONS[0];

  return (
    <div className="dropdown dropdown-end">
      <label
        tabIndex={0}
        className={`btn btn-sm gap-2 ${isSimulating ? 'btn-warning' : 'btn-ghost'}`}
      >
        <span>{currentOption.icon}</span>
        <span className="hidden sm:inline">{currentOption.label}</span>
        <svg
          className="h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </label>
      <ul
        tabIndex={0}
        className="dropdown-content z-[100] menu p-2 shadow-lg bg-base-100 rounded-box w-64 border border-base-300"
      >
        <li className="menu-title">
          <span>View as Role</span>
        </li>
        {ROLE_OPTIONS.map((option) => (
          <li key={option.value}>
            <button
              onClick={() => startSimulation(option.value)}
              className={effectiveRole === option.value ? 'active' : ''}
            >
              <span className="text-lg">{option.icon}</span>
              <div className="flex flex-col items-start">
                <span className="font-medium">{option.label}</span>
                <span className="text-xs opacity-60">{option.description}</span>
              </div>
              {effectiveRole === option.value && (
                <svg className="h-4 w-4 ml-auto" fill="currentColor" viewBox="0 0 20 20">
                  <path
                    fillRule="evenodd"
                    d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                    clipRule="evenodd"
                  />
                </svg>
              )}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
