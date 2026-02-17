/**
 * Gaming Page Tests - Epic 2.1 Task 1.1 (TDD Red)
 * Tests for /gaming placeholder page
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';

vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => vi.fn(),
  useLocation: () => ({ pathname: '/gaming', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

// This import will fail — component doesn't exist yet
import GamingPage from '../gaming/index';

describe('GamingPage (Epic 2.1)', () => {
  it('should render Gaming heading', () => {
    render(<GamingPage />);
    expect(screen.getByRole('heading', { name: /gaming/i })).toBeInTheDocument();
  });

  it('should show Phase 3 message', () => {
    render(<GamingPage />);
    expect(screen.getByText(/coming in phase 3/i)).toBeInTheDocument();
  });
});
