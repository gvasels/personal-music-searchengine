/**
 * Login Page
 * Task 1.5 - Login form with validation and authentication
 */

import { useState, useCallback, useEffect, useRef } from 'react';
import { useNavigate, Link } from '@tanstack/react-router';
import { useAuth } from '../hooks/useAuth';
import { AuthErrorCode } from '../lib/auth';

function getErrorMessage(code: AuthErrorCode): string {
  switch (code) {
    case AuthErrorCode.INVALID_CREDENTIALS:
      return 'Incorrect email or password. Please try again.';
    case AuthErrorCode.USER_NOT_FOUND:
      return 'No account found with this email address.';
    case AuthErrorCode.USER_NOT_CONFIRMED:
      return 'Please verify your email address before signing in.';
    case AuthErrorCode.TOO_MANY_REQUESTS:
      return 'Too many attempts. Please try again later.';
    case AuthErrorCode.NETWORK_ERROR:
      return 'Network error. Please check your connection and try again.';
    case AuthErrorCode.MFA_REQUIRED:
      return 'Multi-factor authentication is required.';
    case AuthErrorCode.NEW_PASSWORD_REQUIRED:
      return 'You must set a new password.';
    default:
      return 'An error occurred. Please try again.';
  }
}

function LoginPage() {
  const navigate = useNavigate();
  const { isLoading: isAuthLoading, isAuthenticated, isSigningIn, signIn, completeNewPassword, error, clearError } = useAuth();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [needsNewPassword, setNeedsNewPassword] = useState(false);
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [validationErrors, setValidationErrors] = useState<{ email?: string; password?: string }>({});
  const [touched, setTouched] = useState<{ email?: boolean; password?: boolean }>({});
  const emailInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    document.title = 'Sign In - Music Search Engine';
  }, []);

  useEffect(() => {
    if (!isAuthenticated && !isAuthLoading) {
      emailInputRef.current?.focus();
    }
  }, [isAuthenticated, isAuthLoading]);

  useEffect(() => {
    if (isAuthenticated) {
      navigate({ to: '/' });
    }
  }, [isAuthenticated, navigate]);

  const validateEmail = (value: string): string | undefined => {
    if (!value.trim()) return 'Email is required';
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return 'Please enter a valid email address';
    return undefined;
  };

  const validatePassword = (value: string): string | undefined => {
    if (!value) return 'Password is required';
    return undefined;
  };

  const handleEmailChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setEmail(value);
    if (touched.email) {
      setValidationErrors((prev) => ({ ...prev, email: validateEmail(value) }));
    }
  }, [touched.email]);

  const handlePasswordChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setPassword(value);
    if (touched.password) {
      setValidationErrors((prev) => ({ ...prev, password: validatePassword(value) }));
    }
  }, [touched.password]);

  const handleEmailBlur = useCallback(() => {
    setTouched((prev) => ({ ...prev, email: true }));
    setValidationErrors((prev) => ({ ...prev, email: validateEmail(email) }));
  }, [email]);

  const handlePasswordBlur = useCallback(() => {
    setTouched((prev) => ({ ...prev, password: true }));
    setValidationErrors((prev) => ({ ...prev, password: validatePassword(password) }));
  }, [password]);

  const handleSubmit = useCallback(async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clearError();

    const emailError = validateEmail(email);
    const passwordError = validatePassword(password);

    setTouched({ email: true, password: true });
    setValidationErrors({ email: emailError, password: passwordError });

    if (emailError || passwordError) return;

    try {
      await signIn(email.trim(), password);
      navigate({ to: '/' });
    } catch (err: unknown) {
      if (err instanceof Error && 'code' in err && (err as { code: string }).code === AuthErrorCode.NEW_PASSWORD_REQUIRED) {
        setNeedsNewPassword(true);
        clearError();
      }
    }
  }, [email, password, signIn, clearError, navigate]);

  const handleNewPasswordSubmit = useCallback(async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clearError();
    if (!newPassword || newPassword.length < 8) return;
    if (newPassword !== confirmPassword) return;
    setIsChangingPassword(true);
    try {
      await completeNewPassword(newPassword);
      navigate({ to: '/' });
    } catch {
      // Error handled by useAuth
    } finally {
      setIsChangingPassword(false);
    }
  }, [newPassword, confirmPassword, completeNewPassword, clearError, navigate]);

  const handleDismissError = useCallback(() => clearError(), [clearError]);

  if (isAuthLoading) {
    return (
      <main className="min-h-screen flex items-center justify-center bg-base-200">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading"></span>
      </main>
    );
  }

  const hasEmailError = touched.email && validationErrors.email;
  const hasPasswordError = touched.password && validationErrors.password;

  return (
    <main id="main-content" className="min-h-screen flex items-center justify-center bg-base-200 p-4">
      <div className="card w-full max-w-md bg-base-100 shadow-xl">
        <div className="card-body">
          <h1 className="card-title text-2xl font-bold text-center justify-center mb-6">
            {needsNewPassword ? 'Set New Password' : 'Sign In'}
          </h1>

          {error && (
            <div role="alert" aria-live="polite" className="alert alert-error mb-4">
              <span>{getErrorMessage(error.code)}</span>
              <button
                type="button"
                className="btn btn-sm btn-ghost"
                onClick={handleDismissError}
                aria-label="Dismiss error"
              >
                ✕
              </button>
            </div>
          )}

          {needsNewPassword ? (
            <form onSubmit={handleNewPasswordSubmit} aria-label="New password form" noValidate>
              <p className="text-sm mb-4">Your account requires a password change.</p>
              <div className="form-control w-full mb-4">
                <label htmlFor="newPassword" className="label">
                  <span className="label-text">New Password</span>
                </label>
                <input
                  type="password"
                  id="newPassword"
                  placeholder="Enter new password"
                  className="input input-bordered w-full"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  disabled={isChangingPassword}
                  autoComplete="new-password"
                  minLength={8}
                  required
                />
              </div>
              <div className="form-control w-full mb-6">
                <label htmlFor="confirmPassword" className="label">
                  <span className="label-text">Confirm Password</span>
                </label>
                <input
                  type="password"
                  id="confirmPassword"
                  placeholder="Confirm new password"
                  className={`input input-bordered w-full ${newPassword && confirmPassword && newPassword !== confirmPassword ? 'input-error' : ''}`}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  disabled={isChangingPassword}
                  autoComplete="new-password"
                  required
                />
                {newPassword && confirmPassword && newPassword !== confirmPassword && (
                  <span className="text-error text-sm mt-1">Passwords do not match</span>
                )}
              </div>
              <button
                type="submit"
                className="btn btn-primary w-full"
                disabled={isChangingPassword || !newPassword || newPassword.length < 8 || newPassword !== confirmPassword}
              >
                {isChangingPassword ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    Changing password...
                  </>
                ) : (
                  'Set New Password'
                )}
              </button>
            </form>
          ) : (
          <>
          <form onSubmit={handleSubmit} aria-label="Login form" noValidate>
            <div className="form-control w-full mb-4">
              <label htmlFor="email" className="label">
                <span className="label-text">Email</span>
              </label>
              <input
                ref={emailInputRef}
                type="email"
                id="email"
                name="email"
                placeholder="Enter your email"
                className={`input input-bordered w-full ${hasEmailError ? 'input-error' : ''}`}
                value={email}
                onChange={handleEmailChange}
                onBlur={handleEmailBlur}
                aria-invalid={hasEmailError ? 'true' : 'false'}
                aria-describedby={hasEmailError ? 'email-error' : undefined}
                disabled={isSigningIn}
                autoComplete="email"
                required
              />
              {hasEmailError && (
                <span className="text-error text-sm mt-1" id="email-error">
                  {validationErrors.email}
                </span>
              )}
            </div>

            <div className="form-control w-full mb-6">
              <label htmlFor="password" className="label">
                <span className="label-text">Password</span>
              </label>
              <input
                type="password"
                id="password"
                name="password"
                placeholder="Enter your password"
                className={`input input-bordered w-full ${hasPasswordError ? 'input-error' : ''}`}
                value={password}
                onChange={handlePasswordChange}
                onBlur={handlePasswordBlur}
                aria-invalid={hasPasswordError ? 'true' : 'false'}
                aria-describedby={hasPasswordError ? 'password-error' : undefined}
                disabled={isSigningIn}
                autoComplete="current-password"
                required
              />
              {hasPasswordError && (
                <span className="text-error text-sm mt-1" id="password-error">
                  {validationErrors.password}
                </span>
              )}
            </div>

            <button
              type="submit"
              className="btn btn-primary w-full"
              disabled={isSigningIn}
            >
              {isSigningIn ? (
                <>
                  <span className="loading loading-spinner loading-sm" role="status"></span>
                  Signing in...
                </>
              ) : (
                'Sign In'
              )}
            </button>
          </form>

          <div className="divider">OR</div>

          <p className="text-center text-sm">
            Don't have an account?{' '}
            <Link to="/signup" className="link link-primary">
              Create account
            </Link>
          </p>
          </>
          )}
        </div>
      </div>
    </main>
  );
}

export default LoginPage;
