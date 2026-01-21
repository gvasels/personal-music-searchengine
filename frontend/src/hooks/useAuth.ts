/**
 * useAuth Hook
 * Task 1.4 - Authentication hook using TanStack Query
 */

import { useState, useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  signIn as authSignIn,
  signOut as authSignOut,
  getCurrentUser,
  AuthUser,
  AuthError,
  AuthErrorCode,
} from '../lib/auth';

export interface UseAuthReturn {
  user: AuthUser | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  isSigningIn: boolean;
  isSigningOut: boolean;
  signIn: (email: string, password: string) => Promise<AuthUser>;
  signOut: () => Promise<void>;
  error: { message: string; code: AuthErrorCode } | null;
  clearError: () => void;
  refetch: () => Promise<unknown>;
}

export function useAuth(): UseAuthReturn {
  const queryClient = useQueryClient();
  const [error, setError] = useState<{ message: string; code: AuthErrorCode } | null>(null);

  const {
    data: user,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ['auth', 'user'],
    queryFn: getCurrentUser,
    staleTime: 5 * 60 * 1000,
    retry: false,
  });

  const signInMutation = useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      authSignIn(email, password),
    onSuccess: (data) => {
      queryClient.setQueryData(['auth', 'user'], data);
      setError(null);
    },
    onError: (err: unknown) => {
      if (err instanceof AuthError) {
        setError({ message: err.message, code: err.code });
      } else {
        setError({ message: 'Sign in failed', code: AuthErrorCode.UNKNOWN });
      }
    },
  });

  const signOutMutation = useMutation({
    mutationFn: () => authSignOut(),
    onSuccess: () => {
      queryClient.setQueryData(['auth', 'user'], null);
      setError(null);
    },
    onError: (err: unknown) => {
      if (err instanceof AuthError) {
        setError({ message: err.message, code: err.code });
      } else {
        setError({ message: 'Sign out failed', code: AuthErrorCode.UNKNOWN });
      }
    },
  });

  const signIn = useCallback(
    async (email: string, password: string): Promise<AuthUser> => {
      setError(null);
      return signInMutation.mutateAsync({ email, password });
    },
    [signInMutation]
  );

  const signOut = useCallback(async (): Promise<void> => {
    return signOutMutation.mutateAsync();
  }, [signOutMutation]);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    user: user ?? null,
    isLoading,
    isAuthenticated: !!user,
    isSigningIn: signInMutation.isPending,
    isSigningOut: signOutMutation.isPending,
    signIn,
    signOut,
    error,
    clearError,
    refetch,
  };
}
