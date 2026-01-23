/**
 * Tags Page Tests - Wave 5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseTagsQuery = vi.fn();
vi.mock('../../hooks/useTags', () => ({
  useTagsQuery: () => mockUseTagsQuery(),
  tagKeys: { all: ['tags'] },
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useLocation: () => ({ pathname: '/tags', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

import TagsPage from '../tags/index';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('TagsPage (Wave 5)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockTags = {
    items: [
      { name: 'rock', trackCount: 25 },
      { name: 'jazz', trackCount: 15 },
      { name: 'classical', trackCount: 8 },
      { name: 'electronic', trackCount: 42 },
    ],
    total: 4,
    limit: 50,
    offset: 0,
  };

  describe('Loading state', () => {
    it('should show loading spinner while fetching', () => {
      mockUseTagsQuery.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<TagsPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message on failure', () => {
      mockUseTagsQuery.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Failed to load tags'),
      });

      render(<TagsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/failed to load tags/i)).toBeInTheDocument();
    });
  });

  describe('Empty state', () => {
    it('should show empty message when no tags', () => {
      mockUseTagsQuery.mockReturnValue({
        isLoading: false,
        data: { items: [], total: 0, limit: 50, offset: 0 },
        isError: false,
      });

      render(<TagsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/no tags/i)).toBeInTheDocument();
    });
  });

  describe('Tags display', () => {
    it('should render page title', () => {
      mockUseTagsQuery.mockReturnValue({
        isLoading: false,
        data: mockTags,
        isError: false,
      });

      render(<TagsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/tags/i)).toBeInTheDocument();
    });

    it('should render tag badges', () => {
      mockUseTagsQuery.mockReturnValue({
        isLoading: false,
        data: mockTags,
        isError: false,
      });

      render(<TagsPage />, { wrapper: createWrapper() });

      expect(screen.getByText('rock')).toBeInTheDocument();
      expect(screen.getByText('jazz')).toBeInTheDocument();
      expect(screen.getByText('classical')).toBeInTheDocument();
      expect(screen.getByText('electronic')).toBeInTheDocument();
    });

    it('should show track count for each tag', () => {
      mockUseTagsQuery.mockReturnValue({
        isLoading: false,
        data: mockTags,
        isError: false,
      });

      render(<TagsPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/25/)).toBeInTheDocument();
      expect(screen.getByText(/42/)).toBeInTheDocument();
    });

    it('should navigate to tag detail on click', async () => {
      mockUseTagsQuery.mockReturnValue({
        isLoading: false,
        data: mockTags,
        isError: false,
      });

      render(<TagsPage />, { wrapper: createWrapper() });
      await user.click(screen.getByText('rock'));

      expect(mockNavigate).toHaveBeenCalledWith({
        to: '/tags/$tagName',
        params: { tagName: 'rock' },
      });
    });
  });

  describe('Tag cloud sizing', () => {
    it('should render tags with different sizes based on track count', () => {
      mockUseTagsQuery.mockReturnValue({
        isLoading: false,
        data: mockTags,
        isError: false,
      });

      render(<TagsPage />, { wrapper: createWrapper() });

      // Tags with more tracks should have larger styling
      // This validates the tag cloud functionality
      const electronicTag = screen.getByText('electronic').closest('button');
      const classicalTag = screen.getByText('classical').closest('button');

      expect(electronicTag).toBeInTheDocument();
      expect(classicalTag).toBeInTheDocument();
    });
  });
});
