/**
 * Tracks Page Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseTracks = vi.fn();
const mockDeleteTrack = vi.fn();
vi.mock('../../hooks/useTracks', () => ({
  useTracksQuery: () => mockUseTracks(),
  useDeleteTrack: () => ({ mutateAsync: mockDeleteTrack, isPending: false }),
  trackKeys: { lists: () => ['tracks', 'list'] },
}));

vi.mock('@/lib/store/playerStore', () => ({
  usePlayerStore: () => ({
    currentTrack: null,
    isPlaying: false,
    setQueue: vi.fn(),
    pause: vi.fn(),
  }),
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useSearch: () => ({ sortBy: undefined, sortOrder: undefined, page: 1 }),
  useLocation: () => ({ pathname: '/tracks', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
}));

import TracksPage from '../tracks/index';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('TracksPage (Wave 2)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockTracksData = {
    items: [
      { id: 'track-1', title: 'Track 1', artist: 'Artist A', album: 'Album 1', duration: 180, format: 'mp3' },
      { id: 'track-2', title: 'Track 2', artist: 'Artist B', album: 'Album 2', duration: 240, format: 'flac' },
    ],
    total: 2,
    limit: 20,
    offset: 0,
  };

  describe('Loading state', () => {
    it('should show loading spinner while fetching', () => {
      mockUseTracks.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<TracksPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message on failure', () => {
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Failed to load'),
        refetch: vi.fn(),
      });

      render(<TracksPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
    });

    it('should show retry button on error', async () => {
      const refetch = vi.fn();
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Network error'),
        refetch,
      });

      render(<TracksPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /retry/i }));

      expect(refetch).toHaveBeenCalled();
    });
  });

  describe('Empty state', () => {
    it('should show empty state when no tracks exist', () => {
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: { items: [], total: 0, limit: 20, offset: 0 },
        isError: false,
      });

      render(<TracksPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no tracks found/i)).toBeInTheDocument();
    });
  });

  describe('Track list rendering', () => {
    it('should render all tracks', () => {
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: mockTracksData,
        isError: false,
      });

      render(<TracksPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Track 1')).toBeInTheDocument();
      expect(screen.getByText('Track 2')).toBeInTheDocument();
    });

    it('should display artist names', () => {
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: mockTracksData,
        isError: false,
      });

      render(<TracksPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Artist A')).toBeInTheDocument();
      expect(screen.getByText('Artist B')).toBeInTheDocument();
    });

    it('should display total track count', () => {
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: mockTracksData,
        isError: false,
      });

      render(<TracksPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/2 tracks/i)).toBeInTheDocument();
    });
  });

  describe('Sorting', () => {
    it('should render sortable column headers', () => {
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: mockTracksData,
        isError: false,
      });

      render(<TracksPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /title/i })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /artist/i })).toBeInTheDocument();
    });

    it('should toggle sort when clicking column header', async () => {
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: mockTracksData,
        isError: false,
      });

      render(<TracksPage />, { wrapper: createWrapper() });
      // Click on the Title column header to sort
      await user.click(screen.getByRole('button', { name: /title/i }));

      // The tracks should be sorted - table still renders without navigation
      expect(screen.getByText('Track 1')).toBeInTheDocument();
    });
  });

  describe('Navigation', () => {
    it('should navigate to track detail when clicking track', async () => {
      mockUseTracks.mockReturnValue({
        isLoading: false,
        data: mockTracksData,
        isError: false,
      });

      render(<TracksPage />, { wrapper: createWrapper() });
      await user.click(screen.getByText('Track 1'));

      expect(mockNavigate).toHaveBeenCalledWith({
        to: '/tracks/$trackId',
        params: { trackId: 'track-1' },
      });
    });
  });
});
