/**
 * useAlbums Hook Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/albums', () => ({
  getAlbums: vi.fn(),
  getAlbum: vi.fn(),
}));

import { useAlbumsQuery, useAlbumQuery, albumKeys } from '../useAlbums';
import * as albumsApi from '../../lib/api/albums';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0 },
      mutations: { retry: false },
    },
  });

  return function Wrapper({ children }: { children: ReactNode }) {
    return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
  };
}

describe('useAlbums Hooks (Wave 2)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockAlbum = {
    id: 'album-1',
    name: 'Test Album',
    artist: 'Test Artist',
    year: 2024,
    trackCount: 10,
    totalDuration: 3600,
    createdAt: '2024-01-01T00:00:00Z',
  };

  describe('albumKeys', () => {
    it('should generate correct key for all albums', () => {
      expect(albumKeys.all).toEqual(['albums']);
    });

    it('should generate correct key for album lists', () => {
      expect(albumKeys.lists()).toEqual(['albums', 'list']);
    });

    it('should generate correct key for album list with params', () => {
      expect(albumKeys.list({ artist: 'Test' })).toEqual(['albums', 'list', { artist: 'Test' }]);
    });

    it('should generate correct key for album details', () => {
      expect(albumKeys.details()).toEqual(['albums', 'detail']);
    });

    it('should generate correct key for specific album', () => {
      expect(albumKeys.detail('album-1')).toEqual(['albums', 'detail', 'album-1']);
    });
  });

  describe('useAlbumsQuery', () => {
    it('should return loading state initially', () => {
      vi.mocked(albumsApi.getAlbums).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useAlbumsQuery(), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(true);
    });

    it('should return albums data on success', async () => {
      vi.mocked(albumsApi.getAlbums).mockResolvedValue({
        items: [mockAlbum],
        total: 1,
        limit: 24,
        offset: 0,
      });

      const { result } = renderHook(() => useAlbumsQuery(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });
    });

    it('should pass params to getAlbums', async () => {
      vi.mocked(albumsApi.getAlbums).mockResolvedValue({
        items: [],
        total: 0,
        limit: 24,
        offset: 0,
      });

      renderHook(() => useAlbumsQuery({ artist: 'Test Artist' }), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(albumsApi.getAlbums).toHaveBeenCalledWith({ artist: 'Test Artist' });
      });
    });

    it('should return error state on failure', async () => {
      vi.mocked(albumsApi.getAlbums).mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => useAlbumsQuery(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useAlbumQuery', () => {
    it('should fetch single album by ID', async () => {
      const albumWithTracks = {
        ...mockAlbum,
        tracks: [{ id: 'track-1', title: 'Track 1', duration: 180 }],
      };
      vi.mocked(albumsApi.getAlbum).mockResolvedValue(albumWithTracks as any);

      const { result } = renderHook(() => useAlbumQuery('album-1'), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data).toEqual(albumWithTracks);
      });
    });

    it('should not fetch when ID is undefined', () => {
      const { result } = renderHook(() => useAlbumQuery(undefined), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(false);
      expect(albumsApi.getAlbum).not.toHaveBeenCalled();
    });
  });
});
