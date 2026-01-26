/**
 * FollowButton Component Tests - Global User Type Feature
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { FollowButton } from '../FollowButton';

// Mock useAuth
const mockUseAuth = vi.fn();
vi.mock('../../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

// Mock useFollowToggle
const mockToggle = vi.fn();
const mockUseFollowToggle = vi.fn();
vi.mock('../../../hooks/useFollow', () => ({
  useFollowToggle: (artistId: string) => mockUseFollowToggle(artistId),
}));

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('FollowButton', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Default: authenticated subscriber
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isSubscriber: true,
    });
    // Default: not following, not loading
    mockUseFollowToggle.mockReturnValue({
      isFollowing: false,
      isLoading: false,
      isToggling: false,
      toggle: mockToggle,
    });
  });

  it('should render Follow button when not following', () => {
    render(<FollowButton artistId="artist-123" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button', { name: 'Follow' })).toBeInTheDocument();
  });

  it('should render Following button when following', () => {
    mockUseFollowToggle.mockReturnValue({
      isFollowing: true,
      isLoading: false,
      isToggling: false,
      toggle: mockToggle,
    });

    render(<FollowButton artistId="artist-123" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button', { name: 'Following' })).toBeInTheDocument();
  });

  it('should call toggle when clicked', () => {
    render(<FollowButton artistId="artist-123" />, { wrapper: createWrapper() });

    fireEvent.click(screen.getByRole('button'));
    expect(mockToggle).toHaveBeenCalled();
  });

  it('should show loading spinner while loading initial state', () => {
    mockUseFollowToggle.mockReturnValue({
      isFollowing: false,
      isLoading: true,
      isToggling: false,
      toggle: mockToggle,
    });

    render(<FollowButton artistId="artist-123" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button')).toBeDisabled();
    expect(screen.getByRole('button').querySelector('.loading')).toBeInTheDocument();
  });

  it('should show loading spinner while toggling', () => {
    mockUseFollowToggle.mockReturnValue({
      isFollowing: false,
      isLoading: false,
      isToggling: true,
      toggle: mockToggle,
    });

    render(<FollowButton artistId="artist-123" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button')).toBeDisabled();
    expect(screen.getByRole('button').querySelector('.loading')).toBeInTheDocument();
  });

  it('should not render when not authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isSubscriber: false,
    });

    const { container } = render(<FollowButton artistId="artist-123" />, {
      wrapper: createWrapper(),
    });
    expect(container).toBeEmptyDOMElement();
  });

  it('should not render when not a subscriber', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isSubscriber: false,
    });

    const { container } = render(<FollowButton artistId="artist-123" />, {
      wrapper: createWrapper(),
    });
    expect(container).toBeEmptyDOMElement();
  });

  it('should apply btn-primary class when not following', () => {
    render(<FollowButton artistId="artist-123" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button')).toHaveClass('btn-primary');
  });

  it('should apply btn-outline class when following', () => {
    mockUseFollowToggle.mockReturnValue({
      isFollowing: true,
      isLoading: false,
      isToggling: false,
      toggle: mockToggle,
    });

    render(<FollowButton artistId="artist-123" />, { wrapper: createWrapper() });
    expect(screen.getByRole('button')).toHaveClass('btn-outline');
  });

  it('should apply size classes correctly', () => {
    const { rerender } = render(<FollowButton artistId="artist-123" size="sm" />, {
      wrapper: createWrapper(),
    });
    expect(screen.getByRole('button')).toHaveClass('btn-sm');

    rerender(<FollowButton artistId="artist-123" size="lg" />);
    expect(screen.getByRole('button')).toHaveClass('btn-lg');
  });

  it('should apply custom className', () => {
    render(<FollowButton artistId="artist-123" className="my-custom-class" />, {
      wrapper: createWrapper(),
    });
    expect(screen.getByRole('button')).toHaveClass('my-custom-class');
  });

  it('should pass artistId to useFollowToggle', () => {
    render(<FollowButton artistId="artist-456" />, { wrapper: createWrapper() });
    expect(mockUseFollowToggle).toHaveBeenCalledWith('artist-456');
  });
});
