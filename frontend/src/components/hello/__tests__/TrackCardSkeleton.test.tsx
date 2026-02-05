/**
 * TrackCardSkeleton Component Tests - Hello World Local Dev Feature
 *
 * Tests for the loading skeleton component.
 * MUST FAIL until TrackCardSkeleton component is implemented.
 */
import { describe, it, expect } from 'vitest';
import { render } from '@/test/test-utils';

describe('TrackCardSkeleton Component', () => {
  it('renders skeleton elements with .skeleton class', async () => {
    const { TrackCardSkeleton } = await import('../TrackCardSkeleton');
    const { container } = render(<TrackCardSkeleton />);

    const skeletonElements = container.querySelectorAll('.skeleton');
    expect(skeletonElements.length).toBeGreaterThan(0);
  });

  it('has multiple skeleton lines', async () => {
    const { TrackCardSkeleton } = await import('../TrackCardSkeleton');
    const { container } = render(<TrackCardSkeleton />);

    // Should have multiple skeleton elements for different track fields
    const skeletonElements = container.querySelectorAll('.skeleton');
    expect(skeletonElements.length).toBeGreaterThanOrEqual(3);
  });

  it('has DaisyUI card class', async () => {
    const { TrackCardSkeleton } = await import('../TrackCardSkeleton');
    const { container } = render(<TrackCardSkeleton />);

    expect(container.querySelector('.card')).toBeInTheDocument();
  });

  it('renders without crashing', async () => {
    const { TrackCardSkeleton } = await import('../TrackCardSkeleton');
    const { container } = render(<TrackCardSkeleton />);

    // Should render something
    expect(container.firstChild).not.toBeNull();
  });
});
