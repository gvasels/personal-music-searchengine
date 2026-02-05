/**
 * TrackCardSkeleton Component Tests - TDD Red Phase
 *
 * Tests for the hello track card skeleton loading component.
 * Source module (TrackCardSkeleton.tsx) does NOT exist yet - these tests MUST fail.
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { TrackCardSkeleton } from '../TrackCardSkeleton';

describe('TrackCardSkeleton', () => {
  it('renders skeleton loading placeholders', () => {
    render(<TrackCardSkeleton />);

    const skeletons = screen.getAllByTestId('skeleton-line');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('has correct number of skeleton elements', () => {
    render(<TrackCardSkeleton />);

    // Should have skeleton lines for: title, artist, album, genre, year, duration
    const skeletons = screen.getAllByTestId('skeleton-line');
    expect(skeletons).toHaveLength(6);
  });
});
