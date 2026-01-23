/**
 * useSearch Hook Tests - Wave 4
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/search', () => ({
  searchTracks: vi.fn(),
  searchAutocomplete: vi.fn(),
}));

import { useSearchQuery, useAutocompleteQuery, searchKeys } from '../useSearch';
import * as searchApi from '../../lib/api/search';

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

describe('useSearch Hooks (Wave 4)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('searchKeys', () => {
    it('should generate correct key for all searches', () => {
      expect(searchKeys.all).toEqual(['search']);
    });

    it('should generate correct key for search results', () => {
      expect(searchKeys.results({ query: 'test' })).toEqual(['search', 'results', { query: 'test' }]);
    });

    it('should generate correct key for autocomplete', () => {
      expect(searchKeys.autocomplete('test')).toEqual(['search', 'autocomplete', 'test']);
    });
  });

  describe('useSearchQuery', () => {
    it('should not fetch when query is empty', () => {
      const { result } = renderHook(() => useSearchQuery({ query: '' }), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(false);
      expect(searchApi.searchTracks).not.toHaveBeenCalled();
    });

    it('should fetch search results', async () => {
      vi.mocked(searchApi.searchTracks).mockResolvedValue({
        query: 'test',
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        tracks: [{ id: 'track-1', title: 'Test' }] as any,
        totalResults: 1,
        limit: 20,
        hasMore: false,
      });

      const { result } = renderHook(() => useSearchQuery({ query: 'test' }), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data?.tracks).toHaveLength(1);
      });
    });

    it('should return error on failure', async () => {
      vi.mocked(searchApi.searchTracks).mockRejectedValue(new Error('Search failed'));

      const { result } = renderHook(() => useSearchQuery({ query: 'test' }), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useAutocompleteQuery', () => {
    it('should not fetch when query is too short', () => {
      const { result } = renderHook(() => useAutocompleteQuery('ab'), { wrapper: createWrapper() });

      expect(result.current.isLoading).toBe(false);
      expect(searchApi.searchAutocomplete).not.toHaveBeenCalled();
    });

    it('should fetch autocomplete suggestions', async () => {
      vi.mocked(searchApi.searchAutocomplete).mockResolvedValue({
        suggestions: [{ type: 'track', value: 'Test Song' }],
      });

      const { result } = renderHook(() => useAutocompleteQuery('test'), { wrapper: createWrapper() });

      await waitFor(() => {
        expect(result.current.data?.suggestions).toHaveLength(1);
      });
    });
  });
});
