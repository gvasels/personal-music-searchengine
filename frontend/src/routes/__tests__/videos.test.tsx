/**
 * Videos Page Tests - Epic 2.1 Task 1.1 (TDD Red)
 * Tests for /videos placeholder page
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';

vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => vi.fn(),
  useLocation: () => ({ pathname: '/videos', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

// This import will fail — component doesn't exist yet
import VideosPage from '../videos/index';

describe('VideosPage (Epic 2.1)', () => {
  it('should render Videos heading', () => {
    render(<VideosPage />);
    expect(screen.getByRole('heading', { name: /videos/i })).toBeInTheDocument();
  });

  it('should show coming soon empty state', () => {
    render(<VideosPage />);
    expect(screen.getByText(/coming soon/i)).toBeInTheDocument();
  });
});
