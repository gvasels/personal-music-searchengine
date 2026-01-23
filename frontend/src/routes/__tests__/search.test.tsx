/**
 * Search Page Tests - Wave 4
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseSearch = vi.fn();
vi.mock('../../hooks/useSearch', () => ({
  useSearchQuery: () => mockUseSearch(),
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  searchKeys: { results: (p: any) => ['search', 'results', p] },
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useSearch: () => ({ q: 'test query', artist: '', album: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
}));

import SearchPage from '../search';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('SearchPage (Wave 4)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockSearchResults = {
    items: [
      { id: 'track-1', title: 'Test Track', artist: 'Test Artist', album: 'Test Album', duration: 180 },
    ],
    total: 1,
    limit: 20,
    offset: 0,
  };

  describe('Loading state', () => {
    it('should show loading spinner while searching', () => {
      mockUseSearch.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<SearchPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message on search failure', () => {
      mockUseSearch.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Search failed'),
      });

      render(<SearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/search failed/i)).toBeInTheDocument();
    });
  });

  describe('Empty state', () => {
    it('should show no results message', () => {
      mockUseSearch.mockReturnValue({
        isLoading: false,
        data: { items: [], total: 0, limit: 20, offset: 0 },
        isError: false,
      });

      render(<SearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no results/i)).toBeInTheDocument();
    });
  });

  describe('Results display', () => {
    it('should show search query', () => {
      mockUseSearch.mockReturnValue({
        isLoading: false,
        data: mockSearchResults,
        isError: false,
      });

      render(<SearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/test query/i)).toBeInTheDocument();
    });

    it('should show result count', () => {
      mockUseSearch.mockReturnValue({
        isLoading: false,
        data: mockSearchResults,
        isError: false,
      });

      render(<SearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/1 result/i)).toBeInTheDocument();
    });

    it('should render search results', () => {
      mockUseSearch.mockReturnValue({
        isLoading: false,
        data: mockSearchResults,
        isError: false,
      });

      render(<SearchPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Test Track')).toBeInTheDocument();
      expect(screen.getByText('Test Artist')).toBeInTheDocument();
    });

    it('should navigate to track when clicking result', async () => {
      mockUseSearch.mockReturnValue({
        isLoading: false,
        data: mockSearchResults,
        isError: false,
      });

      render(<SearchPage />, { wrapper: createWrapper() });
      await user.click(screen.getByText('Test Track'));

      expect(mockNavigate).toHaveBeenCalledWith({
        to: '/tracks/$trackId',
        params: { trackId: 'track-1' },
      });
    });
  });

  describe('Filters', () => {
    it('should render filter inputs', () => {
      mockUseSearch.mockReturnValue({
        isLoading: false,
        data: mockSearchResults,
        isError: false,
      });

      render(<SearchPage />, { wrapper: createWrapper() });

      expect(screen.getByLabelText(/artist/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/album/i)).toBeInTheDocument();
    });

    it('should update filters on input', async () => {
      mockUseSearch.mockReturnValue({
        isLoading: false,
        data: mockSearchResults,
        isError: false,
      });

      render(<SearchPage />, { wrapper: createWrapper() });
      await user.type(screen.getByLabelText(/artist/i), 'Beatles');

      // Wait for debounce (300ms)
      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalled();
      }, { timeout: 500 });
    });
  });
});
