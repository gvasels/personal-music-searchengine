/**
 * useArtistWatch Hook Tests - TDD Red Phase
 *
 * These tests define the contract for the useArtistWatch hooks.
 * They MUST FAIL because useArtistWatch.ts does not exist yet.
 *
 * Tests verify:
 * - watchKeys query key factory
 * - useWatchStatus hook for checking watch state
 * - useWatchToggle hook for watch/unwatch mutations
 * - useWatchedArtists hook for listing watched artists
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/artistWatch', () => ({
  watchArtist: vi.fn(),
  unwatchArtist: vi.fn(),
  getWatchStatus: vi.fn(),
  getWatchedArtists: vi.fn(),
}));

vi.mock('../useAuth', () => ({
  useAuth: () => ({ isAuthenticated: true }),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn(), success: vi.fn() },
}));

// These imports will fail - useArtistWatch.ts does not exist yet
import { useWatchStatus, useWatchToggle, useWatchedArtists, watchKeys } from '../useArtistWatch';
import * as artistWatchApi from '../../lib/api/artistWatch';

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

describe('useArtistWatch Hooks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('watchKeys', () => {
    it('should generate correct key for all watch queries', () => {
      expect(watchKeys.all).toEqual(['artist-watch']);
    });

    it('should generate correct key for watch status', () => {
      expect(watchKeys.status('Daft Punk')).toEqual(['artist-watch', 'status', 'Daft Punk']);
    });

    it('should generate correct key for watched artists list', () => {
      expect(watchKeys.list()).toEqual(['artist-watch', 'list']);
    });
  });

  describe('useWatchStatus', () => {
    it('should return watching=false when not watching', async () => {
      vi.mocked(artistWatchApi.getWatchStatus).mockResolvedValue({
        watching: false,
        artistName: 'Daft Punk',
      });

      const { result } = renderHook(() => useWatchStatus('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.watching).toBe(false);
      });

      expect(artistWatchApi.getWatchStatus).toHaveBeenCalledWith('Daft Punk');
    });

    it('should return watching=true when watching', async () => {
      vi.mocked(artistWatchApi.getWatchStatus).mockResolvedValue({
        watching: true,
        artistName: 'Daft Punk',
      });

      const { result } = renderHook(() => useWatchStatus('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.watching).toBe(true);
      });
    });

    it('should return loading state initially', () => {
      vi.mocked(artistWatchApi.getWatchStatus).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useWatchStatus('Daft Punk'), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(true);
    });

    it('should not fetch when artistName is undefined', () => {
      const { result } = renderHook(() => useWatchStatus(undefined), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(false);
      expect(artistWatchApi.getWatchStatus).not.toHaveBeenCalled();
    });

    it('should return error state on failure', async () => {
      vi.mocked(artistWatchApi.getWatchStatus).mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => useWatchStatus('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useWatchToggle', () => {
    it('should return isWatching, isLoading, and toggle function', async () => {
      vi.mocked(artistWatchApi.getWatchStatus).mockResolvedValue({
        watching: false,
        artistName: 'Daft Punk',
      });

      const { result } = renderHook(() => useWatchToggle('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current).toHaveProperty('isWatching');
        expect(result.current).toHaveProperty('isLoading');
        expect(result.current).toHaveProperty('toggle');
        expect(typeof result.current.toggle).toBe('function');
      });
    });

    it('should call watchArtist when not currently watching', async () => {
      vi.mocked(artistWatchApi.getWatchStatus).mockResolvedValue({
        watching: false,
        artistName: 'Daft Punk',
      });
      vi.mocked(artistWatchApi.watchArtist).mockResolvedValue({
        artistName: 'Daft Punk',
        watching: true,
        watchedAt: '2026-02-25T10:00:00Z',
      });

      const { result } = renderHook(() => useWatchToggle('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isWatching).toBe(false);
      });

      await act(async () => {
        await result.current.toggle();
      });

      expect(artistWatchApi.watchArtist).toHaveBeenCalledWith('Daft Punk');
    });

    it('should call unwatchArtist when currently watching', async () => {
      vi.mocked(artistWatchApi.getWatchStatus).mockResolvedValue({
        watching: true,
        artistName: 'Daft Punk',
      });
      vi.mocked(artistWatchApi.unwatchArtist).mockResolvedValue(undefined);

      const { result } = renderHook(() => useWatchToggle('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isWatching).toBe(true);
      });

      await act(async () => {
        await result.current.toggle();
      });

      expect(artistWatchApi.unwatchArtist).toHaveBeenCalledWith('Daft Punk');
    });
  });

  describe('useWatchedArtists', () => {
    it('should return list of watched artists', async () => {
      vi.mocked(artistWatchApi.getWatchedArtists).mockResolvedValue({
        items: [
          { artistName: 'Daft Punk', watchedAt: '2026-02-25T10:00:00Z' },
          { artistName: 'Aphex Twin', watchedAt: '2026-02-24T08:00:00Z' },
        ],
        hasMore: false,
      });

      const { result } = renderHook(() => useWatchedArtists(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(2);
        expect(result.current.data?.items[0].artistName).toBe('Daft Punk');
      });
    });

    it('should return loading state initially', () => {
      vi.mocked(artistWatchApi.getWatchedArtists).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useWatchedArtists(), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(true);
    });

    it('should support pagination via hasMore and nextCursor', async () => {
      vi.mocked(artistWatchApi.getWatchedArtists).mockResolvedValue({
        items: [{ artistName: 'Daft Punk', watchedAt: '2026-02-25T10:00:00Z' }],
        nextCursor: 'cursor-abc',
        hasMore: true,
      });

      const { result } = renderHook(() => useWatchedArtists(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.hasMore).toBe(true);
        expect(result.current.data?.nextCursor).toBe('cursor-abc');
      });
    });

    it('should return error state on failure', async () => {
      vi.mocked(artistWatchApi.getWatchedArtists).mockRejectedValue(new Error('Forbidden'));

      const { result } = renderHook(() => useWatchedArtists(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });
});
