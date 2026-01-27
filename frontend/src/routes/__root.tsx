/**
 * Root Layout - Access Control Bug Fixes
 * Implements guest user route protection
 * Supports role simulation for admin testing
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
  const { isAuthenticated, isLoading } = useAuth();
  const { role, isSimulating } = useFeatureFlags();
  const location = useLocation();

  // Show loading state while checking auth
  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <span className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  // Check if user should be treated as guest (either not authenticated OR simulating guest)
  const isEffectivelyGuest = !isAuthenticated || (isSimulating && role === 'guest');

  // If effectively guest and trying to access a protected route, redirect to permission-denied
  if (isEffectivelyGuest && !isPublicRoute(location.pathname)) {
    return <Navigate to="/permission-denied" />;
  }

  return <Outlet />;
}

export const Route = createRootRoute({
  component: RootComponent,
});
