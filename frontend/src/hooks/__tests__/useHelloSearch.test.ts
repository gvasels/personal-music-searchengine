import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useHelloSearch, useFeaturedTracks } from '../useHelloSearch';
import * as helloSearchApi from '../../lib/api/helloSearch';
import React from 'react';

vi.mock('../../lib/api/helloSearch');

const mockResponse = {
  tracks: [
    { id: '1', title: 'Midnight Drift', artist: 'Luna Waves', album: 'Waveforms', genre: 'Electronic', year: 2024, duration: 240, durationStr: '4:00', coverArtUrl: '' },
  ],
  total: 1,
  query: 'luna',
};

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return function Wrapper({ children }: { children: React.ReactNode }) {
    return React.createElement(QueryClientProvider, { client: queryClient }, children);
  };
}

describe('useHelloSearch', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('returns isLoading true initially when query is provided', () => {
    vi.mocked(helloSearchApi.searchHello).mockReturnValue(new Promise(() => {}));
    const { result } = renderHook(() => useHelloSearch('luna'), { wrapper: createWrapper() });
    expect(result.current.isLoading).toBe(true);
  });

  it('returns track data on success', async () => {
    vi.mocked(helloSearchApi.searchHello).mockResolvedValue(mockResponse);
    const { result } = renderHook(() => useHelloSearch('luna'), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockResponse);
  });

  it('returns isError true on API failure', async () => {
    vi.mocked(helloSearchApi.searchHello).mockRejectedValue(new Error('Network error'));
    const { result } = renderHook(() => useHelloSearch('luna'), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isError).toBe(true));
  });

  it('skips query when search term is empty', () => {
    const { result } = renderHook(() => useHelloSearch(''), { wrapper: createWrapper() });
    expect(result.current.fetchStatus).toBe('idle');
  });
});

describe('useFeaturedTracks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('returns featured tracks on success', async () => {
    const featured = { ...mockResponse, query: '' };
    vi.mocked(helloSearchApi.getFeaturedTracks).mockResolvedValue(featured);
    const { result } = renderHook(() => useFeaturedTracks(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(featured);
  });
});
