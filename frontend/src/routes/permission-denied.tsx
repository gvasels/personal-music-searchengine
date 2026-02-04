/**
 * Permission Denied Page - Access Control Bug Fixes
 * Displays when guest users try to access protected routes
 * Shows simulation controls when admin is testing guest role
 */
import { Link } from '@tanstack/react-router';
import { useRoleSimulation } from '../hooks/useRoleSimulation';
import { useAuth } from '../hooks/useAuth';

function PermissionDeniedPage() {
  const { isSimulating, stopSimulation, simulatedRole } = useRoleSimulation();
  const { isAuthenticated, role } = useAuth();

  // Check if this is a simulated guest view
  const isSimulatedGuest = isAuthenticated && isSimulating && simulatedRole === 'guest';
  // Check if this is a demoted guest (authenticated but with guest role)
  const isDemotedGuest = isAuthenticated && role === 'guest' && !isSimulating;

  return (
    <main className="min-h-screen bg-base-200 flex items-center justify-center p-4">
      <div className="card bg-base-100 shadow-xl max-w-md w-full">
        <div className="card-body items-center text-center">
          {/* Permission Denied Banner - only for demoted or simulated guests */}
          {(isDemotedGuest || isSimulatedGuest) && (
            <div className="alert alert-error mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
              </svg>
              <span className="font-semibold">Permission Denied</span>
            </div>
          )}

          {/* Simulation indicator */}
          {isSimulatedGuest && (
            <div className="alert alert-warning mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <span>Simulating Guest View</span>
            </div>
          )}

          {/* Lock Icon */}
          <div className="w-20 h-20 rounded-full bg-error/10 flex items-center justify-center mb-4">
            <svg
              className="w-10 h-10 text-error"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
          </div>

          <h1 className="text-2xl font-bold mb-2">
            {isSimulatedGuest
              ? 'Guest View'
              : isDemotedGuest
                ? 'Access Restricted'
                : 'Welcome to Music Library'}
          </h1>

          <p className="text-base-content/70 mb-6">
            {isSimulatedGuest
              ? 'This is what unauthenticated users see. Guest users are immediately redirected here and must sign in or create an account to access the app.'
              : isDemotedGuest
                ? 'Your account does not have permission to access this content. Please contact an administrator if you believe this is an error.'
                : 'Please sign in or create an account to access your personal music library.'}
          </p>

          <div className="flex flex-col gap-3 w-full">
            {isSimulatedGuest ? (
              <button onClick={stopSimulation} className="btn btn-primary w-full">
                Stop Simulation
              </button>
            ) : isDemotedGuest ? (
              <Link to="/" className="btn btn-primary w-full">
                Go to Home
              </Link>
            ) : (
              <>
                <Link to="/login" className="btn btn-primary w-full">
                  Sign In
                </Link>
                <Link to="/signup" className="btn btn-outline w-full">
                  Create Account
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}

export default PermissionDeniedPage;
