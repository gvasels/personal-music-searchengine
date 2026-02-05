/**
 * useHelloSearch Hook Tests - TDD Red Phase
 *
 * Tests for hello-world search hooks using TanStack Query.
 * CRITICAL CONTRACT:
 * - Response uses `items` field (NOT `tracks`)
 * - helloKeys factory generates correct query keys
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/helloSearch', () => ({
  searchHelloTracks: vi.fn(),
  getFeaturedHelloTracks: vi.fn(),
}));

import { useHelloSearch, useHelloFeatured, helloKeys } from '../useHelloSearch';
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
    it('should generate correct key for all hello queries', () => {
      expect(helloKeys.all).toEqual(['hello']);
    });

    it('should generate correct key for search with query', () => {
      expect(helloKeys.search('jazz')).toEqual(['hello', 'search', 'jazz']);
    });

    it('should generate correct key for featured', () => {
      expect(helloKeys.featured()).toEqual(['hello', 'featured']);
    });
  });

  describe('useHelloSearch', () => {
    it('should call searchHelloTracks and return data with items', async () => {
      const mockResponse = {
        items: [
          { id: 'track-1', title: 'Jazz Track', artist: 'Artist', album: 'Album', duration: 200 },
        ],
        total: 1,
      };
      vi.mocked(helloApi.searchHelloTracks).mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useHelloSearch('jazz'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });

      expect(helloApi.searchHelloTracks).toHaveBeenCalledWith('jazz');
    });

    it('should be disabled when query is empty', () => {
      const { result } = renderHook(() => useHelloSearch(''), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(false);
      expect(helloApi.searchHelloTracks).not.toHaveBeenCalled();
    });
  });

  describe('useHelloFeatured', () => {
    it('should call getFeaturedHelloTracks and return data', async () => {
      const mockResponse = {
        items: [
          { id: 'feat-1', title: 'Featured Track', artist: 'Artist', album: 'Album', duration: 240 },
        ],
        total: 1,
      };
      vi.mocked(helloApi.getFeaturedHelloTracks).mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useHelloFeatured(), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
      });

      expect(helloApi.getFeaturedHelloTracks).toHaveBeenCalled();
    });
  });
});
