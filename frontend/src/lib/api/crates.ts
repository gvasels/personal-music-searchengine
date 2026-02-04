/**
 * Crates API
 * Handles DJ crate-related API calls
 */
import { apiClient } from './client';
import type { Crate, CrateWithTracks, Track } from '@/types';

// List all crates
export async function getCrates(): Promise<Crate[]> {
  const response = await apiClient.get<Crate[]>('/crates');
  return response.data;
}

// Get a single crate with its tracks
export async function getCrate(id: string): Promise<CrateWithTracks> {
  const response = await apiClient.get<CrateWithTracks>(`/crates/${id}`);
  return response.data;
}

// Create a new crate
export async function createCrate(data: {
  name: string;
  description?: string;
  color?: string;
}): Promise<Crate> {
  const response = await apiClient.post<Crate>('/crates', data);
  return response.data;
}

// Update a crate
export async function updateCrate(
  id: string,
  data: {
    name?: string;
    description?: string;
    color?: string;
    sortOrder?: string;
  }
): Promise<Crate> {
  const response = await apiClient.patch<Crate>(`/crates/${id}`, data);
  return response.data;
}

// Delete a crate
export async function deleteCrate(id: string): Promise<void> {
  await apiClient.delete(`/crates/${id}`);
}

// Add tracks to a crate
export async function addTracksToCrate(
  crateId: string,
  trackIds: string[],
  position?: number
): Promise<void> {
  await apiClient.post(`/crates/${crateId}/tracks`, { trackIds, position });
}

// Remove tracks from a crate
export async function removeTracksFromCrate(
  crateId: string,
  trackIds: string[]
): Promise<void> {
  await apiClient.delete(`/crates/${crateId}/tracks`, { data: { trackIds } });
}

// Reorder tracks in a crate
export async function reorderCrateTracks(
  crateId: string,
  trackIds: string[]
): Promise<void> {
  await apiClient.put(`/crates/${crateId}/tracks/order`, { trackIds });
}

// Get tracks in a crate
export async function getCrateTracks(crateId: string): Promise<Track[]> {
  const response = await apiClient.get<Track[]>(`/crates/${crateId}/tracks`);
  return response.data;
}
