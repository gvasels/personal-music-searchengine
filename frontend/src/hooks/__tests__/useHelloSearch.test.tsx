/**
 * useHelloSearch Hook Tests - TDD Red Phase
 *
 * These tests define the contract for the useHelloSearch hooks.
 * They MUST FAIL because useHelloSearch.ts does not exist yet.
 *
 * Tests verify:
 * - helloKeys query key factory
 * - useHelloSearch hook with enabled/disabled states
 * - useHelloFeatured hook
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/helloSearch', () => ({
  searchHelloTracks: vi.fn(),
  getFeaturedHelloTracks: vi.fn(),
}));

// This import will fail - useHelloSearch.ts does not exist
import { useHelloSearch, useHelloFeatured, helloKeys } from '../useHelloSearch';
import * as helloSearchApi from '../../lib/api/helloSearch';

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

describe('useHelloSearch Hooks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockHelloTrack = {
    id: 'seed-t1',
    title: 'Aurora Borealis',
    artist: 'Aurora Waves',
    album: 'Dreamscape',
    genre: 'jazz',
    year: 2023,
    duration: 245,
  };

  const mockSearchResponse = {
    items: [mockHelloTrack],
    total: 1,
  };

  describe('helloKeys', () => {
    it('should generate correct key for all hello queries', () => {
      expect(helloKeys.all).toEqual(['hello']);
    });

    it('should generate correct key for search queries', () => {
      expect(helloKeys.search('jazz')).toEqual(['hello', 'search', 'jazz']);
    });

    it('should generate correct key for featured queries', () => {
      expect(helloKeys.featured()).toEqual(['hello', 'featured']);
    });
  });

  describe('useHelloSearch', () => {
    it('should not fetch when query is empty (disabled)', () => {
      const { result } = renderHook(() => useHelloSearch(''), { wrapper: createWrapper() });

      // Query should be disabled - not loading, no request made
      expect(result.current.isLoading).toBe(false);
      expect(result.current.isFetching).toBe(false);
      expect(helloSearchApi.searchHelloTracks).not.toHaveBeenCalled();
    });

    it('should fetch search results when query is provided', async () => {
      vi.mocked(helloSearchApi.searchHelloTracks).mockResolvedValue(mockSearchResponse);

      const { result } = renderHook(() => useHelloSearch('jazz'), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });

      expect(helloSearchApi.searchHelloTracks).toHaveBeenCalledWith('jazz');
    });

    it('should return items array in data (NOT tracks)', async () => {
      vi.mocked(helloSearchApi.searchHelloTracks).mockResolvedValue(mockSearchResponse);

      const { result } = renderHook(() => useHelloSearch('electronic'), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data).toHaveProperty('items');
        expect(result.current.data?.items).toBeInstanceOf(Array);
      });
    });

    it('should return error on failure', async () => {
      vi.mocked(helloSearchApi.searchHelloTracks).mockRejectedValue(new Error('Search failed'));

      const { result } = renderHook(() => useHelloSearch('test'), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });

    it('should return loading state while fetching', () => {
      vi.mocked(helloSearchApi.searchHelloTracks).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useHelloSearch('jazz'), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(true);
    });
  });

  describe('useHelloFeatured', () => {
    it('should fetch featured tracks', async () => {
      vi.mocked(helloSearchApi.getFeaturedHelloTracks).mockResolvedValue(mockSearchResponse);

      const { result } = renderHook(() => useHelloFeatured(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });

      expect(helloSearchApi.getFeaturedHelloTracks).toHaveBeenCalled();
    });

    it('should return items array in data (NOT tracks)', async () => {
      vi.mocked(helloSearchApi.getFeaturedHelloTracks).mockResolvedValue(mockSearchResponse);

      const { result } = renderHook(() => useHelloFeatured(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data).toHaveProperty('items');
        expect(result.current.data?.items).toBeInstanceOf(Array);
      });
    });

    it('should return loading state initially', () => {
      vi.mocked(helloSearchApi.getFeaturedHelloTracks).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useHelloFeatured(), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(true);
    });

    it('should return error on failure', async () => {
      vi.mocked(helloSearchApi.getFeaturedHelloTracks).mockRejectedValue(new Error('Failed to fetch'));

      const { result } = renderHook(() => useHelloFeatured(), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });
});
