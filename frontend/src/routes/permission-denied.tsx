/**
 * Permission Denied Page - Access Control Bug Fixes
 * Displays when guest users try to access protected routes
 */
import { Link } from '@tanstack/react-router';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/permission-denied')({
  component: PermissionDeniedPage,
});

function PermissionDeniedPage() {
  return (
    <main className="min-h-screen bg-base-200 flex items-center justify-center p-4">
      <div className="card bg-base-100 shadow-xl max-w-md w-full">
        <div className="card-body items-center text-center">
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
            You do not have permission to view this page. Please sign in to access your music library.
          </p>

          <div className="flex flex-col sm:flex-row gap-3 w-full">
            <Link to="/" className="btn btn-outline flex-1">
              Go to Dashboard
            </Link>
            <Link to="/login" className="btn btn-primary flex-1">
              Sign In
            </Link>
          </div>
        </div>
      </div>
    </main>
  );
}

export default PermissionDeniedPage;
