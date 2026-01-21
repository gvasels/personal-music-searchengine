/**
 * Search API Tests - Wave 4
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

import { searchTracks, searchAutocomplete } from '../search';
import { apiClient } from '../client';

describe('Search API (Wave 4)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockSearchResults = {
    items: [
      { id: 'track-1', title: 'Test Song', artist: 'Test Artist', album: 'Test Album' },
    ],
    total: 1,
    limit: 20,
    offset: 0,
  };

  describe('searchTracks', () => {
    it('should search tracks with query', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResults });

      const result = await searchTracks({ query: 'test' });

      expect(apiClient.get).toHaveBeenCalledWith('/search', { params: { query: 'test' } });
      expect(result.items).toHaveLength(1);
    });

    it('should search with filters', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResults });

      await searchTracks({ query: 'test', artist: 'Artist', album: 'Album' });

      expect(apiClient.get).toHaveBeenCalledWith('/search', {
        params: { query: 'test', artist: 'Artist', album: 'Album' },
      });
    });

    it('should search with pagination', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSearchResults });

      await searchTracks({ query: 'test', page: 2, limit: 10 });

      expect(apiClient.get).toHaveBeenCalledWith('/search', {
        params: { query: 'test', page: 2, limit: 10 },
      });
    });
  });

  describe('searchAutocomplete', () => {
    it('should return autocomplete suggestions', async () => {
      const mockSuggestions = {
        suggestions: [
          { type: 'track', value: 'Test Song', trackId: 'track-1' },
          { type: 'artist', value: 'Test Artist' },
          { type: 'album', value: 'Test Album', albumId: 'album-1' },
        ],
      };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockSuggestions });

      const result = await searchAutocomplete('test');

      expect(apiClient.get).toHaveBeenCalledWith('/search/autocomplete', { params: { q: 'test' } });
      expect(result.suggestions).toHaveLength(3);
    });
  });
});
