/**
 * Playlist Detail Page Tests - Wave 5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUsePlaylistQuery = vi.fn();
vi.mock('../../hooks/usePlaylists', () => ({
  usePlaylistQuery: () => mockUsePlaylistQuery(),
  playlistKeys: { all: ['playlists'], detail: (id: string) => ['playlists', 'detail', id] },
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useParams: () => ({ playlistId: 'playlist-1' }),
  useLocation: () => ({ pathname: '/playlists/playlist-1', search: '', hash: '' }),
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

import PlaylistDetailPage from '../playlists/$playlistId';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('PlaylistDetailPage (Wave 5)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockPlaylist = {
    id: 'playlist-1',
    name: 'Favorites',
    description: 'My favorite songs',
    trackIds: ['track-1', 'track-2'],
    trackCount: 2,
    tracks: [
      {
        id: 'track-1',
        title: 'Song One',
        artist: 'Artist A',
        album: 'Album 1',
        duration: 180,
      },
      {
        id: 'track-2',
        title: 'Song Two',
        artist: 'Artist B',
        album: 'Album 2',
        duration: 240,
      },
    ],
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-02T00:00:00Z',
  };

  describe('Loading state', () => {
    it('should show loading spinner while fetching', () => {
      mockUsePlaylistQuery.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message on failure', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Playlist not found'),
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/playlist not found/i)).toBeInTheDocument();
    });
  });

  describe('Playlist display', () => {
    it('should render playlist name', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Favorites')).toBeInTheDocument();
    });

    it('should render playlist description', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText('My favorite songs')).toBeInTheDocument();
    });

    it('should render track count', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/2 tracks/i)).toBeInTheDocument();
    });

    it('should render tracks list', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Song One')).toBeInTheDocument();
      expect(screen.getByText('Song Two')).toBeInTheDocument();
    });
  });

  describe('Playback', () => {
    it('should have play all button', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /play all/i })).toBeInTheDocument();
    });

    it('should set queue when play all clicked', async () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /play all/i }));

      expect(mockSetQueue).toHaveBeenCalledWith(mockPlaylist.tracks, 0);
    });
  });

  describe('Navigation', () => {
    it('should have back button', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('link', { name: /back/i })).toBeInTheDocument();
    });
  });

  describe('Edit playlist', () => {
    it('should have edit button', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /edit/i })).toBeInTheDocument();
    });
  });

  describe('Delete playlist', () => {
    it('should have delete button', () => {
      mockUsePlaylistQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylist,
        isError: false,
      });

      render(<PlaylistDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument();
    });
  });
});
