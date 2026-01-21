import { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';

interface UploadDropzoneProps {
  onFilesSelected: (files: File[]) => void;
  onError?: (error: string) => void;
  progress?: number;
  multiple?: boolean;
  maxSize?: number;
  disabled?: boolean;
}

export function UploadDropzone({
  onFilesSelected,
  onError,
  progress,
  multiple = true,
  maxSize = 500 * 1024 * 1024, // 500MB
  disabled = false,
}: UploadDropzoneProps) {
  const [isDragActive, setIsDragActive] = useState(false);

  const onDrop = useCallback(
    (acceptedFiles: File[], rejectedFiles: any[]) => {
      if (rejectedFiles.length > 0) {
        const firstError = rejectedFiles[0].errors[0];
        if (firstError.code === 'file-too-large') {
          onError?.('File size exceeds maximum allowed');
        } else if (firstError.code === 'file-invalid-type') {
          onError?.('Only audio files are accepted');
        } else {
          onError?.(firstError.message);
        }
        return;
      }
      onFilesSelected(acceptedFiles);
    },
    [onFilesSelected, onError]
  );

  const { getRootProps, getInputProps, open } = useDropzone({
    onDrop,
    accept: {
      'audio/*': ['.mp3', '.flac', '.wav', '.ogg', '.m4a', '.aac'],
    },
    multiple,
    maxSize,
    disabled,
    onDragEnter: () => setIsDragActive(true),
    onDragLeave: () => setIsDragActive(false),
    noClick: true,
  });

  return (
    <div
      {...getRootProps()}
      data-testid="dropzone"
      aria-disabled={disabled}
      className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
        isDragActive
          ? 'border-primary bg-primary/10 drag-active'
          : 'border-base-300 hover:border-primary/50'
      } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
    >
      <input {...getInputProps()} data-testid="file-input" />
      <div className="space-y-4">
        <div className="text-4xl">🎵</div>
        <div>
          <p className="text-lg font-medium">Drag and drop audio files here</p>
          <p className="text-base-content/60">or</p>
        </div>
        <button
          type="button"
          onClick={open}
          disabled={disabled}
          className="btn btn-primary"
          aria-label="browse files"
        >
          Browse Files
        </button>
        <p className="text-sm text-base-content/50">
          Supported: MP3, FLAC, WAV, OGG, M4A, AAC (max 500MB)
        </p>
      </div>
      {typeof progress === 'number' && progress > 0 && (
        <div className="mt-4">
          <progress
            className="progress progress-primary w-full"
            value={progress}
            max={100}
            role="progressbar"
          />
          <p className="text-sm text-base-content/60 mt-1">{progress}%</p>
        </div>
      )}
    </div>
  );
}
