import { apiClient } from './client';

export interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;
  durationStr: string;
  coverArtUrl: string;
}

export interface HelloSearchResponse {
  tracks: HelloTrack[];
  total: number;
  query: string;
}

export async function searchHello(query: string): Promise<HelloSearchResponse> {
  const response = await apiClient.get<HelloSearchResponse>('/v1/hello/search', {
    params: { q: query },
  });
  return response.data;
}

export async function getFeaturedTracks(): Promise<HelloSearchResponse> {
  const response = await apiClient.get<HelloSearchResponse>('/v1/hello/featured');
  return response.data;
}
