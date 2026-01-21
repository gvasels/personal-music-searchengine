/**
 * Albums API Module - Wave 2
 */
import { apiClient } from './client';
import type { Album, Track, PaginatedResponse } from '../../types';

export interface GetAlbumsParams {
  page?: number;
  limit?: number;
  sortBy?: 'title' | 'artist' | 'year' | 'createdAt';
  sortOrder?: 'asc' | 'desc';
  artist?: string;
}

export interface AlbumWithTracks extends Album {
  tracks: Track[];
}

export async function getAlbums(params?: GetAlbumsParams): Promise<PaginatedResponse<Album>> {
  const response = await apiClient.get<PaginatedResponse<Album>>('/albums', { params });
  return response.data;
}

export async function getAlbum(id: string): Promise<AlbumWithTracks> {
  const response = await apiClient.get<AlbumWithTracks>(`/albums/${id}`);
  return response.data;
}
