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
  const { isAuthenticated } = useAuth();

  // Check if this is a simulated guest view
  const isSimulatedGuest = isAuthenticated && isSimulating && simulatedRole === 'guest';

  return (
    <main className="min-h-screen bg-base-200 flex items-center justify-center p-4">
      <div className="card bg-base-100 shadow-xl max-w-md w-full">
        <div className="card-body items-center text-center">
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

          <h1 className="text-2xl font-bold mb-2">Access Denied</h1>

          <p className="text-base-content/70 mb-6">
            {isSimulatedGuest
              ? 'Guest users cannot access this page. This is how unauthenticated users experience the app.'
              : 'You do not have permission to view this page. Please sign in to access your music library.'}
          </p>

          <div className="flex flex-col sm:flex-row gap-3 w-full">
            <Link to="/" className="btn btn-outline flex-1">
              Go to Dashboard
            </Link>
            {isSimulatedGuest ? (
              <button onClick={stopSimulation} className="btn btn-primary flex-1">
                Stop Simulation
              </button>
            ) : (
              <Link to="/login" className="btn btn-primary flex-1">
                Sign In
              </Link>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}

export default PermissionDeniedPage;
