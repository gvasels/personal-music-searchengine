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
              className="flex items-center justify-between p-3 bg-base-200 rounded-lg"
            >
              <div className="flex items-center gap-3">
                <span className="text-2xl">🎵</span>
                <div>
                  <p className="font-medium">{item.filename}</p>
                  <p className="text-sm text-base-content/60 capitalize">{item.status}</p>
                </div>
              </div>
              <div>
                {item.status === 'completed' && (
                  <span className="badge badge-success">Completed</span>
                )}
                {item.status === 'processing' && (
                  <span className="badge badge-warning">Processing</span>
                )}
                {item.status === 'uploading' && (
                  <span className="loading loading-spinner loading-sm" />
                )}
                {item.status === 'failed' && (
                  <span className="badge badge-error">Failed</span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
