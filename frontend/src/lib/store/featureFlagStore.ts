/**
 * Feature Flag Store
 * Zustand store for managing feature flags with role-based access
 */
import { create } from 'zustand';
import type { UserRole, FeatureKey } from '@/types';

interface FeatureFlagState {
  // State
  role: UserRole;
  features: Record<string, boolean>;
  isLoaded: boolean;

  // Actions
  setFeatures: (role: UserRole, features: Record<string, boolean>) => void;
  isEnabled: (feature: FeatureKey) => boolean;
  reset: () => void;
}

export const useFeatureFlagStore = create<FeatureFlagState>((set, get) => ({
  role: 'subscriber',
  features: {},
  isLoaded: false,

  setFeatures: (role, features) => {
    set({ role, features, isLoaded: true });
  },

  isEnabled: (feature) => {
    const { features } = get();
    return features[feature] === true;
  },

  reset: () => {
    set({ role: 'subscriber', features: {}, isLoaded: false });
  },
}));
