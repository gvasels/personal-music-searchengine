/**
 * Auth Module Tests - Task 1.1
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('aws-amplify/auth', () => ({
  signIn: vi.fn(),
  signOut: vi.fn(),
  getCurrentUser: vi.fn(),
  fetchAuthSession: vi.fn(),
}));

vi.mock('aws-amplify', () => ({
  Amplify: {
    configure: vi.fn(),
  },
}));

import {
  configureAuth,
  signIn,
  signOut,
  getCurrentUser,
  getToken,
  refreshToken,
  AuthError,
  AuthErrorCode,
} from '../auth';
import * as AmplifyAuth from 'aws-amplify/auth';
import { Amplify } from 'aws-amplify';

describe('Auth Module (Task 1.1)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('configureAuth', () => {
    it('should configure Amplify with Cognito settings', () => {
      configureAuth({ userPoolId: 'us-east-1_test', userPoolClientId: 'clientId' });
      expect(Amplify.configure).toHaveBeenCalledWith(
        expect.objectContaining({
          Auth: expect.objectContaining({
            Cognito: expect.objectContaining({
              userPoolId: 'us-east-1_test',
              userPoolClientId: 'clientId',
            }),
          }),
        })
      );
    });

    it('should throw AuthError for missing userPoolId', () => {
      expect(() => configureAuth({ userPoolId: '', userPoolClientId: 'client' })).toThrow(AuthError);
    });

    it('should throw AuthError for missing userPoolClientId', () => {
      expect(() => configureAuth({ userPoolId: 'pool', userPoolClientId: '' })).toThrow(AuthError);
    });
  });

  describe('signIn', () => {
    it('should sign in user with valid credentials', async () => {
      vi.mocked(AmplifyAuth.signIn).mockResolvedValue({ isSignedIn: true, nextStep: { signInStep: 'DONE' } });
      vi.mocked(AmplifyAuth.getCurrentUser).mockResolvedValue({
        username: 'user-123',
        userId: 'user-123',
        signInDetails: { loginId: 'test@example.com' },
      });
      vi.mocked(AmplifyAuth.fetchAuthSession).mockResolvedValue({
        tokens: {
          idToken: { payload: { 'cognito:groups': ['subscriber'] }, toString: () => 'id-token' },
        },
      } as Awaited<ReturnType<typeof AmplifyAuth.fetchAuthSession>>);

      const result = await signIn('test@example.com', 'password123');

      expect(AmplifyAuth.signIn).toHaveBeenCalledWith({ username: 'test@example.com', password: 'password123' });
      expect(result).toEqual({ userId: 'user-123', email: 'test@example.com', role: 'subscriber', groups: ['subscriber'] });
    });

    it('should throw INVALID_CREDENTIALS for NotAuthorizedException', async () => {
      vi.mocked(AmplifyAuth.signIn).mockRejectedValue({ name: 'NotAuthorizedException' });

      try {
        await signIn('test@example.com', 'wrong');
        expect.fail('Should throw');
      } catch (error) {
        expect(error).toBeInstanceOf(AuthError);
        expect((error as AuthError).code).toBe(AuthErrorCode.INVALID_CREDENTIALS);
      }
    });

    it('should throw USER_NOT_FOUND for UserNotFoundException', async () => {
      vi.mocked(AmplifyAuth.signIn).mockRejectedValue({ name: 'UserNotFoundException' });

      try {
        await signIn('nonexistent@example.com', 'password');
        expect.fail('Should throw');
      } catch (error) {
        expect((error as AuthError).code).toBe(AuthErrorCode.USER_NOT_FOUND);
      }
    });

    it('should throw for empty email', async () => {
      await expect(signIn('', 'password')).rejects.toThrow(AuthError);
    });

    it('should throw for empty password', async () => {
      await expect(signIn('test@example.com', '')).rejects.toThrow(AuthError);
    });

    it('should throw MFA_REQUIRED for MFA challenge', async () => {
      vi.mocked(AmplifyAuth.signIn).mockResolvedValue({
        isSignedIn: false,
        nextStep: { signInStep: 'CONFIRM_SIGN_IN_WITH_SMS_CODE' },
      } as Awaited<ReturnType<typeof AmplifyAuth.signIn>>);

      try {
        await signIn('test@example.com', 'password');
        expect.fail('Should throw');
      } catch (error) {
        expect((error as AuthError).code).toBe(AuthErrorCode.MFA_REQUIRED);
      }
    });
  });

  describe('signOut', () => {
    it('should sign out the current user', async () => {
      vi.mocked(AmplifyAuth.signOut).mockResolvedValue(undefined);
      await signOut();
      expect(AmplifyAuth.signOut).toHaveBeenCalled();
    });

    it('should sign out globally when option is true', async () => {
      vi.mocked(AmplifyAuth.signOut).mockResolvedValue(undefined);
      await signOut({ global: true });
      expect(AmplifyAuth.signOut).toHaveBeenCalledWith({ global: true });
    });

    it('should throw SIGN_OUT_FAILED on error', async () => {
      vi.mocked(AmplifyAuth.signOut).mockRejectedValue(new Error('Failed'));

      try {
        await signOut();
        expect.fail('Should throw');
      } catch (error) {
        expect((error as AuthError).code).toBe(AuthErrorCode.SIGN_OUT_FAILED);
      }
    });
  });

  describe('getCurrentUser', () => {
    it('should return current authenticated user', async () => {
      vi.mocked(AmplifyAuth.getCurrentUser).mockResolvedValue({
        username: 'user-123',
        userId: 'user-123',
        signInDetails: { loginId: 'test@example.com' },
      });
      vi.mocked(AmplifyAuth.fetchAuthSession).mockResolvedValue({
        tokens: {
          idToken: { payload: { 'cognito:groups': ['artist'] }, toString: () => 'id-token' },
        },
      } as Awaited<ReturnType<typeof AmplifyAuth.fetchAuthSession>>);

      const user = await getCurrentUser();
      expect(user).toEqual({ userId: 'user-123', email: 'test@example.com', role: 'artist', groups: ['artist'] });
    });

    it('should return null when not authenticated', async () => {
      vi.mocked(AmplifyAuth.getCurrentUser).mockRejectedValue({ name: 'UserUnAuthenticatedException' });

      const user = await getCurrentUser();
      expect(user).toBeNull();
    });
  });

  describe('getToken', () => {
    it('should return access token', async () => {
      vi.mocked(AmplifyAuth.fetchAuthSession).mockResolvedValue({
        tokens: {
          accessToken: { toString: () => 'access-token', payload: {} } as unknown as ReturnType<typeof AmplifyAuth.fetchAuthSession> extends Promise<infer T> ? T extends { tokens?: { accessToken?: infer A } } ? A : never : never,
          idToken: { toString: () => 'id-token', payload: {} } as unknown as ReturnType<typeof AmplifyAuth.fetchAuthSession> extends Promise<infer T> ? T extends { tokens?: { idToken?: infer I } } ? I : never : never,
        },
      } as Awaited<ReturnType<typeof AmplifyAuth.fetchAuthSession>>);

      const token = await getToken();
      expect(token).toBe('access-token');
    });

    it('should return id token when useIdToken is true', async () => {
      vi.mocked(AmplifyAuth.fetchAuthSession).mockResolvedValue({
        tokens: {
          accessToken: { toString: () => 'access-token', payload: {} },
          idToken: { toString: () => 'id-token', payload: {} },
        },
      } as Awaited<ReturnType<typeof AmplifyAuth.fetchAuthSession>>);

      const token = await getToken({ useIdToken: true });
      expect(token).toBe('id-token');
    });

    it('should return null when no tokens', async () => {
      vi.mocked(AmplifyAuth.fetchAuthSession).mockResolvedValue({ tokens: undefined } as Awaited<ReturnType<typeof AmplifyAuth.fetchAuthSession>>);

      const token = await getToken();
      expect(token).toBeNull();
    });

    it('should force refresh when option is true', async () => {
      vi.mocked(AmplifyAuth.fetchAuthSession).mockResolvedValue({
        tokens: { accessToken: { toString: () => 'refreshed', payload: {} } },
      } as Awaited<ReturnType<typeof AmplifyAuth.fetchAuthSession>>);

      await getToken({ forceRefresh: true });
      expect(AmplifyAuth.fetchAuthSession).toHaveBeenCalledWith({ forceRefresh: true });
    });
  });

  describe('refreshToken', () => {
    it('should refresh and return new tokens', async () => {
      vi.mocked(AmplifyAuth.fetchAuthSession).mockResolvedValue({
        tokens: {
          accessToken: { toString: () => 'new-access', payload: {} },
          idToken: { toString: () => 'new-id', payload: {} },
        },
      } as Awaited<ReturnType<typeof AmplifyAuth.fetchAuthSession>>);

      const tokens = await refreshToken();

      expect(AmplifyAuth.fetchAuthSession).toHaveBeenCalledWith({ forceRefresh: true });
      expect(tokens).toEqual({ accessToken: 'new-access', idToken: 'new-id' });
    });

    it('should throw TOKEN_REFRESH_FAILED when no tokens', async () => {
      vi.mocked(AmplifyAuth.fetchAuthSession).mockResolvedValue({ tokens: undefined } as Awaited<ReturnType<typeof AmplifyAuth.fetchAuthSession>>);

      try {
        await refreshToken();
        expect.fail('Should throw');
      } catch (error) {
        expect((error as AuthError).code).toBe(AuthErrorCode.TOKEN_REFRESH_FAILED);
      }
    });
  });

  describe('AuthError', () => {
    it('should be instance of Error', () => {
      const error = new AuthError('Test', AuthErrorCode.UNKNOWN);
      expect(error).toBeInstanceOf(Error);
    });

    it('should have correct name', () => {
      const error = new AuthError('Test', AuthErrorCode.UNKNOWN);
      expect(error.name).toBe('AuthError');
    });

    it('should contain error code', () => {
      const error = new AuthError('Test', AuthErrorCode.INVALID_CREDENTIALS);
      expect(error.code).toBe(AuthErrorCode.INVALID_CREDENTIALS);
    });

    it('should contain original error', () => {
      const original = new Error('Original');
      const error = new AuthError('Test', AuthErrorCode.UNKNOWN, original);
      expect(error.originalError).toBe(original);
    });
  });

  describe('AuthErrorCode enum', () => {
    it('should have all required codes', () => {
      expect(AuthErrorCode.INVALID_CREDENTIALS).toBeDefined();
      expect(AuthErrorCode.USER_NOT_FOUND).toBeDefined();
      expect(AuthErrorCode.USER_NOT_CONFIRMED).toBeDefined();
      expect(AuthErrorCode.TOO_MANY_REQUESTS).toBeDefined();
      expect(AuthErrorCode.NETWORK_ERROR).toBeDefined();
      expect(AuthErrorCode.TOKEN_EXPIRED).toBeDefined();
      expect(AuthErrorCode.TOKEN_REFRESH_FAILED).toBeDefined();
      expect(AuthErrorCode.AUTH_CONFIG_INVALID).toBeDefined();
      expect(AuthErrorCode.SIGN_OUT_FAILED).toBeDefined();
      expect(AuthErrorCode.MFA_REQUIRED).toBeDefined();
      expect(AuthErrorCode.UNKNOWN).toBeDefined();
    });
  });
});
