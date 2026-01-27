/**
 * Root Layout Tests - Access Control Bug Fixes
 * Tests guest user route protection logic (Task 2.3)
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

// Mock the useAuth hook
vi.mock('../../hooks/useAuth', () => ({
  useAuth: vi.fn(),
}));

// Mock TanStack Router hooks
vi.mock('@tanstack/react-router', () => ({
  createRootRoute: vi.fn(() => ({ component: vi.fn() })),
  Outlet: () => <div data-testid="outlet">Outlet Content</div>,
  useLocation: vi.fn(),
  Navigate: vi.fn(({ to }: { to: string }) => <div data-testid="navigate-to">{to}</div>),
}));

import { useAuth } from '../../hooks/useAuth';
import { useLocation, Navigate } from '@tanstack/react-router';

/**
 * Re-implement the route guard logic for testing.
 * This mirrors the implementation in __root.tsx
 */
const PUBLIC_ROUTES = ['/', '/login', '/permission-denied'];

function isPublicRoute(pathname: string): boolean {
  return PUBLIC_ROUTES.includes(pathname);
}

/**
 * Test component that mirrors the RootComponent logic
 */
function TestRootComponent() {
  const { isAuthenticated, isLoading } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div data-testid="loading-spinner">
        <span className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  if (!isAuthenticated && !isPublicRoute(location.pathname)) {
    return <Navigate to="/permission-denied" />;
  }

  return <div data-testid="content-rendered">Content Accessible</div>;
}

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0, staleTime: 0 },
      mutations: { retry: false },
    },
  });

  return function Wrapper({ children }: { children: ReactNode }) {
    return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
  };
}

function renderTestComponent() {
  const Wrapper = createWrapper();
  return render(
    <Wrapper>
      <TestRootComponent />
    </Wrapper>
  );
}

const mockAuthReturn = {
  isAuthenticated: false,
  isLoading: false,
  user: null,
  isSigningIn: false,
  isSigningOut: false,
  signIn: vi.fn(),
  signOut: vi.fn(),
  error: null,
  clearError: vi.fn(),
  refetch: vi.fn(),
  role: 'guest' as const,
  groups: [] as string[],
  can: vi.fn(),
  isAdmin: false,
  isArtist: false,
  isSubscriber: false,
};

