/**
 * WatchButton Component Tests - Artist Events Feature (Group 4)
 *
 * TDD Red Phase: These tests define the contract for WatchButton.
 * They MUST FAIL because the component does not exist yet.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { WatchButton } from '../WatchButton';

// Mock useAuth
const mockUseAuth = vi.fn();
vi.mock('../../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

// Mock useWatchToggle
const mockToggle = vi.fn();
const mockUseWatchToggle = vi.fn();
vi.mock('../../../hooks/useArtistWatch', () => ({
  useWatchToggle: (artistName: string) => mockUseWatchToggle(artistName),
}));

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('WatchButton', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Default: authenticated subscriber
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isSubscriber: true,
    });
    // Default: not watching, not loading
    mockUseWatchToggle.mockReturnValue({
      isWatching: false,
      isLoading: false,
      toggle: mockToggle,
    });
  });

  it('should render "Watch" button when not watching', () => {
    render(<WatchButton artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button', { name: 'Watch' })).toBeInTheDocument();
  });

  it('should render "Watching" button when watching', () => {
    mockUseWatchToggle.mockReturnValue({
      isWatching: true,
      isLoading: false,
      toggle: mockToggle,
    });

    render(<WatchButton artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button', { name: 'Watching' })).toBeInTheDocument();
  });

  it('should call toggle on click', () => {
    render(<WatchButton artistName="Daft Punk" />, { wrapper: createWrapper() });

    fireEvent.click(screen.getByRole('button'));
    expect(mockToggle).toHaveBeenCalled();
  });

  it('should show loading spinner when isLoading is true', () => {
    mockUseWatchToggle.mockReturnValue({
      isWatching: false,
      isLoading: true,
      toggle: mockToggle,
    });

    render(<WatchButton artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button')).toBeDisabled();
    expect(screen.getByRole('button').querySelector('.loading')).toBeInTheDocument();
  });

  it('should return null when not authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isSubscriber: false,
    });

    const { container } = render(<WatchButton artistName="Daft Punk" />, {
      wrapper: createWrapper(),
    });
    expect(container).toBeEmptyDOMElement();
  });

  it('should return null when not a subscriber', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isSubscriber: false,
    });

    const { container } = render(<WatchButton artistName="Daft Punk" />, {
      wrapper: createWrapper(),
    });
    expect(container).toBeEmptyDOMElement();
  });

  it('should pass artistName to useWatchToggle', () => {
    render(<WatchButton artistName="Bonobo" />, { wrapper: createWrapper() });
    expect(mockUseWatchToggle).toHaveBeenCalledWith('Bonobo');
  });

  it('should apply btn-primary class when not watching', () => {
    render(<WatchButton artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button')).toHaveClass('btn-primary');
  });

  it('should apply btn-outline class when watching', () => {
    mockUseWatchToggle.mockReturnValue({
      isWatching: true,
      isLoading: false,
      toggle: mockToggle,
    });

    render(<WatchButton artistName="Daft Punk" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button')).toHaveClass('btn-outline');
  });

  it('should apply btn-sm class for size="sm"', () => {
    render(<WatchButton artistName="Daft Punk" size="sm" />, {
      wrapper: createWrapper(),
    });
    expect(screen.getByRole('button')).toHaveClass('btn-sm');
  });

  it('should apply custom className', () => {
    render(<WatchButton artistName="Daft Punk" className="my-custom-class" />, {
      wrapper: createWrapper(),
    });
    expect(screen.getByRole('button')).toHaveClass('my-custom-class');
  });
});
