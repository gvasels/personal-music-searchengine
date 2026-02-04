/**
 * useRoleSimulation Hook Tests - Admin Role Switching Feature
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useRoleSimulation } from '../useRoleSimulation';
import { useRoleSimulationStore } from '../../lib/store/roleSimulationStore';

// Mock useAuth hook
const mockUseAuth = vi.fn();
vi.mock('../useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

describe('useRoleSimulation', () => {
  beforeEach(() => {
    // Reset the store
    useRoleSimulationStore.setState({
      simulatedRole: null,
      startedAt: null,
    });
    // Default: non-admin user
    mockUseAuth.mockReturnValue({
      role: 'subscriber',
      isAdmin: false,
    });
  });

  describe('canSimulate', () => {
    it('should return true for admin users', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.canSimulate).toBe(true);
    });

    it('should return false for non-admin users', () => {
      mockUseAuth.mockReturnValue({ role: 'subscriber', isAdmin: false });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.canSimulate).toBe(false);
    });

    it('should return false for artist users', () => {
      mockUseAuth.mockReturnValue({ role: 'artist', isAdmin: false });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.canSimulate).toBe(false);
    });
  });

  describe('isSimulating', () => {
    it('should return false when no simulated role is set', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.isSimulating).toBe(false);
    });

    it('should return true when simulating a non-admin role', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });
      useRoleSimulationStore.setState({ simulatedRole: 'subscriber', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.isSimulating).toBe(true);
    });

    it('should return false when simulating admin role', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });
      useRoleSimulationStore.setState({ simulatedRole: 'admin', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.isSimulating).toBe(false);
    });

    it('should return false for non-admin even with simulated role', () => {
      mockUseAuth.mockReturnValue({ role: 'subscriber', isAdmin: false });
      useRoleSimulationStore.setState({ simulatedRole: 'guest', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.isSimulating).toBe(false);
    });
  });

  describe('effectiveRole', () => {
    it('should return actual role when not simulating', () => {
      mockUseAuth.mockReturnValue({ role: 'subscriber', isAdmin: false });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.effectiveRole).toBe('subscriber');
    });

    it('should return simulated role when simulating', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });
      useRoleSimulationStore.setState({ simulatedRole: 'guest', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.effectiveRole).toBe('guest');
    });

    it('should return actual role when simulating admin', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });
      useRoleSimulationStore.setState({ simulatedRole: 'admin', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.effectiveRole).toBe('admin');
    });
  });

  describe('actualRole', () => {
    it('should always return the actual role from auth', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });
      useRoleSimulationStore.setState({ simulatedRole: 'guest', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.actualRole).toBe('admin');
    });
  });

  describe('startSimulation', () => {
    it('should set simulated role when user can simulate', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });

      const { result } = renderHook(() => useRoleSimulation());

      act(() => {
        result.current.startSimulation('subscriber');
      });

      expect(result.current.simulatedRole).toBe('subscriber');
      expect(result.current.isSimulating).toBe(true);
    });

    it('should not set simulated role when user cannot simulate', () => {
      mockUseAuth.mockReturnValue({ role: 'subscriber', isAdmin: false });

      const { result } = renderHook(() => useRoleSimulation());

      act(() => {
        result.current.startSimulation('guest');
      });

      expect(result.current.simulatedRole).toBeNull();
      expect(result.current.isSimulating).toBe(false);
    });

    it('should update the effective role when starting simulation', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.effectiveRole).toBe('admin');

      act(() => {
        result.current.startSimulation('artist');
      });

      expect(result.current.effectiveRole).toBe('artist');
    });
  });

  describe('stopSimulation', () => {
    it('should clear simulated role', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });
      useRoleSimulationStore.setState({ simulatedRole: 'guest', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.isSimulating).toBe(true);

      act(() => {
        result.current.stopSimulation();
      });

      expect(result.current.isSimulating).toBe(false);
      expect(result.current.simulatedRole).toBeNull();
    });

    it('should return effective role to actual role', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });
      useRoleSimulationStore.setState({ simulatedRole: 'subscriber', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.effectiveRole).toBe('subscriber');

      act(() => {
        result.current.stopSimulation();
      });

      expect(result.current.effectiveRole).toBe('admin');
    });
  });

  describe('simulatedRole', () => {
    it('should expose the raw simulated role from store', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });
      useRoleSimulationStore.setState({ simulatedRole: 'artist', startedAt: Date.now() });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.simulatedRole).toBe('artist');
    });

    it('should be null when no simulation is active', () => {
      mockUseAuth.mockReturnValue({ role: 'admin', isAdmin: true });

      const { result } = renderHook(() => useRoleSimulation());

      expect(result.current.simulatedRole).toBeNull();
    });
  });
});
