/**
 * useRoleSimulation Hook - Admin Role Switching Feature
 *
 * Combines auth state with simulation state to provide
 * an effective role for UI rendering.
 */
import { useCallback, useMemo } from 'react';
import { useRoleSimulationStore } from '../lib/store/roleSimulationStore';
import { useAuth } from './useAuth';
import type { UserRole } from '../types';

interface UseRoleSimulationReturn {
  // State
  isSimulating: boolean;
  simulatedRole: UserRole | null;
  effectiveRole: UserRole;
  actualRole: UserRole;

  // Actions
  startSimulation: (role: UserRole) => void;
  stopSimulation: () => void;

  // Utilities
  canSimulate: boolean;
}

export function useRoleSimulation(): UseRoleSimulationReturn {
  const { simulatedRole, setSimulatedRole, clearSimulation } = useRoleSimulationStore();
  const { role: actualRole, isAdmin } = useAuth();
  const canSimulate = isAdmin;

  const isSimulating = canSimulate && simulatedRole !== null && simulatedRole !== 'admin';

  const effectiveRole = useMemo(() => {
    if (isSimulating && simulatedRole) {
      return simulatedRole;
    }
    return actualRole;
  }, [isSimulating, simulatedRole, actualRole]);

  // Note: Write blocking during simulation is a future enhancement.
  // When implemented, mutation hooks should check effectiveRole before executing.

  const startSimulation = useCallback((role: UserRole) => {
    if (canSimulate) {
      setSimulatedRole(role);
    }
  }, [canSimulate, setSimulatedRole]);

  const stopSimulation = useCallback(() => {
    clearSimulation();
  }, [clearSimulation]);

  return {
    isSimulating,
    simulatedRole,
    effectiveRole,
    actualRole,
    startSimulation,
    stopSimulation,
    canSimulate,
  };
}
