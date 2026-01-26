/**
 * AWS Amplify Authentication Module
 * Task 1.1 - Amplify Auth Configuration
 */

import { signIn as amplifySignIn, signOut as amplifySignOut, getCurrentUser as amplifyGetCurrentUser, fetchAuthSession, signUp as amplifySignUp, confirmSignUp as amplifyConfirmSignUp, resendSignUpCode as amplifyResendSignUpCode } from 'aws-amplify/auth';
import { Amplify } from 'aws-amplify';
import type { UserRole } from '../types';

export enum AuthErrorCode {
  INVALID_CREDENTIALS = 'INVALID_CREDENTIALS',
  USER_NOT_FOUND = 'USER_NOT_FOUND',
  USER_NOT_CONFIRMED = 'USER_NOT_CONFIRMED',
  USER_ALREADY_EXISTS = 'USER_ALREADY_EXISTS',
  INVALID_PASSWORD = 'INVALID_PASSWORD',
  INVALID_CODE = 'INVALID_CODE',
  CODE_EXPIRED = 'CODE_EXPIRED',
  TOO_MANY_REQUESTS = 'TOO_MANY_REQUESTS',
  NETWORK_ERROR = 'NETWORK_ERROR',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
  TOKEN_REFRESH_FAILED = 'TOKEN_REFRESH_FAILED',
  AUTH_CONFIG_INVALID = 'AUTH_CONFIG_INVALID',
  SIGN_OUT_FAILED = 'SIGN_OUT_FAILED',
  SIGN_UP_FAILED = 'SIGN_UP_FAILED',
  CONFIRMATION_FAILED = 'CONFIRMATION_FAILED',
  MFA_REQUIRED = 'MFA_REQUIRED',
  UNKNOWN = 'UNKNOWN',
}

export class AuthError extends Error {
  code: AuthErrorCode;
  originalError?: Error;

  constructor(message: string, code: AuthErrorCode, originalError?: Error) {
    super(message);
    this.name = 'AuthError';
    this.code = code;
    this.originalError = originalError;
  }
}

export interface AuthConfig {
  userPoolId: string;
  userPoolClientId: string;
  identityPoolId?: string;
}

export interface AuthUser {
  userId: string;
  email: string;
  name?: string;
  role: UserRole;
  groups: string[];
}

export interface TokenOptions {
  forceRefresh?: boolean;
  useIdToken?: boolean;
}

export interface SignOutOptions {
  global?: boolean;
}

export interface Tokens {
  accessToken: string;
  idToken: string;
}

/**
 * Extract user role from Cognito groups.
 * Priority: admin > artist > subscriber > guest
 */
function roleFromGroups(groups: string[]): UserRole {
  const lowerGroups = groups.map((g) => g.toLowerCase());

  if (lowerGroups.includes('admin') || lowerGroups.includes('admins')) {
    return 'admin';
  }
  if (lowerGroups.includes('artist') || lowerGroups.includes('artists')) {
    return 'artist';
  }
  if (lowerGroups.includes('subscriber') || lowerGroups.includes('subscribers')) {
    return 'subscriber';
  }
  return 'guest';
}

/**
 * Extract Cognito groups from session tokens.
 */
async function getGroupsFromSession(): Promise<string[]> {
  try {
    const session = await fetchAuthSession();
    if (!session.tokens?.idToken) {
      console.log('[Auth] No ID token in session');
      return [];
    }
    // Groups are in the cognito:groups claim of the ID token
    const payload = session.tokens.idToken.payload;
    const groups = payload['cognito:groups'];
    console.log('[Auth] Token groups claim:', groups);
    if (Array.isArray(groups)) {
      return groups as string[];
    }
    if (typeof groups === 'string') {
      return groups.split(' ');
    }
    return [];
  } catch (err) {
    console.error('[Auth] Error getting groups:', err);
    return [];
  }
}

function mapAmplifyError(error: unknown): AuthError {
  const err = error as { name?: string; message?: string };
  const name = err?.name || '';
  const message = err?.message || 'Unknown error';

  switch (name) {
    case 'NotAuthorizedException':
      return new AuthError('Invalid credentials', AuthErrorCode.INVALID_CREDENTIALS, error as Error);
    case 'UserNotFoundException':
      return new AuthError('User not found', AuthErrorCode.USER_NOT_FOUND, error as Error);
    case 'UserNotConfirmedException':
      return new AuthError('User not confirmed', AuthErrorCode.USER_NOT_CONFIRMED, error as Error);
    case 'UsernameExistsException':
      return new AuthError('An account with this email already exists', AuthErrorCode.USER_ALREADY_EXISTS, error as Error);
    case 'InvalidPasswordException':
      return new AuthError('Password does not meet requirements', AuthErrorCode.INVALID_PASSWORD, error as Error);
    case 'CodeMismatchException':
      return new AuthError('Invalid verification code', AuthErrorCode.INVALID_CODE, error as Error);
    case 'ExpiredCodeException':
      return new AuthError('Verification code has expired', AuthErrorCode.CODE_EXPIRED, error as Error);
    case 'TooManyRequestsException':
      return new AuthError('Too many requests', AuthErrorCode.TOO_MANY_REQUESTS, error as Error);
    case 'NetworkError':
      return new AuthError('Network error', AuthErrorCode.NETWORK_ERROR, error as Error);
    case 'TokenExpiredException':
      return new AuthError('Token expired', AuthErrorCode.TOKEN_EXPIRED, error as Error);
    default:
      return new AuthError(message, AuthErrorCode.UNKNOWN, error as Error);
  }
}

