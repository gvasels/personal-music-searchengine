/**
 * TrackCard Component Tests - Hello World Feature (Red Phase)
 *
 * Tests for a card component that displays track metadata with formatted duration.
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { TrackCard } from '../TrackCard';

interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;
}

const mockTrack: HelloTrack = {
  id: 'track-1',
  title: 'Blue in Green',
  artist: 'Miles Davis',
  album: 'Kind of Blue',
  genre: 'Jazz',
  year: 1959,
  duration: 240,
};

describe('TrackCard', () => {
  it('renders track title', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Blue in Green')).toBeInTheDocument();
  });

  it('renders track artist', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Miles Davis')).toBeInTheDocument();
  });

  it('renders track album', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Kind of Blue')).toBeInTheDocument();
  });

  it('renders track genre', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Jazz')).toBeInTheDocument();
  });

  it('renders track year', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('1959')).toBeInTheDocument();
  });

  it('formats duration 240 seconds as "4:00"', () => {
    render(<TrackCard track={{ ...mockTrack, duration: 240 }} />);
    expect(screen.getByText('4:00')).toBeInTheDocument();
  });

  it('formats duration 185 seconds as "3:05" with zero padding', () => {
    render(<TrackCard track={{ ...mockTrack, duration: 185 }} />);
    expect(screen.getByText('3:05')).toBeInTheDocument();
  });

  it('formats duration 60 seconds as "1:00"', () => {
    render(<TrackCard track={{ ...mockTrack, duration: 60 }} />);
    expect(screen.getByText('1:00')).toBeInTheDocument();
  });

  it('renders with DaisyUI card classes', () => {
    const { container } = render(<TrackCard track={mockTrack} />);
    const cardElement = container.querySelector('.card');
    expect(cardElement).toBeInTheDocument();
  });
});
