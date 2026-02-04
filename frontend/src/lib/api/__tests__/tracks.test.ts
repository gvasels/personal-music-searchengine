/**
 * Tracks API Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}));

import { getTracks, getTrack, updateTrack, deleteTrack } from '../tracks';
import { apiClient } from '../client';

describe('Tracks API (Wave 2)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockTrack = {
    id: 'track-1',
    title: 'Test Track',
    artist: 'Test Artist',
    album: 'Test Album',
    duration: 180,
    format: 'mp3',
    fileSize: 5000000,
    tags: ['rock'],
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  };

  describe('getTracks', () => {
    it('should fetch tracks with default params', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockTrack], total: 1, limit: 20, offset: 0 },
      });

      const result = await getTracks();

      expect(apiClient.get).toHaveBeenCalledWith('/tracks', { params: undefined });
      expect(result.items).toHaveLength(1);
    });

    it('should fetch tracks with pagination', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [], total: 100, limit: 20, offset: 20 },
      });

      await getTracks({ page: 2, limit: 20 });

      expect(apiClient.get).toHaveBeenCalledWith('/tracks', {
        params: { page: 2, limit: 20 },
      });
    });

    it('should fetch tracks with sorting', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [], total: 0, limit: 20, offset: 0 },
      });

      await getTracks({ sortBy: 'title', sortOrder: 'asc' });

      expect(apiClient.get).toHaveBeenCalledWith('/tracks', {
        params: { sortBy: 'title', sortOrder: 'asc' },
      });
    });

    it('should fetch tracks with artist filter', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockTrack], total: 1, limit: 20, offset: 0 },
      });

      await getTracks({ artist: 'Test Artist' });

      expect(apiClient.get).toHaveBeenCalledWith('/tracks', {
        params: { artist: 'Test Artist' },
      });
    });

    it('should fetch tracks with album filter', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockTrack], total: 1, limit: 20, offset: 0 },
      });

      await getTracks({ album: 'Test Album' });

      expect(apiClient.get).toHaveBeenCalledWith('/tracks', {
        params: { album: 'Test Album' },
      });
    });

    it('should fetch tracks with search query', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockTrack], total: 1, limit: 20, offset: 0 },
      });

      await getTracks({ search: 'test' });

      expect(apiClient.get).toHaveBeenCalledWith('/tracks', {
        params: { search: 'test' },
      });
    });
  });

  describe('getTrack', () => {
    it('should fetch a single track by ID', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockTrack });

      const result = await getTrack('track-1');

      expect(apiClient.get).toHaveBeenCalledWith('/tracks/track-1');
      expect(result).toEqual(mockTrack);
    });
  });

  describe('updateTrack', () => {
    it('should update track metadata', async () => {
      const updatedTrack = { ...mockTrack, title: 'Updated Title' };
      vi.mocked(apiClient.patch).mockResolvedValue({ data: updatedTrack });

      const result = await updateTrack('track-1', { title: 'Updated Title' });

      expect(apiClient.patch).toHaveBeenCalledWith('/tracks/track-1', { title: 'Updated Title' });
      expect(result.title).toBe('Updated Title');
    });

    it('should update track tags', async () => {
      const updatedTrack = { ...mockTrack, tags: ['jazz', 'classical'] };
      vi.mocked(apiClient.patch).mockResolvedValue({ data: updatedTrack });

      const result = await updateTrack('track-1', { tags: ['jazz', 'classical'] });

      expect(apiClient.patch).toHaveBeenCalledWith('/tracks/track-1', { tags: ['jazz', 'classical'] });
      expect(result.tags).toEqual(['jazz', 'classical']);
    });
  });

  describe('deleteTrack', () => {
    it('should delete a track by ID', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({ data: undefined });

      await deleteTrack('track-1');

      expect(apiClient.delete).toHaveBeenCalledWith('/tracks/track-1');
    });
  });
});