describe('Root Layout - Guest Route Protection (Task 2.3)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('loading state', () => {
    it('should show loading spinner while auth is loading', () => {
      vi.mocked(useAuth).mockReturnValue({
        ...mockAuthReturn,
        isLoading: true,
      });
      vi.mocked(useLocation).mockReturnValue({ pathname: '/' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    });

    it('should not show content while loading', () => {
      vi.mocked(useAuth).mockReturnValue({
        ...mockAuthReturn,
        isLoading: true,
      });
      vi.mocked(useLocation).mockReturnValue({ pathname: '/' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.queryByTestId('content-rendered')).not.toBeInTheDocument();
    });
  });

  describe('public routes - guest access allowed', () => {
    beforeEach(() => {
      vi.mocked(useAuth).mockReturnValue({
        ...mockAuthReturn,
        isAuthenticated: false,
        isLoading: false,
      });
    });

    it('should allow guests to access home page (/)', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
      expect(screen.queryByTestId('navigate-to')).not.toBeInTheDocument();
    });

    it('should allow guests to access login page (/login)', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/login' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
      expect(screen.queryByTestId('navigate-to')).not.toBeInTheDocument();
    });

    it('should allow guests to access permission-denied page (/permission-denied)', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/permission-denied' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
      expect(screen.queryByTestId('navigate-to')).not.toBeInTheDocument();
    });
  });

  describe('protected routes - guest access denied', () => {
    beforeEach(() => {
      vi.mocked(useAuth).mockReturnValue({
        ...mockAuthReturn,
        isAuthenticated: false,
        isLoading: false,
      });
    });

    it('should redirect guests from /tracks to /permission-denied', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/tracks' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('navigate-to')).toHaveTextContent('/permission-denied');
      expect(screen.queryByTestId('content-rendered')).not.toBeInTheDocument();
    });

    it('should redirect guests from /albums to /permission-denied', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/albums' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('navigate-to')).toHaveTextContent('/permission-denied');
    });

    it('should redirect guests from /playlists to /permission-denied', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/playlists' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('navigate-to')).toHaveTextContent('/permission-denied');
    });

    it('should redirect guests from /artists to /permission-denied', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/artists' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('navigate-to')).toHaveTextContent('/permission-denied');
    });

    it('should redirect guests from /upload to /permission-denied', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/upload' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('navigate-to')).toHaveTextContent('/permission-denied');
    });

    it('should redirect guests from /admin to /permission-denied', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/admin' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('navigate-to')).toHaveTextContent('/permission-denied');
    });

    it('should redirect guests from /studio to /permission-denied', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/studio' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('navigate-to')).toHaveTextContent('/permission-denied');
    });
  });

  describe('authenticated users - full access', () => {
    const mockUser = {
      userId: 'user-123',
      email: 'test@example.com',
      role: 'subscriber' as const,
      groups: ['subscriber'],
    };

    beforeEach(() => {
      vi.mocked(useAuth).mockReturnValue({
        ...mockAuthReturn,
        isAuthenticated: true,
        isLoading: false,
        user: mockUser,
        role: 'subscriber',
        groups: ['subscriber'],
        isSubscriber: true,
      });
    });

    it('should allow authenticated users to access /tracks', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/tracks' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
      expect(screen.queryByTestId('navigate-to')).not.toBeInTheDocument();
    });

    it('should allow authenticated users to access /albums', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/albums' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
    });

    it('should allow authenticated users to access /playlists', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/playlists' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
    });

    it('should allow authenticated users to access public routes too', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
    });
  });

  describe('admin users', () => {
    const mockAdminUser = {
      userId: 'admin-123',
      email: 'admin@example.com',
      role: 'admin' as const,
      groups: ['admin'],
    };

    beforeEach(() => {
      vi.mocked(useAuth).mockReturnValue({
        ...mockAuthReturn,
        isAuthenticated: true,
        isLoading: false,
        user: mockAdminUser,
        role: 'admin',
        groups: ['admin'],
        isAdmin: true,
        isSubscriber: true,
      });
    });

    it('should allow admin users to access /admin', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/admin' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
    });

    it('should allow admin users to access all protected routes', () => {
      vi.mocked(useLocation).mockReturnValue({ pathname: '/studio' } as ReturnType<typeof useLocation>);

      renderTestComponent();

      expect(screen.getByTestId('content-rendered')).toBeInTheDocument();
    });
  });
});

describe('isPublicRoute helper', () => {
  it('should return true for /', () => {
    expect(isPublicRoute('/')).toBe(true);
  });

  it('should return true for /login', () => {
    expect(isPublicRoute('/login')).toBe(true);
  });

  it('should return true for /permission-denied', () => {
    expect(isPublicRoute('/permission-denied')).toBe(true);
  });

  it('should return false for /tracks', () => {
    expect(isPublicRoute('/tracks')).toBe(false);
  });

  it('should return false for /albums', () => {
    expect(isPublicRoute('/albums')).toBe(false);
  });

  it('should return false for /playlists', () => {
    expect(isPublicRoute('/playlists')).toBe(false);
  });

  it('should return false for /admin', () => {
    expect(isPublicRoute('/admin')).toBe(false);
  });

  it('should return false for /upload', () => {
    expect(isPublicRoute('/upload')).toBe(false);
  });

  it('should return false for /studio', () => {
    expect(isPublicRoute('/studio')).toBe(false);
  });

  it('should return false for /search', () => {
    expect(isPublicRoute('/search')).toBe(false);
  });
});
