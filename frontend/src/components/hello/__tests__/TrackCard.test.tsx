/**
 * TrackCard Component Tests - Hello World Local Dev Feature
 *
 * Tests for the track card display component.
 * MUST FAIL until TrackCard component is implemented.
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import type { HelloTrack } from '@/lib/api/helloSearch';

const mockTrack: HelloTrack = {
  id: 'seed-t1',
  title: 'Aurora Borealis',
  artist: 'Aurora Waves',
  album: 'Dreamscape',
  genre: 'jazz',
  year: 2023,
  duration: 245,
};

describe('TrackCard Component', () => {
  it('renders track title', async () => {
    const { TrackCard } = await import('../TrackCard');
    render(<TrackCard track={mockTrack} />);

    expect(screen.getByText('Aurora Borealis')).toBeInTheDocument();
  });

  it('renders track artist', async () => {
    const { TrackCard } = await import('../TrackCard');
    render(<TrackCard track={mockTrack} />);

    expect(screen.getByText('Aurora Waves')).toBeInTheDocument();
  });

  it('renders track album', async () => {
    const { TrackCard } = await import('../TrackCard');
    render(<TrackCard track={mockTrack} />);

    expect(screen.getByText('Dreamscape')).toBeInTheDocument();
  });

  it('renders track genre', async () => {
    const { TrackCard } = await import('../TrackCard');
    render(<TrackCard track={mockTrack} />);

    expect(screen.getByText('jazz')).toBeInTheDocument();
  });

  it('renders track year', async () => {
    const { TrackCard } = await import('../TrackCard');
    render(<TrackCard track={mockTrack} />);

    expect(screen.getByText('2023')).toBeInTheDocument();
  });

  it('formats duration 240s as "4:00"', async () => {
    const { TrackCard } = await import('../TrackCard');
    const trackWith240s = { ...mockTrack, duration: 240 };
    render(<TrackCard track={trackWith240s} />);

    expect(screen.getByText('4:00')).toBeInTheDocument();
  });

  it('formats duration 185s as "3:05" (zero padding)', async () => {
    const { TrackCard } = await import('../TrackCard');
    const trackWith185s = { ...mockTrack, duration: 185 };
    render(<TrackCard track={trackWith185s} />);

    expect(screen.getByText('3:05')).toBeInTheDocument();
  });

  it('has DaisyUI card class', async () => {
    const { TrackCard } = await import('../TrackCard');
    const { container } = render(<TrackCard track={mockTrack} />);

    expect(container.querySelector('.card')).toBeInTheDocument();
  });
});
