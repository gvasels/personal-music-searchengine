/**
 * PlayerBar Component Tests - REQ-5.6
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@/test/test-utils';

vi.mock('@/lib/store/playerStore', () => ({
  usePlayerStore: vi.fn(),
}));

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(() => ({
    isAuthenticated: true,
    user: { id: 'user-1', email: 'test@test.com' },
    isLoading: false,
    signIn: vi.fn(),
    signOut: vi.fn(),
    signUp: vi.fn(),
  })),
}));

const mockTrack = {
  id: 'track-1',
  title: 'Test Track',
  artist: 'Test Artist',
  album: 'Test Album',
  duration: 180,
  format: 'mp3',
  fileSize: 5000000,
  s3Key: 'tracks/1.mp3',
  tags: [],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};

describe('PlayerBar Component', () => {
  it('REQ-5.6: should render player bar', async () => {
    const { usePlayerStore } = await import('@/lib/store/playerStore');
    vi.mocked(usePlayerStore).mockReturnValue({
      currentTrack: null,
      isPlaying: false,
      volume: 1,
      progress: 0,
      shuffle: false,
      repeat: 'off',
      play: vi.fn(),
      pause: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seek: vi.fn(),
      setVolume: vi.fn(),
      toggleShuffle: vi.fn(),
      cycleRepeat: vi.fn(),
    });

    const { PlayerBar } = await import('@/components/player/PlayerBar');
    render(<PlayerBar />);
    expect(screen.getByTestId('player-bar')).toBeInTheDocument();
  });

  it('should display current track info', async () => {
    const { usePlayerStore } = await import('@/lib/store/playerStore');
    vi.mocked(usePlayerStore).mockReturnValue({
      currentTrack: mockTrack,
      isPlaying: true,
      volume: 1,
      progress: 60,
      shuffle: false,
      repeat: 'off',
      play: vi.fn(),
      pause: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seek: vi.fn(),
      setVolume: vi.fn(),
      toggleShuffle: vi.fn(),
      cycleRepeat: vi.fn(),
    });

    const { PlayerBar } = await import('@/components/player/PlayerBar');
    render(<PlayerBar />);
    expect(screen.getByText('Test Track')).toBeInTheDocument();
    expect(screen.getByText('Test Artist')).toBeInTheDocument();
  });

  it('REQ-5.6: should show play/pause button', async () => {
    const { usePlayerStore } = await import('@/lib/store/playerStore');
    vi.mocked(usePlayerStore).mockReturnValue({
      currentTrack: mockTrack,
      isPlaying: false,
      volume: 1,
      progress: 0,
      shuffle: false,
      repeat: 'off',
      play: vi.fn(),
      pause: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seek: vi.fn(),
      setVolume: vi.fn(),
      toggleShuffle: vi.fn(),
      cycleRepeat: vi.fn(),
    });

    const { PlayerBar } = await import('@/components/player/PlayerBar');
    render(<PlayerBar />);
    expect(screen.getByRole('button', { name: /play/i })).toBeInTheDocument();
  });

  it('should show next/previous buttons', async () => {
    const { usePlayerStore } = await import('@/lib/store/playerStore');
    vi.mocked(usePlayerStore).mockReturnValue({
      currentTrack: mockTrack,
      isPlaying: false,
      volume: 1,
      progress: 0,
      shuffle: false,
      repeat: 'off',
      play: vi.fn(),
      pause: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seek: vi.fn(),
      setVolume: vi.fn(),
      toggleShuffle: vi.fn(),
      cycleRepeat: vi.fn(),
    });

    const { PlayerBar } = await import('@/components/player/PlayerBar');
    render(<PlayerBar />);
    expect(screen.getByRole('button', { name: /next/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /previous/i })).toBeInTheDocument();
  });

  it('REQ-5.6: should show volume control', async () => {
    const { usePlayerStore } = await import('@/lib/store/playerStore');
    vi.mocked(usePlayerStore).mockReturnValue({
      currentTrack: mockTrack,
      isPlaying: false,
      volume: 0.7,
      progress: 0,
      shuffle: false,
      repeat: 'off',
      play: vi.fn(),
      pause: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seek: vi.fn(),
      setVolume: vi.fn(),
      toggleShuffle: vi.fn(),
      cycleRepeat: vi.fn(),
    });

    const { PlayerBar } = await import('@/components/player/PlayerBar');
    render(<PlayerBar />);
    expect(screen.getByRole('slider', { name: /volume/i })).toBeInTheDocument();
  });
});
