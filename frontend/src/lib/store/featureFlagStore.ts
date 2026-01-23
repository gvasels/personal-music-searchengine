/**
 * Feature Flag Store
 * Zustand store for managing feature flags and subscription state
 */
import { create } from 'zustand';
import type { SubscriptionTier, FeatureKey } from '@/types';

interface FeatureFlagState {
  // State
  tier: SubscriptionTier;
  features: Record<string, boolean>;
  isLoaded: boolean;

  // Actions
  setFeatures: (tier: SubscriptionTier, features: Record<string, boolean>) => void;
  isEnabled: (feature: FeatureKey) => boolean;
  reset: () => void;
}

export const useFeatureFlagStore = create<FeatureFlagState>((set, get) => ({
  tier: 'free',
  features: {},
  isLoaded: false,

  setFeatures: (tier, features) => {
    set({ tier, features, isLoaded: true });
  },

  isEnabled: (feature) => {
    const { features } = get();
    return features[feature] === true;
  },

  reset: () => {
    set({ tier: 'free', features: {}, isLoaded: false });
  },
}));
