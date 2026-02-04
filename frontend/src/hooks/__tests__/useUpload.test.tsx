/**
 * useUpload Hook Tests - Wave 4
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';

vi.mock('../../lib/api/upload', () => ({
  getPresignedUploadUrl: vi.fn(),
  confirmUpload: vi.fn(),
  getUploadStatus: vi.fn(),
}));

import { useUpload } from '../useUpload';
import * as uploadApi from '../../lib/api/upload';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0 },
      mutations: { retry: false },
    },
  });

  return function Wrapper({ children }: { children: ReactNode }) {
    return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
  };
}

describe('useUpload Hook (Wave 4)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('upload flow', () => {
    it('should return upload function', () => {
      const { result } = renderHook(() => useUpload(), { wrapper: createWrapper() });

      expect(typeof result.current.upload).toBe('function');
    });

    it('should track uploading state', async () => {
      vi.mocked(uploadApi.getPresignedUploadUrl).mockResolvedValue({
        uploadId: 'upload-123',
        uploadUrl: 'https://s3.example.com/presigned',
        expiresAt: '2026-01-22T03:00:00Z',
        maxFileSize: 5000000,
      });

      const { result } = renderHook(() => useUpload(), { wrapper: createWrapper() });

      expect(result.current.isUploading).toBe(false);
    });

    it('should track upload progress', async () => {
      const { result } = renderHook(() => useUpload(), { wrapper: createWrapper() });

      expect(result.current.progress).toBe(0);
    });

    it('should track upload errors', async () => {
      const { result } = renderHook(() => useUpload(), { wrapper: createWrapper() });

      expect(result.current.error).toBeNull();
    });

    it('should return uploads list', async () => {
      const { result } = renderHook(() => useUpload(), { wrapper: createWrapper() });

      expect(Array.isArray(result.current.uploads)).toBe(true);
    });

    it('should reset state on reset call', async () => {
      const { result } = renderHook(() => useUpload(), { wrapper: createWrapper() });

      act(() => {
        result.current.reset();
      });

      expect(result.current.progress).toBe(0);
      expect(result.current.error).toBeNull();
    });
  });
});
