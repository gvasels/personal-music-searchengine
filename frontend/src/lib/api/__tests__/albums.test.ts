/**
 * Albums API Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

import { getAlbums, getAlbum } from '../albums';
import { apiClient } from '../client';

describe('Albums API (Wave 2)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockAlbum = {
    id: 'album-1',
    title: 'Test Album',
    artist: 'Test Artist',
    year: 2024,
    coverArtUrl: 'https://example.com/cover.jpg',
    trackCount: 10,
    totalDuration: 3600,
    createdAt: '2024-01-01T00:00:00Z',
  };

  describe('getAlbums', () => {
    it('should fetch albums with default params', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockAlbum], total: 1, limit: 20, offset: 0 },
      });

      const result = await getAlbums();

      expect(apiClient.get).toHaveBeenCalledWith('/albums', { params: undefined });
      expect(result.items).toHaveLength(1);
    });

    it('should fetch albums with pagination', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [], total: 100, limit: 24, offset: 24 },
      });

      await getAlbums({ page: 2, limit: 24 });

      expect(apiClient.get).toHaveBeenCalledWith('/albums', {
        params: { page: 2, limit: 24 },
      });
    });

    it('should fetch albums with sorting', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [], total: 0, limit: 20, offset: 0 },
      });

      await getAlbums({ sortBy: 'year', sortOrder: 'desc' });

      expect(apiClient.get).toHaveBeenCalledWith('/albums', {
        params: { sortBy: 'year', sortOrder: 'desc' },
      });
    });

    it('should fetch albums with artist filter', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockAlbum], total: 1, limit: 20, offset: 0 },
      });

      await getAlbums({ artist: 'Test Artist' });

      expect(apiClient.get).toHaveBeenCalledWith('/albums', {
        params: { artist: 'Test Artist' },
      });
    });
  });

  describe('getAlbum', () => {
    it('should fetch a single album with tracks', async () => {
      const albumWithTracks = {
        ...mockAlbum,
        tracks: [
          { id: 'track-1', title: 'Track 1', duration: 180, trackNumber: 1 },
          { id: 'track-2', title: 'Track 2', duration: 200, trackNumber: 2 },
        ],
      };
      vi.mocked(apiClient.get).mockResolvedValue({ data: albumWithTracks });

      const result = await getAlbum('album-1');

      expect(apiClient.get).toHaveBeenCalledWith('/albums/album-1');
      expect(result.tracks).toHaveLength(2);
    });
  });
});
