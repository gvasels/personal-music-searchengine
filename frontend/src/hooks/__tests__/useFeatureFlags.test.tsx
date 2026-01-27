/**
 * useFeatureFlags Hook Tests - WS4 Creator Studio
 * Updated for role-based access control
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useFeatureFlags, useFeatureGate, featureKeys } from '../useFeatureFlags';
import * as featuresApi from '../../lib/api/features';
import { useFeatureFlagStore } from '../../lib/store/featureFlagStore';

// Mock the API
vi.mock('../../lib/api/features', () => ({
  getUserFeatures: vi.fn(),
}));

// Mock useAuth hook
vi.mock('../useAuth', () => ({
  useAuth: vi.fn(() => ({
    isAuthenticated: true,
    user: { id: 'user-1', role: 'subscriber' },
  })),
}));

// Mock useRoleSimulation hook
const mockUseRoleSimulation = vi.fn(() => ({
  effectiveRole: 'subscriber',
  isSimulating: false,
  canSimulate: true,
  startSimulation: vi.fn(),
  stopSimulation: vi.fn(),
}));

vi.mock('../useRoleSimulation', () => ({
  useRoleSimulation: () => mockUseRoleSimulation(),
}));

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('useFeatureFlags (WS4 Creator Studio)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Reset the store before each test
    useFeatureFlagStore.getState().reset();
    // Reset useRoleSimulation mock to default (not simulating)
    mockUseRoleSimulation.mockReturnValue({
      effectiveRole: 'subscriber',
      isSimulating: false,
      canSimulate: true,
      startSimulation: vi.fn(),
      stopSimulation: vi.fn(),
          });
  });

  afterEach(() => {
    useFeatureFlagStore.getState().reset();
  });

  describe('featureKeys', () => {
    it('should generate correct query keys', () => {
      expect(featureKeys.all).toEqual(['features']);
      expect(featureKeys.user()).toEqual(['features', 'user']);
    });
  });

  describe('useFeatureFlags', () => {
    it('should fetch user features and update store', async () => {
      const mockFeatures = {
        tier: 'creator' as const,
        features: {
          CRATES: true,
          HOT_CUES: true,
          DJ_MODULE: true,
          BPM_MATCHING: true,
          KEY_MATCHING: false,
        },
      };
      vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

      const { result } = renderHook(() => useFeatureFlags(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isLoaded).toBe(true);
      });

      // Role is mapped from tier (creator -> artist)
      expect(result.current.role).toBe('subscriber'); // User role from auth mock
      expect(result.current.features).toEqual(mockFeatures.features);
      expect(featuresApi.getUserFeatures).toHaveBeenCalled();
    });

    it('should handle loading state', () => {
      vi.mocked(featuresApi.getUserFeatures).mockReturnValue(new Promise(() => {}));

      const { result } = renderHook(() => useFeatureFlags(), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(true);
      expect(result.current.isLoaded).toBe(false);
    });

    it('should handle error state', async () => {
      vi.mocked(featuresApi.getUserFeatures).mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => useFeatureFlags(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });

    describe('isEnabled', () => {
      it('should return true for enabled features', async () => {
        const mockFeatures = {
          tier: 'pro' as const,
          features: { CRATES: true, HOT_CUES: true },
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        expect(result.current.isEnabled('CRATES')).toBe(true);
        expect(result.current.isEnabled('HOT_CUES')).toBe(true);
      });

      it('should return false for disabled features', async () => {
        const mockFeatures = {
          tier: 'free' as const,
          features: { CRATES: false, HOT_CUES: false },
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        expect(result.current.isEnabled('CRATES')).toBe(false);
        expect(result.current.isEnabled('HOT_CUES')).toBe(false);
      });

      it('should return false when not loaded', () => {
        vi.mocked(featuresApi.getUserFeatures).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        expect(result.current.isEnabled('CRATES')).toBe(false);
      });
    });

    describe('hasRole', () => {
      it('should return true when user has required role', async () => {
        const mockFeatures = {
          tier: 'pro' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        // Set the store directly to test hasRole
        useFeatureFlagStore.getState().setFeatures('admin', {});

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        expect(result.current.hasRole('guest')).toBe(true);
        expect(result.current.hasRole('subscriber')).toBe(true);
        expect(result.current.hasRole('artist')).toBe(true);
        expect(result.current.hasRole('admin')).toBe(true);
      });

      it('should return false when user lacks required role', async () => {
        const mockFeatures = {
          tier: 'free' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        // Set the store directly to test hasRole
        useFeatureFlagStore.getState().setFeatures('subscriber', {});

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        expect(result.current.hasRole('guest')).toBe(true);
        expect(result.current.hasRole('subscriber')).toBe(true);
        expect(result.current.hasRole('artist')).toBe(false);
        expect(result.current.hasRole('admin')).toBe(false);
      });

      it('should handle artist role correctly', async () => {
        const mockFeatures = {
          tier: 'creator' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        // Set the store directly to test hasRole
        useFeatureFlagStore.getState().setFeatures('artist', {});

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        expect(result.current.hasRole('guest')).toBe(true);
        expect(result.current.hasRole('subscriber')).toBe(true);
        expect(result.current.hasRole('artist')).toBe(true);
        expect(result.current.hasRole('admin')).toBe(false);
      });
    });

    describe('simulation mode', () => {
      it('should use simulated role when isSimulating is true', async () => {
        // Mock simulation active with guest role
        mockUseRoleSimulation.mockReturnValue({
          effectiveRole: 'guest',
          isSimulating: true,
          canSimulate: true,
          startSimulation: vi.fn(),
          stopSimulation: vi.fn(),
                  });

        const mockFeatures = {
          tier: 'pro' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        // Set the store with admin role
        useFeatureFlagStore.getState().setFeatures('admin', {});

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        // Should return simulated role, not actual role
        expect(result.current.role).toBe('guest');
        expect(result.current.actualRole).toBe('admin');
        expect(result.current.isSimulating).toBe(true);
      });

      it('should use actual role when not simulating', async () => {
        // Mock simulation inactive
        mockUseRoleSimulation.mockReturnValue({
          effectiveRole: 'subscriber',
          isSimulating: false,
          canSimulate: true,
          startSimulation: vi.fn(),
          stopSimulation: vi.fn(),
                  });

        const mockFeatures = {
          tier: 'pro' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        // Set the store with admin role
        useFeatureFlagStore.getState().setFeatures('admin', {});

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        // Should return actual role
        expect(result.current.role).toBe('admin');
        expect(result.current.actualRole).toBe('admin');
        expect(result.current.isSimulating).toBe(false);
      });

      it('should check hasRole against simulated role when simulating', async () => {
        // Mock simulation with subscriber role
        mockUseRoleSimulation.mockReturnValue({
          effectiveRole: 'subscriber',
          isSimulating: true,
          canSimulate: true,
          startSimulation: vi.fn(),
          stopSimulation: vi.fn(),
                  });

        const mockFeatures = {
          tier: 'pro' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        // Set the store with admin role
        useFeatureFlagStore.getState().setFeatures('admin', {});

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        // hasRole should check against simulated role (subscriber), not actual role (admin)
        expect(result.current.hasRole('guest')).toBe(true);
        expect(result.current.hasRole('subscriber')).toBe(true);
        expect(result.current.hasRole('artist')).toBe(false); // subscriber < artist
        expect(result.current.hasRole('admin')).toBe(false);  // subscriber < admin
      });

      it('should check hasRole against actual role when not simulating', async () => {
        // Mock simulation inactive
        mockUseRoleSimulation.mockReturnValue({
          effectiveRole: 'subscriber',
          isSimulating: false,
          canSimulate: true,
          startSimulation: vi.fn(),
          stopSimulation: vi.fn(),
                  });

        const mockFeatures = {
          tier: 'pro' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        // Set the store with admin role
        useFeatureFlagStore.getState().setFeatures('admin', {});

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        // hasRole should check against actual role (admin)
        expect(result.current.hasRole('guest')).toBe(true);
        expect(result.current.hasRole('subscriber')).toBe(true);
        expect(result.current.hasRole('artist')).toBe(true);
        expect(result.current.hasRole('admin')).toBe(true);
      });
    });

    describe('invalidate', () => {
      it('should invalidate the features query', async () => {
        const mockFeatures = {
          tier: 'free' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        // Should not throw
        act(() => {
          result.current.invalidate();
        });

        // After invalidation, query should refetch
        await waitFor(() => {
          expect(featuresApi.getUserFeatures).toHaveBeenCalledTimes(2);
        });
      });
    });
  });

  describe('useFeatureGate', () => {
    it('should return isEnabled true when feature is enabled', async () => {
      const mockFeatures = {
        tier: 'pro' as const,
        features: { CRATES: true },
      };
      vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

      const { result } = renderHook(() => useFeatureGate('CRATES'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isEnabled).toBe(true);
      });

      expect(result.current.isLocked).toBe(false);
    });

    it('should return isLocked true when feature is disabled', async () => {
      const mockFeatures = {
        tier: 'free' as const,
        features: { CRATES: false },
      };
      vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

      const { result } = renderHook(() => useFeatureGate('CRATES'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.isEnabled).toBe(false);
      expect(result.current.isLocked).toBe(true);
    });

    it('should not show locked state while loading', () => {
      vi.mocked(featuresApi.getUserFeatures).mockReturnValue(new Promise(() => {}));

      const { result } = renderHook(() => useFeatureGate('CRATES'), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(true);
      expect(result.current.isLocked).toBe(false);
    });

    it('should return the user role', async () => {
      const mockFeatures = {
        tier: 'creator' as const,
        features: { CRATES: true },
      };
      vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

      const { result } = renderHook(() => useFeatureGate('CRATES'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.role).toBe('subscriber'); // From auth mock
      });
    });
  });
});

describe('featureFlagStore', () => {
  beforeEach(() => {
    useFeatureFlagStore.getState().reset();
  });

  it('should start with default state', () => {
    const state = useFeatureFlagStore.getState();
    expect(state.role).toBe('subscriber');
    expect(state.features).toEqual({});
    expect(state.isLoaded).toBe(false);
  });

  it('should setFeatures correctly', () => {
    const store = useFeatureFlagStore.getState();
    store.setFeatures('artist', { CRATES: true, HOT_CUES: true });

    const newState = useFeatureFlagStore.getState();
    expect(newState.role).toBe('artist');
    expect(newState.features).toEqual({ CRATES: true, HOT_CUES: true });
    expect(newState.isLoaded).toBe(true);
  });

  it('should check isEnabled correctly', () => {
    const store = useFeatureFlagStore.getState();
    store.setFeatures('artist', { CRATES: true, HOT_CUES: false });

    const state = useFeatureFlagStore.getState();
    expect(state.isEnabled('CRATES')).toBe(true);
    expect(state.isEnabled('HOT_CUES')).toBe(false);
    expect(state.isEnabled('UNKNOWN_FEATURE')).toBe(false);
  });

  it('should reset correctly', () => {
    const store = useFeatureFlagStore.getState();
    store.setFeatures('admin', { CRATES: true });

    store.reset();

    const newState = useFeatureFlagStore.getState();
    expect(newState.role).toBe('subscriber');
    expect(newState.features).toEqual({});
    expect(newState.isLoaded).toBe(false);
  });
});
