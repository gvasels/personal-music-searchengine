import { apiClient } from './client';

export interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre?: string;
  year?: number;
  duration: number;
}

export interface HelloSearchResponse {
  items: HelloTrack[];
  total: number;
}

export interface HelloHealthResponse {
  status: string;
  service: string;
}

export async function searchHelloTracks(query: string): Promise<HelloSearchResponse> {
  const { data } = await apiClient.get('/v1/hello/search', {
    params: { q: query },
  });
  return data;
}

export async function getFeaturedHelloTracks(limit?: number): Promise<HelloSearchResponse> {
  const { data } = await apiClient.get(
    '/v1/hello/featured',
    limit ? { params: { limit } } : {}
  );
  return data;
}

export async function getHelloHealth(): Promise<HelloHealthResponse> {
  const { data } = await apiClient.get('/v1/hello/health');
  return data;
}
