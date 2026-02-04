/**
 * AuthGuard Component
 * Protects routes by redirecting unauthenticated users or simulated guests
 * Supports role simulation for admin testing
 */

import { ReactNode, useEffect } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { useAuth } from '../../hooks/useAuth';
import { useFeatureFlags } from '../../hooks/useFeatureFlags';

interface AuthGuardProps {
  children: ReactNode;
}

export function AuthGuard({ children }: AuthGuardProps) {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading } = useAuth();
  const { role: effectiveRole, isLoaded } = useFeatureFlags();

  // Check if user should be treated as guest:
  // 1. Not authenticated at all
  // 2. Guest role from API (real-time from DB) - includes demoted users and simulated guests
  const isEffectivelyGuest =
    !isAuthenticated ||
    effectiveRole === 'guest';

  useEffect(() => {
    // Wait for both auth and feature flags to load
    if (isLoading || (isAuthenticated && !isLoaded)) {
      return;
    }

    if (isEffectivelyGuest) {
      // Redirect to permission-denied page for guest users or simulated guests
      navigate({ to: '/permission-denied' });
    }
  }, [isLoading, isLoaded, isAuthenticated, isEffectivelyGuest, navigate]);

  // Show loading spinner while checking auth or loading feature flags
  if (isLoading || (isAuthenticated && !isLoaded)) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading" />
      </div>
    );
  }

  // Don't render children if effectively a guest
  if (isEffectivelyGuest) {
    return null;
  }

  return <>{children}</>;
}
