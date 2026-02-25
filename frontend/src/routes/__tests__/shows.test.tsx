/**
 * My Shows Page Tests - TDD Red Phase
 *
 * These tests define the contract for the /shows route.
 * They MUST FAIL because shows.tsx does not exist yet.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseWatchedArtists = vi.fn();
const mockUseSearchArtistEvents = vi.fn();

vi.mock('../../hooks/useArtistWatch', () => ({
  useWatchedArtists: () => mockUseWatchedArtists(),
  watchKeys: { all: ['artist-watch'], list: () => ['artist-watch', 'list'] },
}));

vi.mock('../../hooks/useArtistEvents', () => ({
  useSearchArtistEvents: (q: string) => mockUseSearchArtistEvents(q),
  useArtistEvents: vi.fn(),
  eventKeys: { all: ['events'] },
}));

vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => vi.fn(),
  useSearch: () => ({}),
  useLocation: () => ({ pathname: '/shows', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
  createFileRoute: () => () => ({ component: undefined }),
}));

// This import will fail - shows.tsx does not exist yet
import ShowsPage from '../shows';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('ShowsPage', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
    mockUseSearchArtistEvents.mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: false,
    });
  });

  describe('Loading state', () => {
    it('should show loading spinner while fetching', () => {
      mockUseWatchedArtists.mockReturnValue({
        isLoading: true,
        data: undefined,
        isError: false,
      });

      render(<ShowsPage />, { wrapper: createWrapper() });
      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Empty state', () => {
    it('should show empty state when no watched artists', () => {
      mockUseWatchedArtists.mockReturnValue({
        isLoading: false,
        data: { items: [], hasMore: false },
        isError: false,
      });

      render(<ShowsPage />, { wrapper: createWrapper() });
      expect(screen.getByText(/no watched artists/i)).toBeInTheDocument();
    });
  });

  describe('Data rendering', () => {
    it('should render watched artist names', () => {
      mockUseWatchedArtists.mockReturnValue({
        isLoading: false,
        data: {
          items: [
            { artistName: 'Daft Punk', watchedAt: '2026-02-25T10:00:00Z' },
            { artistName: 'Aphex Twin', watchedAt: '2026-02-24T08:00:00Z' },
          ],
          hasMore: false,
        },
        isError: false,
      });

      render(<ShowsPage />, { wrapper: createWrapper() });
      expect(screen.getByText('Daft Punk')).toBeInTheDocument();
      expect(screen.getByText('Aphex Twin')).toBeInTheDocument();
    });
  });

  describe('Search', () => {
    it('should have a search input', () => {
      mockUseWatchedArtists.mockReturnValue({
        isLoading: false,
        data: { items: [], hasMore: false },
        isError: false,
      });

      render(<ShowsPage />, { wrapper: createWrapper() });
      expect(screen.getByRole('searchbox')).toBeInTheDocument();
    });

    it('should show search results when query entered', async () => {
      mockUseWatchedArtists.mockReturnValue({
        isLoading: false,
        data: { items: [], hasMore: false },
        isError: false,
      });
      mockUseSearchArtistEvents.mockReturnValue({
        data: {
          items: [
            { name: 'Daft Punk', upcomingEvents: 3, source: 'mock' },
          ],
        },
        isLoading: false,
        isError: false,
      });

      render(<ShowsPage />, { wrapper: createWrapper() });

      const searchInput = screen.getByRole('searchbox');
      await user.type(searchInput, 'Daft');

      expect(mockUseSearchArtistEvents).toHaveBeenCalledWith('Daft');
    });
  });

  describe('Error state', () => {
    it('should show error state on failure', () => {
      mockUseWatchedArtists.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Failed'),
      });

      render(<ShowsPage />, { wrapper: createWrapper() });
      expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
    });
  });
});
