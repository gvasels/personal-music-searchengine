/**
 * Feature Flags API
 * Handles feature flag and subscription-related API calls
 */
import { apiClient } from './client';
import type {
  UserFeaturesResponse,
  TierConfig,
  SubscriptionResponse,
  StorageUsageResponse
} from '@/types';

// Get user's enabled features and tier
export async function getUserFeatures(): Promise<UserFeaturesResponse> {
  const response = await apiClient.get<UserFeaturesResponse>('/features');
  return response.data;
}

// Get subscription details
export async function getSubscription(): Promise<SubscriptionResponse> {
  const response = await apiClient.get<SubscriptionResponse>('/subscription');
  return response.data;
}

// Get available tier configurations
export async function getTierConfigs(): Promise<TierConfig[]> {
  const response = await apiClient.get<TierConfig[]>('/subscription/tiers');
  return response.data;
}

// Get storage usage
export async function getStorageUsage(): Promise<StorageUsageResponse> {
  const response = await apiClient.get<StorageUsageResponse>('/subscription/storage');
  return response.data;
}

// Create checkout session for upgrading subscription
export async function createCheckoutSession(
  tier: 'creator' | 'pro',
  interval: 'monthly' | 'yearly',
  successUrl?: string,
  cancelUrl?: string
): Promise<{ checkoutUrl: string; sessionId: string }> {
  const params = new URLSearchParams();
  if (successUrl) params.set('successUrl', successUrl);
  if (cancelUrl) params.set('cancelUrl', cancelUrl);

  const response = await apiClient.post<{ checkoutUrl: string; sessionId: string }>(
    `/subscription/checkout?${params.toString()}`,
    { tier, interval }
  );
  return response.data;
}

// Create customer portal session for managing subscription
export async function createPortalSession(returnUrl?: string): Promise<{ portalUrl: string }> {
  const params = new URLSearchParams();
  if (returnUrl) params.set('returnUrl', returnUrl);

  const response = await apiClient.post<{ portalUrl: string }>(
    `/subscription/portal?${params.toString()}`
  );
  return response.data;
}
