/**
 * Login Page Tests - Task 1.5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RouterProvider, createRouter, createRootRoute, createRoute, createMemoryHistory } from '@tanstack/react-router';

// Mock useAuth
const mockSignIn = vi.fn();
const mockClearError = vi.fn();
const mockUseAuth = vi.fn((): {
  user: { userId: string; email: string } | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  isSigningIn: boolean;
  signIn: typeof mockSignIn;
  signOut: ReturnType<typeof vi.fn>;
  error: { message: string; code: string } | null;
  clearError: typeof mockClearError;
} => ({
  user: null,
  isLoading: false,
  isAuthenticated: false,
  isSigningIn: false,
  signIn: mockSignIn,
  signOut: vi.fn(),
  error: null,
  clearError: mockClearError,
}));

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock('../../lib/auth', () => ({
  AuthErrorCode: {
    INVALID_CREDENTIALS: 'INVALID_CREDENTIALS',
    USER_NOT_FOUND: 'USER_NOT_FOUND',
    USER_NOT_CONFIRMED: 'USER_NOT_CONFIRMED',
    TOO_MANY_REQUESTS: 'TOO_MANY_REQUESTS',
    NETWORK_ERROR: 'NETWORK_ERROR',
    MFA_REQUIRED: 'MFA_REQUIRED',
    UNKNOWN: 'UNKNOWN',
  },
}));

// Import after mocks
import LoginPage from '../login';

function renderLogin() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });

  const rootRoute = createRootRoute();
  const loginRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/login',
    component: LoginPage,
  });
  const homeRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/',
    component: () => <div>Home</div>,
  });

  const routeTree = rootRoute.addChildren([loginRoute, homeRoute]);
  const router = createRouter({ routeTree, history: createMemoryHistory({ initialEntries: ['/login'] }) });

  return {
    user: userEvent.setup(),
    ...render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    ),
  };
}

describe('Login Page (Task 1.5)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue({
      user: null,
      isLoading: false,
      isAuthenticated: false,
      isSigningIn: false,
      signIn: mockSignIn,
      signOut: vi.fn(),
      error: null,
      clearError: mockClearError,
    });
  });

  describe('form rendering', () => {
    it('should render email input', async () => {
      renderLogin();
      await waitFor(() => {
        expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
      });
    });

    it('should render password input', async () => {
      renderLogin();
      await waitFor(() => {
        expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
      });
    });

    it('should render submit button', async () => {
      renderLogin();
      await waitFor(() => {
        expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();
      });
    });

    it('should render page title', async () => {
      renderLogin();
      await waitFor(() => {
        expect(screen.getByRole('heading', { name: /sign in/i })).toBeInTheDocument();
      });
    });

    it('should render signup link', async () => {
      renderLogin();
      await waitFor(() => {
        expect(screen.getByRole('link', { name: /create account/i })).toBeInTheDocument();
      });
    });
  });

  describe('form validation', () => {
    it('should show error for empty email', async () => {
      const { user } = renderLogin();

      await waitFor(() => expect(screen.getByLabelText(/email/i)).toBeInTheDocument());

      const passwordInput = screen.getByLabelText(/password/i);
      await user.type(passwordInput, 'ValidPassword123!');

      const submitButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/email.*required/i)).toBeInTheDocument();
      });
    });

    it('should show error for empty password', async () => {
      const { user } = renderLogin();

      await waitFor(() => expect(screen.getByLabelText(/email/i)).toBeInTheDocument());

      const emailInput = screen.getByLabelText(/email/i);
      await user.type(emailInput, 'test@example.com');

      const submitButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/password.*required/i)).toBeInTheDocument();
      });
    });

    it('should show error for invalid email', async () => {
      const { user } = renderLogin();

      await waitFor(() => expect(screen.getByLabelText(/email/i)).toBeInTheDocument());

      const emailInput = screen.getByLabelText(/email/i);
      await user.type(emailInput, 'invalid-email');

      const passwordInput = screen.getByLabelText(/password/i);
      await user.type(passwordInput, 'password123');

      const submitButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/valid email/i)).toBeInTheDocument();
      });
    });

    it('should not submit when validation fails', async () => {
      const { user } = renderLogin();

      await waitFor(() => expect(screen.getByLabelText(/email/i)).toBeInTheDocument());

      const submitButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(submitButton);

      expect(mockSignIn).not.toHaveBeenCalled();
    });
  });

  describe('form submission', () => {
    it('should call signIn with email and password', async () => {
      mockSignIn.mockResolvedValue({ userId: '123', email: 'test@example.com' });

      const { user } = renderLogin();

      await waitFor(() => expect(screen.getByLabelText(/email/i)).toBeInTheDocument());

      const emailInput = screen.getByLabelText(/email/i);
      await user.type(emailInput, 'test@example.com');

      const passwordInput = screen.getByLabelText(/password/i);
      await user.type(passwordInput, 'ValidPassword123!');

      const submitButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(mockSignIn).toHaveBeenCalledWith('test@example.com', 'ValidPassword123!');
      });
    });

    it('should trim email whitespace', async () => {
      mockSignIn.mockResolvedValue({ userId: '123', email: 'test@example.com' });

      const { user } = renderLogin();

      await waitFor(() => expect(screen.getByLabelText(/email/i)).toBeInTheDocument());

      const emailInput = screen.getByLabelText(/email/i);
      await user.type(emailInput, '  test@example.com  ');

      const passwordInput = screen.getByLabelText(/password/i);
      await user.type(passwordInput, 'password');

      const submitButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(mockSignIn).toHaveBeenCalledWith('test@example.com', 'password');
      });
    });
  });

  describe('loading state', () => {
    it('should disable button during authentication', async () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: false,
        isAuthenticated: false,
        isSigningIn: true,
        signIn: mockSignIn,
        signOut: vi.fn(),
        error: null,
        clearError: mockClearError,
      });

      renderLogin();

      await waitFor(() => {
        const submitButton = screen.getByRole('button', { name: /signing in/i });
        expect(submitButton).toBeDisabled();
      });
    });

    it('should disable inputs during authentication', async () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: false,
        isAuthenticated: false,
        isSigningIn: true,
        signIn: mockSignIn,
        signOut: vi.fn(),
        error: null,
        clearError: mockClearError,
      });

      renderLogin();

      await waitFor(() => {
        expect(screen.getByLabelText(/email/i)).toBeDisabled();
        expect(screen.getByLabelText(/password/i)).toBeDisabled();
      });
    });
  });

  describe('error handling', () => {
    it('should display error message', async () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: false,
        isAuthenticated: false,
        isSigningIn: false,
        signIn: mockSignIn,
        signOut: vi.fn(),
        error: { message: 'Invalid credentials', code: 'INVALID_CREDENTIALS' as const },
        clearError: mockClearError,
      });

      renderLogin();

      await waitFor(() => {
        expect(screen.getByRole('alert')).toBeInTheDocument();
        expect(screen.getByText(/incorrect email or password/i)).toBeInTheDocument();
      });
    });

    it('should call clearError when dismiss clicked', async () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: false,
        isAuthenticated: false,
        isSigningIn: false,
        signIn: mockSignIn,
        signOut: vi.fn(),
        error: { message: 'Error', code: 'UNKNOWN' as const },
        clearError: mockClearError,
      });

      const { user } = renderLogin();

      await waitFor(() => expect(screen.getByRole('alert')).toBeInTheDocument());

      const dismissButton = screen.getByRole('button', { name: /dismiss/i });
      await user.click(dismissButton);

      expect(mockClearError).toHaveBeenCalled();
    });
  });

  describe('accessibility', () => {
    it('should have accessible form labels', async () => {
      renderLogin();

      await waitFor(() => {
        const emailInput = screen.getByLabelText(/email/i);
        const passwordInput = screen.getByLabelText(/password/i);

        expect(emailInput).toHaveAccessibleName();
        expect(passwordInput).toHaveAccessibleName();
      });
    });

    it('should have required attribute on inputs', async () => {
      renderLogin();

      await waitFor(() => {
        expect(screen.getByLabelText(/email/i)).toHaveAttribute('required');
        expect(screen.getByLabelText(/password/i)).toHaveAttribute('required');
      });
    });

    it('should have main landmark', async () => {
      renderLogin();

      await waitFor(() => {
        expect(screen.getByRole('main')).toBeInTheDocument();
      });
    });
  });

  describe('initial auth loading', () => {
    it('should show loading state', async () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: true,
        isAuthenticated: false,
        isSigningIn: false,
        signIn: mockSignIn,
        signOut: vi.fn(),
        error: null,
        clearError: mockClearError,
      });

      renderLogin();

      await waitFor(() => {
        expect(screen.getByRole('status')).toBeInTheDocument();
      });
    });
  });
});
