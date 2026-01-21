/**
 * Artists API Module - Wave 2
 */
import { apiClient } from './client';
import type { Artist, Album, Track, PaginatedResponse } from '../../types';

export interface GetArtistsParams {
  page?: number;
  limit?: number;
  sortBy?: 'name' | 'trackCount' | 'albumCount';
  sortOrder?: 'asc' | 'desc';
  search?: string;
}

export interface ArtistWithDetails extends Artist {
  albums: Album[];
  recentTracks: Track[];
}

export async function getArtists(params?: GetArtistsParams): Promise<PaginatedResponse<Artist>> {
  const response = await apiClient.get<PaginatedResponse<Artist>>('/artists', { params });
  return response.data;
}

export async function getArtist(name: string): Promise<ArtistWithDetails> {
  const response = await apiClient.get<ArtistWithDetails>(`/artists/${name}`);
  return response.data;
}
