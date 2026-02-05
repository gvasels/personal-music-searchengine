/**
 * HelloSearchPage Tests - TDD Red Phase
 *
 * Tests for the /hello-search route page component.
 * These tests define the contract for the HelloSearchPage before implementation.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mock hooks
const mockUseHelloSearch = vi.fn();
const mockUseHelloFeatured = vi.fn();
vi.mock('../../hooks/useHelloSearch', () => ({
  useHelloSearch: (q: string) => mockUseHelloSearch(q),
  useHelloFeatured: () => mockUseHelloFeatured(),
  helloKeys: {
    all: ['hello'],
    search: (q: string) => ['hello', 'search', q],
    featured: () => ['hello', 'featured'],
  },
}));

// Mock TanStack Router
vi.mock('@tanstack/react-router', async () => ({
  createFileRoute: () => () => ({ component: undefined }),
  useNavigate: () => vi.fn(),
  Link: ({
    children,
    to,
  }: {
    children: React.ReactNode;
    to: string;
  }) => <a href={to}>{children}</a>,
}));

// Import the page component (Red phase: this module does not exist yet)
import HelloSearchPage from '../hello-search';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

// Mock data matching the API contract
const mockFeaturedData = {
  items: [
    {
      id: '1',
      title: 'Test Track',
      artist: 'Test Artist',
      album: 'Test Album',
      genre: 'jazz',
      year: 2023,
      duration: 240,
    },
    {
      id: '2',
      title: 'Another Track',
      artist: 'Another Artist',
      album: 'Another Album',
      genre: 'rock',
      year: 2022,
      duration: 180,
    },
  ],
  total: 2,
};

const _mockSearchData = {
  items: [
    {
      id: '3',
      title: 'Search Result Track',
      artist: 'Found Artist',
      album: 'Found Album',
      genre: 'electronic',
      year: 2024,
      duration: 300,
    },
  ],
  total: 1,
};

const mockEmptySearchData = {
  items: [],
  total: 0,
};

describe('HelloSearchPage', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('renders search input and featured tracks on load', () => {
    it('should render a search input and featured track cards', () => {
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

      // Should render a search input
      expect(
        screen.getByPlaceholderText(/search/i),
      ).toBeInTheDocument();

      // Should render featured track titles
      expect(screen.getByText('Test Track')).toBeInTheDocument();
      expect(screen.getByText('Another Track')).toBeInTheDocument();
    });
  });

  describe('shows loading skeletons during fetch', () => {
    it('should display skeleton elements when featured data is loading', () => {
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

      // Should render skeleton loading elements
      const skeletons = document.querySelectorAll('.skeleton');
      expect(skeletons.length).toBeGreaterThan(0);
    });
  });

  describe('shows error message on API failure', () => {
    it('should display an error alert when the API fails', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('API Error'),
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('alert')).toBeInTheDocument();
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
  });

  describe('shows "No results found" for empty search results', () => {
    it('should display a no-results message when search returns empty items', () => {
      mockUseHelloFeatured.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: false,
      });
      mockUseHelloSearch.mockReturnValue({
        isLoading: false,
        data: mockEmptySearchData,
        isError: false,
      });

      render(<HelloSearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no results found/i)).toBeInTheDocument();
    });
  });

  describe('search input triggers search', () => {
    it('should call useHelloSearch with the typed query', async () => {
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

      const searchInput = screen.getByPlaceholderText(/search/i);
      await user.type(searchInput, 'jazz');

      await waitFor(() => {
        expect(mockUseHelloSearch).toHaveBeenCalledWith('jazz');
      });
    });
  });

  describe('displays track cards from featured endpoint', () => {
    it('should render track cards with title and artist from featured data', () => {
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

      // Verify first track card
      expect(screen.getByText('Test Track')).toBeInTheDocument();
      expect(screen.getByText('Test Artist')).toBeInTheDocument();

      // Verify second track card
      expect(screen.getByText('Another Track')).toBeInTheDocument();
      expect(screen.getByText('Another Artist')).toBeInTheDocument();
    });

    it('should render the correct number of track cards', () => {
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

      // There should be 2 track cards matching the mock data
      const trackTitles = ['Test Track', 'Another Track'];
      trackTitles.forEach((title) => {
        expect(screen.getByText(title)).toBeInTheDocument();
      });
    });
  });
});
