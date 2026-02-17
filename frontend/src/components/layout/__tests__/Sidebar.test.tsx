/**
 * Sidebar Tests - Epic 2.1 Task 2.1 (TDD Red → Green)
 * Tests section-based navigation with Music sub-nav
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';

const mockUseAuth = vi.fn();
vi.mock('../../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

const mockUseFeatureFlags = vi.fn();
vi.mock('../../../hooks/useFeatureFlags', () => ({
  useFeatureFlags: () => mockUseFeatureFlags(),
}));

const mockUseLocation = vi.fn();
vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to, ...props }: { children: React.ReactNode; to: string; [key: string]: unknown }) => (
    <a href={to} {...props}>{children}</a>
  ),
  useLocation: () => mockUseLocation(),
}));

import { Sidebar } from '../Sidebar';

describe('Sidebar (Epic 2.1 - Section Navigation)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue({ isAuthenticated: true, isLoading: false });
    mockUseFeatureFlags.mockReturnValue({
      role: 'subscriber',
      hasRole: (r: string) => ['guest', 'subscriber'].includes(r),
      isSimulating: false,
    });
  });

  it('should render Home, Music, Videos, Gaming sections', () => {
    mockUseLocation.mockReturnValue({ pathname: '/', search: '', hash: '' });
    render(<Sidebar />);
    expect(screen.getByText('Home')).toBeInTheDocument();
    expect(screen.getByText('Music')).toBeInTheDocument();
    expect(screen.getByText('Videos')).toBeInTheDocument();
    expect(screen.getByText('Gaming')).toBeInTheDocument();
  });

  it('should show music sub-nav when on /music path', () => {
    mockUseLocation.mockReturnValue({ pathname: '/music/tracks', search: '', hash: '' });
    render(<Sidebar />);
    expect(screen.getByText('Tracks')).toBeInTheDocument();
    expect(screen.getByText('Albums')).toBeInTheDocument();
    expect(screen.getByText('Artists')).toBeInTheDocument();
    expect(screen.getByText('Playlists')).toBeInTheDocument();
    expect(screen.getByText('Tags')).toBeInTheDocument();
  });

  it('should hide music sub-nav when on /videos path', () => {
    mockUseLocation.mockReturnValue({ pathname: '/videos', search: '', hash: '' });
    render(<Sidebar />);
    expect(screen.getByText('Music')).toBeInTheDocument();
    expect(screen.queryByText('Tracks')).not.toBeInTheDocument();
  });

  it('should show admin section for admin role', () => {
    mockUseFeatureFlags.mockReturnValue({
      role: 'admin',
      hasRole: () => true,
      isSimulating: false,
    });
    mockUseLocation.mockReturnValue({ pathname: '/', search: '', hash: '' });
    render(<Sidebar />);
    expect(screen.getByText('User Management')).toBeInTheDocument();
  });

  it('should not render when not authenticated', () => {
    mockUseAuth.mockReturnValue({ isAuthenticated: false, isLoading: false });
    mockUseLocation.mockReturnValue({ pathname: '/', search: '', hash: '' });
    const { container } = render(<Sidebar />);
    expect(container.innerHTML).toBe('');
  });
});
