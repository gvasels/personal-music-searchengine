import { apiClient } from './client';
import type { WatchResponse, WatchStatusResponse, WatchedArtistsResponse } from '../../types';

export interface WatchedArtistsParams {
  cursor?: string;
  limit?: number;
}

export async function watchArtist(artistName: string): Promise<WatchResponse> {
  const response = await apiClient.post(`/v1/artists/${encodeURIComponent(artistName)}/watch`);
  return response.data;
}

export async function unwatchArtist(artistName: string): Promise<void> {
  await apiClient.delete(`/v1/artists/${encodeURIComponent(artistName)}/watch`);
}

export async function getWatchStatus(artistName: string): Promise<WatchStatusResponse> {
  const response = await apiClient.get(`/v1/artists/${encodeURIComponent(artistName)}/watch`);
  return response.data;
}

export async function getWatchedArtists(params?: WatchedArtistsParams): Promise<WatchedArtistsResponse> {
  const response = await apiClient.get('/v1/users/me/watched-artists', { params });
  return response.data;
}
