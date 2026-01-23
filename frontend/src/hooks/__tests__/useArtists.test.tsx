/**
 * useArtists Hook Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/artists', () => ({
  getArtists: vi.fn(),
  getArtist: vi.fn(),
}));

import { useArtistsQuery, useArtistQuery, artistKeys } from '../useArtists';
import * as artistsApi from '../../lib/api/artists';

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

describe('useArtists Hooks (Wave 2)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockArtist = {
    id: 'artist-1',
    name: 'Test Artist',
    trackCount: 25,
    albumCount: 3,
  };

  describe('artistKeys', () => {
    it('should generate correct key for all artists', () => {
      expect(artistKeys.all).toEqual(['artists']);
    });

    it('should generate correct key for artist lists', () => {
      expect(artistKeys.lists()).toEqual(['artists', 'list']);
    });

    it('should generate correct key for artist list with params', () => {
      expect(artistKeys.list({ search: 'Test' })).toEqual(['artists', 'list', { search: 'Test' }]);
    });

    it('should generate correct key for artist details', () => {
      expect(artistKeys.details()).toEqual(['artists', 'detail']);
    });

    it('should generate correct key for specific artist', () => {
      expect(artistKeys.detail('Test Artist')).toEqual(['artists', 'detail', 'Test Artist']);
    });
  });

  describe('useArtistsQuery', () => {
    it('should return loading state initially', () => {
      vi.mocked(artistsApi.getArtists).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useArtistsQuery(), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(true);
    });

    it('should return artists data on success', async () => {
      vi.mocked(artistsApi.getArtists).mockResolvedValue({
        items: [mockArtist],
        total: 1,
        limit: 20,
        offset: 0,
      });

      const { result } = renderHook(() => useArtistsQuery(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });
    });

    it('should pass params to getArtists', async () => {
      vi.mocked(artistsApi.getArtists).mockResolvedValue({
        items: [],
        total: 0,
        limit: 20,
        offset: 0,
      });

      renderHook(() => useArtistsQuery({ search: 'Beatles' }), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(artistsApi.getArtists).toHaveBeenCalledWith({ search: 'Beatles' });
      });
    });

    it('should return error state on failure', async () => {
      vi.mocked(artistsApi.getArtists).mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => useArtistsQuery(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useArtistQuery', () => {
    it('should fetch single artist by name', async () => {
      const artistWithDetails = {
        ...mockArtist,
        albums: [{ id: 'album-1', name: 'Album 1', year: 2024 }],
        recentTracks: [{ id: 'track-1', title: 'Track 1', duration: 180 }],
      };
      vi.mocked(artistsApi.getArtist).mockResolvedValue(artistWithDetails as unknown as ReturnType<typeof artistsApi.getArtist>);

      const { result } = renderHook(() => useArtistQuery('Test Artist'), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data).toEqual(artistWithDetails);
      });
    });

    it('should not fetch when name is undefined', () => {
      const { result } = renderHook(() => useArtistQuery(undefined), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(false);
      expect(artistsApi.getArtist).not.toHaveBeenCalled();
    });
  });
});
