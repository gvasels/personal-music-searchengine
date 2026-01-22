/**
 * Signup Page
 * User registration with email verification
 */

import { useState, useCallback, useEffect, useRef } from 'react';
import { useNavigate, Link } from '@tanstack/react-router';
import { signUp, confirmSignUp, resendSignUpCode, AuthError, AuthErrorCode } from '../lib/auth';

function getErrorMessage(code: AuthErrorCode): string {
  switch (code) {
    case AuthErrorCode.USER_ALREADY_EXISTS:
      return 'An account with this email already exists.';
    case AuthErrorCode.INVALID_PASSWORD:
      return 'Password must be at least 8 characters with uppercase, lowercase, and numbers.';
    case AuthErrorCode.INVALID_CODE:
      return 'Invalid verification code. Please try again.';
    case AuthErrorCode.CODE_EXPIRED:
      return 'Verification code has expired. Please request a new one.';
    case AuthErrorCode.TOO_MANY_REQUESTS:
      return 'Too many attempts. Please try again later.';
    case AuthErrorCode.NETWORK_ERROR:
      return 'Network error. Please check your connection and try again.';
    default:
      return 'An error occurred. Please try again.';
  }
}

type Step = 'signup' | 'verify';

function SignupPage() {
  const navigate = useNavigate();
  const [step, setStep] = useState<Step>('signup');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<{ message: string; code: AuthErrorCode } | null>(null);
  const [validationErrors, setValidationErrors] = useState<{
    email?: string;
    password?: string;
    confirmPassword?: string;
    verificationCode?: string;
  }>({});
  const [touched, setTouched] = useState<{
    email?: boolean;
    password?: boolean;
    confirmPassword?: boolean;
    verificationCode?: boolean;
  }>({});
  const emailInputRef = useRef<HTMLInputElement>(null);
  const codeInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    document.title = 'Create Account - Music Search Engine';
  }, []);

  useEffect(() => {
    if (step === 'signup') {
      emailInputRef.current?.focus();
    } else {
      codeInputRef.current?.focus();
    }
  }, [step]);

  const validateEmail = (value: string): string | undefined => {
    if (!value.trim()) return 'Email is required';
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return 'Please enter a valid email address';
    return undefined;
  };

  const validatePassword = (value: string): string | undefined => {
    if (!value) return 'Password is required';
    if (value.length < 8) return 'Password must be at least 8 characters';
    if (!/[A-Z]/.test(value)) return 'Password must contain an uppercase letter';
    if (!/[a-z]/.test(value)) return 'Password must contain a lowercase letter';
    if (!/[0-9]/.test(value)) return 'Password must contain a number';
    return undefined;
  };

  const validateConfirmPassword = (value: string): string | undefined => {
    if (!value) return 'Please confirm your password';
    if (value !== password) return 'Passwords do not match';
    return undefined;
  };

  const validateVerificationCode = (value: string): string | undefined => {
    if (!value.trim()) return 'Verification code is required';
    if (!/^\d{6}$/.test(value.trim())) return 'Code must be 6 digits';
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
    if (touched.confirmPassword && confirmPassword) {
      setValidationErrors((prev) => ({
        ...prev,
        confirmPassword: value !== confirmPassword ? 'Passwords do not match' : undefined,
      }));
    }
  }, [touched.password, touched.confirmPassword, confirmPassword]);

  const handleConfirmPasswordChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setConfirmPassword(value);
    if (touched.confirmPassword) {
      setValidationErrors((prev) => ({
        ...prev,
        confirmPassword: value !== password ? 'Passwords do not match' : undefined,
      }));
    }
  }, [touched.confirmPassword, password]);

  const handleVerificationCodeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 6);
    setVerificationCode(value);
    if (touched.verificationCode) {
      setValidationErrors((prev) => ({ ...prev, verificationCode: validateVerificationCode(value) }));
    }
  }, [touched.verificationCode]);

  const handleBlur = useCallback((field: keyof typeof touched) => {
    setTouched((prev) => ({ ...prev, [field]: true }));
    if (field === 'email') {
      setValidationErrors((prev) => ({ ...prev, email: validateEmail(email) }));
    } else if (field === 'password') {
      setValidationErrors((prev) => ({ ...prev, password: validatePassword(password) }));
    } else if (field === 'confirmPassword') {
      setValidationErrors((prev) => ({ ...prev, confirmPassword: validateConfirmPassword(confirmPassword) }));
    } else if (field === 'verificationCode') {
      setValidationErrors((prev) => ({ ...prev, verificationCode: validateVerificationCode(verificationCode) }));
    }
  }, [email, password, confirmPassword, verificationCode]);

  const handleSignupSubmit = useCallback(async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);

    const emailError = validateEmail(email);
    const passwordError = validatePassword(password);
    const confirmPasswordError = validateConfirmPassword(confirmPassword);

    setTouched({ email: true, password: true, confirmPassword: true });
    setValidationErrors({ email: emailError, password: passwordError, confirmPassword: confirmPasswordError });

    if (emailError || passwordError || confirmPasswordError) return;

    setIsSubmitting(true);
    try {
      const result = await signUp(email.trim(), password);
      if (result.nextStep === 'CONFIRM_SIGN_UP') {
        setStep('verify');
      } else if (result.isSignUpComplete) {
        navigate({ to: '/login' });
      }
    } catch (err) {
      if (err instanceof AuthError) {
        setError({ message: err.message, code: err.code });
      } else {
        setError({ message: 'Sign up failed', code: AuthErrorCode.UNKNOWN });
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [email, password, confirmPassword, navigate]);

  const handleVerifySubmit = useCallback(async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);

    const codeError = validateVerificationCode(verificationCode);
    setTouched((prev) => ({ ...prev, verificationCode: true }));
    setValidationErrors((prev) => ({ ...prev, verificationCode: codeError }));

    if (codeError) return;

    setIsSubmitting(true);
    try {
      await confirmSignUp(email.trim(), verificationCode.trim());
      navigate({ to: '/login' });
    } catch (err) {
      if (err instanceof AuthError) {
        setError({ message: err.message, code: err.code });
      } else {
        setError({ message: 'Verification failed', code: AuthErrorCode.UNKNOWN });
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [email, verificationCode, navigate]);

  const handleResendCode = useCallback(async () => {
    setError(null);
    setIsSubmitting(true);
    try {
      await resendSignUpCode(email.trim());
      setError({ message: 'Verification code sent!', code: AuthErrorCode.UNKNOWN });
    } catch (err) {
      if (err instanceof AuthError) {
        setError({ message: err.message, code: err.code });
      } else {
        setError({ message: 'Failed to resend code', code: AuthErrorCode.UNKNOWN });
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [email]);

  const handleDismissError = useCallback(() => setError(null), []);

  const hasEmailError = touched.email && validationErrors.email;
  const hasPasswordError = touched.password && validationErrors.password;
  const hasConfirmPasswordError = touched.confirmPassword && validationErrors.confirmPassword;
  const hasVerificationCodeError = touched.verificationCode && validationErrors.verificationCode;

  return (
    <main id="main-content" className="min-h-screen flex items-center justify-center bg-base-200 p-4">
      <div className="card w-full max-w-md bg-base-100 shadow-xl">
        <div className="card-body">
          <h1 className="card-title text-2xl font-bold text-center justify-center mb-6">
            {step === 'signup' ? 'Create Account' : 'Verify Email'}
          </h1>

          {error && (
            <div
              role="alert"
              aria-live="polite"
              className={`alert ${error.message.includes('sent') ? 'alert-success' : 'alert-error'} mb-4`}
            >
              <span>{error.message.includes('sent') ? error.message : getErrorMessage(error.code)}</span>
              <button
                type="button"
                className="btn btn-sm btn-ghost"
                onClick={handleDismissError}
                aria-label="Dismiss message"
              >
                ✕
              </button>
            </div>
          )}

          {step === 'signup' ? (
            <form onSubmit={handleSignupSubmit} aria-label="Signup form" noValidate>
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
                  onBlur={() => handleBlur('email')}
                  aria-invalid={hasEmailError ? 'true' : 'false'}
                  aria-describedby={hasEmailError ? 'email-error' : undefined}
                  disabled={isSubmitting}
                  autoComplete="email"
                  required
                />
                {hasEmailError && (
                  <span className="text-error text-sm mt-1" id="email-error">
                    {validationErrors.email}
                  </span>
                )}
              </div>

              <div className="form-control w-full mb-4">
                <label htmlFor="password" className="label">
                  <span className="label-text">Password</span>
                </label>
                <input
                  type="password"
                  id="password"
                  name="password"
                  placeholder="Create a password"
                  className={`input input-bordered w-full ${hasPasswordError ? 'input-error' : ''}`}
                  value={password}
                  onChange={handlePasswordChange}
                  onBlur={() => handleBlur('password')}
                  aria-invalid={hasPasswordError ? 'true' : 'false'}
                  aria-describedby={hasPasswordError ? 'password-error' : undefined}
                  disabled={isSubmitting}
                  autoComplete="new-password"
                  required
                />
                {hasPasswordError && (
                  <span className="text-error text-sm mt-1" id="password-error">
                    {validationErrors.password}
                  </span>
                )}
              </div>

              <div className="form-control w-full mb-6">
                <label htmlFor="confirmPassword" className="label">
                  <span className="label-text">Confirm Password</span>
                </label>
                <input
                  type="password"
                  id="confirmPassword"
                  name="confirmPassword"
                  placeholder="Confirm your password"
                  className={`input input-bordered w-full ${hasConfirmPasswordError ? 'input-error' : ''}`}
                  value={confirmPassword}
                  onChange={handleConfirmPasswordChange}
                  onBlur={() => handleBlur('confirmPassword')}
                  aria-invalid={hasConfirmPasswordError ? 'true' : 'false'}
                  aria-describedby={hasConfirmPasswordError ? 'confirmPassword-error' : undefined}
                  disabled={isSubmitting}
                  autoComplete="new-password"
                  required
                />
                {hasConfirmPasswordError && (
                  <span className="text-error text-sm mt-1" id="confirmPassword-error">
                    {validationErrors.confirmPassword}
                  </span>
                )}
              </div>

              <button type="submit" className="btn btn-primary w-full" disabled={isSubmitting}>
                {isSubmitting ? (
                  <>
                    <span className="loading loading-spinner loading-sm" role="status"></span>
                    Creating account...
                  </>
                ) : (
                  'Create Account'
                )}
              </button>
            </form>
          ) : (
            <form onSubmit={handleVerifySubmit} aria-label="Verification form" noValidate>
              <p className="text-base-content/70 mb-4 text-center">
                We sent a verification code to <strong>{email}</strong>
              </p>

              <div className="form-control w-full mb-6">
                <label htmlFor="verificationCode" className="label">
                  <span className="label-text">Verification Code</span>
                </label>
                <input
                  ref={codeInputRef}
                  type="text"
                  id="verificationCode"
                  name="verificationCode"
                  placeholder="Enter 6-digit code"
                  className={`input input-bordered w-full text-center text-2xl tracking-widest ${hasVerificationCodeError ? 'input-error' : ''}`}
                  value={verificationCode}
                  onChange={handleVerificationCodeChange}
                  onBlur={() => handleBlur('verificationCode')}
                  aria-invalid={hasVerificationCodeError ? 'true' : 'false'}
                  aria-describedby={hasVerificationCodeError ? 'verificationCode-error' : undefined}
                  disabled={isSubmitting}
                  maxLength={6}
                  inputMode="numeric"
                  autoComplete="one-time-code"
                  required
                />
                {hasVerificationCodeError && (
                  <span className="text-error text-sm mt-1" id="verificationCode-error">
                    {validationErrors.verificationCode}
                  </span>
                )}
              </div>

              <button type="submit" className="btn btn-primary w-full mb-4" disabled={isSubmitting}>
                {isSubmitting ? (
                  <>
                    <span className="loading loading-spinner loading-sm" role="status"></span>
                    Verifying...
                  </>
                ) : (
                  'Verify Email'
                )}
              </button>

              <button
                type="button"
                className="btn btn-ghost w-full"
                onClick={handleResendCode}
                disabled={isSubmitting}
              >
                Resend Code
              </button>
            </form>
          )}

          <div className="divider">OR</div>

          <p className="text-center text-sm">
            Already have an account?{' '}
            <Link to="/login" className="link link-primary">
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </main>
  );
}

export default SignupPage;
