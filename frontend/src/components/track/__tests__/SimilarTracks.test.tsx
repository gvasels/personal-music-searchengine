import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { SimilarTracks } from '../SimilarTracks';

vi.mock('../../../hooks/useSimilarTracks', () => ({
  useSimilarTracks: vi.fn(),
}));

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, ...props }: { children: React.ReactNode; to: string; params?: Record<string, string> }) => (
    <a href={props.to} {...props}>{children}</a>
  ),
}));

import { useSimilarTracks } from '../../../hooks/useSimilarTracks';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('SimilarTracks', () => {
  it('shows loading state', () => {
    vi.mocked(useSimilarTracks).mockReturnValue({
      data: undefined,
      isLoading: true,
      error: null,
    } as ReturnType<typeof useSimilarTracks>);

    render(<SimilarTracks trackId="track-1" />, { wrapper: createWrapper() });
    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('shows similar tracks', () => {
    vi.mocked(useSimilarTracks).mockReturnValue({
      data: {
        similar: [
          { track: { id: 'track-2', title: 'Track Two', artist: 'Artist A', coverArtUrl: '', bpm: 120, keyCamelot: '8B' }, score: 0.05 },
          { track: { id: 'track-3', title: 'Track Three', artist: 'Artist B', coverArtUrl: '', bpm: 0, keyCamelot: '' }, score: 0.15 },
        ],
      },
      isLoading: false,
      error: null,
    } as ReturnType<typeof useSimilarTracks>);

    render(<SimilarTracks trackId="track-1" />, { wrapper: createWrapper() });
    expect(screen.getByText('Track Two')).toBeInTheDocument();
    expect(screen.getByText('Track Three')).toBeInTheDocument();
  });

  it('shows empty state when no similar tracks', () => {
    vi.mocked(useSimilarTracks).mockReturnValue({
      data: { similar: [] },
      isLoading: false,
      error: null,
    } as ReturnType<typeof useSimilarTracks>);

    render(<SimilarTracks trackId="track-1" />, { wrapper: createWrapper() });
    expect(screen.getByText(/no similar tracks/i)).toBeInTheDocument();
  });
});
