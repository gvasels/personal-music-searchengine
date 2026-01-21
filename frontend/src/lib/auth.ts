/**
 * AWS Amplify Authentication Module
 * Task 1.1 - Amplify Auth Configuration
 */

import { signIn as amplifySignIn, signOut as amplifySignOut, getCurrentUser as amplifyGetCurrentUser, fetchAuthSession } from 'aws-amplify/auth';
import { Amplify } from 'aws-amplify';

export enum AuthErrorCode {
  INVALID_CREDENTIALS = 'INVALID_CREDENTIALS',
  USER_NOT_FOUND = 'USER_NOT_FOUND',
  USER_NOT_CONFIRMED = 'USER_NOT_CONFIRMED',
  TOO_MANY_REQUESTS = 'TOO_MANY_REQUESTS',
  NETWORK_ERROR = 'NETWORK_ERROR',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
  TOKEN_REFRESH_FAILED = 'TOKEN_REFRESH_FAILED',
  AUTH_CONFIG_INVALID = 'AUTH_CONFIG_INVALID',
  SIGN_OUT_FAILED = 'SIGN_OUT_FAILED',
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
    return {
      userId: user.userId,
      email: user.signInDetails?.loginId || email.trim(),
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
    return {
      userId: user.userId,
      email: user.signInDetails?.loginId || '',
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
