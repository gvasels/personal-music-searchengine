import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { UserDetailModal } from '../UserDetailModal';
import type { UserRole } from '../../../types';

// Mock the useAdmin hooks
vi.mock('../../../hooks/useAdmin', () => ({
  useUserDetails: vi.fn(),
  useUpdateUserRole: vi.fn(),
  useUpdateUserStatus: vi.fn(),
}));

import { useUserDetails, useUpdateUserRole, useUpdateUserStatus } from '../../../hooks/useAdmin';

interface UserDetails {
  id: string;
  email: string;
  displayName: string;
  role: UserRole;
  disabled: boolean;
  createdAt: string;
  lastLoginAt?: string;
  trackCount?: number;
  playlistCount?: number;
  albumCount?: number;
  storageUsed?: number;
  followerCount?: number;
  followingCount?: number;
}

const createMockUserDetails = (overrides: Partial<UserDetails> = {}): UserDetails => ({
  id: 'user-123',
  email: 'test@example.com',
  displayName: 'Test User',
  role: 'subscriber',
  disabled: false,
  createdAt: '2024-01-15T10:00:00Z',
  lastLoginAt: '2024-06-15T14:30:00Z',
  trackCount: 42,
  playlistCount: 5,
  albumCount: 3,
  storageUsed: 1073741824, // 1 GB
  followerCount: 100,
  followingCount: 50,
  ...overrides,
});

const createQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

const renderModal = (props: { isOpen: boolean; userId: string | null; onClose: () => void }) => {
  const queryClient = createQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <UserDetailModal {...props} />
    </QueryClientProvider>
  );
};

