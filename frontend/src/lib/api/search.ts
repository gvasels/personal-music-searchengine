/**
 * Search API Module - Wave 4
 */
import { apiClient } from './client';
import type { Track, Playlist } from '../../types';

export interface SearchParams {
  query: string;
  artist?: string;
  album?: string;
  page?: number;
  limit?: number;
}

export interface AutocompleteSuggestion {
  type: 'track' | 'artist' | 'album' | 'playlist';
  value: string;
  trackId?: string;
  albumId?: string;
  playlistId?: string;
}

export interface AutocompleteResponse {
  suggestions: AutocompleteSuggestion[];
}

export interface SearchResponse {
  query: string;
  totalResults: number;
  tracks: Track[];
  playlists?: Playlist[];
  limit: number;
  nextCursor?: string;
  hasMore: boolean;
}

export async function searchTracks(params: SearchParams): Promise<SearchResponse> {
  const response = await apiClient.get<SearchResponse>('/search', { params });
  return response.data;
}

export async function searchAutocomplete(query: string): Promise<AutocompleteResponse> {
  const response = await apiClient.get<AutocompleteResponse>('/search/autocomplete', { params: { q: query } });
  return response.data;
}
