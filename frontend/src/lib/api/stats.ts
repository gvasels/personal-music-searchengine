/**
 * Library Stats API
 * Handles library statistics API calls with role-based scope
 */
import { apiClient } from './client';
import type { LibraryStats, StatsScope } from '@/types';

/**
 * Get library statistics
 * @param scope - 'own' (user's tracks + public), 'public' (public only), 'all' (admin: all tracks)
 * @returns Library stats including track, album, artist counts and total duration
 */
export async function getLibraryStats(scope: StatsScope = 'own'): Promise<LibraryStats> {
  const response = await apiClient.get<LibraryStats>('/stats', {
    params: { scope },
  });
  return response.data;
}
