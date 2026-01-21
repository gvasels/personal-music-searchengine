/**
 * Tracks API Module - Wave 2
 */
import { apiClient } from './client';
import type { Track, PaginatedResponse } from '../../types';

export interface GetTracksParams {
  page?: number;
  limit?: number;
  sortBy?: 'title' | 'artist' | 'album' | 'createdAt' | 'duration';
  sortOrder?: 'asc' | 'desc';
  search?: string;
  artist?: string;
  album?: string;
  tags?: string[];
}

export interface UpdateTrackData {
  title?: string;
  artist?: string;
  album?: string;
  tags?: string[];
}

export async function getTracks(params?: GetTracksParams): Promise<PaginatedResponse<Track>> {
  const response = await apiClient.get<PaginatedResponse<Track>>('/tracks', { params });
  return response.data;
}

export async function getTrack(id: string): Promise<Track> {
  const response = await apiClient.get<Track>(`/tracks/${id}`);
  return response.data;
}

export async function updateTrack(id: string, data: UpdateTrackData): Promise<Track> {
  const response = await apiClient.patch<Track>(`/tracks/${id}`, data);
  return response.data;
}

export async function deleteTrack(id: string): Promise<void> {
  await apiClient.delete(`/tracks/${id}`);
}
