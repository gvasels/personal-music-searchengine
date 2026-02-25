/**
 * ReprocessButton Component Tests - Admin Track Reprocess Feature (TDD Red Phase)
 *
 * Tests for the admin-only button that triggers AI reprocessing of track metadata.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@/test/test-utils';
import { useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';

// Mock useMutation
vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual('@tanstack/react-query');
  return {
    ...actual,
    useMutation: vi.fn(() => ({
      mutate: vi.fn(),
      isPending: false,
      isError: false,
      error: null,
    })),
  };
});

// Mock toast
vi.mock('react-hot-toast', () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
  },
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

// Mock reprocessTrack API function
vi.mock('@/lib/api/tracks', () => ({
  reprocessTrack: vi.fn(),
}));

describe('ReprocessButton Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders reprocess button with icon and label', async () => {
    vi.mocked(useMutation).mockReturnValue({
      mutate: vi.fn(),
      mutateAsync: vi.fn(),
      isPending: false,
      isError: false,
      error: null,
      isIdle: true,
      isSuccess: false,
      isPaused: false,
      data: undefined,
      variables: undefined,
      failureCount: 0,
      failureReason: null,
      reset: vi.fn(),
      status: 'idle',
      context: undefined,
      submittedAt: 0,
    });

    // Import the component dynamically to ensure mocks are in place
    const { ReprocessButton } = await import('../ReprocessButton');
    render(<ReprocessButton trackId="track-1" />);

    // Should have a button element
    expect(screen.getByRole('button')).toBeInTheDocument();
    // Should have accessible label for reprocess action
    expect(screen.getByLabelText(/reprocess/i)).toBeInTheDocument();
  });

  it('shows loading spinner during reprocessing', async () => {
    vi.mocked(useMutation).mockReturnValue({
      mutate: vi.fn(),
      mutateAsync: vi.fn(),
      isPending: true,
      isError: false,
      error: null,
      isIdle: false,
      isSuccess: false,
      isPaused: false,
      data: undefined,
      variables: undefined,
      failureCount: 0,
      failureReason: null,
      reset: vi.fn(),
      status: 'pending',
      context: undefined,
      submittedAt: Date.now(),
    });

    const { ReprocessButton } = await import('../ReprocessButton');
    render(<ReprocessButton trackId="track-1" />);

    const button = screen.getByRole('button');
    // DaisyUI loading state adds 'loading' class or shows spinner
    expect(button).toHaveClass('loading');
  });

  it('disables button during reprocessing', async () => {
    vi.mocked(useMutation).mockReturnValue({
      mutate: vi.fn(),
      mutateAsync: vi.fn(),
      isPending: true,
      isError: false,
      error: null,
      isIdle: false,
      isSuccess: false,
      isPaused: false,
      data: undefined,
      variables: undefined,
      failureCount: 0,
      failureReason: null,
      reset: vi.fn(),
      status: 'pending',
      context: undefined,
      submittedAt: Date.now(),
    });

    const { ReprocessButton } = await import('../ReprocessButton');
    render(<ReprocessButton trackId="track-1" />);

    const button = screen.getByRole('button');
    expect(button).toBeDisabled();
  });

  it('calls onSuccess callback after successful reprocess', async () => {
    const onSuccessMock = vi.fn();
    const mutateMock = vi.fn();

    // Capture the onSuccess callback from useMutation
    let capturedOnSuccess: (() => void) | undefined;
    vi.mocked(useMutation).mockImplementation((options: { onSuccess?: () => void }) => {
      capturedOnSuccess = options?.onSuccess;
      return {
        mutate: mutateMock,
        mutateAsync: vi.fn(),
        isPending: false,
        isError: false,
        error: null,
        isIdle: true,
        isSuccess: false,
        isPaused: false,
        data: undefined,
        variables: undefined,
        failureCount: 0,
        failureReason: null,
        reset: vi.fn(),
        status: 'idle',
        context: undefined,
        submittedAt: 0,
      };
    });

    const { ReprocessButton } = await import('../ReprocessButton');
    const { user } = render(<ReprocessButton trackId="track-1" onSuccess={onSuccessMock} />);

    // Click the button to trigger mutation
    const button = screen.getByRole('button');
    await user.click(button);

    // Verify mutate was called
    expect(mutateMock).toHaveBeenCalled();

    // Simulate success callback
    if (capturedOnSuccess) {
      capturedOnSuccess();
    }

    // onSuccess prop should be called
    await waitFor(() => {
      expect(onSuccessMock).toHaveBeenCalled();
    });
  });

  it('shows error toast on failure', async () => {
    const mutateMock = vi.fn();

    // Capture the onError callback from useMutation
    let capturedOnError: ((error: Error) => void) | undefined;
    vi.mocked(useMutation).mockImplementation((options: { onError?: (error: Error) => void }) => {
      capturedOnError = options?.onError;
      return {
        mutate: mutateMock,
        mutateAsync: vi.fn(),
        isPending: false,
        isError: false,
        error: null,
        isIdle: true,
        isSuccess: false,
        isPaused: false,
        data: undefined,
        variables: undefined,
        failureCount: 0,
        failureReason: null,
        reset: vi.fn(),
        status: 'idle',
        context: undefined,
        submittedAt: 0,
      };
    });

    const { ReprocessButton } = await import('../ReprocessButton');
    const { user } = render(<ReprocessButton trackId="track-1" />);

    // Click the button to trigger mutation
    const button = screen.getByRole('button');
    await user.click(button);

    // Simulate error callback
    if (capturedOnError) {
      capturedOnError(new Error('Reprocess failed'));
    }

    // Error toast should be shown
    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith(expect.stringMatching(/failed|error/i));
    });
  });

  it('shows success toast on successful reprocess', async () => {
    const mutateMock = vi.fn();

    // Capture the onSuccess callback from useMutation
    let capturedOnSuccess: (() => void) | undefined;
    vi.mocked(useMutation).mockImplementation((options: { onSuccess?: () => void }) => {
      capturedOnSuccess = options?.onSuccess;
      return {
        mutate: mutateMock,
        mutateAsync: vi.fn(),
        isPending: false,
        isError: false,
        error: null,
        isIdle: true,
        isSuccess: false,
        isPaused: false,
        data: undefined,
        variables: undefined,
        failureCount: 0,
        failureReason: null,
        reset: vi.fn(),
        status: 'idle',
        context: undefined,
        submittedAt: 0,
      };
    });

    const { ReprocessButton } = await import('../ReprocessButton');
    const { user } = render(<ReprocessButton trackId="track-1" />);

    // Click the button to trigger mutation
    const button = screen.getByRole('button');
    await user.click(button);

    // Simulate success callback
    if (capturedOnSuccess) {
      capturedOnSuccess();
    }

    // Success toast should be shown
    await waitFor(() => {
      expect(toast.success).toHaveBeenCalledWith(expect.stringMatching(/reprocess|success|started/i));
    });
  });

  it('calls reprocessTrack API with correct trackId', async () => {
    await import('@/lib/api/tracks');
    const mutateMock = vi.fn();

    vi.mocked(useMutation).mockImplementation(() => ({
      mutate: mutateMock,
      mutateAsync: vi.fn(),
      isPending: false,
      isError: false,
      error: null,
      isIdle: true,
      isSuccess: false,
      isPaused: false,
      data: undefined,
      variables: undefined,
      failureCount: 0,
      failureReason: null,
      reset: vi.fn(),
      status: 'idle',
      context: undefined,
      submittedAt: 0,
    }));

    const { ReprocessButton } = await import('../ReprocessButton');
    const { user } = render(<ReprocessButton trackId="track-123" />);

    // Click the button
    const button = screen.getByRole('button');
    await user.click(button);

    // Verify mutate was called (the mutation should call reprocessTrack internally)
    expect(mutateMock).toHaveBeenCalled();
  });
});
