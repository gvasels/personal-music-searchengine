/**
 * TrackCardSkeleton Component Tests - Hello World Feature (Red Phase)
 *
 * Tests for a skeleton loading placeholder that mimics the TrackCard layout.
 */
import { describe, it, expect } from 'vitest';
import { render } from '@/test/test-utils';
import { TrackCardSkeleton } from '../TrackCardSkeleton';

describe('TrackCardSkeleton', () => {
  it('renders skeleton elements', () => {
    const { container } = render(<TrackCardSkeleton />);
    const skeletonElements = container.querySelectorAll('.skeleton');
    expect(skeletonElements.length).toBeGreaterThan(0);
  });

  it('has multiple skeleton lines for layout', () => {
    const { container } = render(<TrackCardSkeleton />);
    const skeletonElements = container.querySelectorAll('.skeleton');
    expect(skeletonElements.length).toBeGreaterThanOrEqual(2);
  });
});
