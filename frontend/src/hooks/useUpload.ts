/**
 * useUpload Hook - Wave 4
 */
import { useState, useCallback } from 'react';
import { getPresignedUploadUrl, confirmUpload, getUploadStatus, UploadSteps } from '../lib/api/upload';

export interface UploadItem {
  id: string;
  filename: string;
  status: 'pending' | 'uploading' | 'processing' | 'completed' | 'failed';
  progress: number;
  currentStep: string;
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

// Calculate progress percentage based on completed steps (50-100%)
function calculateProcessingProgress(steps: UploadSteps): { progress: number; currentStep: string } {
  const stepList = [
    { key: 'metadataExtracted', label: 'Extracting metadata', weight: 10 },
    { key: 'coverArtExtracted', label: 'Processing cover art', weight: 10 },
    { key: 'trackCreated', label: 'Creating track record', weight: 10 },
    { key: 'fileMoved', label: 'Moving file to storage', weight: 10 },
    { key: 'indexed', label: 'Indexing for search', weight: 10 },
  ] as const;

  let completedWeight = 0;
  let currentStep = 'Processing...';

  for (const step of stepList) {
    if (steps[step.key]) {
      completedWeight += step.weight;
    } else {
      currentStep = step.label;
      break;
    }
  }

  // If all steps completed
  if (completedWeight === 50) {
    currentStep = 'Completed';
  }

  // Processing is 50-100%, so add 50 base
  return { progress: 50 + completedWeight, currentStep };
}

// Upload file to S3 with progress tracking using XHR
function uploadToS3WithProgress(
  url: string,
  file: File,
  onProgress: (percent: number) => void
): Promise<void> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();

    xhr.upload.addEventListener('progress', (event) => {
      if (event.lengthComputable) {
        // Upload progress is 0-50% of total
        const percent = Math.round((event.loaded / event.total) * 50);
        onProgress(percent);
      }
    });

    xhr.addEventListener('load', () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        resolve();
      } else {
        reject(new Error(`Upload to S3 failed: ${xhr.status}`));
      }
    });

    xhr.addEventListener('error', () => {
      reject(new Error('Network error during upload'));
    });

    xhr.addEventListener('abort', () => {
      reject(new Error('Upload aborted'));
    });

    xhr.open('PUT', url);
    xhr.setRequestHeader('Content-Type', file.type);
    xhr.send(file);
  });
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
          currentStep: 'Uploading...',
          trackId: null,
          error: null,
        };

        setUploads((prev) => [...prev, uploadItem]);

        // Upload to S3 with progress tracking
        await uploadToS3WithProgress(presigned.uploadUrl, file, (percent) => {
          setProgress(percent);
          setUploads((prev) =>
            prev.map((u) =>
              u.id === presigned.uploadId
                ? { ...u, progress: percent, currentStep: `Uploading... ${percent}%` }
                : u
            )
          );
        });

        // Update to processing state
        setProgress(50);
        setUploads((prev) =>
          prev.map((u) =>
            u.id === presigned.uploadId
              ? { ...u, progress: 50, status: 'processing', currentStep: 'Starting processing...' }
              : u
          )
        );

        // Confirm upload (triggers Step Functions)
        await confirmUpload(presigned.uploadId);

        // Poll for status with step progress
        let status = await getUploadStatus(presigned.uploadId);
        while (status.status === 'PROCESSING') {
          const { progress: stepProgress, currentStep } = calculateProcessingProgress(status.steps);

          setProgress(stepProgress);
          setUploads((prev) =>
            prev.map((u) =>
              u.id === presigned.uploadId
                ? { ...u, progress: stepProgress, currentStep }
                : u
            )
          );

          await new Promise((resolve) => setTimeout(resolve, 1500));
          status = await getUploadStatus(presigned.uploadId);
        }

        // Final status update
        const finalStatus = status.status === 'COMPLETED' ? 'completed' : 'failed';
        setUploads((prev) =>
          prev.map((u) =>
            u.id === presigned.uploadId
              ? {
                  ...u,
                  status: finalStatus,
                  progress: finalStatus === 'completed' ? 100 : u.progress,
                  currentStep: finalStatus === 'completed' ? 'Completed' : 'Failed',
                  trackId: status.trackId,
                  error: status.errorMsg,
                }
              : u
          )
        );
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Upload failed';
        setError(message);
        setUploads((prev) =>
          prev.map((u) =>
            u.filename === file.name
              ? { ...u, status: 'failed', error: message, currentStep: 'Failed' }
              : u
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
