/**
 * usePlaylists Hook Tests - Wave 5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { usePlaylistsQuery, usePlaylistQuery, playlistKeys } from '../usePlaylists';
import * as playlistsApi from '../../lib/api/playlists';

vi.mock('../../lib/api/playlists', () => ({
  getPlaylists: vi.fn(),
  getPlaylist: vi.fn(),
  createPlaylist: vi.fn(),
  updatePlaylist: vi.fn(),
  deletePlaylist: vi.fn(),
  addTrackToPlaylist: vi.fn(),
  removeTrackFromPlaylist: vi.fn(),
}));

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('usePlaylists (Wave 5)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('playlistKeys', () => {
    it('should generate correct query keys', () => {
      expect(playlistKeys.all).toEqual(['playlists']);
      expect(playlistKeys.lists()).toEqual(['playlists', 'list']);
      expect(playlistKeys.list({ limit: 20 })).toEqual(['playlists', 'list', { limit: 20 }]);
      expect(playlistKeys.details()).toEqual(['playlists', 'detail']);
      expect(playlistKeys.detail('playlist-1')).toEqual(['playlists', 'detail', 'playlist-1']);
    });
  });

  describe('usePlaylistsQuery', () => {
    it('should fetch playlists', async () => {
      const mockPlaylists = {
        items: [
          { id: 'playlist-1', name: 'Favorites', trackIds: [], trackCount: 5, createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
        ],
        total: 1,
        limit: 20,
        offset: 0,
      };
      vi.mocked(playlistsApi.getPlaylists).mockResolvedValue(mockPlaylists);

      const { result } = renderHook(() => usePlaylistsQuery(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockPlaylists);
      expect(playlistsApi.getPlaylists).toHaveBeenCalled();
    });

    it('should pass parameters to API', async () => {
      const params = { limit: 10, offset: 20 };
      const mockPlaylists = { items: [], total: 0, limit: 10, offset: 20 };
      vi.mocked(playlistsApi.getPlaylists).mockResolvedValue(mockPlaylists);

      const { result } = renderHook(() => usePlaylistsQuery(params), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(playlistsApi.getPlaylists).toHaveBeenCalledWith(params);
    });

    it('should handle error state', async () => {
      vi.mocked(playlistsApi.getPlaylists).mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => usePlaylistsQuery(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toBeDefined();
    });
  });

  describe('usePlaylistQuery', () => {
    it('should fetch a single playlist with tracks', async () => {
      const mockPlaylistWithTracks = {
        playlist: {
          id: 'playlist-1',
          name: 'Favorites',
          trackIds: ['track-1', 'track-2'],
          trackCount: 2,
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        tracks: [
          { id: 'track-1', title: 'Song 1', artist: 'Artist 1', album: 'Album 1', duration: 180, format: 'mp3', fileSize: 1000, s3Key: 'key1', tags: [], createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
          { id: 'track-2', title: 'Song 2', artist: 'Artist 2', album: 'Album 2', duration: 200, format: 'mp3', fileSize: 1200, s3Key: 'key2', tags: [], createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
        ],
      };
      vi.mocked(playlistsApi.getPlaylist).mockResolvedValue(mockPlaylistWithTracks);

      const { result } = renderHook(() => usePlaylistQuery('playlist-1'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockPlaylistWithTracks);
      expect(result.current.data?.playlist.name).toBe('Favorites');
      expect(result.current.data?.tracks).toHaveLength(2);
      expect(playlistsApi.getPlaylist).toHaveBeenCalledWith('playlist-1');
    });

    it('should not fetch when id is undefined', () => {
      const { result } = renderHook(() => usePlaylistQuery(undefined), {
        wrapper: createWrapper(),
      });

      expect(result.current.isFetching).toBe(false);
      expect(playlistsApi.getPlaylist).not.toHaveBeenCalled();
    });
  });
});
