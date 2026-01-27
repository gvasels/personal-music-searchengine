/**
 * Role Simulation Store - Admin Role Switching Feature
 *
 * Zustand store for managing role simulation state.
 * Allows admins to view the app as different user roles.
 */
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import type { UserRole } from '../../types';

interface RoleSimulationState {
  simulatedRole: UserRole | null;
  startedAt: number | null;
  setSimulatedRole: (role: UserRole | null) => void;
  clearSimulation: () => void;
}

export const useRoleSimulationStore = create<RoleSimulationState>()(
  persist(
    (set) => ({
      simulatedRole: null,
      startedAt: null,
      setSimulatedRole: (role) =>
        set({
          simulatedRole: role,
          startedAt: role ? Date.now() : null,
        }),
      clearSimulation: () =>
        set({
          simulatedRole: null,
          startedAt: null,
        }),
    }),
    {
      name: 'role-simulation',
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);
