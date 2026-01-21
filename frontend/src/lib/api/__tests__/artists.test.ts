/**
 * Artists API Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

import { getArtists, getArtist } from '../artists';
import { apiClient } from '../client';

describe('Artists API (Wave 2)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockArtist = {
    name: 'Test Artist',
    trackCount: 25,
    albumCount: 3,
  };

  describe('getArtists', () => {
    it('should fetch artists with default params', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockArtist], total: 1, limit: 20, offset: 0 },
      });

      const result = await getArtists();

      expect(apiClient.get).toHaveBeenCalledWith('/artists', { params: undefined });
      expect(result.items).toHaveLength(1);
    });

    it('should fetch artists with pagination', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [], total: 100, limit: 20, offset: 20 },
      });

      await getArtists({ page: 2, limit: 20 });

      expect(apiClient.get).toHaveBeenCalledWith('/artists', {
        params: { page: 2, limit: 20 },
      });
    });

    it('should fetch artists with sorting', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [], total: 0, limit: 20, offset: 0 },
      });

      await getArtists({ sortBy: 'trackCount', sortOrder: 'desc' });

      expect(apiClient.get).toHaveBeenCalledWith('/artists', {
        params: { sortBy: 'trackCount', sortOrder: 'desc' },
      });
    });

    it('should fetch artists with search', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockArtist], total: 1, limit: 20, offset: 0 },
      });

      await getArtists({ search: 'test' });

      expect(apiClient.get).toHaveBeenCalledWith('/artists', {
        params: { search: 'test' },
      });
    });
  });

  describe('getArtist', () => {
    it('should fetch a single artist with details', async () => {
      const artistWithDetails = {
        ...mockArtist,
        albums: [
          { id: 'album-1', title: 'Album 1', year: 2024, trackCount: 10 },
        ],
        recentTracks: [
          { id: 'track-1', title: 'Track 1', duration: 180 },
        ],
      };
      vi.mocked(apiClient.get).mockResolvedValue({ data: artistWithDetails });

      const result = await getArtist('Test Artist');

      expect(apiClient.get).toHaveBeenCalledWith('/artists/Test Artist');
      expect(result.albums).toHaveLength(1);
      expect(result.recentTracks).toHaveLength(1);
    });
  });
});
