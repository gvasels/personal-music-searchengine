import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '../../test/test-utils';
import React from 'react';

// Mock the hooks
vi.mock('../../hooks/useHelloSearch', () => ({
  useHelloSearch: vi.fn(),
  useFeaturedTracks: vi.fn(),
}));

// Mock TanStack Router - createFileRoute returns a function that returns a function
vi.mock('@tanstack/react-router', () => ({
  createFileRoute: () => (opts: { component: unknown }) => opts,
  Link: ({
    children,
    ...props
  }: {
    children: React.ReactNode;
    to: string;
  }) => <a {...props}>{children}</a>,
}));

import { useHelloSearch, useFeaturedTracks } from '../../hooks/useHelloSearch';

// Import component after mocks are set up
const { Route } = await import('../hello-search');
const HelloSearchPage = Route.component;

describe('HelloSearchPage', () => {
  const mockUseHelloSearch = vi.mocked(useHelloSearch);
  const mockUseFeaturedTracks = vi.mocked(useFeaturedTracks);

  const mockTracks = [
    {
      id: '1',
      title: 'Midnight Drift',
      artist: 'Luna Waves',
      album: 'Waveforms',
      genre: 'Electronic',
      year: 2024,
      duration: 240,
      durationStr: '4:00',
      coverArtUrl: '',
    },
    {
      id: '2',
      title: 'Ghost Protocol',
      artist: 'DJ Phantom',
      album: 'Spectre',
      genre: 'House',
      year: 2025,
      duration: 330,
      durationStr: '5:30',
      coverArtUrl: '',
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    // Default: not searching, featured tracks loaded
    mockUseHelloSearch.mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: false,
      error: null,
    } as ReturnType<typeof useHelloSearch>);
    mockUseFeaturedTracks.mockReturnValue({
      data: { tracks: mockTracks, total: 2, query: '' },
      isLoading: false,
      isError: false,
      error: null,
    } as ReturnType<typeof useFeaturedTracks>);
  });

  it('renders hero section with "Music Search" heading', () => {
    render(<HelloSearchPage />);
    expect(screen.getByText('Music Search')).toBeInTheDocument();
  });

  it('renders search input', () => {
    render(<HelloSearchPage />);
    expect(screen.getAllByRole('textbox').length).toBeGreaterThan(0);
  });

  it('shows "Featured Tracks" section when no search query', () => {
    render(<HelloSearchPage />);
    expect(screen.getByText('Featured Tracks')).toBeInTheDocument();
  });

  it('renders track cards from featured data', () => {
    render(<HelloSearchPage />);
    expect(screen.getByText('Midnight Drift')).toBeInTheDocument();
    expect(screen.getByText('Ghost Protocol')).toBeInTheDocument();
  });

  it('shows subtitle text', () => {
    render(<HelloSearchPage />);
    expect(
      screen.getByText('Find your next favorite track')
    ).toBeInTheDocument();
  });
});
