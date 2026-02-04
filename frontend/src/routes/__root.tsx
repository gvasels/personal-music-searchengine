/**
 * Root Layout - Access Control Bug Fixes
 * Implements guest user route protection
 * Supports role simulation for admin testing
 * Updated: 2026-01-27 - Added demoted guest detection
 */
import { createRootRoute, Outlet, useLocation, Navigate } from '@tanstack/react-router';
import { useAuth } from '../hooks/useAuth';
import { useFeatureFlags } from '../hooks/useFeatureFlags';

/**
 * Routes accessible to unauthenticated (guest) users.
 * All other routes require authentication.
 */
const PUBLIC_ROUTES = ['/', '/login', '/permission-denied'];

/**
 * Check if a path matches any public route.
 * Handles exact matches for simple routes.
 */
function isPublicRoute(pathname: string): boolean {
  return PUBLIC_ROUTES.includes(pathname);
}

function RootComponent() {
  const { isAuthenticated, isLoading, role: authRole } = useAuth();
  const { role: effectiveRole, isSimulating } = useFeatureFlags();
  const location = useLocation();

  // Show loading state while checking auth
  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <span className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  // Check if user should be treated as guest:
  // 1. Not authenticated at all
  // 2. Authenticated but with guest role (demoted user)
  // 3. Admin simulating guest role
  const isEffectivelyGuest =
    !isAuthenticated ||
    authRole === 'guest' ||
    (isSimulating && effectiveRole === 'guest');

  // Debug: log access control check
  if (typeof window !== 'undefined') {
    (window as unknown as Record<string, unknown>).__routeDebug = {
      pathname: location.pathname,
      isAuthenticated,
      authRole,
      isSimulating,
      effectiveRole,
      isEffectivelyGuest,
      isPublicRoute: isPublicRoute(location.pathname),
    };
    console.log('[Route] Access check v2:', (window as unknown as Record<string, unknown>).__routeDebug);
  }

  // If effectively guest and trying to access a protected route, redirect to permission-denied
  if (isEffectivelyGuest && !isPublicRoute(location.pathname)) {
    console.log('[Route] Redirecting guest to permission-denied');
    return <Navigate to="/permission-denied" />;
  }

  return <Outlet />;
}

export const Route = createRootRoute({
  component: RootComponent,
});
