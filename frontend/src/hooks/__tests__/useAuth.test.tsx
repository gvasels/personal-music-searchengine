/**
 * useAuth Hook Tests - Task 1.4
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/auth', () => ({
  signIn: vi.fn(),
  signOut: vi.fn(),
  getCurrentUser: vi.fn(),
  AuthError: class AuthError extends Error {
    code: string;
    constructor(message: string, code: string) {
      super(message);
      this.code = code;
      this.name = 'AuthError';
    }
  },
  AuthErrorCode: {
    INVALID_CREDENTIALS: 'INVALID_CREDENTIALS',
    USER_NOT_FOUND: 'USER_NOT_FOUND',
    NETWORK_ERROR: 'NETWORK_ERROR',
    UNKNOWN: 'UNKNOWN',
  },
}));

import { useAuth } from '../useAuth';
import * as auth from '../../lib/auth';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0, staleTime: 0 },
      mutations: { retry: false },
    },
  });

  return function Wrapper({ children }: { children: ReactNode }) {
    return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
  };
}

describe('useAuth Hook (Task 1.4)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('initial state', () => {
    it('should return isLoading true initially', async () => {
      vi.mocked(auth.getCurrentUser).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(true);
    });

    it('should return user as null initially', () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      expect(result.current.user).toBeNull();
    });
  });

  describe('authenticated state', () => {
    const mockUser = { userId: 'user-123', email: 'test@example.com' };

    it('should return user when authenticated', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(mockUser);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.user).toEqual(mockUser);
      });
    });

    it('should return isAuthenticated true when user exists', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(mockUser);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isAuthenticated).toBe(true);
      });
    });

    it('should return isLoading false after user is fetched', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(mockUser);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });
    });
  });

  describe('unauthenticated state', () => {
    it('should return user as null when not authenticated', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.user).toBeNull();
      });
    });

    it('should return isAuthenticated false when no user', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isAuthenticated).toBe(false);
      });
    });
  });

  describe('signIn function', () => {
    const mockUser = { userId: 'user-123', email: 'test@example.com' };

    it('should be a function', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      expect(typeof result.current.signIn).toBe('function');
    });

    it('should call auth.signIn with email and password', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);
      vi.mocked(auth.signIn).mockResolvedValue(mockUser);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      await act(async () => {
        await result.current.signIn('test@example.com', 'password123');
      });

      expect(auth.signIn).toHaveBeenCalledWith('test@example.com', 'password123');
    });

    it('should update user state after successful sign in', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);
      vi.mocked(auth.signIn).mockResolvedValue(mockUser);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      await act(async () => {
        await result.current.signIn('test@example.com', 'password123');
      });

      await waitFor(() => {
        expect(result.current.user).toEqual(mockUser);
        expect(result.current.isAuthenticated).toBe(true);
      });
    });

    it('should set error on sign in failure', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);
      vi.mocked(auth.signIn).mockRejectedValue(new auth.AuthError('Invalid', auth.AuthErrorCode.INVALID_CREDENTIALS));

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      await act(async () => {
        try {
          await result.current.signIn('test@example.com', 'wrong');
        } catch {}
      });

      await waitFor(() => {
        expect(result.current.error).toBeTruthy();
      });
    });
  });

  describe('signOut function', () => {
    const mockUser = { userId: 'user-123', email: 'test@example.com' };

    it('should be a function', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(mockUser);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      expect(typeof result.current.signOut).toBe('function');
    });

    it('should call auth.signOut', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(mockUser);
      vi.mocked(auth.signOut).mockResolvedValue(undefined);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.isAuthenticated).toBe(true));

      await act(async () => {
        await result.current.signOut();
      });

      expect(auth.signOut).toHaveBeenCalled();
    });

    it('should clear user state after sign out', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(mockUser);
      vi.mocked(auth.signOut).mockResolvedValue(undefined);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.isAuthenticated).toBe(true));

      await act(async () => {
        await result.current.signOut();
      });

      await waitFor(() => {
        expect(result.current.user).toBeNull();
        expect(result.current.isAuthenticated).toBe(false);
      });
    });
  });

  describe('error handling', () => {
    it('should clear error when clearError is called', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);
      vi.mocked(auth.signIn).mockRejectedValue(new Error('Error'));

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      await act(async () => {
        try {
          await result.current.signIn('test@example.com', 'wrong');
        } catch {}
      });

      await waitFor(() => expect(result.current.error).toBeTruthy());

      act(() => {
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe('refetch functionality', () => {
    it('should provide refetch function', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      expect(typeof result.current.refetch).toBe('function');
    });
  });

  describe('isSigningIn state', () => {
    it('should return isSigningIn false initially', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue(null);

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.isSigningIn).toBe(false);
    });
  });

  describe('isSigningOut state', () => {
    it('should return isSigningOut false initially', async () => {
      vi.mocked(auth.getCurrentUser).mockResolvedValue({ userId: '123', email: 'test@example.com' });

      const { result } = renderHook(() => useAuth(), { wrapper: createWrapper() });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.isSigningOut).toBe(false);
    });
  });
});
