/**
 * Hello Search Page Tests - TDD Red Phase
 *
 * Tests for the hello search route page component.
 * Source module (hello-search.tsx) does NOT exist yet - these tests MUST fail.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';

const mockUseHelloSearch = vi.fn();
const mockUseFeaturedTracks = vi.fn();

vi.mock('../../hooks/useHelloSearch', () => ({
  useHelloSearch: (...args: unknown[]) => mockUseHelloSearch(...args),
  useFeaturedTracks: (...args: unknown[]) => mockUseFeaturedTracks(...args),
}));

vi.mock('@tanstack/react-router', () => ({
  createFileRoute: () => (opts: Record<string, unknown>) => opts,
  Link: ({ children, ...props }: { children: React.ReactNode; [key: string]: unknown }) => (
    <a {...props}>{children}</a>
  ),
}));

import HelloSearchPage from '../hello-search';

describe('HelloSearchPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseHelloSearch.mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: false,
      error: null,
    });
    mockUseFeaturedTracks.mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: false,
      error: null,
    });
  });

  it('renders hero section with title text', () => {
    render(<HelloSearchPage />);

    expect(screen.getByText(/hello/i)).toBeInTheDocument();
  });

  it('shows featured tracks when loaded', () => {
    mockUseFeaturedTracks.mockReturnValue({
      data: {
        items: [
          {
            id: 'feat-1',
            title: 'Featured Hit',
            artist: 'Star Artist',
            album: 'Top Album',
            genre: 'Pop',
            year: 2024,
            duration: 200,
          },
          {
            id: 'feat-2',
            title: 'Another Featured',
            artist: 'Cool Artist',
            album: 'Cool Album',
            genre: 'Rock',
            year: 2023,
            duration: 180,
          },
        ],
        total: 2,
      },
      isLoading: false,
      isError: false,
      error: null,
    });

    render(<HelloSearchPage />);

    expect(screen.getByText('Featured Hit')).toBeInTheDocument();
    expect(screen.getByText('Another Featured')).toBeInTheDocument();
  });

  it('shows loading skeletons while fetching', () => {
    mockUseFeaturedTracks.mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
      error: null,
    });

    render(<HelloSearchPage />);

    const skeletons = screen.getAllByTestId('skeleton-line');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('shows error message when API fails', () => {
    mockUseFeaturedTracks.mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
      error: new Error('API is down'),
    });

    render(<HelloSearchPage />);

    expect(screen.getByText(/error/i)).toBeInTheDocument();
  });

  it('performs search when query is entered', () => {
    mockUseHelloSearch.mockReturnValue({
      data: {
        items: [
          {
            id: 'search-1',
            title: 'Search Result',
            artist: 'Found Artist',
            album: 'Found Album',
            genre: 'Jazz',
            year: 2024,
            duration: 300,
          },
        ],
        total: 1,
      },
      isLoading: false,
      isError: false,
      error: null,
    });

    render(<HelloSearchPage />);

    const input = screen.getByPlaceholderText(/search/i);
    fireEvent.change(input, { target: { value: 'jazz' } });

    expect(screen.getByText('Search Result')).toBeInTheDocument();
  });

  it('shows search results instead of featured tracks when searching', () => {
    mockUseFeaturedTracks.mockReturnValue({
      data: {
        items: [
          {
            id: 'feat-1',
            title: 'Featured Track',
            artist: 'Featured Artist',
            album: 'Featured Album',
            genre: 'Pop',
            year: 2024,
            duration: 200,
          },
        ],
        total: 1,
      },
      isLoading: false,
      isError: false,
      error: null,
    });

    mockUseHelloSearch.mockReturnValue({
      data: {
        items: [
          {
            id: 'search-1',
            title: 'Search Result Track',
            artist: 'Search Artist',
            album: 'Search Album',
            genre: 'Electronic',
            year: 2024,
            duration: 250,
          },
        ],
        total: 1,
      },
      isLoading: false,
      isError: false,
      error: null,
    });

    render(<HelloSearchPage />);

    const input = screen.getByPlaceholderText(/search/i);
    fireEvent.change(input, { target: { value: 'electronic' } });

    expect(screen.getByText('Search Result Track')).toBeInTheDocument();
  });
});
