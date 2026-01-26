import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { UserCard } from '../UserCard';
import type { UserRole } from '../../../types';

interface UserSummary {
  id: string;
  email: string;
  displayName: string;
  role: UserRole;
  disabled: boolean;
  createdAt: string;
}

const createMockUser = (overrides: Partial<UserSummary> = {}): UserSummary => ({
  id: 'user-123',
  email: 'test@example.com',
  displayName: 'Test User',
  role: 'subscriber',
  disabled: false,
  createdAt: '2024-01-15T10:00:00Z',
  ...overrides,
});

describe('UserCard', () => {
  const mockOnSelect = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('rendering', () => {
    it('renders user display name', () => {
      render(<UserCard user={createMockUser({ displayName: 'John Doe' })} onSelect={mockOnSelect} />);
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    it('renders fallback when no display name', () => {
      render(<UserCard user={createMockUser({ displayName: '' })} onSelect={mockOnSelect} />);
      expect(screen.getByText('No display name')).toBeInTheDocument();
    });

    it('renders user email', () => {
      render(<UserCard user={createMockUser({ email: 'john@example.com' })} onSelect={mockOnSelect} />);
      expect(screen.getByText('john@example.com')).toBeInTheDocument();
    });

    it('renders formatted join date', () => {
      render(<UserCard user={createMockUser({ createdAt: '2024-06-15T10:00:00Z' })} onSelect={mockOnSelect} />);
      expect(screen.getByText(/joined/i)).toBeInTheDocument();
    });
  });

  describe('role badges', () => {
    it('renders guest badge with ghost style', () => {
      render(<UserCard user={createMockUser({ role: 'guest' })} onSelect={mockOnSelect} />);
      const badge = screen.getByText('Guest');
      expect(badge).toHaveClass('badge-ghost');
    });

    it('renders subscriber badge with info style', () => {
      render(<UserCard user={createMockUser({ role: 'subscriber' })} onSelect={mockOnSelect} />);
      const badge = screen.getByText('Subscriber');
      expect(badge).toHaveClass('badge-info');
    });

    it('renders artist badge with secondary style', () => {
      render(<UserCard user={createMockUser({ role: 'artist' })} onSelect={mockOnSelect} />);
      const badge = screen.getByText('Artist');
      expect(badge).toHaveClass('badge-secondary');
    });

    it('renders admin badge with primary style', () => {
      render(<UserCard user={createMockUser({ role: 'admin' })} onSelect={mockOnSelect} />);
      const badge = screen.getByText('Admin');
      expect(badge).toHaveClass('badge-primary');
    });
  });

  describe('disabled state', () => {
    it('shows disabled badge when user is disabled', () => {
      render(<UserCard user={createMockUser({ disabled: true })} onSelect={mockOnSelect} />);
      expect(screen.getByText('Disabled')).toBeInTheDocument();
    });

    it('does not show disabled badge when user is enabled', () => {
      render(<UserCard user={createMockUser({ disabled: false })} onSelect={mockOnSelect} />);
      expect(screen.queryByText('Disabled')).not.toBeInTheDocument();
    });

    it('applies opacity class when disabled', () => {
      render(<UserCard user={createMockUser({ disabled: true })} onSelect={mockOnSelect} />);
      const button = screen.getByRole('button');
      expect(button).toHaveClass('opacity-60');
    });
  });

  describe('selection', () => {
    it('applies ring style when selected', () => {
      render(<UserCard user={createMockUser()} onSelect={mockOnSelect} isSelected={true} />);
      const button = screen.getByRole('button');
      expect(button).toHaveClass('ring-2', 'ring-primary');
    });

    it('does not apply ring style when not selected', () => {
      render(<UserCard user={createMockUser()} onSelect={mockOnSelect} isSelected={false} />);
      const button = screen.getByRole('button');
      expect(button).not.toHaveClass('ring-2');
    });
  });

  describe('click handling', () => {
    it('calls onSelect with user id when clicked', async () => {
      const user = userEvent.setup();
      render(<UserCard user={createMockUser({ id: 'user-xyz' })} onSelect={mockOnSelect} />);

      await user.click(screen.getByRole('button'));
      expect(mockOnSelect).toHaveBeenCalledWith('user-xyz');
    });

    it('calls onSelect only once per click', async () => {
      const user = userEvent.setup();
      render(<UserCard user={createMockUser()} onSelect={mockOnSelect} />);

      await user.click(screen.getByRole('button'));
      expect(mockOnSelect).toHaveBeenCalledTimes(1);
    });
  });

  describe('accessibility', () => {
    it('renders as a button', () => {
      render(<UserCard user={createMockUser()} onSelect={mockOnSelect} />);
      expect(screen.getByRole('button')).toBeInTheDocument();
    });

    it('has type button to prevent form submission', () => {
      render(<UserCard user={createMockUser()} onSelect={mockOnSelect} />);
      expect(screen.getByRole('button')).toHaveAttribute('type', 'button');
    });
  });
});
