/**
 * Hello Search API Client
 *
 * API functions for the Hello World search feature.
 * Uses the seed track data from LocalStack DynamoDB.
 *
 * CONTRACT CRITICAL:
 * - searchHelloTracks uses param `q` (NOT `query`)
 * - Response type uses `items` field (NOT `tracks`)
 */
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

/**
 * Search for tracks by query string.
 * @param query - Search term (matched against title, artist, album, genre)
 * @returns Search results with items array and total count
 */
export async function searchHelloTracks(query: string): Promise<HelloSearchResponse> {
  const response = await apiClient.get('/v1/hello/search', {
    params: { q: query },
  });
  return response.data;
}

/**
 * Get featured tracks.
 * @param limit - Optional limit on number of tracks returned
 * @returns Featured tracks with items array and total count
 */
export async function getFeaturedHelloTracks(limit?: number): Promise<HelloSearchResponse> {
  const response = await apiClient.get('/v1/hello/featured', {
    params: limit ? { limit } : undefined,
  });
  return response.data;
}

/**
 * Check the hello service health status.
 * @returns Health status object
 */
export async function getHelloHealth(): Promise<{ status: string; service: string }> {
  const response = await apiClient.get('/v1/hello/health');
  return response.data;
}
