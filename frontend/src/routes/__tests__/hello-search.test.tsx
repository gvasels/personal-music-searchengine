/**
 * Hello Search Page Tests - TDD Red Phase
 *
 * Tests for the hello-search.tsx route page.
 * These tests MUST FAIL because the route doesn't exist yet.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mock the useHelloSearch hook - CONTRACT CRITICAL: response uses "items" not "tracks"
const mockUseHelloSearch = vi.fn();
const mockUseHelloFeatured = vi.fn();
vi.mock('../../hooks/useHelloSearch', () => ({
  useHelloSearch: (q: string) => mockUseHelloSearch(q),
  useHelloFeatured: () => mockUseHelloFeatured(),
  helloKeys: {
    all: ['hello'] as const,
    search: (q: string) => ['hello', 'search', q] as const,
    featured: () => ['hello', 'featured'] as const,
  },
}));

// Mock TanStack Router
const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  createFileRoute: () => () => ({}),
  useNavigate: () => mockNavigate,
  useSearch: () => ({ q: '' }),
  useLocation: () => ({ pathname: '/hello-search', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

// Import after mocks
import HelloSearchPage from '../hello-search';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

// Mock data - CONTRACT CRITICAL: uses "items" not "tracks"
const mockFeaturedData = {
  items: [
    {
      id: 'seed-t1',
      title: 'Aurora Borealis',
      artist: 'Aurora Waves',
      album: 'Dreamscape',
      genre: 'jazz',
      year: 2023,
      duration: 245,
    },
    {
      id: 'seed-t2',
      title: 'Midnight Jazz',
      artist: 'Jazz Collective',
      album: 'Night Sessions',
      genre: 'jazz',
      year: 2022,
      duration: 312,
    },
  ],
  total: 2,
};

const mockSearchData = {
  items: [
    {
      id: 'seed-t1',
      title: 'Aurora Borealis',
      artist: 'Aurora Waves',
      album: 'Dreamscape',
      genre: 'jazz',
      year: 2023,
      duration: 245,
    },
  ],
  total: 1,
};

const mockEmptyData = {
  items: [],
  total: 0,
};

describe('HelloSearchPage', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Initial render with featured tracks', () => {
    it('should render search input on load', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      expect(
        screen.getByPlaceholderText(/search.*tracks/i)
      ).toBeInTheDocument();
    });

    it('should render featured tracks on load', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Aurora Borealis')).toBeInTheDocument();
      expect(screen.getByText('Midnight Jazz')).toBeInTheDocument();
    });
  });

  describe('Loading state', () => {
    it('should show loading skeletons during fetch (.skeleton class)', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: true,
        data: undefined,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      // Check for skeleton elements with .skeleton class
      const skeletons = document.querySelectorAll('.skeleton');
      expect(skeletons.length).toBeGreaterThan(0);
    });

    it('should show loading skeletons during search', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: true,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      const skeletons = document.querySelectorAll('.skeleton');
      expect(skeletons.length).toBeGreaterThan(0);
    });
  });

  describe('Error state', () => {
    it('should show error message on API failure (role="alert")', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Failed to fetch featured tracks'),
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      const alert = screen.getByRole('alert');
      expect(alert).toBeInTheDocument();
      expect(alert).toHaveTextContent(/failed|error/i);
    });

    it('should show error message on search failure', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Search failed'),
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      const alert = screen.getByRole('alert');
      expect(alert).toBeInTheDocument();
      expect(alert).toHaveTextContent(/failed|error/i);
    });
  });

  describe('Empty results state', () => {
    it('should show "No results found" for empty search results', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: mockEmptyData,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no results found/i)).toBeInTheDocument();
    });
  });

  describe('Search functionality', () => {
    it('should trigger search when typing (calls useHelloSearch with query)', async () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      const searchInput = screen.getByPlaceholderText(/search.*tracks/i);
      await user.type(searchInput, 'jazz');

      // Verify useHelloSearch was called with the query
      await waitFor(() => {
        expect(mockUseHelloSearch).toHaveBeenCalledWith('jazz');
      });
    });

    it('should display search results when query matches', async () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: mockSearchData,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      // With search results, should show the matching track
      expect(screen.getByText('Aurora Borealis')).toBeInTheDocument();
      expect(screen.getByText('Aurora Waves')).toBeInTheDocument();
    });
  });

  describe('Track card display', () => {
    it('should display track cards from featured endpoint', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      // Verify track card content from featured tracks
      expect(screen.getByText('Aurora Borealis')).toBeInTheDocument();
      expect(screen.getByText('Aurora Waves')).toBeInTheDocument();
      expect(screen.getByText('Dreamscape')).toBeInTheDocument();

      expect(screen.getByText('Midnight Jazz')).toBeInTheDocument();
      expect(screen.getByText('Jazz Collective')).toBeInTheDocument();
      expect(screen.getByText('Night Sessions')).toBeInTheDocument();
    });

    it('should display track duration in MM:SS format', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      // 245 seconds = 4:05, 312 seconds = 5:12
      expect(screen.getByText('4:05')).toBeInTheDocument();
      expect(screen.getByText('5:12')).toBeInTheDocument();
    });

    it('should display track genre', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      // Both tracks have jazz genre
      const jazzElements = screen.getAllByText(/jazz/i);
      expect(jazzElements.length).toBeGreaterThanOrEqual(1);
    });

    it('should display track year', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: mockFeaturedData,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText('2023')).toBeInTheDocument();
      expect(screen.getByText('2022')).toBeInTheDocument();
    });
  });
});
