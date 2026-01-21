/**
 * Layout Component Tests - REQ-5.2
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@/test/test-utils';

vi.mock('@tanstack/react-router', () => ({
  Link: vi.fn(({ children, to, ...props }) => (
    <a href={to} {...props}>{children}</a>
  )),
  useNavigate: vi.fn(() => vi.fn()),
}));

vi.mock('@/lib/store/themeStore', () => ({
  useThemeStore: vi.fn(() => ({ theme: 'dark', toggleTheme: vi.fn() })),
}));

vi.mock('@/lib/store/playerStore', () => ({
  usePlayerStore: vi.fn(() => ({
    currentTrack: null,
    isPlaying: false,
    volume: 1,
    progress: 0,
    play: vi.fn(),
    pause: vi.fn(),
    next: vi.fn(),
    previous: vi.fn(),
    setVolume: vi.fn(),
    seek: vi.fn(),
  })),
}));

describe('Layout Component', () => {
  it('REQ-5.2: should render app shell', async () => {
    const { Layout } = await import('@/components/layout/Layout');
    render(<Layout><div>Content</div></Layout>);
    expect(screen.getByText('Content')).toBeInTheDocument();
  });

  it('should render Header', async () => {
    const { Layout } = await import('@/components/layout/Layout');
    render(<Layout><div>Content</div></Layout>);
    expect(screen.getByRole('banner')).toBeInTheDocument();
  });

  it('should render Sidebar', async () => {
    const { Layout } = await import('@/components/layout/Layout');
    render(<Layout><div>Content</div></Layout>);
    expect(screen.getByRole('navigation')).toBeInTheDocument();
  });

  it('should render PlayerBar', async () => {
    const { Layout } = await import('@/components/layout/Layout');
    render(<Layout><div>Content</div></Layout>);
    expect(screen.getByTestId('player-bar')).toBeInTheDocument();
  });

  it('REQ-5.2: should be responsive', async () => {
    const { Layout } = await import('@/components/layout/Layout');
    const { container } = render(<Layout><div>Content</div></Layout>);
    expect(container.querySelector('.min-h-screen')).toBeInTheDocument();
  });
});
