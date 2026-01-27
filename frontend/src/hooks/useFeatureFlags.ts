/**
 * useFeatureFlags Hook
 * Fetches and caches user features with role-based access
 * Respects role simulation for admin testing
 */
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect } from 'react';
import { useAuth } from './useAuth';
import { useRoleSimulation } from './useRoleSimulation';
import { getUserFeatures } from '@/lib/api/features';
import { useFeatureFlagStore } from '@/lib/store/featureFlagStore';
import type { FeatureKey, UserRole } from '@/types';

// Query key factory
export const featureKeys = {
  all: ['features'] as const,
  user: () => [...featureKeys.all, 'user'] as const,
};

export function useFeatureFlags() {
  const { isAuthenticated, user } = useAuth();
  const { effectiveRole, isSimulating } = useRoleSimulation();
  const queryClient = useQueryClient();
  const { role, features, isLoaded, setFeatures, isEnabled: storeIsEnabled, reset } = useFeatureFlagStore();

  // Fetch features when authenticated
  const query = useQuery({
    queryKey: featureKeys.user(),
    queryFn: getUserFeatures,
    enabled: isAuthenticated,
    staleTime: 1000 * 60 * 5, // 5 minutes
    gcTime: 1000 * 60 * 30, // 30 minutes
  });

  // Update store when data changes - map tier to role for backward compatibility
  useEffect(() => {
    if (query.data) {
      // Map subscription tier to role (backward compatibility during migration)
      const roleFromTier: UserRole =
        query.data.tier === 'pro' ? 'artist' :
        query.data.tier === 'creator' ? 'artist' :
        'subscriber';
      // Use actual role from user if available, otherwise infer from tier
      const actualRole = (user?.role as UserRole) || roleFromTier;
      setFeatures(actualRole, query.data.features);
    }
  }, [query.data, setFeatures, user?.role]);

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

  // Check if user has at least a certain role level
  // Uses effective role (simulated or actual) for UI rendering
  const hasRole = (minRole: UserRole): boolean => {
    const roleOrder: UserRole[] = ['guest', 'subscriber', 'artist', 'admin'];
    const currentRole = isSimulating ? effectiveRole : role;
    const userRoleIndex = roleOrder.indexOf(currentRole);
    const requiredRoleIndex = roleOrder.indexOf(minRole);
    return userRoleIndex >= requiredRoleIndex;
  };

  // Invalidate features cache (call after role change)
  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: featureKeys.all });
  };

  return {
    role: isSimulating ? effectiveRole : role,
    actualRole: role,
    features,
    isLoading: query.isLoading,
    isError: query.isError,
    isLoaded,
    isEnabled,
    hasRole,
    invalidate,
    refetch: query.refetch,
    isSimulating,
  };
}

// Hook for feature-gated components
export function useFeatureGate(feature: FeatureKey) {
  const { isEnabled, isLoading, role } = useFeatureFlags();

  return {
    isEnabled: isEnabled(feature),
    isLoading,
    role,
    // Feature is locked if not enabled and not loading
    isLocked: !isEnabled(feature) && !isLoading,
  };
}