export function configureAuth(config: AuthConfig): void {
  if (!config.userPoolId || !config.userPoolClientId) {
    throw new AuthError('Invalid auth configuration', AuthErrorCode.AUTH_CONFIG_INVALID);
  }

  Amplify.configure({
    Auth: {
      Cognito: {
        userPoolId: config.userPoolId,
        userPoolClientId: config.userPoolClientId,
        identityPoolId: config.identityPoolId || '',
      },
    },
  } as Parameters<typeof Amplify.configure>[0]);
}

export async function signIn(email: string, password: string): Promise<AuthUser> {
  if (!email || !email.trim()) {
    throw new AuthError('Email is required', AuthErrorCode.INVALID_CREDENTIALS);
  }
  if (!password) {
    throw new AuthError('Password is required', AuthErrorCode.INVALID_CREDENTIALS);
  }

  try {
    const result = await amplifySignIn({
      username: email.trim(),
      password,
    });

    if (!result.isSignedIn) {
      const step = result.nextStep?.signInStep || '';
      if (step.includes('MFA') || step.includes('SMS') || step.includes('TOTP')) {
        throw new AuthError('MFA required', AuthErrorCode.MFA_REQUIRED);
      }
      throw new AuthError('Sign in incomplete', AuthErrorCode.UNKNOWN);
    }

    const user = await amplifyGetCurrentUser();
    const groups = await getGroupsFromSession();
    return {
      userId: user.userId,
      email: user.signInDetails?.loginId || email.trim(),
      role: roleFromGroups(groups),
      groups,
    };
  } catch (error) {
    if (error instanceof AuthError) throw error;
    throw mapAmplifyError(error);
  }
}

export async function signOut(options?: SignOutOptions): Promise<void> {
  try {
    if (options?.global) {
      await amplifySignOut({ global: true });
    } else {
      await amplifySignOut();
    }
  } catch (error) {
    throw new AuthError('Sign out failed', AuthErrorCode.SIGN_OUT_FAILED, error as Error);
  }
}

export async function getCurrentUser(): Promise<AuthUser | null> {
  try {
    const user = await amplifyGetCurrentUser();
    const groups = await getGroupsFromSession();
    return {
      userId: user.userId,
      email: user.signInDetails?.loginId || '',
      role: roleFromGroups(groups),
      groups,
    };
  } catch (error) {
    const err = error as { name?: string };
    if (err?.name === 'UserUnAuthenticatedException') {
      return null;
    }
    throw mapAmplifyError(error);
  }
}

export async function getToken(options?: TokenOptions): Promise<string | null> {
  try {
    const session = await fetchAuthSession(
      options?.forceRefresh ? { forceRefresh: true } : undefined
    );

    if (!session.tokens) {
      return null;
    }

    if (options?.useIdToken) {
      return session.tokens.idToken?.toString() || null;
    }

    return session.tokens.accessToken?.toString() || null;
  } catch (error) {
    const err = error as { name?: string };
    if (err?.name === 'TokenExpiredException') {
      throw new AuthError('Token expired', AuthErrorCode.TOKEN_EXPIRED, error as Error);
    }
    throw mapAmplifyError(error);
  }
}

export async function refreshToken(): Promise<Tokens> {
  try {
    const session = await fetchAuthSession({ forceRefresh: true });

    if (!session.tokens) {
      throw new AuthError('No tokens available', AuthErrorCode.TOKEN_REFRESH_FAILED);
    }

    return {
      accessToken: session.tokens.accessToken?.toString() || '',
      idToken: session.tokens.idToken?.toString() || '',
    };
  } catch (error) {
    if (error instanceof AuthError) throw error;
    throw new AuthError('Token refresh failed', AuthErrorCode.TOKEN_REFRESH_FAILED, error as Error);
  }
}

export interface SignUpResult {
  isSignUpComplete: boolean;
  userId?: string;
  nextStep: string;
}

export async function signUp(email: string, password: string): Promise<SignUpResult> {
  if (!email || !email.trim()) {
    throw new AuthError('Email is required', AuthErrorCode.INVALID_CREDENTIALS);
  }
  if (!password) {
    throw new AuthError('Password is required', AuthErrorCode.INVALID_CREDENTIALS);
  }
  if (password.length < 8) {
    throw new AuthError('Password must be at least 8 characters', AuthErrorCode.INVALID_PASSWORD);
  }

  try {
    const result = await amplifySignUp({
      username: email.trim(),
      password,
      options: {
        userAttributes: {
          email: email.trim(),
        },
      },
    });

    return {
      isSignUpComplete: result.isSignUpComplete,
      userId: result.userId,
      nextStep: result.nextStep?.signUpStep || 'DONE',
    };
  } catch (error) {
    if (error instanceof AuthError) throw error;
    throw mapAmplifyError(error);
  }
}

export async function confirmSignUp(email: string, code: string): Promise<boolean> {
  if (!email || !email.trim()) {
    throw new AuthError('Email is required', AuthErrorCode.INVALID_CREDENTIALS);
  }
  if (!code || !code.trim()) {
    throw new AuthError('Verification code is required', AuthErrorCode.INVALID_CODE);
  }

  try {
    const result = await amplifyConfirmSignUp({
      username: email.trim(),
      confirmationCode: code.trim(),
    });

    return result.isSignUpComplete;
  } catch (error) {
    if (error instanceof AuthError) throw error;
    throw mapAmplifyError(error);
  }
}

export async function resendSignUpCode(email: string): Promise<void> {
  if (!email || !email.trim()) {
    throw new AuthError('Email is required', AuthErrorCode.INVALID_CREDENTIALS);
  }

  try {
    await amplifyResendSignUpCode({
      username: email.trim(),
    });
  } catch (error) {
    if (error instanceof AuthError) throw error;
    throw mapAmplifyError(error);
  }
}
