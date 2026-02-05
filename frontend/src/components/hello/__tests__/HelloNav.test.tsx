/**
 * HelloNav Component Tests - Hello World Feature (Red Phase)
 *
 * Tests for a navigation bar that displays the page title and a link to home.
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { HelloNav } from '../HelloNav';

describe('HelloNav', () => {
  it('renders "Hello Music Search" text', () => {
    render(<HelloNav />);
    expect(screen.getByText('Hello Music Search')).toBeInTheDocument();
  });

  it('has a link to home page', () => {
    render(<HelloNav />);
    const link = screen.getByRole('link', { name: /home/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute('href', '/');
  });
});
