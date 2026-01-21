/**
 * Tags API Tests - Wave 5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { apiClient } from '../client';
import { getTags, getTracksByTag, addTagToTrack, removeTagFromTrack } from '../tags';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('Tags API (Wave 5)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getTags', () => {
    it('should fetch all tags', async () => {
      const mockResponse = {
        data: { items: [{ name: 'rock', trackCount: 10 }], total: 1, limit: 50, offset: 0 },
      };
      vi.mocked(apiClient.get).mockResolvedValue(mockResponse);

      const result = await getTags();

      expect(apiClient.get).toHaveBeenCalledWith('/tags', { params: undefined });
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('getTracksByTag', () => {
    it('should fetch tracks for a specific tag', async () => {
      const mockResponse = {
        data: { items: [], total: 0, limit: 20, offset: 0 },
      };
      vi.mocked(apiClient.get).mockResolvedValue(mockResponse);

      const result = await getTracksByTag('rock', { limit: 20, offset: 0 });

      expect(apiClient.get).toHaveBeenCalledWith('/tags/rock/tracks', {
        params: { limit: 20, offset: 0 },
      });
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('addTagToTrack', () => {
    it('should add a tag to a track', async () => {
      const mockTrack = { id: 'track-1', tags: ['rock'] };
      vi.mocked(apiClient.post).mockResolvedValue({ data: mockTrack });

      const result = await addTagToTrack('track-1', 'rock');

      expect(apiClient.post).toHaveBeenCalledWith('/tracks/track-1/tags', { tagName: 'rock' });
      expect(result).toEqual(mockTrack);
    });
  });

  describe('removeTagFromTrack', () => {
    it('should remove a tag from a track', async () => {
      const mockTrack = { id: 'track-1', tags: [] };
      vi.mocked(apiClient.delete).mockResolvedValue({ data: mockTrack });

      const result = await removeTagFromTrack('track-1', 'rock');

      expect(apiClient.delete).toHaveBeenCalledWith('/tracks/track-1/tags/rock');
      expect(result).toEqual(mockTrack);
    });
  });
});
