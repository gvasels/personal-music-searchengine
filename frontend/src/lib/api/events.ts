import { apiClient } from './client';
import type { ArtistEventsResponse, ArtistSearchResponse } from '../../types';

export async function getArtistEvents(artistName: string): Promise<ArtistEventsResponse> {
  const response = await apiClient.get(`/v1/artists/${encodeURIComponent(artistName)}/events`);
  return response.data;
}

export async function searchArtistEvents(query: string, limit: number = 10): Promise<ArtistSearchResponse> {
  const response = await apiClient.get('/v1/events/search', {
    params: { q: query, limit },
  });
  return response.data;
}
