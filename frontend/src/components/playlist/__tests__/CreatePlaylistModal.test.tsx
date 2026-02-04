/**
 * CreatePlaylistModal Component Tests - REQ-5.9
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@/test/test-utils';

describe('CreatePlaylistModal Component', () => {
  it('REQ-5.9: should render modal when open', async () => {
    const { CreatePlaylistModal } = await import('@/components/playlist/CreatePlaylistModal');
    render(<CreatePlaylistModal isOpen onClose={vi.fn()} />);
    expect(screen.getByRole('dialog')).toBeInTheDocument();
  });

  it('should not render when closed', async () => {
    const { CreatePlaylistModal } = await import('@/components/playlist/CreatePlaylistModal');
    render(<CreatePlaylistModal isOpen={false} onClose={vi.fn()} />);
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });

  it('should have name input', async () => {
    const { CreatePlaylistModal } = await import('@/components/playlist/CreatePlaylistModal');
    render(<CreatePlaylistModal isOpen onClose={vi.fn()} />);
    expect(screen.getByLabelText(/name/i)).toBeInTheDocument();
  });

  it('should have description textarea', async () => {
    const { CreatePlaylistModal } = await import('@/components/playlist/CreatePlaylistModal');
    render(<CreatePlaylistModal isOpen onClose={vi.fn()} />);
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument();
  });

  it('REQ-5.9: should require name', async () => {
    const { CreatePlaylistModal } = await import('@/components/playlist/CreatePlaylistModal');
    const { user } = render(<CreatePlaylistModal isOpen onClose={vi.fn()} />);

    await user.click(screen.getByRole('button', { name: /create/i }));
    expect(screen.getByText(/name.*required/i)).toBeInTheDocument();
  });

  it('should close on cancel', async () => {
    const onClose = vi.fn();
    const { CreatePlaylistModal } = await import('@/components/playlist/CreatePlaylistModal');
    const { user } = render(<CreatePlaylistModal isOpen onClose={onClose} />);

    await user.click(screen.getByRole('button', { name: /cancel/i }));
    expect(onClose).toHaveBeenCalled();
  });

  it('should close on Escape key', async () => {
    const onClose = vi.fn();
    const { CreatePlaylistModal } = await import('@/components/playlist/CreatePlaylistModal');
    const { user } = render(<CreatePlaylistModal isOpen onClose={onClose} />);

    await user.keyboard('{Escape}');
    expect(onClose).toHaveBeenCalled();
  });
});
