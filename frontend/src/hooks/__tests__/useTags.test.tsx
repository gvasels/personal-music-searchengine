/**
 * useTags Hook Tests - Wave 5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useTagsQuery, useTracksByTagQuery, tagKeys } from '../useTags';
import * as tagsApi from '../../lib/api/tags';

vi.mock('../../lib/api/tags', () => ({
  getTags: vi.fn(),
  getTracksByTag: vi.fn(),
  addTagToTrack: vi.fn(),
  removeTagFromTrack: vi.fn(),
}));

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('useTags (Wave 5)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('tagKeys', () => {
    it('should generate correct query keys', () => {
      expect(tagKeys.all).toEqual(['tags']);
      expect(tagKeys.lists()).toEqual(['tags', 'list']);
      expect(tagKeys.list({})).toEqual(['tags', 'list', {}]);
      expect(tagKeys.tracks('rock')).toEqual(['tags', 'tracks', 'rock']);
      expect(tagKeys.trackList('rock', { limit: 20 })).toEqual([
        'tags',
        'tracks',
        'rock',
        { limit: 20 },
      ]);
    });
  });

  describe('useTagsQuery', () => {
    it('should fetch all tags', async () => {
      const mockTags = {
        items: [
          { name: 'rock', trackCount: 15 },
          { name: 'jazz', trackCount: 8 },
        ],
        total: 2,
        limit: 50,
        offset: 0,
      };
      vi.mocked(tagsApi.getTags).mockResolvedValue(mockTags);

      const { result } = renderHook(() => useTagsQuery(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockTags);
      expect(tagsApi.getTags).toHaveBeenCalled();
    });

    it('should handle error state', async () => {
      vi.mocked(tagsApi.getTags).mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => useTagsQuery(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toBeDefined();
    });
  });

  describe('useTracksByTagQuery', () => {
    it('should fetch tracks for a tag', async () => {
      const mockTracks = {
        items: [
          { id: 'track-1', title: 'Rock Song', tags: ['rock'] },
        ],
        total: 1,
        limit: 20,
        offset: 0,
      };
      vi.mocked(tagsApi.getTracksByTag).mockResolvedValue(mockTracks as any);

      const { result } = renderHook(() => useTracksByTagQuery('rock'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockTracks);
      expect(tagsApi.getTracksByTag).toHaveBeenCalledWith('rock', undefined);
    });

    it('should pass pagination params', async () => {
      const params = { limit: 10, offset: 20 };
      const mockTracks = { items: [], total: 0, limit: 10, offset: 20 };
      vi.mocked(tagsApi.getTracksByTag).mockResolvedValue(mockTracks as any);

      const { result } = renderHook(() => useTracksByTagQuery('rock', params), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(tagsApi.getTracksByTag).toHaveBeenCalledWith('rock', params);
    });

    it('should not fetch when tagName is empty', () => {
      const { result } = renderHook(() => useTracksByTagQuery(''), {
        wrapper: createWrapper(),
      });

      expect(result.current.isFetching).toBe(false);
      expect(tagsApi.getTracksByTag).not.toHaveBeenCalled();
    });
  });
});
