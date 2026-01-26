/**
 * Tracks API Module - Wave 2
 */
import { apiClient } from './client';
import type { Track, PaginatedResponse, TrackVisibility } from '../../types';

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

export async function updateTrackVisibility(
  id: string,
  visibility: TrackVisibility
): Promise<{ trackId: string; visibility: TrackVisibility }> {
  const response = await apiClient.put<{ trackId: string; visibility: TrackVisibility }>(
    `/tracks/${id}/visibility`,
    { visibility }
  );
  return response.data;
}
