import { describe, it, expect } from 'vitest';
import { render, screen } from '../../../test/test-utils';
import { TrackCard } from '../TrackCard';

const mockTrack = {
  id: '1',
  title: 'Midnight Drift',
  artist: 'Luna Waves',
  album: 'Waveforms',
  genre: 'Electronic',
  year: 2024,
  duration: 240,
  durationStr: '4:00',
  coverArtUrl: '',
};

describe('TrackCard', () => {
  it('renders track title', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Midnight Drift')).toBeInTheDocument();
  });

  it('renders artist name', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Luna Waves')).toBeInTheDocument();
  });

  it('renders album name', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Waveforms')).toBeInTheDocument();
  });

  it('renders genre as badge', () => {
    render(<TrackCard track={mockTrack} />);
    const badge = screen.getByText('Electronic');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('badge');
  });

  it('renders formatted duration', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('4:00')).toBeInTheDocument();
  });

  it('renders cover art image with alt text', () => {
    render(<TrackCard track={mockTrack} />);
    const img = screen.getByAltText('Midnight Drift by Luna Waves');
    expect(img).toBeInTheDocument();
  });
});
