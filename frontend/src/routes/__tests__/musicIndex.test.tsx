/**
 * Music Index Page Tests - Epic 2.1 Task 1.1 (TDD Red)
 * Tests for /music landing page
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseAuth = vi.fn();
vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock('../../lib/api/stats', () => ({
  getLibraryStats: vi.fn(() =>
    Promise.resolve({
      totalTracks: 10,
      totalAlbums: 3,
      totalArtists: 5,
      totalDuration: 3600,
    })
  ),
}));

vi.mock('../../lib/api/client', () => ({
  getTracks: vi.fn(() =>
    Promise.resolve({
      items: [
        { id: 't1', title: 'Track One', artist: 'Artist A', duration: 180 },
      ],
      total: 1,
      limit: 5,
      offset: 0,
    })
  ),
}));

vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => vi.fn(),
  useLocation: () => ({ pathname: '/music', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

// This import will fail — component doesn't exist yet
import MusicIndexPage from '../music/index';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('MusicIndexPage (Epic 2.1)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue({
      user: { sub: 'user-1', role: 'subscriber' },
      isAuthenticated: true,
    });
  });

  it('should render music library heading', () => {
    render(<MusicIndexPage />, { wrapper: createWrapper() });
    expect(screen.getByRole('heading', { name: /music library/i })).toBeInTheDocument();
  });

  it('should show recent tracks section', () => {
    render(<MusicIndexPage />, { wrapper: createWrapper() });
    expect(screen.getByText(/recent tracks/i)).toBeInTheDocument();
  });

  it('should show library stats', async () => {
    render(<MusicIndexPage />, { wrapper: createWrapper() });
    await waitFor(() => {
      expect(screen.getByText('Albums')).toBeInTheDocument();
    });
    expect(screen.getByText('Artists')).toBeInTheDocument();
  });
});
