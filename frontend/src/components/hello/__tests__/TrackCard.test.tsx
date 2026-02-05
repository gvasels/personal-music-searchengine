/**
 * TrackCard Component Tests - TDD Red Phase
 *
 * Tests for the hello track card component.
 * Source module (TrackCard.tsx) does NOT exist yet - these tests MUST fail.
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { TrackCard } from '../TrackCard';

const mockTrack = {
  id: 'track-1',
  title: 'Test Song',
  artist: 'Test Artist',
  album: 'Test Album',
  genre: 'Rock',
  year: 2024,
  duration: 234,
};

describe('TrackCard', () => {
  it('renders track title', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Test Song')).toBeInTheDocument();
  });

  it('renders artist name', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Test Artist')).toBeInTheDocument();
  });

  it('renders album name', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('Test Album')).toBeInTheDocument();
  });

  it('renders genre as badge', () => {
    render(<TrackCard track={mockTrack} />);
    const badge = screen.getByText('Rock');
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass('badge');
  });

  it('renders year', () => {
    render(<TrackCard track={mockTrack} />);
    expect(screen.getByText('2024')).toBeInTheDocument();
  });

  it('formats duration as MM:SS', () => {
    render(<TrackCard track={mockTrack} />);
    // 234 seconds = 3 minutes 54 seconds = "3:54"
    expect(screen.getByText('3:54')).toBeInTheDocument();
  });
});
