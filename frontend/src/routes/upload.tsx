/**
 * Upload Page - Wave 4
 */
import { UploadDropzone } from '../components/upload';
import { useUpload } from '../hooks/useUpload';

export default function UploadPage() {
  const { upload, isUploading, progress, error, uploads, reset } = useUpload();

  const handleFilesSelected = (files: File[]) => {
    void upload(files);
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Upload Music</h1>
        {uploads.length > 0 && (
          <button className="btn btn-ghost btn-sm" onClick={reset}>
            Clear
          </button>
        )}
      </div>

      <UploadDropzone
        onFilesSelected={handleFilesSelected}
        onError={(err) => console.error(err)}
        progress={isUploading ? progress : undefined}
        disabled={isUploading}
      />

      {error && (
        <div className="alert alert-error">
          <span>{error}</span>
        </div>
      )}

      {uploads.length > 0 && (
        <div className="space-y-3">
          <h2 className="text-lg font-semibold">Uploads</h2>
          {uploads.map((item) => (
            <div
              key={item.id}
              className="p-4 bg-base-200 rounded-lg space-y-3"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <span className="text-2xl">🎵</span>
                  <div>
                    <p className="font-medium">{item.filename}</p>
                    <p className="text-sm text-base-content/60">{item.currentStep}</p>
                  </div>
                </div>
                <div>
                  {item.status === 'completed' && (
                    <span className="badge badge-success">Completed</span>
                  )}
                  {item.status === 'failed' && (
                    <span className="badge badge-error">Failed</span>
                  )}
                  {(item.status === 'uploading' || item.status === 'processing') && (
                    <span className="text-sm font-medium">{item.progress}%</span>
                  )}
                </div>
              </div>
              {(item.status === 'uploading' || item.status === 'processing') && (
                <progress
                  className="progress progress-primary w-full"
                  value={item.progress}
                  max={100}
                />
              )}
              {item.error && (
                <p className="text-sm text-error">{item.error}</p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
