/**
 * Upload Page Tests - Wave 4
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockUseUpload = vi.fn();
vi.mock('../../hooks/useUpload', () => ({
  useUpload: () => mockUseUpload(),
}));

vi.mock('@tanstack/react-router', async () => ({
  useNavigate: () => vi.fn(),
  useLocation: () => ({ pathname: '/upload', search: '', hash: '' }),
  Link: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
}));

import UploadPage from '../upload';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('UploadPage (Wave 4)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseUpload.mockReturnValue({
      upload: vi.fn(),
      isUploading: false,
      progress: 0,
      error: null,
      uploads: [],
      reset: vi.fn(),
    });
  });

  describe('Rendering', () => {
    it('should render page title', () => {
      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('heading', { name: /upload/i })).toBeInTheDocument();
    });

    it('should render upload dropzone', () => {
      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByTestId('dropzone')).toBeInTheDocument();
    });

    it('should render browse button', () => {
      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('button', { name: /browse/i })).toBeInTheDocument();
    });
  });

  describe('Upload state', () => {
    it('should show progress bar when uploading', () => {
      mockUseUpload.mockReturnValue({
        upload: vi.fn(),
        isUploading: true,
        progress: 50,
        error: null,
        uploads: [],
        reset: vi.fn(),
      });

      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('should show error message on upload failure', () => {
      mockUseUpload.mockReturnValue({
        upload: vi.fn(),
        isUploading: false,
        progress: 0,
        error: 'Upload failed',
        uploads: [],
        reset: vi.fn(),
      });

      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByText(/upload failed/i)).toBeInTheDocument();
    });

    it('should disable dropzone while uploading', () => {
      mockUseUpload.mockReturnValue({
        upload: vi.fn(),
        isUploading: true,
        progress: 50,
        error: null,
        uploads: [],
        reset: vi.fn(),
      });

      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByTestId('dropzone')).toHaveAttribute('aria-disabled', 'true');
    });
  });

  describe('Upload list', () => {
    it('should show upload list when there are uploads', () => {
      mockUseUpload.mockReturnValue({
        upload: vi.fn(),
        isUploading: false,
        progress: 0,
        error: null,
        uploads: [
          { id: 'upload-1', filename: 'song1.mp3', status: 'completed', trackId: 'track-1' },
          { id: 'upload-2', filename: 'song2.mp3', status: 'processing', trackId: null },
        ],
        reset: vi.fn(),
      });

      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByText('song1.mp3')).toBeInTheDocument();
      expect(screen.getByText('song2.mp3')).toBeInTheDocument();
    });

    it('should show completed status for finished uploads', () => {
      mockUseUpload.mockReturnValue({
        upload: vi.fn(),
        isUploading: false,
        progress: 0,
        error: null,
        uploads: [
          { id: 'upload-1', filename: 'song.mp3', status: 'completed', trackId: 'track-1' },
        ],
        reset: vi.fn(),
      });

      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Completed')).toBeInTheDocument();
    });

    it('should show processing status for uploads in progress', () => {
      mockUseUpload.mockReturnValue({
        upload: vi.fn(),
        isUploading: false,
        progress: 0,
        error: null,
        uploads: [
          { id: 'upload-1', filename: 'song.mp3', status: 'processing', trackId: null, currentStep: 'Processing metadata', progress: 45 },
        ],
        reset: vi.fn(),
      });

      render(<UploadPage />, { wrapper: createWrapper() });

      expect(screen.getByText('Processing metadata')).toBeInTheDocument();
      expect(screen.getByText('45%')).toBeInTheDocument();
    });
  });
});
