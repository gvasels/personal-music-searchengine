/**
 * Upload API Tests - Wave 4
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    post: vi.fn(),
    get: vi.fn(),
  },
}));

import { getPresignedUploadUrl, confirmUpload, getUploadStatus } from '../upload';
import { apiClient } from '../client';

describe('Upload API (Wave 4)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getPresignedUploadUrl', () => {
    it('should request presigned URL for file upload', async () => {
      const mockResponse = {
        uploadId: 'upload-123',
        uploadUrl: 'https://s3.amazonaws.com/presigned-url',
        expiresAt: '2026-01-22T03:00:00Z',
        maxFileSize: 5000000,
      };
      vi.mocked(apiClient.post).mockResolvedValue({ data: mockResponse });

      const result = await getPresignedUploadUrl({
        filename: 'song.mp3',
        contentType: 'audio/mpeg',
        fileSize: 5000000,
      });

      expect(apiClient.post).toHaveBeenCalledWith('/upload/presigned', {
        filename: 'song.mp3',
        contentType: 'audio/mpeg',
        fileSize: 5000000,
      });
      expect(result.uploadId).toBe('upload-123');
      expect(result.uploadUrl).toBeTruthy();
    });
  });

  describe('confirmUpload', () => {
    it('should confirm upload completion', async () => {
      const mockResponse = { status: 'processing', trackId: null };
      vi.mocked(apiClient.post).mockResolvedValue({ data: mockResponse });

      const result = await confirmUpload('upload-123');

      expect(apiClient.post).toHaveBeenCalledWith('/upload/confirm', { uploadId: 'upload-123' });
      expect(result.status).toBe('processing');
    });
  });

  describe('getUploadStatus', () => {
    it('should get upload processing status', async () => {
      const mockResponse = {
        uploadId: 'upload-123',
        status: 'completed',
        trackId: 'track-456',
        error: null,
      };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse });

      const result = await getUploadStatus('upload-123');

      expect(apiClient.get).toHaveBeenCalledWith('/uploads/upload-123');
      expect(result.status).toBe('completed');
      expect(result.trackId).toBe('track-456');
    });

    it('should return error status on failure', async () => {
      const mockResponse = {
        uploadId: 'upload-123',
        status: 'failed',
        trackId: null,
        error: 'Invalid audio format',
      };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse });

      const result = await getUploadStatus('upload-123');

      expect(result.status).toBe('failed');
      expect(result.error).toBe('Invalid audio format');
    });
  });
});
