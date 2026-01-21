/**
 * Search API Module - Wave 4
 */
import { apiClient } from './client';
import type { Track, PaginatedResponse } from '../../types';

export interface SearchParams {
  query: string;
  artist?: string;
  album?: string;
  page?: number;
  limit?: number;
}

export interface AutocompleteSuggestion {
  type: 'track' | 'artist' | 'album';
  value: string;
  trackId?: string;
  albumId?: string;
}

export interface AutocompleteResponse {
  suggestions: AutocompleteSuggestion[];
}

export async function searchTracks(params: SearchParams): Promise<PaginatedResponse<Track>> {
  const response = await apiClient.get<PaginatedResponse<Track>>('/search', { params });
  return response.data;
}

export async function searchAutocomplete(query: string): Promise<AutocompleteResponse> {
  const response = await apiClient.get<AutocompleteResponse>('/search/autocomplete', { params: { q: query } });
  return response.data;
}
