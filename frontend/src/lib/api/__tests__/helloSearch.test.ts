/**
 * Hello Search API Tests - TDD Red Phase
 *
 * Tests for the hello-world search API client functions.
 * CRITICAL CONTRACT:
 * - Search uses param key `q` (NOT `query`)
 * - Response uses `items` field (NOT `tracks`)
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

import { searchHelloTracks, getFeaturedHelloTracks, getHelloHealth } from '../helloSearch';
import { apiClient } from '../client';

describe('Hello Search API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockHelloTrack = {
    id: 'hello-track-1',
    title: 'Hello World',
    artist: 'Test Artist',
    album: 'Test Album',
    duration: 180,
  };

  const mockSearchResponse = {
    items: [mockHelloTrack],
    total: 1,
  };

  describe('searchHelloTracks', () => {
    it('should call apiClient.get with /v1/hello/search and q param', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      await searchHelloTracks('jazz');

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/search', {
        params: { q: 'jazz' },
      });
    });

    it('should return response with items array (not tracks)', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      const result = await searchHelloTracks('jazz');

      expect(result).toHaveProperty('items');
      expect(result).not.toHaveProperty('tracks');
      expect(result.items).toHaveLength(1);
      expect(result.total).toBe(1);
    });
  });

  describe('getFeaturedHelloTracks', () => {
    it('should call apiClient.get with /v1/hello/featured', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      await getFeaturedHelloTracks();

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/featured', {});
    });

    it('should call with limit param when provided', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResponse });

      await getFeaturedHelloTracks(5);

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/featured', {
        params: { limit: 5 },
      });
    });
  });

  describe('getHelloHealth', () => {
    it('should call apiClient.get with /v1/hello/health', async () => {
      const mockHealth = { status: 'ok', service: 'hello' };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockHealth });

      await getHelloHealth();

      expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/health');
    });

    it('should return health status with service name', async () => {
      const mockHealth = { status: 'ok', service: 'hello' };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockHealth });

      const result = await getHelloHealth();

      expect(result.status).toBe('ok');
      expect(result.service).toBe('hello');
    });
  });
});
