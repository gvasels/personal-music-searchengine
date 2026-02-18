/**
 * TrackList Component Tests - REQ-5.3
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { Track } from '@/types';

vi.mock('@/lib/store/playerStore', () => ({
  usePlayerStore: vi.fn(() => ({
    currentTrack: null,
    isPlaying: false,
    setQueue: vi.fn(),
  })),
}));

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(() => ({
    isAdmin: false,
    isLoading: false,
    isAuthenticated: true,
    role: 'subscriber',
  })),
}));

vi.mock('@/lib/store/preferencesStore', () => ({
  useShowUploadedBy: vi.fn(() => true),
}));

const mockTracks: Track[] = [
  {
    id: 'track-1',
    title: 'Track One',
    artist: 'Artist A',
    album: 'Album X',
    duration: 180,
    format: 'mp3',
    fileSize: 5000000,
    s3Key: 'tracks/1.mp3',
    tags: [],
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'track-2',
    title: 'Track Two',
    artist: 'Artist B',
    album: 'Album Y',
    duration: 240,
    format: 'mp3',
    fileSize: 6000000,
    s3Key: 'tracks/2.mp3',
    tags: [],
    createdAt: '2024-01-02T00:00:00Z',
    updatedAt: '2024-01-02T00:00:00Z',
  },
];

describe('TrackList Component', () => {
  it('REQ-5.3: should render all tracks', async () => {
    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={mockTracks} />);
    expect(screen.getByText('Track One')).toBeInTheDocument();
    expect(screen.getByText('Track Two')).toBeInTheDocument();
  });

  it('should render column headers', async () => {
    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={mockTracks} />);
    // Use getAllByText since both header and data cells contain these words
    const titleElements = screen.getAllByText(/title/i);
    const artistElements = screen.getAllByText(/artist/i);
    expect(titleElements.length).toBeGreaterThan(0);
    expect(artistElements.length).toBeGreaterThan(0);
  });

  it('should show empty state when no tracks', async () => {
    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={[]} />);
    expect(screen.getByText(/no tracks/i)).toBeInTheDocument();
  });

  it('should show loading skeleton when loading', async () => {
    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={[]} isLoading />);
    expect(screen.getByTestId('track-list-skeleton')).toBeInTheDocument();
  });

  it('REQ-5.3: should play track on click', async () => {
    const { TrackList } = await import('@/components/library/TrackList');
    const { usePlayerStore } = await import('@/lib/store/playerStore');
    const setQueue = vi.fn();
    vi.mocked(usePlayerStore).mockReturnValue({
      currentTrack: null,
      isPlaying: false,
      setQueue,
    });

    const { user } = render(<TrackList tracks={mockTracks} />);
    await user.click(screen.getByText('Track One'));
    expect(setQueue).toHaveBeenCalled();
  });

  it('should have table role for accessibility', async () => {
    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={mockTracks} />);
    expect(screen.getByRole('table')).toBeInTheDocument();
  });
});

/**
 * Admin Track Reprocess Tests - TDD Red Phase
 *
 * Tests for the reprocess button visibility based on admin role.
 */
