/**
 * Track Detail Page - Wave 2
 */
import { useState } from 'react';
import { useNavigate, useParams, Link } from '@tanstack/react-router';
import { useTrackQuery, useUpdateTrack, useDeleteTrack } from '../../hooks/useTracks';

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export default function TrackDetailPage() {
  const navigate = useNavigate();
  const { trackId } = useParams({ from: '/tracks/$trackId' });
  const { data: track, isLoading, isError, error } = useTrackQuery(trackId);
  const updateTrack = useUpdateTrack();
  const deleteTrack = useDeleteTrack();

  const [isEditing, setIsEditing] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [editForm, setEditForm] = useState({ title: '', artist: '', album: '' });

  const handleBack = () => {
    void navigate({ to: '/tracks' });
  };

  const handleEdit = () => {
    if (track) {
      setEditForm({ title: track.title, artist: track.artist, album: track.album });
      setIsEditing(true);
    }
  };

  const handleSave = async () => {
    if (!track) return;
    await updateTrack.mutateAsync({ id: track.id, data: editForm });
    setIsEditing(false);
  };

  const handleDelete = async () => {
    if (!track) return;
    await deleteTrack.mutateAsync(track.id);
    void navigate({ to: '/tracks' });
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading" />
      </div>
    );
  }

  if (isError || !track) {
    return (
      <div className="space-y-4">
        <div className="alert alert-error">
          <span>Track not found: {error?.message}</span>
        </div>
        <button className="btn" onClick={handleBack}>
          Back
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <button className="btn btn-ghost btn-sm" onClick={handleBack}>
        ← Back
      </button>

      <div className="flex gap-6">
        <div className="flex-1 space-y-4">
          {isEditing ? (
            <div className="space-y-4">
              <div className="form-control">
                <label className="label" htmlFor="title">
                  <span className="label-text">Title</span>
                </label>
                <input
                  id="title"
                  type="text"
                  className="input input-bordered"
                  value={editForm.title}
                  onChange={(e) => setEditForm({ ...editForm, title: e.target.value })}
                />
              </div>
              <div className="form-control">
                <label className="label" htmlFor="artist">
                  <span className="label-text">Artist</span>
                </label>
                <input
                  id="artist"
                  type="text"
                  className="input input-bordered"
                  value={editForm.artist}
                  onChange={(e) => setEditForm({ ...editForm, artist: e.target.value })}
                />
              </div>
              <div className="form-control">
                <label className="label" htmlFor="album">
                  <span className="label-text">Album</span>
                </label>
                <input
                  id="album"
                  type="text"
                  className="input input-bordered"
                  value={editForm.album}
                  onChange={(e) => setEditForm({ ...editForm, album: e.target.value })}
                />
              </div>
              <div className="flex gap-2">
                <button className="btn btn-primary" onClick={() => void handleSave()}>
                  Save
                </button>
                <button className="btn btn-ghost" onClick={() => setIsEditing(false)}>
                  Cancel
                </button>
              </div>
            </div>
          ) : (
            <>
              <h1 className="text-3xl font-bold">{track.title}</h1>
              <p className="text-lg">
                <Link to="/artists/$artistName" params={{ artistName: track.artist }} className="link">
                  {track.artist}
                </Link>
              </p>
              {track.albumId && (
                <p>
                  <Link to="/albums/$albumId" params={{ albumId: track.albumId }} className="link">
                    {track.album}
                  </Link>
                </p>
              )}

              <div className="stats shadow">
                <div className="stat">
                  <div className="stat-title">Duration</div>
                  <div className="stat-value text-lg">{formatDuration(track.duration)}</div>
                </div>
                <div className="stat">
                  <div className="stat-title">Format</div>
                  <div className="stat-value text-lg uppercase">{track.format}</div>
                </div>
              </div>

              {track.tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {track.tags.map((tag) => (
                    <span key={tag} className="badge badge-outline">
                      {tag}
                    </span>
                  ))}
                </div>
              )}

              <div className="flex gap-2">
                <button className="btn btn-primary">Play</button>
                <button className="btn btn-outline" onClick={handleEdit}>
                  Edit
                </button>
                <button className="btn btn-error btn-outline" onClick={() => setShowDeleteConfirm(true)}>
                  Delete
                </button>
              </div>
            </>
          )}
        </div>
      </div>

      {showDeleteConfirm && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">Delete Track</h3>
            <p className="py-4">Are you sure you want to delete "{track.title}"?</p>
            <div className="modal-action">
              <button className="btn btn-error" onClick={() => void handleDelete()}>
                Confirm
              </button>
              <button className="btn" onClick={() => setShowDeleteConfirm(false)}>
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
