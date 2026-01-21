/**
 * Tags API - Wave 5
 */
import { apiClient } from './client';
import type { Tag, Track, PaginatedResponse } from '../../types';

export interface GetTagsParams {
  limit?: number;
  offset?: number;
}

export interface GetTracksByTagParams {
  limit?: number;
  offset?: number;
}

export async function getTags(params?: GetTagsParams): Promise<PaginatedResponse<Tag>> {
  const response = await apiClient.get<PaginatedResponse<Tag>>('/tags', { params });
  return response.data;
}

export async function getTracksByTag(
  tagName: string,
  params?: GetTracksByTagParams
): Promise<PaginatedResponse<Track>> {
  const response = await apiClient.get<PaginatedResponse<Track>>(`/tags/${tagName}/tracks`, {
    params,
  });
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
