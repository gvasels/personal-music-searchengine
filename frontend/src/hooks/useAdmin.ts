/**
 * useAdmin Hook - Admin Panel & Track Visibility Feature
 * Provides TanStack Query hooks for admin user management operations
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useState, useEffect, useMemo } from 'react';
import {
  searchUsers,
  getUserDetails,
  updateUserRole,
  updateUserStatus,
  type SearchUsersParams,
  type UserSummary,
  type UserDetails,
} from '../lib/api/admin';
import type { UserRole } from '../types';

// Query key factory for admin operations
export const adminKeys = {
  all: ['admin'] as const,
  users: () => [...adminKeys.all, 'users'] as const,
  userSearch: (params: SearchUsersParams) => [...adminKeys.users(), 'search', params] as const,
  userDetails: () => [...adminKeys.users(), 'details'] as const,
  userDetail: (id: string) => [...adminKeys.userDetails(), id] as const,
};

/**
 * Hook for debounced search value
 */
function useDebouncedValue<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

/**
 * Hook for searching users with debounce
 * @param search - Search query string
 * @param options - Query options including limit and cursor
 * @param debounceMs - Debounce delay in milliseconds (default: 300)
 */
export function useSearchUsers(
  search: string,
  options?: { limit?: number; cursor?: string },
  debounceMs: number = 300
) {
  const debouncedSearch = useDebouncedValue(search, debounceMs);

  const params: SearchUsersParams = useMemo(
    () => ({
      search: debouncedSearch,
      limit: options?.limit,
      cursor: options?.cursor,
    }),
    [debouncedSearch, options?.limit, options?.cursor]
  );

  return useQuery({
    queryKey: adminKeys.userSearch(params),
    queryFn: () => searchUsers(params),
    enabled: debouncedSearch.length >= 1, // Only search when there's input
    staleTime: 1000 * 60, // 1 minute
  });
}

/**
 * Hook for fetching user details by ID
 * @param userId - User ID to fetch details for
 */
export function useUserDetails(userId: string | undefined) {
  return useQuery({
    queryKey: adminKeys.userDetail(userId!),
    queryFn: () => getUserDetails(userId!),
    enabled: !!userId,
    staleTime: 1000 * 60, // 1 minute
  });
}

/**
 * Hook for updating a user's role
 * Invalidates user search and detail queries on success
 */
export function useUpdateUserRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: UserRole }) =>
      updateUserRole(userId, role),
    onSuccess: (updatedUser, { userId }) => {
      // Update the user details cache with the new data
      queryClient.setQueryData<UserDetails>(adminKeys.userDetail(userId), updatedUser);
      // Invalidate search results to reflect role change
      void queryClient.invalidateQueries({ queryKey: adminKeys.users() });
    },
  });
}

/**
 * Hook for updating a user's status (enable/disable)
 * Invalidates user search and detail queries on success
 */
export function useUpdateUserStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userId, disabled }: { userId: string; disabled: boolean }) =>
      updateUserStatus(userId, disabled),
    onSuccess: (updatedUser, { userId }) => {
      // Update the user details cache with the new data
      queryClient.setQueryData<UserDetails>(adminKeys.userDetail(userId), updatedUser);
      // Invalidate search results to reflect status change
      void queryClient.invalidateQueries({ queryKey: adminKeys.users() });
    },
  });
}

// Re-export types for convenience
export type { UserSummary, UserDetails, SearchUsersParams };
