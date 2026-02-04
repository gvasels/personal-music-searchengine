/**
 * Admin API - Admin Panel & Track Visibility Feature
 * Provides functions for admin user management operations
 */
import { apiClient } from './client';
import type { UserRole } from '../../types';

// User summary returned from search results
export interface UserSummary {
  id: string;
  email: string;
  displayName: string;
  role: UserRole;
  disabled: boolean;
  createdAt: string;
}

// Full user details for admin management
export interface UserDetails extends UserSummary {
  lastLoginAt?: string;
  trackCount: number;
  playlistCount: number;
  albumCount: number;
  storageUsed: number;
  followerCount: number;
  followingCount: number;
}

// Search users response
export interface SearchUsersResponse {
  items: UserSummary[];
  nextCursor?: string;
}

// Search users parameters
export interface SearchUsersParams {
  search: string;
  limit?: number;
  cursor?: string;
}

// Update user role request
export interface UpdateUserRoleRequest {
  role: UserRole;
}

// Update user status request
export interface UpdateUserStatusRequest {
  disabled: boolean;
}

/**
 * Search users by email or display name
 */
export async function searchUsers(params: SearchUsersParams): Promise<SearchUsersResponse> {
  const response = await apiClient.get<SearchUsersResponse>('/admin/users', { params });
  return response.data;
}

/**
 * Get detailed user information by ID
 */
export async function getUserDetails(userId: string): Promise<UserDetails> {
  const response = await apiClient.get<UserDetails>(`/admin/users/${userId}`);
  return response.data;
}

/**
 * Update a user's role
 */
export async function updateUserRole(userId: string, role: UserRole): Promise<UserDetails> {
  const response = await apiClient.put<UserDetails>(`/admin/users/${userId}/role`, {
    role,
  } as UpdateUserRoleRequest);
  return response.data;
}

/**
 * Update a user's status (enable/disable)
 */
export async function updateUserStatus(userId: string, disabled: boolean): Promise<UserDetails> {
  const response = await apiClient.put<UserDetails>(`/admin/users/${userId}/status`, {
    disabled,
  } as UpdateUserStatusRequest);
  return response.data;
}
