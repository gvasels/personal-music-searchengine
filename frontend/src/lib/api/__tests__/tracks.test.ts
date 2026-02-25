/**
 * Tracks API Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
    post: vi.fn(),
  },
}));

import { getTracks, getTrack, updateTrack, deleteTrack, reprocessTrack } from '../tracks';
import { apiClient } from '../client';
import type { ReprocessResult } from '../../../types';

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

  /**
   * Admin Track Reprocess API Tests - TDD Red Phase
   *
   * Tests for the reprocessTrack API function that triggers AI reanalysis.
   */
  describe('reprocessTrack', () => {
    const mockReprocessResult: ReprocessResult = {
      trackId: 'track-1',
      status: 'processing',
      processedAt: '2026-02-17T10:30:00Z',
    };

    it('sends POST request to /tracks/:id/reprocess', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ data: mockReprocessResult });

      await reprocessTrack('track-1');

      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/tracks/track-1/reprocess');
    });

    it('returns ReprocessResult on success', async () => {
      const completeResult: ReprocessResult = {
        trackId: 'track-1',
        status: 'complete',
        bpm: 128.5,
        bpmConfidence: 0.92,
        musicalKey: 'G minor',
        keyCamelot: '6A',
        embeddingStatus: 'updated',
        processedAt: '2026-02-17T10:30:05Z',
      };
      vi.mocked(apiClient.post).mockResolvedValue({ data: completeResult });

      const result = await reprocessTrack('track-1');

      expect(result).toEqual(completeResult);
      expect(result.trackId).toBe('track-1');
      expect(result.status).toBe('complete');
      expect(result.bpm).toBe(128.5);
    });

    it('throws error on API failure', async () => {
      const apiError = new Error('Forbidden');
      vi.mocked(apiClient.post).mockRejectedValue(apiError);

      await expect(reprocessTrack('track-1')).rejects.toThrow('Forbidden');
    });

    it('handles different track IDs correctly', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({
        data: { ...mockReprocessResult, trackId: 'different-track-id' },
      });

      await reprocessTrack('different-track-id');

      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/tracks/different-track-id/reprocess');
    });
  });
});
