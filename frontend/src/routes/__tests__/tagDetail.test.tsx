/**
 * Tag Detail Page Tests - Wave 5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseTracksByTagQuery = vi.fn();
vi.mock('../../hooks/useTags', () => ({
  useTracksByTagQuery: () => mockUseTracksByTagQuery(),
  tagKeys: { all: ['tags'], tracks: (name: string) => ['tags', 'tracks', name] },
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useParams: () => ({ tagName: 'rock' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

// Mock player store
const mockSetQueue = vi.fn();
vi.mock('../../lib/store/playerStore', () => ({
  usePlayerStore: () => ({
    setQueue: mockSetQueue,
    currentTrack: null,
    isPlaying: false,
  }),
}));

import TagDetailPage from '../tags/$tagName';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('TagDetailPage (Wave 5)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockTracks = {
    items: [
      {
        id: 'track-1',
        title: 'Rock Anthem',
        artist: 'Rock Band',
        album: 'Greatest Hits',
        duration: 240,
        tags: ['rock', 'classic'],
      },
      {
        id: 'track-2',
        title: 'Rock Ballad',
        artist: 'Rock Band',
        album: 'Greatest Hits',
        duration: 300,
        tags: ['rock', 'ballad'],
      },
    ],
    total: 2,
    limit: 20,
    offset: 0,
  };

  describe('Loading state', () => {
    it('should show loading spinner while fetching', () => {
      mockUseTracksByTagQuery.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<TagDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message on failure', () => {
      mockUseTracksByTagQuery.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Failed to load tracks'),
      });

      render(<TagDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/failed to load tracks/i)).toBeInTheDocument();
    });
  });

  describe('Empty state', () => {
    it('should show empty message when no tracks', () => {
      mockUseTracksByTagQuery.mockReturnValue({
        isLoading: false,
        data: { items: [], total: 0, limit: 20, offset: 0 },
        isError: false,
      });

      render(<TagDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no tracks/i)).toBeInTheDocument();
    });
  });

  describe('Tag detail display', () => {
    it('should render tag name as title', () => {
      mockUseTracksByTagQuery.mockReturnValue({
        isLoading: false,
        data: mockTracks,
        isError: false,
      });

      render(<TagDetailPage />, { wrapper: createWrapper() });

      // The tag name appears as heading text
      expect(screen.getByRole('heading', { name: /rock/i })).toBeInTheDocument();
    });

    it('should render track count', () => {
      mockUseTracksByTagQuery.mockReturnValue({
        isLoading: false,
        data: mockTracks,
        isError: false,
      });

      render(<TagDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/2 tracks/i)).toBeInTheDocument();
    });

    it('should render tracks list', () => {
      mockUseTracksByTagQuery.mockReturnValue({
        isLoading: false,
        data: mockTracks,
        isError: false,
      });

      render(<TagDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Rock Anthem')).toBeInTheDocument();
      expect(screen.getByText('Rock Ballad')).toBeInTheDocument();
    });
  });

  describe('Playback', () => {
    it('should have play all button', () => {
      mockUseTracksByTagQuery.mockReturnValue({
        isLoading: false,
        data: mockTracks,
        isError: false,
      });

      render(<TagDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /play all/i })).toBeInTheDocument();
    });

    it('should set queue when play all clicked', async () => {
      mockUseTracksByTagQuery.mockReturnValue({
        isLoading: false,
        data: mockTracks,
        isError: false,
      });

      render(<TagDetailPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /play all/i }));

      expect(mockSetQueue).toHaveBeenCalledWith(mockTracks.items, 0);
    });
  });

  describe('Navigation', () => {
    it('should have back button to tags', () => {
      mockUseTracksByTagQuery.mockReturnValue({
        isLoading: false,
        data: mockTracks,
        isError: false,
      });

      render(<TagDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('link', { name: /back/i })).toBeInTheDocument();
    });
  });
});