describe('TrackList Admin Reprocess Button', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Reset to default mocks before each test
    vi.resetModules();
  });

  it('shows reprocess button for admin users when showUploadedBy is true', async () => {
    // Mock useAuth to return admin user
    const { useAuth } = await import('@/hooks/useAuth');
    vi.mocked(useAuth).mockReturnValue({
      isAdmin: true,
      isLoading: false,
      isAuthenticated: true,
      role: 'admin',
      user: null,
      isSigningIn: false,
      isSigningOut: false,
      signIn: vi.fn(),
      signOut: vi.fn(),
      completeNewPassword: vi.fn(),
      error: null,
      clearError: vi.fn(),
      refetch: vi.fn(),
      groups: ['admin'],
      can: vi.fn(() => true),
      isArtist: true,
      isSubscriber: true,
    });

    // Mock preferences to show uploaded by column
    const { useShowUploadedBy } = await import('@/lib/store/preferencesStore');
    vi.mocked(useShowUploadedBy).mockReturnValue(true);

    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={mockTracks} showUploadedBy />);

    // Should show reprocess button(s) for admin with showUploadedBy enabled
    const reprocessButtons = screen.getAllByLabelText(/reprocess/i);
    expect(reprocessButtons.length).toBeGreaterThan(0);
  });

  it('hides reprocess button for non-admin users', async () => {
    // Mock useAuth to return non-admin user
    const { useAuth } = await import('@/hooks/useAuth');
    vi.mocked(useAuth).mockReturnValue({
      isAdmin: false,
      isLoading: false,
      isAuthenticated: true,
      role: 'subscriber',
      user: null,
      isSigningIn: false,
      isSigningOut: false,
      signIn: vi.fn(),
      signOut: vi.fn(),
      completeNewPassword: vi.fn(),
      error: null,
      clearError: vi.fn(),
      refetch: vi.fn(),
      groups: ['subscriber'],
      can: vi.fn(() => false),
      isArtist: false,
      isSubscriber: true,
    });

    // Mock preferences to show uploaded by column
    const { useShowUploadedBy } = await import('@/lib/store/preferencesStore');
    vi.mocked(useShowUploadedBy).mockReturnValue(true);

    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={mockTracks} showUploadedBy />);

    // Should NOT show reprocess button for non-admin
    expect(screen.queryByLabelText(/reprocess/i)).not.toBeInTheDocument();
  });

  it('hides reprocess button for admin users when showUploadedBy is false', async () => {
    // Mock useAuth to return admin user
    const { useAuth } = await import('@/hooks/useAuth');
    vi.mocked(useAuth).mockReturnValue({
      isAdmin: true,
      isLoading: false,
      isAuthenticated: true,
      role: 'admin',
      user: null,
      isSigningIn: false,
      isSigningOut: false,
      signIn: vi.fn(),
      signOut: vi.fn(),
      completeNewPassword: vi.fn(),
      error: null,
      clearError: vi.fn(),
      refetch: vi.fn(),
      groups: ['admin'],
      can: vi.fn(() => true),
      isArtist: true,
      isSubscriber: true,
    });

    // Mock preferences to hide uploaded by column
    const { useShowUploadedBy } = await import('@/lib/store/preferencesStore');
    vi.mocked(useShowUploadedBy).mockReturnValue(false);

    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={mockTracks} showUploadedBy={false} />);

    // Should NOT show reprocess button when showUploadedBy is explicitly false
    expect(screen.queryByLabelText(/reprocess/i)).not.toBeInTheDocument();
  });

  it('renders reprocess button in actions column for each track when admin', async () => {
    // Mock useAuth to return admin user
    const { useAuth } = await import('@/hooks/useAuth');
    vi.mocked(useAuth).mockReturnValue({
      isAdmin: true,
      isLoading: false,
      isAuthenticated: true,
      role: 'admin',
      user: null,
      isSigningIn: false,
      isSigningOut: false,
      signIn: vi.fn(),
      signOut: vi.fn(),
      completeNewPassword: vi.fn(),
      error: null,
      clearError: vi.fn(),
      refetch: vi.fn(),
      groups: ['admin'],
      can: vi.fn(() => true),
      isArtist: true,
      isSubscriber: true,
    });

    // Mock preferences to show uploaded by column
    const { useShowUploadedBy } = await import('@/lib/store/preferencesStore');
    vi.mocked(useShowUploadedBy).mockReturnValue(true);

    const { TrackList } = await import('@/components/library/TrackList');
    render(<TrackList tracks={mockTracks} showUploadedBy />);

    // Should show one reprocess button per track
    const reprocessButtons = screen.getAllByLabelText(/reprocess/i);
    expect(reprocessButtons).toHaveLength(mockTracks.length);
  });
});
