/**
 * useFeatureFlags Hook
 * Fetches and caches user features with subscription tier awareness
 */
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect } from 'react';
import { useAuth } from './useAuth';
import { getUserFeatures } from '@/lib/api/features';
import { useFeatureFlagStore } from '@/lib/store/featureFlagStore';
import type { FeatureKey, SubscriptionTier } from '@/types';

// Query key factory
export const featureKeys = {
  all: ['features'] as const,
  user: () => [...featureKeys.all, 'user'] as const,
};

export function useFeatureFlags() {
  const { isAuthenticated } = useAuth();
  const queryClient = useQueryClient();
  const { tier, features, isLoaded, setFeatures, isEnabled: storeIsEnabled, reset } = useFeatureFlagStore();

  // Fetch features when authenticated
  const query = useQuery({
    queryKey: featureKeys.user(),
    queryFn: getUserFeatures,
    enabled: isAuthenticated,
    staleTime: 1000 * 60 * 5, // 5 minutes
    gcTime: 1000 * 60 * 30, // 30 minutes
  });

  // Update store when data changes
  useEffect(() => {
    if (query.data) {
      setFeatures(query.data.tier, query.data.features);
    }
  }, [query.data, setFeatures]);

  // Reset store on logout
  useEffect(() => {
    if (!isAuthenticated) {
      reset();
    }
  }, [isAuthenticated, reset]);

  // Check if a specific feature is enabled
  const isEnabled = (feature: FeatureKey): boolean => {
    if (!isLoaded) return false;
    return storeIsEnabled(feature);
  };

  // Check if user has at least a certain tier
  const hasTier = (minTier: SubscriptionTier): boolean => {
    const tierOrder: SubscriptionTier[] = ['free', 'creator', 'pro'];
    const userTierIndex = tierOrder.indexOf(tier);
    const requiredTierIndex = tierOrder.indexOf(minTier);
    return userTierIndex >= requiredTierIndex;
  };

  // Invalidate features cache (call after subscription change)
  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: featureKeys.all });
  };

  return {
    tier,
    features,
    isLoading: query.isLoading,
    isError: query.isError,
    isLoaded,
    isEnabled,
    hasTier,
    invalidate,
    refetch: query.refetch,
  };
}

// Hook for feature-gated components
export function useFeatureGate(feature: FeatureKey) {
  const { isEnabled, isLoading, tier } = useFeatureFlags();

  return {
    isEnabled: isEnabled(feature),
    isLoading,
    tier,
    showUpgrade: !isEnabled(feature) && !isLoading,
  };
}
