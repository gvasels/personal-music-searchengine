/**
 * Artists API Module - Wave 2
 *
 * This module handles the LEGACY /artists endpoint which returns aggregated
 * artist data from tracks/albums. For the new Artist entity endpoints,
 * see the functions in client.ts (getArtistEntities, getArtistEntity, etc.)
 */
import { apiClient } from './client';
import type { ArtistSummary, Album, Track, PaginatedResponse } from '../../types';

export interface GetArtistsParams {
  page?: number;
  limit?: number;
  sortBy?: 'name' | 'trackCount' | 'albumCount';
  sortOrder?: 'asc' | 'desc';
  search?: string;
}

export interface ArtistWithDetails extends ArtistSummary {
  albums: Album[];
  recentTracks: Track[];
}

export async function getArtists(params?: GetArtistsParams): Promise<PaginatedResponse<ArtistSummary>> {
  const response = await apiClient.get<PaginatedResponse<ArtistSummary>>('/artists', { params });
  return response.data;
}

export async function getArtist(name: string): Promise<ArtistWithDetails> {
  const response = await apiClient.get<ArtistWithDetails>(`/artists/${name}`);
  return response.data;
}
