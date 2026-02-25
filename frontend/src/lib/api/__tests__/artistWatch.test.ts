/**
 * Artist Watch API Tests - TDD Red Phase
 *
 * These tests define the contract for the artistWatch API client.
 * They MUST FAIL because artistWatch.ts does not exist yet.
 *
 * Contract critical:
 * - watchArtist sends POST to /v1/artists/{name}/watch
 * - unwatchArtist sends DELETE to /v1/artists/{name}/watch
 * - getWatchStatus sends GET to /v1/artists/{name}/watch
 * - getWatchedArtists sends GET to /v1/users/me/watched-artists
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn(),
  },
}));

// These imports will fail - artistWatch.ts does not exist yet
import { watchArtist, unwatchArtist, getWatchStatus, getWatchedArtists } from '../artistWatch';
import { apiClient } from '../client';

describe('Artist Watch API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockWatchResponse = {
    artistName: 'Daft Punk',
    watching: true,
    watchedAt: '2026-02-25T10:00:00Z',
  };

  const mockWatchStatusResponse = {
    watching: true,
    artistName: 'Daft Punk',
  };

  const mockWatchedArtistsResponse = {
    items: [
      { artistName: 'Daft Punk', watchedAt: '2026-02-25T10:00:00Z' },
      { artistName: 'Aphex Twin', watchedAt: '2026-02-24T08:00:00Z' },
    ],
    nextCursor: 'cursor-abc',
    hasMore: true,
  };

  describe('watchArtist', () => {
    it('should send POST to /v1/artists/{name}/watch', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ data: mockWatchResponse });

      await watchArtist('Daft Punk');

      expect(apiClient.post).toHaveBeenCalledWith('/v1/artists/Daft Punk/watch');
    });

    it('should return watch confirmation with artistName, watching, and watchedAt', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ data: mockWatchResponse });

      const result = await watchArtist('Daft Punk');

      expect(result.artistName).toBe('Daft Punk');
      expect(result.watching).toBe(true);
      expect(result.watchedAt).toBe('2026-02-25T10:00:00Z');
    });

    it('should throw error on API failure', async () => {
      vi.mocked(apiClient.post).mockRejectedValue(new Error('Unauthorized'));

      await expect(watchArtist('Daft Punk')).rejects.toThrow('Unauthorized');
    });

    it('should handle artist names with special characters', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({
        data: { ...mockWatchResponse, artistName: 'M83' },
      });

      await watchArtist('M83');

      expect(apiClient.post).toHaveBeenCalledWith('/v1/artists/M83/watch');
    });
  });

  describe('unwatchArtist', () => {
    it('should send DELETE to /v1/artists/{name}/watch', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({ status: 204 });

      await unwatchArtist('Daft Punk');

      expect(apiClient.delete).toHaveBeenCalledWith('/v1/artists/Daft Punk/watch');
    });

    it('should throw error on API failure', async () => {
      vi.mocked(apiClient.delete).mockRejectedValue(new Error('Not Found'));

      await expect(unwatchArtist('Unknown Artist')).rejects.toThrow('Not Found');
    });
  });

  describe('getWatchStatus', () => {
    it('should send GET to /v1/artists/{name}/watch', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockWatchStatusResponse });

      await getWatchStatus('Daft Punk');

      expect(apiClient.get).toHaveBeenCalledWith('/v1/artists/Daft Punk/watch');
    });

    it('should return watching boolean and artistName', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockWatchStatusResponse });

      const result = await getWatchStatus('Daft Punk');

      expect(result.watching).toBe(true);
      expect(result.artistName).toBe('Daft Punk');
    });

    it('should return watching=false when not watching', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { watching: false, artistName: 'Aphex Twin' },
      });

      const result = await getWatchStatus('Aphex Twin');

      expect(result.watching).toBe(false);
    });

    it('should throw error on API failure', async () => {
      vi.mocked(apiClient.get).mockRejectedValue(new Error('Network error'));

      await expect(getWatchStatus('Daft Punk')).rejects.toThrow('Network error');
    });
  });

  describe('getWatchedArtists', () => {
    it('should send GET to /v1/users/me/watched-artists', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockWatchedArtistsResponse });

      await getWatchedArtists();

      expect(apiClient.get).toHaveBeenCalledWith('/v1/users/me/watched-artists', { params: undefined });
    });

    it('should return items array with watched artists', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockWatchedArtistsResponse });

      const result = await getWatchedArtists();

      expect(result.items).toHaveLength(2);
      expect(result.items[0].artistName).toBe('Daft Punk');
      expect(result.items[0].watchedAt).toBe('2026-02-25T10:00:00Z');
    });

    it('should pass pagination params', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockWatchedArtistsResponse });

      await getWatchedArtists({ cursor: 'cursor-abc', limit: 10 });

      expect(apiClient.get).toHaveBeenCalledWith('/v1/users/me/watched-artists', {
        params: { cursor: 'cursor-abc', limit: 10 },
      });
    });

    it('should return hasMore and nextCursor for pagination', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockWatchedArtistsResponse });

      const result = await getWatchedArtists();

      expect(result.hasMore).toBe(true);
      expect(result.nextCursor).toBe('cursor-abc');
    });

    it('should return empty items when no artists are watched', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [], hasMore: false },
      });

      const result = await getWatchedArtists();

      expect(result.items).toEqual([]);
      expect(result.hasMore).toBe(false);
    });

    it('should throw error on API failure', async () => {
      vi.mocked(apiClient.get).mockRejectedValue(new Error('Forbidden'));

      await expect(getWatchedArtists()).rejects.toThrow('Forbidden');
    });
  });
});
