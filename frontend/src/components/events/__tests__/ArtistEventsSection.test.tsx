/**
 * ArtistEventsSection Component Tests - TDD Red Phase
 *
 * These tests define the contract for the ArtistEventsSection component.
 * They MUST FAIL because ArtistEventsSection.tsx does not exist yet.
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ArtistEventsSection } from '../ArtistEventsSection';

vi.mock('../../../hooks/useArtistEvents', () => ({
  useArtistEvents: vi.fn(),
}));

import { useArtistEvents } from '../../../hooks/useArtistEvents';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('ArtistEventsSection', () => {
  it('should show loading state', () => {
    vi.mocked(useArtistEvents).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
      error: null,
    } as ReturnType<typeof useArtistEvents>);

    render(<ArtistEventsSection artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('should render events when data available', () => {
    vi.mocked(useArtistEvents).mockReturnValue({
      data: {
        artistName: 'Daft Punk',
        events: [
          {
            id: 'evt-1',
            artistName: 'Daft Punk',
            title: 'Alive 2027 Tour',
            date: '2027-06-15T20:00:00Z',
            venue: 'Madison Square Garden',
            city: 'New York',
            region: 'NY',
            country: 'US',
            status: 'scheduled',
            source: 'mock',
          },
          {
            id: 'evt-2',
            artistName: 'Daft Punk',
            title: 'Summer Festival',
            date: '2027-07-20T18:00:00Z',
            venue: 'Coachella',
            city: 'Indio',
            region: 'CA',
            country: 'US',
            status: 'scheduled',
            source: 'mock',
          },
        ],
        totalCount: 2,
        source: 'mock',
      },
      isLoading: false,
      isError: false,
      error: null,
    } as ReturnType<typeof useArtistEvents>);

    render(<ArtistEventsSection artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByText('Alive 2027 Tour')).toBeInTheDocument();
    expect(screen.getByText('Summer Festival')).toBeInTheDocument();
  });

  it('should show empty state when no events', () => {
    vi.mocked(useArtistEvents).mockReturnValue({
      data: {
        artistName: 'Daft Punk',
        events: [],
        totalCount: 0,
        source: 'mock',
      },
      isLoading: false,
      isError: false,
      error: null,
    } as ReturnType<typeof useArtistEvents>);

    render(<ArtistEventsSection artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByText(/no upcoming/i)).toBeInTheDocument();
  });

  it('should show error state on failure', () => {
    vi.mocked(useArtistEvents).mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
      error: new Error('Failed to load events'),
    } as ReturnType<typeof useArtistEvents>);

    render(<ArtistEventsSection artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
  });
});
