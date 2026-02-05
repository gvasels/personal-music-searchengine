import { apiClient } from './client';

export interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;
}

export interface HelloSearchResponse {
  items: HelloTrack[];
  total: number;
}

export async function searchHelloTracks(query: string): Promise<HelloSearchResponse> {
  const response = await apiClient.get<HelloSearchResponse>('/v1/hello/search', {
    params: { q: query },
  });
  return response.data;
}

export async function getFeaturedTracks(limit?: number): Promise<HelloSearchResponse> {
  const params: Record<string, unknown> = {};
  if (limit !== undefined) {
    params.limit = limit;
  }
  const response = await apiClient.get<HelloSearchResponse>('/v1/hello/featured', { params });
  return response.data;
}

export async function getHelloHealth(): Promise<{ status: string; service: string }> {
  const response = await apiClient.get<{ status: string; service: string }>('/v1/hello/health');
  return response.data;
}
