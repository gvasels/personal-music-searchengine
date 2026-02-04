/**
 * Playlists API - Wave 5
 */
import { apiClient } from './client';
import type { Playlist, PlaylistWithTracks, PaginatedResponse, PlaylistVisibility } from '../../types';

export interface GetPlaylistsParams {
  limit?: number;
  offset?: number;
}

export interface CreatePlaylistData {
  name: string;
  description?: string;
  visibility?: PlaylistVisibility;
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

export interface ReorderTracksData {
  trackIds: string[];
}

export async function reorderPlaylistTracks(
  playlistId: string,
  data: ReorderTracksData
): Promise<Playlist> {
  const response = await apiClient.put<Playlist>(`/playlists/${playlistId}/reorder`, data);
  return response.data;
}

export interface UpdateVisibilityResponse {
  playlistId: string;
  visibility: PlaylistVisibility;
}

export async function updatePlaylistVisibility(
  playlistId: string,
  visibility: PlaylistVisibility
): Promise<UpdateVisibilityResponse> {
  const response = await apiClient.put<UpdateVisibilityResponse>(
    `/playlists/${playlistId}/visibility`,
    { visibility }
  );
  return response.data;
}

export interface GetPublicPlaylistsParams {
  limit?: number;
  cursor?: string;
}

export interface PublicPlaylistsResponse {
  items: Playlist[];
  nextCursor?: string;
}

export async function getPublicPlaylists(
  params?: GetPublicPlaylistsParams
): Promise<PublicPlaylistsResponse> {
  const response = await apiClient.get<PublicPlaylistsResponse>('/playlists/public', { params });
  return response.data;
}