describe('UserDetailModal', () => {
  const mockOnClose = vi.fn();
  const mockMutateFn = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();

    // Default mock implementations
    vi.mocked(useUserDetails).mockReturnValue({
      data: undefined,
      isLoading: false,
      error: null,
    } as ReturnType<typeof useUserDetails>);

    vi.mocked(useUpdateUserRole).mockReturnValue({
      mutateAsync: mockMutateFn,
      isPending: false,
    } as unknown as ReturnType<typeof useUpdateUserRole>);

    vi.mocked(useUpdateUserStatus).mockReturnValue({
      mutateAsync: mockMutateFn,
      isPending: false,
    } as unknown as ReturnType<typeof useUpdateUserStatus>);
  });

  describe('modal visibility', () => {
    it('returns null when isOpen is false', () => {
      renderModal({ isOpen: false, userId: 'user-123', onClose: mockOnClose });
      expect(screen.queryByText('User Details')).not.toBeInTheDocument();
    });

    it('renders dialog when isOpen is true', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails(),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('User Details')).toBeInTheDocument();
    });
  });

  describe('loading state', () => {
    it('shows loading spinner when data is loading', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: undefined,
        isLoading: true,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('User Details')).toBeInTheDocument();
      expect(document.querySelector('.loading-spinner')).toBeInTheDocument();
    });
  });

  describe('error state', () => {
    it('displays error alert when API fails', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: undefined,
        isLoading: false,
        error: new Error('Network error'),
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('Failed to load user details')).toBeInTheDocument();
    });
  });

  describe('user information display', () => {
    it('renders user display name', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ displayName: 'John Doe' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    it('renders "No display name" when displayName is empty', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ displayName: '' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('No display name')).toBeInTheDocument();
    });

    it('renders user email', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ email: 'john@example.com' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('john@example.com')).toBeInTheDocument();
    });

    it('renders role badge', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ role: 'admin' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('Admin')).toBeInTheDocument();
    });

    it('renders disabled badge when user is disabled', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ disabled: true }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('Disabled')).toBeInTheDocument();
    });
  });

  describe('stats display', () => {
    it('renders all stat values correctly', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({
          trackCount: 42,
          playlistCount: 5,
          albumCount: 3,
          followerCount: 100,
          followingCount: 50,
          storageUsed: 1073741824, // 1 GB = 1024 MB
        }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });

      expect(screen.getByText('Tracks')).toBeInTheDocument();
      expect(screen.getByText('42')).toBeInTheDocument();
      expect(screen.getByText('Playlists')).toBeInTheDocument();
      expect(screen.getByText('5')).toBeInTheDocument();
      expect(screen.getByText('Albums')).toBeInTheDocument();
      expect(screen.getByText('3')).toBeInTheDocument();
      expect(screen.getByText('Followers')).toBeInTheDocument();
      expect(screen.getByText('100')).toBeInTheDocument();
      expect(screen.getByText('Following')).toBeInTheDocument();
      expect(screen.getByText('50')).toBeInTheDocument();
      expect(screen.getByText('Storage')).toBeInTheDocument();
      // Storage value may be rendered with separate text nodes
      expect(screen.getByText(/1024\.0/)).toBeInTheDocument();
    });

    it('handles null/undefined stats gracefully with fallback to 0', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({
          trackCount: undefined,
          playlistCount: undefined,
          albumCount: undefined,
          followerCount: undefined,
          followingCount: undefined,
          storageUsed: undefined,
        }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });

      // All stats should show 0 for undefined values
      const statValues = screen.getAllByText('0');
      expect(statValues.length).toBeGreaterThanOrEqual(5);

      // Storage should show 0.0 MB
      expect(screen.getByText('0.0 MB')).toBeInTheDocument();
    });

    it('handles zero values correctly', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({
          trackCount: 0,
          playlistCount: 0,
          albumCount: 0,
          followerCount: 0,
          followingCount: 0,
          storageUsed: 0,
        }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });

      // Should show 0 values (minimum 5 stats with 0)
      const statValues = screen.getAllByText('0');
      expect(statValues.length).toBeGreaterThanOrEqual(5);
      expect(screen.getByText('0.0 MB')).toBeInTheDocument();
    });
  });

  describe('date display', () => {
    it('renders join date', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ createdAt: '2024-01-15T10:00:00Z' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText(/joined/i)).toBeInTheDocument();
    });

    it('renders last login date when available', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ lastLoginAt: '2024-06-15T14:30:00Z' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText(/last login/i)).toBeInTheDocument();
    });

    it('does not render last login when not available', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ lastLoginAt: undefined }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.queryByText(/last login/i)).not.toBeInTheDocument();
    });
  });

  describe('close behavior', () => {
    it('calls onClose when close button is clicked', async () => {
      const user = userEvent.setup();
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails(),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });

      // Use getByText to find the Close button in modal-action footer
      const closeButton = screen.getByText('Close');
      await user.click(closeButton);
      expect(mockOnClose).toHaveBeenCalledTimes(1);
    });

    it('calls onClose when X button is clicked', async () => {
      const user = userEvent.setup();
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails(),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });

      // Find the X button by aria-label
      const xButton = screen.getByLabelText('Close');
      await user.click(xButton);
      expect(mockOnClose).toHaveBeenCalledTimes(1);
    });
  });

  describe('no user selected', () => {
    it('shows placeholder when no userId provided', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: undefined,
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: null, onClose: mockOnClose });
      expect(screen.getByText('Select a user to view details')).toBeInTheDocument();
    });
  });

  describe('role selector', () => {
    it('renders role selector with current role', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ role: 'subscriber' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });

      // Find the select by the label text
      expect(screen.getByText('User Role')).toBeInTheDocument();
      const roleSelector = document.querySelector('select.select');
      expect(roleSelector).toBeInTheDocument();
      expect(roleSelector).toHaveValue('subscriber');
    });
  });

  describe('status toggle', () => {
    it('renders status toggle for enabled user', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ disabled: false }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('Account Active')).toBeInTheDocument();
    });

    it('renders status toggle for disabled user', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ disabled: true }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('Account Disabled')).toBeInTheDocument();
    });
  });

  describe('avatar', () => {
    it('renders avatar with first letter of display name', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ displayName: 'John Doe' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('J')).toBeInTheDocument();
    });

    it('renders avatar with first letter of email when no display name', () => {
      vi.mocked(useUserDetails).mockReturnValue({
        data: createMockUserDetails({ displayName: '', email: 'test@example.com' }),
        isLoading: false,
        error: null,
      } as ReturnType<typeof useUserDetails>);

      renderModal({ isOpen: true, userId: 'user-123', onClose: mockOnClose });
      expect(screen.getByText('T')).toBeInTheDocument();
    });
  });
});
