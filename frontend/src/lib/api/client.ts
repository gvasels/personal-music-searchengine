import axios from 'axios';
import type { Track, Album, Artist, Playlist, PaginatedResponse } from '@/types';

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: { 'Content-Type': 'application/json' },
});

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export async function getTracks(params?: { limit?: number; offset?: number }): Promise<PaginatedResponse<Track>> {
  const response = await apiClient.get<PaginatedResponse<Track>>('/tracks', { params });
  return response.data;
}

export async function getTrack(id: string): Promise<Track> {
  const response = await apiClient.get<Track>(`/tracks/${id}`);
  return response.data;
}

export async function getAlbums(params?: { limit?: number }): Promise<PaginatedResponse<Album>> {
  const response = await apiClient.get<PaginatedResponse<Album>>('/albums', { params });
  return response.data;
}

export async function getArtists(params?: { limit?: number }): Promise<PaginatedResponse<Artist>> {
  const response = await apiClient.get<PaginatedResponse<Artist>>('/artists', { params });
  return response.data;
}

export async function getPlaylists(params?: { limit?: number }): Promise<PaginatedResponse<Playlist>> {
  const response = await apiClient.get<PaginatedResponse<Playlist>>('/playlists', { params });
  return response.data;
}

export async function getPlaylist(id: string): Promise<Playlist> {
  const response = await apiClient.get<Playlist>(`/playlists/${id}`);
  return response.data;
}

export async function createPlaylist(data: { name: string; description?: string }): Promise<Playlist> {
  const response = await apiClient.post<Playlist>('/playlists', data);
  return response.data;
}

export async function deletePlaylist(id: string): Promise<void> {
  await apiClient.delete(`/playlists/${id}`);
}

export async function addTrackToPlaylist(playlistId: string, trackId: string): Promise<Playlist> {
  const response = await apiClient.post<Playlist>(`/playlists/${playlistId}/tracks`, { trackId });
  return response.data;
}

export async function removeTrackFromPlaylist(playlistId: string, trackId: string): Promise<Playlist> {
  const response = await apiClient.delete<Playlist>(`/playlists/${playlistId}/tracks/${trackId}`);
  return response.data;
}

export async function addTagToTrack(trackId: string, tagName: string): Promise<Track> {
  const response = await apiClient.post<Track>(`/tracks/${trackId}/tags`, { tagName });
  return response.data;
}

export async function removeTagFromTrack(trackId: string, tagName: string): Promise<Track> {
  const response = await apiClient.delete<Track>(`/tracks/${trackId}/tags/${tagName}`);
  return response.data;
}

export async function searchTracks(query: string): Promise<Track[]> {
  const response = await apiClient.get<Track[]>('/search', { params: { q: query } });
  return response.data;
}

export async function getPresignedUploadUrl(data: { fileName: string; contentType: string; fileSize: number }): Promise<{ uploadId: string; uploadUrl: string }> {
  const response = await apiClient.post('/upload/presigned', data);
  return response.data;
}

export async function getStreamUrl(trackId: string): Promise<{ streamUrl: string }> {
  const response = await apiClient.get(`/tracks/${trackId}/stream`);
  return response.data;
}
