/**
 * useTracks Hook Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/tracks', () => ({
  getTracks: vi.fn(),
  getTrack: vi.fn(),
  updateTrack: vi.fn(),
  deleteTrack: vi.fn(),
}));

import { useTracksQuery, useTrackQuery, useUpdateTrack, useDeleteTrack, trackKeys } from '../useTracks';
import * as tracksApi from '../../lib/api/tracks';

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

describe('useTracks Hooks (Wave 2)', () => {
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
    s3Key: 'tracks/track-1.mp3',
    tags: [] as string[],
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  };

  describe('trackKeys', () => {
    it('should generate correct key for all tracks', () => {
      expect(trackKeys.all).toEqual(['tracks']);
    });

    it('should generate correct key for track lists', () => {
      expect(trackKeys.lists()).toEqual(['tracks', 'list']);
    });

    it('should generate correct key for track list with params', () => {
      expect(trackKeys.list({ artist: 'Test' })).toEqual(['tracks', 'list', { artist: 'Test' }]);
    });

    it('should generate correct key for track details', () => {
      expect(trackKeys.details()).toEqual(['tracks', 'detail']);
    });

    it('should generate correct key for specific track', () => {
      expect(trackKeys.detail('track-1')).toEqual(['tracks', 'detail', 'track-1']);
    });
  });

  describe('useTracksQuery', () => {
    it('should return loading state initially', () => {
      vi.mocked(tracksApi.getTracks).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useTracksQuery(), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(true);
    });

    it('should return tracks data on success', async () => {
      vi.mocked(tracksApi.getTracks).mockResolvedValue({
        items: [mockTrack],
        total: 1,
        limit: 20,
        offset: 0,
      });

      const { result } = renderHook(() => useTracksQuery(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });
    });

    it('should pass params to getTracks', async () => {
      vi.mocked(tracksApi.getTracks).mockResolvedValue({
        items: [],
        total: 0,
        limit: 20,
        offset: 0,
      });

      renderHook(() => useTracksQuery({ artist: 'Test Artist' }), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(tracksApi.getTracks).toHaveBeenCalledWith({ artist: 'Test Artist' });
      });
    });

    it('should return error state on failure', async () => {
      vi.mocked(tracksApi.getTracks).mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => useTracksQuery(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useTrackQuery', () => {
    it('should fetch single track by ID', async () => {
      vi.mocked(tracksApi.getTrack).mockResolvedValue(mockTrack);

      const { result } = renderHook(() => useTrackQuery('track-1'), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data).toEqual(mockTrack);
      });
    });

    it('should not fetch when ID is undefined', () => {
      const { result } = renderHook(() => useTrackQuery(undefined), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(false);
      expect(tracksApi.getTrack).not.toHaveBeenCalled();
    });
  });

  describe('useUpdateTrack', () => {
    it('should update track and return updated data', async () => {
      const updatedTrack = { ...mockTrack, title: 'Updated Title' };
      vi.mocked(tracksApi.updateTrack).mockResolvedValue(updatedTrack);

      const { result } = renderHook(() => useUpdateTrack(), { wrapper: createWrapper() });

      await act(async () => {
        const data = await result.current.mutateAsync({ id: 'track-1', data: { title: 'Updated Title' } });
        expect(data.title).toBe('Updated Title');
      });

      expect(tracksApi.updateTrack).toHaveBeenCalledWith('track-1', { title: 'Updated Title' });
    });
  });

  describe('useDeleteTrack', () => {
    it('should delete track', async () => {
      vi.mocked(tracksApi.deleteTrack).mockResolvedValue(undefined);

      const { result } = renderHook(() => useDeleteTrack(), { wrapper: createWrapper() });

      await act(async () => {
        await result.current.mutateAsync('track-1');
      });

      expect(tracksApi.deleteTrack).toHaveBeenCalledWith('track-1');
    });
  });
});
