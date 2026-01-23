/**
 * Track Detail Page Tests - Wave 2
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseTrack = vi.fn();
const mockUpdateTrack = vi.fn();
const mockDeleteTrack = vi.fn();
vi.mock('../../hooks/useTracks', () => ({
  useTrackQuery: () => mockUseTrack(),
  useUpdateTrack: () => ({ mutateAsync: mockUpdateTrack, isPending: false }),
  useDeleteTrack: () => ({ mutateAsync: mockDeleteTrack, isPending: false }),
  trackKeys: { detail: (id: string) => ['tracks', 'detail', id] },
}));

const mockNavigate = vi.fn();
vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => mockNavigate,
  useParams: () => ({ trackId: 'track-1' }),
  useLocation: () => ({ pathname: '/tracks/track-1', search: '', hash: '' }),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
}));

import TrackDetailPage from '../tracks/$trackId';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('TrackDetailPage (Wave 2)', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockTrack = {
    id: 'track-1',
    title: 'Test Track',
    artist: 'Test Artist',
    album: 'Test Album',
    albumId: 'album-1',
    duration: 180,
    format: 'mp3',
    fileSize: 5000000,
    tags: ['rock', 'indie'],
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  };

  describe('Loading state', () => {
    it('should show loading spinner while fetching', () => {
      mockUseTrack.mockReturnValue({ isLoading: true, data: undefined, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Error state', () => {
    it('should show error message when track not found', () => {
      mockUseTrack.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Track not found'),
      });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/track not found/i)).toBeInTheDocument();
    });

    it('should show back button on error', () => {
      mockUseTrack.mockReturnValue({
        isLoading: false,
        data: undefined,
        isError: true,
        error: new Error('Not found'),
      });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /back/i })).toBeInTheDocument();
    });
  });

  describe('Track details display', () => {
    it('should render track title', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('heading', { name: 'Test Track' })).toBeInTheDocument();
    });

    it('should render artist name with link', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('link', { name: 'Test Artist' })).toBeInTheDocument();
    });

    it('should render album name with link', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('link', { name: 'Test Album' })).toBeInTheDocument();
    });

    it('should render duration', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText('3:00')).toBeInTheDocument();
    });

    it('should render format', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/mp3/i)).toBeInTheDocument();
    });

    it('should render tags', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByText('rock')).toBeInTheDocument();
      expect(screen.getByText('indie')).toBeInTheDocument();
    });
  });

  describe('Edit mode', () => {
    it('should show edit button', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /edit/i })).toBeInTheDocument();
    });

    it('should show edit form when edit clicked', async () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /edit/i }));

      expect(screen.getByLabelText(/title/i)).toBeInTheDocument();
    });

    it('should submit updates', async () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });
      mockUpdateTrack.mockResolvedValue({ ...mockTrack, title: 'Updated Title' });

      render(<TrackDetailPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /edit/i }));
      await user.clear(screen.getByLabelText(/title/i));
      await user.type(screen.getByLabelText(/title/i), 'Updated Title');
      await user.click(screen.getByRole('button', { name: /save/i }));

      await waitFor(() => {
        expect(mockUpdateTrack).toHaveBeenCalledWith({ id: 'track-1', data: expect.objectContaining({ title: 'Updated Title' }) });
      });
    });
  });

  describe('Delete functionality', () => {
    it('should show delete button', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument();
    });

    it('should show confirmation dialog when delete clicked', async () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /delete/i }));

      expect(screen.getByText(/are you sure/i)).toBeInTheDocument();
    });

    it('should delete track and navigate on confirm', async () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });
      mockDeleteTrack.mockResolvedValue(undefined);

      render(<TrackDetailPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /delete/i }));
      await user.click(screen.getByRole('button', { name: /confirm/i }));

      await waitFor(() => {
        expect(mockDeleteTrack).toHaveBeenCalledWith('track-1');
        expect(mockNavigate).toHaveBeenCalledWith({ to: '/tracks' });
      });
    });
  });

  describe('Player integration', () => {
    it('should show play button', () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /play/i })).toBeInTheDocument();
    });
  });

  describe('Navigation', () => {
    it('should have back button', async () => {
      mockUseTrack.mockReturnValue({ isLoading: false, data: mockTrack, isError: false });

      render(<TrackDetailPage />, { wrapper: createWrapper() });
      await user.click(screen.getByRole('button', { name: /back/i }));

      expect(mockNavigate).toHaveBeenCalledWith({ to: '/tracks' });
    });
  });
});
