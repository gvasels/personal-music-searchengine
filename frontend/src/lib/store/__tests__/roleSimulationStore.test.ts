/**
 * roleSimulationStore Tests - Admin Role Switching Feature
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useRoleSimulationStore } from '../roleSimulationStore';

describe('roleSimulationStore', () => {
  beforeEach(() => {
    // Clear the store state before each test
    useRoleSimulationStore.setState({
      simulatedRole: null,
      startedAt: null,
    });
    // Clear session storage
    sessionStorage.clear();
  });

  describe('initial state', () => {
    it('should start with null simulatedRole', () => {
      const state = useRoleSimulationStore.getState();
      expect(state.simulatedRole).toBeNull();
    });

    it('should start with null startedAt', () => {
      const state = useRoleSimulationStore.getState();
      expect(state.startedAt).toBeNull();
    });
  });

  describe('setSimulatedRole', () => {
    it('should set simulatedRole to the provided role', () => {
      const { setSimulatedRole } = useRoleSimulationStore.getState();

      setSimulatedRole('subscriber');

      const state = useRoleSimulationStore.getState();
      expect(state.simulatedRole).toBe('subscriber');
    });

    it('should set startedAt when setting a role', () => {
      const now = Date.now();
      vi.spyOn(Date, 'now').mockReturnValue(now);

      const { setSimulatedRole } = useRoleSimulationStore.getState();
      setSimulatedRole('guest');

      const state = useRoleSimulationStore.getState();
      expect(state.startedAt).toBe(now);

      vi.restoreAllMocks();
    });

    it('should clear startedAt when setting role to null', () => {
      const { setSimulatedRole } = useRoleSimulationStore.getState();

      // First set a role
      setSimulatedRole('artist');
      expect(useRoleSimulationStore.getState().startedAt).not.toBeNull();

      // Then clear it
      setSimulatedRole(null);
      expect(useRoleSimulationStore.getState().startedAt).toBeNull();
    });

    it('should update startedAt when changing roles', () => {
      const firstTime = 1000;
      const secondTime = 2000;

      vi.spyOn(Date, 'now').mockReturnValueOnce(firstTime).mockReturnValueOnce(secondTime);

      const { setSimulatedRole } = useRoleSimulationStore.getState();

      setSimulatedRole('subscriber');
      expect(useRoleSimulationStore.getState().startedAt).toBe(firstTime);

      setSimulatedRole('artist');
      expect(useRoleSimulationStore.getState().startedAt).toBe(secondTime);

      vi.restoreAllMocks();
    });

    it('should set all role types correctly', () => {
      const { setSimulatedRole } = useRoleSimulationStore.getState();

      const roles = ['guest', 'subscriber', 'artist', 'admin'] as const;

      for (const role of roles) {
        setSimulatedRole(role);
        expect(useRoleSimulationStore.getState().simulatedRole).toBe(role);
      }
    });
  });

  describe('clearSimulation', () => {
    it('should reset simulatedRole to null', () => {
      const { setSimulatedRole, clearSimulation } = useRoleSimulationStore.getState();

      setSimulatedRole('subscriber');
      expect(useRoleSimulationStore.getState().simulatedRole).toBe('subscriber');

      clearSimulation();
      expect(useRoleSimulationStore.getState().simulatedRole).toBeNull();
    });

    it('should reset startedAt to null', () => {
      const { setSimulatedRole, clearSimulation } = useRoleSimulationStore.getState();

      setSimulatedRole('artist');
      expect(useRoleSimulationStore.getState().startedAt).not.toBeNull();

      clearSimulation();
      expect(useRoleSimulationStore.getState().startedAt).toBeNull();
    });

    it('should be idempotent when called multiple times', () => {
      const { clearSimulation } = useRoleSimulationStore.getState();

      clearSimulation();
      clearSimulation();
      clearSimulation();

      const state = useRoleSimulationStore.getState();
      expect(state.simulatedRole).toBeNull();
      expect(state.startedAt).toBeNull();
    });
  });

  describe('persistence', () => {
    it('should use sessionStorage for persistence', () => {
      const { setSimulatedRole } = useRoleSimulationStore.getState();

      setSimulatedRole('artist');

      // Check sessionStorage was used
      const stored = sessionStorage.getItem('role-simulation');
      expect(stored).not.toBeNull();

      const parsed = JSON.parse(stored!);
      expect(parsed.state.simulatedRole).toBe('artist');
    });
  });
});
