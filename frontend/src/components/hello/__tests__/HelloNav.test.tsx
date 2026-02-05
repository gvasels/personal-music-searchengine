/**
 * HelloNav Component Tests - Hello World Local Dev Feature
 *
 * Tests for the hello search navigation component.
 * MUST FAIL until HelloNav component is implemented.
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';

describe('HelloNav Component', () => {
  it('renders "Hello Music Search" text', async () => {
    const { HelloNav } = await import('../HelloNav');
    render(<HelloNav />);

    expect(screen.getByText('Hello Music Search')).toBeInTheDocument();
  });

  it('has a link to home page (href="/")', async () => {
    const { HelloNav } = await import('../HelloNav');
    render(<HelloNav />);

    const homeLink = screen.getByRole('link', { name: /home/i });
    expect(homeLink).toHaveAttribute('href', '/');
  });

  it('renders navigation element', async () => {
    const { HelloNav } = await import('../HelloNav');
    render(<HelloNav />);

    expect(screen.getByRole('navigation')).toBeInTheDocument();
  });

  it('has DaisyUI navbar class', async () => {
    const { HelloNav } = await import('../HelloNav');
    const { container } = render(<HelloNav />);

    expect(container.querySelector('.navbar')).toBeInTheDocument();
  });
});
