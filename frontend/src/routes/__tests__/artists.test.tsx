/**
 * Artists Page Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseArtists = vi.fn();
vi.mock('../../hooks/useArtists', () => ({
  useArtistsQuery: () => mockUseArtists(),
  artistKeys: { lists: () => ['artists', 'list'] },
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useSearch: () => ({ sortBy: undefined, sortOrder: undefined, page: 1, search: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
}));

import ArtistsPage from '../artists/index';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('ArtistsPage (Wave 2)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockArtistsData = {
    items: [
      { name: 'Artist A', trackCount: 25, albumCount: 3 },
      { name: 'Artist B', trackCount: 15, albumCount: 2 },
    ],
    total: 2,
    limit: 20,
    offset: 0,
  };

  describe('Loading state', () => {
    it('should show loading spinner while fetching', () => {
      mockUseArtists.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<ArtistsPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message on failure', () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Failed'),
        refetch: vi.fn(),
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
    });

    it('should show retry button on error', async () => {
      const refetch = vi.fn();
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Error'),
        refetch,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /retry/i }));

      expect(refetch).toHaveBeenCalled();
    });
  });

  describe('Empty state', () => {
    it('should show empty state when no artists exist', () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: { items: [], total: 0, limit: 20, offset: 0 },
        isError: false,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no artists found/i)).toBeInTheDocument();
    });
  });

  describe('Artist list rendering', () => {
    it('should render all artists', () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: mockArtistsData,
        isError: false,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Artist A')).toBeInTheDocument();
      expect(screen.getByText('Artist B')).toBeInTheDocument();
    });

    it('should display track counts', () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: mockArtistsData,
        isError: false,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/25 tracks/i)).toBeInTheDocument();
      expect(screen.getByText(/15 tracks/i)).toBeInTheDocument();
    });

    it('should display album counts', () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: mockArtistsData,
        isError: false,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/3 albums/i)).toBeInTheDocument();
      expect(screen.getByText(/2 albums/i)).toBeInTheDocument();
    });

    it('should display total artist count', () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: mockArtistsData,
        isError: false,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/2 artists/i)).toBeInTheDocument();
    });
  });

  describe('Sorting', () => {
    it('should render sort dropdown', () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: mockArtistsData,
        isError: false,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('combobox', { name: /sort/i })).toBeInTheDocument();
    });

    it('should update URL when sort changes', async () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: mockArtistsData,
        isError: false,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });
      await user.selectOptions(screen.getByRole('combobox', { name: /sort/i }), 'trackCount');

      expect(mockNavigate).toHaveBeenCalled();
    });
  });

  describe('Navigation', () => {
    it('should navigate to artist detail when clicking artist', async () => {
      mockUseArtists.mockReturnValue({
        isLoading: false,
        data: mockArtistsData,
        isError: false,
      });

      render(<ArtistsPage />, { wrapper: createWrapper() });
      await user.click(screen.getByText('Artist A'));

      expect(mockNavigate).toHaveBeenCalledWith({
        to: '/artists/$artistName',
        params: { artistName: 'Artist A' },
      });
    });
  });
});
