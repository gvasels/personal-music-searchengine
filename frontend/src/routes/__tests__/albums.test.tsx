/**
 * Albums Page Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseAlbums = vi.fn();
vi.mock('../../hooks/useAlbums', () => ({
  useAlbumsQuery: () => mockUseAlbums(),
  albumKeys: { lists: () => ['albums', 'list'] },
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useSearch: () => ({ sortBy: undefined, sortOrder: undefined, page: 1 }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
}));

import AlbumsPage from '../albums/index';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('AlbumsPage (Wave 2)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockAlbumsData = {
    items: [
      { id: 'album-1', name: 'Album 1', artist: 'Artist A', year: 2024, trackCount: 10 },
      { id: 'album-2', name: 'Album 2', artist: 'Artist B', year: 2023, trackCount: 12 },
    ],
    total: 2,
    limit: 24,
    offset: 0,
  };

  describe('Loading state', () => {
    it('should show loading skeleton while fetching', () => {
      mockUseAlbums.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<AlbumsPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message on failure', () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Failed'),
        refetch: vi.fn(),
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
    });

    it('should show retry button on error', async () => {
      const refetch = vi.fn();
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Error'),
        refetch,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /retry/i }));

      expect(refetch).toHaveBeenCalled();
    });
  });

  describe('Empty state', () => {
    it('should show empty state when no albums exist', () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: { items: [], total: 0, limit: 24, offset: 0 },
        isError: false,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no albums found/i)).toBeInTheDocument();
    });
  });

  describe('Album grid rendering', () => {
    it('should render all albums', () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: mockAlbumsData,
        isError: false,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Album 1')).toBeInTheDocument();
      expect(screen.getByText('Album 2')).toBeInTheDocument();
    });

    it('should display artist names', () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: mockAlbumsData,
        isError: false,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Artist A')).toBeInTheDocument();
      expect(screen.getByText('Artist B')).toBeInTheDocument();
    });

    it('should display years', () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: mockAlbumsData,
        isError: false,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });

      expect(screen.getByText('2024')).toBeInTheDocument();
      expect(screen.getByText('2023')).toBeInTheDocument();
    });

    it('should display total album count', () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: mockAlbumsData,
        isError: false,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/2 albums/i)).toBeInTheDocument();
    });
  });

  describe('Sorting', () => {
    it('should render sort dropdown', () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: mockAlbumsData,
        isError: false,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('combobox', { name: /sort/i })).toBeInTheDocument();
    });

    it('should update URL when sort changes', async () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: mockAlbumsData,
        isError: false,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });
      await user.selectOptions(screen.getByRole('combobox', { name: /sort/i }), 'year');

      expect(mockNavigate).toHaveBeenCalled();
    });
  });

  describe('Navigation', () => {
    it('should navigate to album detail when clicking album', async () => {
      mockUseAlbums.mockReturnValue({
        isLoading: false,
        data: mockAlbumsData,
        isError: false,
      });

      render(<AlbumsPage />, { wrapper: createWrapper() });
      await user.click(screen.getByText('Album 1'));

      expect(mockNavigate).toHaveBeenCalledWith({
        to: '/albums/$albumId',
        params: { albumId: 'album-1' },
      });
    });
  });
});
