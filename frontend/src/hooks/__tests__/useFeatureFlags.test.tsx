/**
 * useFeatureFlags Hook Tests - WS4 Creator Studio
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
    user: { id: 'user-1' },
  })),
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

      expect(result.current.tier).toBe('creator');
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

    describe('hasTier', () => {
      it('should return true when user has required tier', async () => {
        const mockFeatures = {
          tier: 'pro' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        expect(result.current.hasTier('free')).toBe(true);
        expect(result.current.hasTier('creator')).toBe(true);
        expect(result.current.hasTier('pro')).toBe(true);
      });

      it('should return false when user lacks required tier', async () => {
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

        expect(result.current.hasTier('free')).toBe(true);
        expect(result.current.hasTier('creator')).toBe(false);
        expect(result.current.hasTier('pro')).toBe(false);
      });

      it('should handle creator tier correctly', async () => {
        const mockFeatures = {
          tier: 'creator' as const,
          features: {},
        };
        vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

        const { result } = renderHook(() => useFeatureFlags(), {
          wrapper: createWrapper(),
        });

        await waitFor(() => {
          expect(result.current.isLoaded).toBe(true);
        });

        expect(result.current.hasTier('free')).toBe(true);
        expect(result.current.hasTier('creator')).toBe(true);
        expect(result.current.hasTier('pro')).toBe(false);
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

      expect(result.current.showUpgrade).toBe(false);
    });

    it('should return showUpgrade true when feature is disabled', async () => {
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
      expect(result.current.showUpgrade).toBe(true);
    });

    it('should not show upgrade while loading', () => {
      vi.mocked(featuresApi.getUserFeatures).mockReturnValue(new Promise(() => {}));

      const { result } = renderHook(() => useFeatureGate('CRATES'), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(true);
      expect(result.current.showUpgrade).toBe(false);
    });

    it('should return the user tier', async () => {
      const mockFeatures = {
        tier: 'creator' as const,
        features: { CRATES: true },
      };
      vi.mocked(featuresApi.getUserFeatures).mockResolvedValue(mockFeatures);

      const { result } = renderHook(() => useFeatureGate('CRATES'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.tier).toBe('creator');
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
    expect(state.tier).toBe('free');
    expect(state.features).toEqual({});
    expect(state.isLoaded).toBe(false);
  });

  it('should setFeatures correctly', () => {
    const store = useFeatureFlagStore.getState();
    store.setFeatures('pro', { CRATES: true, HOT_CUES: true });

    const newState = useFeatureFlagStore.getState();
    expect(newState.tier).toBe('pro');
    expect(newState.features).toEqual({ CRATES: true, HOT_CUES: true });
    expect(newState.isLoaded).toBe(true);
  });

  it('should check isEnabled correctly', () => {
    const store = useFeatureFlagStore.getState();
    store.setFeatures('creator', { CRATES: true, HOT_CUES: false });

    const state = useFeatureFlagStore.getState();
    expect(state.isEnabled('CRATES')).toBe(true);
    expect(state.isEnabled('HOT_CUES')).toBe(false);
    expect(state.isEnabled('UNKNOWN_FEATURE')).toBe(false);
  });

  it('should reset correctly', () => {
    const store = useFeatureFlagStore.getState();
    store.setFeatures('pro', { CRATES: true });

    store.reset();

    const newState = useFeatureFlagStore.getState();
    expect(newState.tier).toBe('free');
    expect(newState.features).toEqual({});
    expect(newState.isLoaded).toBe(false);
  });
});
