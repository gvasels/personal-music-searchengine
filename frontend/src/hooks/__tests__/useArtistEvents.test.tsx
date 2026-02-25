/**
 * useArtistEvents Hook Tests - TDD Red Phase
 *
 * These tests define the contract for the useArtistEvents hooks.
 * They MUST FAIL because useArtistEvents.ts does not exist yet.
 *
 * Tests verify:
 * - eventKeys query key factory
 * - useArtistEvents hook for fetching artist events
 * - useSearchArtistEvents hook for searching events
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/events', () => ({
  getArtistEvents: vi.fn(),
  searchArtistEvents: vi.fn(),
}));

// These imports will fail - useArtistEvents.ts does not exist yet
import { useArtistEvents, useSearchArtistEvents, eventKeys } from '../useArtistEvents';
import * as eventsApi from '../../lib/api/events';

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

describe('useArtistEvents Hooks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockEvent = {
    id: 'evt-1',
    artistName: 'Daft Punk',
    title: 'Alive 2027 Tour',
    date: '2027-06-15T20:00:00Z',
    venue: 'Madison Square Garden',
    city: 'New York',
    region: 'NY',
    country: 'US',
    ticketUrl: 'https://tickets.example.com/daft-punk',
    status: 'scheduled' as const,
    source: 'ticketmaster',
  };

  const mockEventsResponse = {
    artistName: 'Daft Punk',
    events: [mockEvent],
    totalCount: 1,
    source: 'ticketmaster',
  };

  const mockSearchResult = {
    name: 'Daft Punk',
    imageUrl: 'https://img.example.com/daft-punk.jpg',
    upcomingEvents: 3,
    source: 'ticketmaster',
  };

  describe('eventKeys', () => {
    it('should generate correct key for all event queries', () => {
      expect(eventKeys.all).toEqual(['events']);
    });

    it('should generate correct key for artist events', () => {
      expect(eventKeys.artist('Daft Punk')).toEqual(['events', 'artist', 'Daft Punk']);
    });

    it('should generate correct key for event search', () => {
      expect(eventKeys.search('Daft Punk')).toEqual(['events', 'search', 'Daft Punk']);
    });
  });

  describe('useArtistEvents', () => {
    it('should return events data when available', async () => {
      vi.mocked(eventsApi.getArtistEvents).mockResolvedValue(mockEventsResponse);

      const { result } = renderHook(() => useArtistEvents('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.events).toHaveLength(1);
        expect(result.current.data?.artistName).toBe('Daft Punk');
        expect(result.current.data?.totalCount).toBe(1);
      });

      expect(eventsApi.getArtistEvents).toHaveBeenCalledWith('Daft Punk');
    });

    it('should return loading state initially', () => {
      vi.mocked(eventsApi.getArtistEvents).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useArtistEvents('Daft Punk'), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(true);
    });

    it('should return empty events for unknown artist', async () => {
      vi.mocked(eventsApi.getArtistEvents).mockResolvedValue({
        artistName: 'Unknown Artist',
        events: [],
        totalCount: 0,
        source: 'ticketmaster',
      });

      const { result } = renderHook(() => useArtistEvents('Unknown Artist'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.events).toEqual([]);
        expect(result.current.data?.totalCount).toBe(0);
      });
    });

    it('should not fetch when artistName is undefined', () => {
      const { result } = renderHook(() => useArtistEvents(undefined), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(false);
      expect(eventsApi.getArtistEvents).not.toHaveBeenCalled();
    });

    it('should return error state on failure', async () => {
      vi.mocked(eventsApi.getArtistEvents).mockRejectedValue(new Error('Service unavailable'));

      const { result } = renderHook(() => useArtistEvents('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useSearchArtistEvents', () => {
    it('should return search results', async () => {
      vi.mocked(eventsApi.searchArtistEvents).mockResolvedValue({
        items: [mockSearchResult],
      });

      const { result } = renderHook(() => useSearchArtistEvents('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.items).toHaveLength(1);
        expect(result.current.data?.items[0].name).toBe('Daft Punk');
        expect(result.current.data?.items[0].upcomingEvents).toBe(3);
      });

      expect(eventsApi.searchArtistEvents).toHaveBeenCalledWith('Daft Punk');
    });

    it('should not execute when query is empty', () => {
      const { result } = renderHook(() => useSearchArtistEvents(''), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.isFetching).toBe(false);
      expect(eventsApi.searchArtistEvents).not.toHaveBeenCalled();
    });

    it('should return loading state while fetching', () => {
      vi.mocked(eventsApi.searchArtistEvents).mockImplementation(() => new Promise(() => {}));

      const { result } = renderHook(() => useSearchArtistEvents('Daft Punk'), {
        wrapper: createWrapper(),
      });

      expect(result.current.isLoading).toBe(true);
    });

    it('should return error state on failure', async () => {
      vi.mocked(eventsApi.searchArtistEvents).mockRejectedValue(new Error('Search failed'));

      const { result } = renderHook(() => useSearchArtistEvents('Daft Punk'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });

    it('should return empty results for no matches', async () => {
      vi.mocked(eventsApi.searchArtistEvents).mockResolvedValue({ items: [] });

      const { result } = renderHook(() => useSearchArtistEvents('nonexistent'), {
        wrapper: createWrapper(),
      });

      await waitFor(() => {
        expect(result.current.data?.items).toEqual([]);
      });
    });
  });
});
