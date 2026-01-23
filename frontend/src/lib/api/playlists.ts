/**
 * Playlists API - Wave 5
 */
import { apiClient } from './client';
import type { Playlist, PlaylistWithTracks, PaginatedResponse } from '../../types';

export interface GetPlaylistsParams {
  limit?: number;
  offset?: number;
}

export interface CreatePlaylistData {
  name: string;
  description?: string;
}

export interface UpdatePlaylistData {
  name?: string;
  description?: string;
}

export async function getPlaylists(
  params?: GetPlaylistsParams
): Promise<PaginatedResponse<Playlist>> {
  const response = await apiClient.get<PaginatedResponse<Playlist>>('/playlists', { params });
  return response.data;
}

export async function getPlaylist(id: string): Promise<PlaylistWithTracks> {
  const response = await apiClient.get<PlaylistWithTracks>(`/playlists/${id}`);
  return response.data;
}

export async function createPlaylist(data: CreatePlaylistData): Promise<Playlist> {
  const response = await apiClient.post<Playlist>('/playlists', data);
  return response.data;
}

export async function updatePlaylist(id: string, data: UpdatePlaylistData): Promise<Playlist> {
  const response = await apiClient.put<Playlist>(`/playlists/${id}`, data);
  return response.data;
}

export async function deletePlaylist(id: string): Promise<void> {
  await apiClient.delete(`/playlists/${id}`);
}

export async function addTrackToPlaylist(playlistId: string, trackId: string): Promise<Playlist> {
  const response = await apiClient.post<Playlist>(`/playlists/${playlistId}/tracks`, { trackIds: [trackId] });
  return response.data;
}

export async function removeTrackFromPlaylist(
  playlistId: string,
  trackId: string
): Promise<Playlist> {
  const response = await apiClient.delete<Playlist>(`/playlists/${playlistId}/tracks`, {
    data: { trackIds: [trackId] },
  });
  return response.data;
}
