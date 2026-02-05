/**
 * useHelloSearch Hook Tests - TDD Red Phase
 *
 * Tests for the hello search hooks (useHelloSearch, useFeaturedTracks).
 * Source module (useHelloSearch.ts) does NOT exist yet - these tests MUST fail.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/helloSearch', () => ({
  searchHelloTracks: vi.fn(),
  getFeaturedTracks: vi.fn(),
}));

import { useHelloSearch, useFeaturedTracks, helloKeys } from '../useHelloSearch';
import * as helloApi from '../../lib/api/helloSearch';

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

  describe('helloKeys', () => {
    it('should return correct base key for all', () => {
      expect(helloKeys.all).toEqual(['hello']);
    });

    it('should return correct key for search with query', () => {
      expect(helloKeys.search('test')).toEqual(['hello', 'search', 'test']);
    });

    it('should return correct key for featured', () => {
      expect(helloKeys.featured()).toEqual(['hello', 'featured']);
    });
  });

  describe('useHelloSearch', () => {
    it('should be disabled when query is empty', () => {
      const { result } = renderHook(() => useHelloSearch(''), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(false);
      expect(helloApi.searchHelloTracks).not.toHaveBeenCalled();
    });

    it('should fetch when query is non-empty', async () => {
      vi.mocked(helloApi.searchHelloTracks).mockResolvedValue({
        items: [
          {
            id: 'track-1',
            title: 'Test Song',
            artist: 'Test Artist',
            album: 'Test Album',
            genre: 'Rock',
            year: 2024,
            duration: 200,
          },
        ],
        total: 1,
      });

      const { result } = renderHook(() => useHelloSearch('test'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });

      expect(helloApi.searchHelloTracks).toHaveBeenCalledWith('test');
    });
  });

  describe('useFeaturedTracks', () => {
    it('should fetch featured tracks on mount', async () => {
      vi.mocked(helloApi.getFeaturedTracks).mockResolvedValue({
        items: [
          {
            id: 'featured-1',
            title: 'Featured Track',
            artist: 'Featured Artist',
            album: 'Featured Album',
            genre: 'Pop',
            year: 2024,
            duration: 180,
          },
        ],
        total: 1,
      });

      const { result } = renderHook(() => useFeaturedTracks(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });

      expect(helloApi.getFeaturedTracks).toHaveBeenCalled();
    });
  });
});
