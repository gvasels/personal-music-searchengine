/**
 * Home Page Tests - Task 1.11
 * Updated to mock new /stats API and useFeatureFlags
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RouterProvider, createRouter, createRootRoute, createRoute, createMemoryHistory } from '@tanstack/react-router';

// Mock useAuth
const mockUseAuth = vi.fn();
vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

// Mock useFeatureFlags
const mockUseFeatureFlags = vi.fn();
vi.mock('../../hooks/useFeatureFlags', () => ({
  useFeatureFlags: () => mockUseFeatureFlags(),
}));

// Mock tracks API
vi.mock('../../lib/api/client', () => ({
  getTracks: vi.fn(() =>
    Promise.resolve({
      items: [
        { id: 'track-1', title: 'Track One', artist: 'Artist A', album: 'Album 1', duration: 180 },
        { id: 'track-2', title: 'Track Two', artist: 'Artist B', album: 'Album 2', duration: 240 },
      ],
      total: 2,
      limit: 5,
      offset: 0,
    })
  ),
}));

// Mock stats API
vi.mock('../../lib/api/stats', () => ({
  getLibraryStats: vi.fn(() =>
    Promise.resolve({
      totalTracks: 2,
      totalAlbums: 2,
      totalArtists: 2,
      totalDuration: 420,
    })
  ),
}));

// Import after mocks
import HomePage from '../index';

function renderHome() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });

  const rootRoute = createRootRoute();
  const homeRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/',
    component: HomePage,
  });
  const loginRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/login',
    component: () => <div>Login Page</div>,
  });
  const uploadRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/upload',
    component: () => <div>Upload Page</div>,
  });
  const libraryRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/library',
    component: () => <div>Library Page</div>,
  });

  const routeTree = rootRoute.addChildren([homeRoute, loginRoute, uploadRoute, libraryRoute]);
  const router = createRouter({ routeTree, history: createMemoryHistory({ initialEntries: ['/'] }) });

  return {
    user: userEvent.setup(),
    router,
    ...render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    ),
  };
}

const mockUser = {
  userId: 'user-123',
  email: 'test@example.com',
  name: 'Test User',
};

describe('Home Page (Task 1.11)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue({
      user: mockUser,
      isLoading: false,
      isAuthenticated: true,
    });
    mockUseFeatureFlags.mockReturnValue({
      role: 'admin',
      isSimulating: false,
      isLoaded: true,
      features: {},
      tier: 'pro',
    });
  });

  describe('authentication', () => {
    it('should show loading state while checking auth', async () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: true,
        isAuthenticated: false,
      });

      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('status')).toBeInTheDocument();
      });
    });

    it('should render page when authenticated', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('main')).toBeInTheDocument();
      });
    });
  });

  describe('welcome message', () => {
    it('should display welcome with user name', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByText(/welcome back.*test user/i)).toBeInTheDocument();
      });
    });

    it('should display welcome with email if no name', async () => {
      mockUseAuth.mockReturnValue({
        user: { userId: 'user-123', email: 'test@example.com' },
        isLoading: false,
        isAuthenticated: true,
      });

      renderHome();

      await waitFor(() => {
        expect(screen.getByText(/welcome back.*test@example.com/i)).toBeInTheDocument();
      });
    });

    it('should have welcome as heading', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: /welcome/i })).toBeInTheDocument();
      });
    });
  });

  describe('library stats', () => {
    it('should display stats section', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByText(/library stats/i)).toBeInTheDocument();
      });
    });

    it('should display track count', async () => {
      renderHome();

      await waitFor(() => {
        // Look specifically for the stat title, not any tracks text
        const statTitles = screen.getAllByText(/tracks/i);
        // Should have at least the stat title "Tracks"
        expect(statTitles.length).toBeGreaterThan(0);
        // Verify the stat title specifically
        expect(screen.getByText('Tracks', { selector: '.stat-title' })).toBeInTheDocument();
      });
    });

    it('should display albums count', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByText(/albums/i)).toBeInTheDocument();
      });
    });

    it('should display artists count', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByText(/artists/i)).toBeInTheDocument();
      });
    });
  });

  describe('recent tracks', () => {
    it('should display recent tracks heading', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: /recent tracks/i })).toBeInTheDocument();
      });
    });

    it('should display track titles', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByText('Track One')).toBeInTheDocument();
        expect(screen.getByText('Track Two')).toBeInTheDocument();
      });
    });

    it('should display track artists', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByText(/artist a/i)).toBeInTheDocument();
        expect(screen.getByText(/artist b/i)).toBeInTheDocument();
      });
    });

    it('should display formatted durations', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByText('3:00')).toBeInTheDocument();
        expect(screen.getByText('4:00')).toBeInTheDocument();
      });
    });
  });

  describe('upload CTA', () => {
    it('should display upload button', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /upload/i })).toBeInTheDocument();
      });
    });

    it('should have primary styling', async () => {
      renderHome();

      await waitFor(() => {
        const uploadButton = screen.getByRole('button', { name: /upload/i });
        expect(uploadButton).toHaveClass('btn-primary');
      });
    });
  });

  describe('quick actions', () => {
    it('should display library link', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('link', { name: /library|browse/i })).toBeInTheDocument();
      });
    });

    it('should display playlists link', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('link', { name: /playlists/i })).toBeInTheDocument();
      });
    });

    it('should display search link', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('link', { name: /search/i })).toBeInTheDocument();
      });
    });
  });

  describe('accessibility', () => {
    it('should have main landmark', async () => {
      renderHome();

      await waitFor(() => {
        expect(screen.getByRole('main')).toBeInTheDocument();
      });
    });

    it('should have proper heading hierarchy', async () => {
      renderHome();

      await waitFor(() => {
        const h1 = screen.getByRole('heading', { level: 1 });
        expect(h1).toBeInTheDocument();

        const h2s = screen.getAllByRole('heading', { level: 2 });
        expect(h2s.length).toBeGreaterThanOrEqual(1);
      });
    });

    it('should have accessible names for buttons', async () => {
      renderHome();

      await waitFor(() => {
        const buttons = screen.getAllByRole('button');
        buttons.forEach((button) => {
          expect(button).toHaveAccessibleName();
        });
      });
    });

    it('should have accessible names for links', async () => {
      renderHome();

      await waitFor(() => {
        const links = screen.getAllByRole('link');
        links.forEach((link) => {
          expect(link).toHaveAccessibleName();
        });
      });
    });
  });
});
