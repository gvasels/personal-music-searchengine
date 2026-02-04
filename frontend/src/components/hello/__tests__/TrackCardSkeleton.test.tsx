import { describe, it, expect } from 'vitest';
import { render } from '../../../test/test-utils';
import { TrackCardSkeleton } from '../TrackCardSkeleton';

describe('TrackCardSkeleton', () => {
  it('renders skeleton elements', () => {
    const { container } = render(<TrackCardSkeleton />);
    const skeletons = container.querySelectorAll('.skeleton');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders correct number of skeleton lines', () => {
    const { container } = render(<TrackCardSkeleton />);
    const skeletons = container.querySelectorAll('.skeleton');
    // image skeleton + title + artist + album + badge + duration = 6
    expect(skeletons.length).toBe(6);
  });
});
