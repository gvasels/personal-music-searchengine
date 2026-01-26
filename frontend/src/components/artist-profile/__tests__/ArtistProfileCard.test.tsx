/**
 * ArtistProfileCard Component Tests - Global User Type Feature
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ArtistProfileCard } from '../ArtistProfileCard';
import type { ArtistProfile } from '../../../types';

// Mock useNavigate
const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', () => ({
  useNavigate: () => mockNavigate,
}));

// Mock FollowButton
vi.mock('../../follow/FollowButton', () => ({
  FollowButton: ({ artistId }: { artistId: string }) => (
    <button data-testid="follow-button">Follow {artistId}</button>
  ),
}));

const mockProfile: ArtistProfile = {
  userId: 'user-123',
  displayName: 'Test Artist',
  bio: 'This is a test artist bio that describes their music.',
  avatarUrl: 'https://example.com/avatar.jpg',
  headerImageUrl: 'https://example.com/header.jpg',
  location: 'New York, USA',
  website: 'https://example.com',
  socialLinks: { twitter: 'https://twitter.com/testartist' },
  isVerified: true,
  followerCount: 1000,
  followingCount: 50,
  trackCount: 25,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('ArtistProfileCard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render artist display name', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    expect(screen.getByText('Test Artist')).toBeInTheDocument();
  });

  it('should render verified badge for verified artists', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    expect(screen.getByText('Verified')).toBeInTheDocument();
  });

  it('should not render verified badge for unverified artists', () => {
    const unverifiedProfile = { ...mockProfile, isVerified: false };
    render(<ArtistProfileCard profile={unverifiedProfile} />, { wrapper: createWrapper() });
    expect(screen.queryByText('Verified')).not.toBeInTheDocument();
  });

  it('should render bio text', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    expect(screen.getByText(/This is a test artist bio/)).toBeInTheDocument();
  });

  it('should render location when provided', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    expect(screen.getByText('New York, USA')).toBeInTheDocument();
  });

  it('should not render location when not provided', () => {
    const profileWithoutLocation = { ...mockProfile, location: undefined };
    render(<ArtistProfileCard profile={profileWithoutLocation} />, { wrapper: createWrapper() });
    expect(screen.queryByText('New York, USA')).not.toBeInTheDocument();
  });

  it('should render follower count', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    expect(screen.getByText('1000')).toBeInTheDocument();
    expect(screen.getByText('followers')).toBeInTheDocument();
  });

  it('should render track count', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    expect(screen.getByText('25')).toBeInTheDocument();
    expect(screen.getByText('tracks')).toBeInTheDocument();
  });

  it('should render avatar image when avatarUrl is provided', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    const avatar = screen.getByAltText('Test Artist');
    expect(avatar).toBeInTheDocument();
    expect(avatar).toHaveAttribute('src', 'https://example.com/avatar.jpg');
  });

  it('should render initial when no avatarUrl is provided', () => {
    const profileWithoutAvatar = { ...mockProfile, avatarUrl: undefined };
    render(<ArtistProfileCard profile={profileWithoutAvatar} />, { wrapper: createWrapper() });
    expect(screen.getByText('T')).toBeInTheDocument();
  });

  it('should render header image when provided', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    const headerImage = document.querySelector('figure img');
    expect(headerImage).toHaveAttribute('src', 'https://example.com/header.jpg');
  });

  it('should show FollowButton by default', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });
    expect(screen.getByTestId('follow-button')).toBeInTheDocument();
  });

  it('should hide FollowButton when showFollowButton is false', () => {
    render(<ArtistProfileCard profile={mockProfile} showFollowButton={false} />, {
      wrapper: createWrapper(),
    });
    expect(screen.queryByTestId('follow-button')).not.toBeInTheDocument();
  });

  it('should navigate to artist profile on click', () => {
    render(<ArtistProfileCard profile={mockProfile} />, { wrapper: createWrapper() });

    const nameElement = screen.getByText('Test Artist');
    fireEvent.click(nameElement);

    expect(mockNavigate).toHaveBeenCalledWith({
      to: '/artists/entity/user-123',
    });
  });
});
