/**
 * Hello Search API Tests - TDD Red Phase
 *
 * These tests define the contract for the helloSearch API client.
 * They MUST FAIL because helloSearch.ts does not exist yet.
 *
 * Contract critical:
 * - searchHelloTracks uses param `q` (NOT `query`)
 * - Response type expects `items` field (NOT `tracks`)
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

// This import will fail - helloSearch.ts does not exist
import { searchHelloTracks, getFeaturedHelloTracks, getHelloHealth } from '../helloSearch';
import { apiClient } from '../client';

describe('Hello Search API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockHelloTrack = {
    id: 'seed-t1',
    title: 'Aurora Borealis',
    artist: 'Aurora Waves',
    album: 'Dreamscape',
    genre: 'jazz',
    year: 2023,
    duration: 245,
  };

  const mockSearchResponse = {
    items: [mockHelloTrack],
    total: 1,
  };

  describe('searchHelloTracks', () => {
    it('should search tracks with q param (NOT query)', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      const result = await searchHelloTracks('jazz');

      // CRITICAL: param is `q`, NOT `query`
      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/search', { params: { q: 'jazz' } });
      expect(result.items).toHaveLength(1);
      expect(result.total).toBe(1);
    });

    it('should return items array (NOT tracks)', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      const result = await searchHelloTracks('electronic');

      // CRITICAL: response uses `items`, NOT `tracks`
      expect(result).toHaveProperty('items');
      expect(result.items).toBeInstanceOf(Array);
      expect(result.items[0]).toEqual(mockHelloTrack);
    });

    it('should return empty items for no matches', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: { items: [], total: 0 } });

      const result = await searchHelloTracks('nonexistent');

      expect(result.items).toEqual([]);
      expect(result.total).toBe(0);
    });
  });

  describe('getFeaturedHelloTracks', () => {
    it('should fetch featured tracks without limit', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      const result = await getFeaturedHelloTracks();

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/featured', { params: undefined });
      expect(result.items).toHaveLength(1);
    });

    it('should fetch featured tracks with limit param', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      await getFeaturedHelloTracks(5);

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/featured', { params: { limit: 5 } });
    });

    it('should return items array (NOT tracks)', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      const result = await getFeaturedHelloTracks();

      expect(result).toHaveProperty('items');
      expect(result.items).toBeInstanceOf(Array);
    });
  });

  describe('getHelloHealth', () => {
    it('should fetch health status', async () => {
      const mockHealth = { status: 'ok', service: 'hello' };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockHealth });

      const result = await getHelloHealth();

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/health');
      expect(result.status).toBe('ok');
      expect(result.service).toBe('hello');
    });
  });
});
