/**
 * useUpload Hook - Wave 4
 */
import { useState, useCallback } from 'react';
import { getPresignedUploadUrl, confirmUpload, getUploadStatus } from '../lib/api/upload';

export interface UploadItem {
  id: string;
  filename: string;
  status: 'pending' | 'uploading' | 'processing' | 'completed' | 'failed';
  progress: number;
  trackId: string | null;
  error: string | null;
}

export interface UseUploadReturn {
  upload: (files: File[]) => Promise<void>;
  isUploading: boolean;
  progress: number;
  error: string | null;
  uploads: UploadItem[];
  reset: () => void;
}

export function useUpload(): UseUploadReturn {
  const [isUploading, setIsUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [uploads, setUploads] = useState<UploadItem[]>([]);

  const upload = useCallback(async (files: File[]) => {
    setIsUploading(true);
    setError(null);

    for (const file of files) {
      try {
        // Get presigned URL
        const presigned = await getPresignedUploadUrl({
          filename: file.name,
          contentType: file.type,
          fileSize: file.size,
        });

        const uploadItem: UploadItem = {
          id: presigned.uploadId,
          filename: file.name,
          status: 'uploading',
          progress: 0,
          trackId: null,
          error: null,
        };

        setUploads((prev) => [...prev, uploadItem]);

        // Upload to S3
        const uploadResponse = await fetch(presigned.uploadUrl, {
          method: 'PUT',
          body: file,
          headers: {
            'Content-Type': file.type,
          },
        });

        if (!uploadResponse.ok) {
          throw new Error(`Upload to S3 failed: ${uploadResponse.status}`);
        }

        // Update progress
        setProgress(50);
        setUploads((prev) =>
          prev.map((u) =>
            u.id === presigned.uploadId ? { ...u, progress: 50, status: 'processing' } : u
          )
        );

        // Confirm upload
        await confirmUpload(presigned.uploadId);

        // Poll for status
        let status = await getUploadStatus(presigned.uploadId);
        while (status.status === 'processing') {
          await new Promise((resolve) => setTimeout(resolve, 2000));
          status = await getUploadStatus(presigned.uploadId);
        }

        setUploads((prev) =>
          prev.map((u) =>
            u.id === presigned.uploadId
              ? {
                  ...u,
                  status: status.status as UploadItem['status'],
                  progress: 100,
                  trackId: status.trackId,
                  error: status.error,
                }
              : u
          )
        );
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Upload failed';
        setError(message);
        setUploads((prev) =>
          prev.map((u) =>
            u.filename === file.name ? { ...u, status: 'failed', error: message } : u
          )
        );
      }
    }

    setIsUploading(false);
    setProgress(0);
  }, []);

  const reset = useCallback(() => {
    setProgress(0);
    setError(null);
    setUploads([]);
  }, []);

  return {
    upload,
    isUploading,
    progress,
    error,
    uploads,
    reset,
  };
}
