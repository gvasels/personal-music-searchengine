/**
 * Hello Search API Tests - TDD Red Phase
 *
 * Tests for the hello search API client functions.
 * Source module (helloSearch.ts) does NOT exist yet - these tests MUST fail.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

import { searchHelloTracks, getFeaturedTracks, getHelloHealth } from '../helloSearch';
import { apiClient } from '../client';

describe('Hello Search API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('searchHelloTracks', () => {
    it('should call GET /v1/hello/search with query param', async () => {
      const mockResponse = {
        items: [
          { id: 'track-1', title: 'Test Song', artist: 'Test Artist', album: 'Test Album', genre: 'Rock', year: 2024, duration: 234 },
        ],
        total: 1,
      };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse });

      const result = await searchHelloTracks('test query');

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/search', {
        params: { q: 'test query' },
      });
      expect(result.items).toHaveLength(1);
      expect(result.items[0].title).toBe('Test Song');
    });
  });

  describe('getFeaturedTracks', () => {
    it('should call GET /v1/hello/featured', async () => {
      const mockResponse = {
        items: [
          { id: 'track-1', title: 'Featured Song', artist: 'Featured Artist', album: 'Featured Album', genre: 'Pop', year: 2024, duration: 180 },
        ],
        total: 1,
      };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse });

      const result = await getFeaturedTracks();

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/featured', {
        params: {},
      });
      expect(result.items).toHaveLength(1);
    });

    it('should pass limit parameter when provided', async () => {
      const mockResponse = { items: [], total: 0 };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse });

      await getFeaturedTracks(5);

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/featured', {
        params: { limit: 5 },
      });
    });
  });

  describe('getHelloHealth', () => {
    it('should call GET /v1/hello/health', async () => {
      const mockResponse = { status: 'ok', service: 'hello' };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse });

      const result = await getHelloHealth();

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/health');
      expect(result.status).toBe('ok');
    });
  });
});
