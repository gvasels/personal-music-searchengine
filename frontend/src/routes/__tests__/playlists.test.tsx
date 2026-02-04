/**
 * Playlists Page Tests - Wave 5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUsePlaylistsQuery = vi.fn();
vi.mock('../../hooks/usePlaylists', () => ({
  usePlaylistsQuery: () => mockUsePlaylistsQuery(),
  playlistKeys: { all: ['playlists'] },
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useLocation: () => ({ pathname: '/playlists', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

import PlaylistsPage from '../playlists/index';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('PlaylistsPage (Wave 5)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockPlaylists = {
    items: [
      {
        id: 'playlist-1',
        name: 'Favorites',
        description: 'My favorite songs',
        trackIds: ['track-1', 'track-2'],
        trackCount: 2,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-02T00:00:00Z',
      },
      {
        id: 'playlist-2',
        name: 'Chill Vibes',
        description: '',
        trackIds: [],
        trackCount: 0,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      },
    ],
    total: 2,
    limit: 20,
    offset: 0,
  };

  describe('Loading state', () => {
    it('should show loading spinner while fetching', () => {
      mockUsePlaylistsQuery.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<PlaylistsPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message on failure', () => {
      mockUsePlaylistsQuery.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Failed to load playlists'),
      });

      render(<PlaylistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/failed to load playlists/i)).toBeInTheDocument();
    });
  });

  describe('Empty state', () => {
    it('should show empty message when no playlists', () => {
      mockUsePlaylistsQuery.mockReturnValue({
        isLoading: false,
        data: { items: [], total: 0, limit: 20, offset: 0 },
        isError: false,
      });

      render(<PlaylistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no playlists/i)).toBeInTheDocument();
    });
  });

  describe('Playlists display', () => {
    it('should render page title', () => {
      mockUsePlaylistsQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylists,
        isError: false,
      });

      render(<PlaylistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/playlists/i)).toBeInTheDocument();
    });

    it('should render playlist cards', () => {
      mockUsePlaylistsQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylists,
        isError: false,
      });

      render(<PlaylistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Favorites')).toBeInTheDocument();
      expect(screen.getByText('Chill Vibes')).toBeInTheDocument();
    });

    it('should show track count', () => {
      mockUsePlaylistsQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylists,
        isError: false,
      });

      render(<PlaylistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/2 tracks/i)).toBeInTheDocument();
    });

    it('should navigate to playlist on click', async () => {
      mockUsePlaylistsQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylists,
        isError: false,
      });

      render(<PlaylistsPage />, { wrapper: createWrapper() });
      await user.click(screen.getByText('Favorites'));

      expect(mockNavigate).toHaveBeenCalledWith({
        to: '/playlists/$playlistId',
        params: { playlistId: 'playlist-1' },
      });
    });
  });

  describe('Create playlist', () => {
    it('should render create playlist button', () => {
      mockUsePlaylistsQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylists,
        isError: false,
      });

      render(<PlaylistsPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /create playlist/i })).toBeInTheDocument();
    });

    it('should open create modal when button clicked', async () => {
      mockUsePlaylistsQuery.mockReturnValue({
        isLoading: false,
        data: mockPlaylists,
        isError: false,
      });

      render(<PlaylistsPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /create playlist/i }));

      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });
});
