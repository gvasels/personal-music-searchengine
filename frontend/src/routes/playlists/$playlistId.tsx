/**
 * Playlist Detail Page - Wave 5
 */
import { useState } from 'react';
import { Link, useParams, useNavigate } from '@tanstack/react-router';
import { usePlaylistQuery, useUpdatePlaylist, useDeletePlaylist, useRemoveTracksFromPlaylist } from '../../hooks/usePlaylists';
import { usePlayerStore } from '../../lib/store/playerStore';

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export default function PlaylistDetailPage() {
  const navigate = useNavigate();
  const { playlistId } = useParams({ strict: false });
  const { data, isLoading, isError, error } = usePlaylistQuery(playlistId);
  const updatePlaylist = useUpdatePlaylist();
  const deletePlaylist = useDeletePlaylist();
  const removeTracksFromPlaylist = useRemoveTracksFromPlaylist();
  const { setQueue, currentTrack, isPlaying } = usePlayerStore();

  const [isEditing, setIsEditing] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [editForm, setEditForm] = useState({ name: '', description: '' });

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-64">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="alert alert-error">
        <span>{error?.message || 'Playlist not found'}</span>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="alert alert-error">
        <span>Playlist not found</span>
      </div>
    );
  }

  const { playlist, tracks } = data;

  const handlePlayAll = () => {
    if (tracks.length > 0) {
      setQueue(tracks, 0);
    }
  };

  const handleEdit = () => {
    setEditForm({ name: playlist.name, description: playlist.description || '' });
    setIsEditing(true);
  };

  const handleSave = async () => {
    if (!playlistId) return;
    await updatePlaylist.mutateAsync({ id: playlistId, data: editForm });
    setIsEditing(false);
  };

  const handleDelete = async () => {
    if (!playlistId) return;
    await deletePlaylist.mutateAsync(playlistId);
    void navigate({ to: '/playlists' });
  };

  const handleRemoveTrack = async (trackId: string) => {
    if (!playlistId) return;
    await removeTracksFromPlaylist.mutateAsync({ id: playlistId, trackIds: [trackId] });
  };

  const handlePlayTrack = (index: number) => {
    setQueue(tracks, index);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/playlists" className="btn btn-ghost btn-sm">
          ← Back
        </Link>
      </div>

      <div className="flex flex-col md:flex-row md:items-end gap-6">
        <div className="w-48 h-48 bg-base-300 rounded-lg flex items-center justify-center">
          <span className="text-6xl">🎵</span>
        </div>

        <div className="space-y-2 flex-1">
          {isEditing ? (
            <div className="space-y-4 max-w-md">
              <div className="form-control">
                <label className="label" htmlFor="name">
                  <span className="label-text">Name</span>
                </label>
                <input
                  id="name"
                  type="text"
                  className="input input-bordered"
                  value={editForm.name}
                  onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                />
              </div>
              <div className="form-control">
                <label className="label" htmlFor="description">
                  <span className="label-text">Description</span>
                </label>
                <textarea
                  id="description"
                  className="textarea textarea-bordered"
                  value={editForm.description}
                  onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                />
              </div>
              <div className="flex gap-2">
                <button
                  className="btn btn-primary"
                  onClick={() => void handleSave()}
                  disabled={updatePlaylist.isPending}
                >
                  {updatePlaylist.isPending ? 'Saving...' : 'Save'}
                </button>
                <button className="btn btn-ghost" onClick={() => setIsEditing(false)}>
                  Cancel
                </button>
              </div>
            </div>
          ) : (
            <>
              <p className="text-sm uppercase tracking-wide text-base-content/60">Playlist</p>
              <h1 className="text-4xl font-bold">{playlist.name}</h1>
              {playlist.description && (
                <p className="text-base-content/70">{playlist.description}</p>
              )}
              <p className="text-sm text-base-content/60">{tracks.length} tracks</p>
            </>
          )}
        </div>
      </div>

      {!isEditing && (
        <div className="flex gap-2">
          <button
            className="btn btn-primary"
            onClick={handlePlayAll}
            disabled={tracks.length === 0}
          >
            Play All
          </button>
          <button className="btn btn-outline" onClick={handleEdit}>
            Edit
          </button>
          <button
            className="btn btn-outline btn-error"
            onClick={() => setShowDeleteConfirm(true)}
          >
            Delete
          </button>
        </div>
      )}

      {tracks.length === 0 ? (
        <div className="text-center py-12 text-base-content/60">
          <p>No tracks in this playlist</p>
          <p className="text-sm mt-2">Add tracks from your library</p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="table table-zebra">
            <thead>
              <tr>
                <th className="w-12">#</th>
                <th>Title</th>
                <th>Artist</th>
                <th>Album</th>
                <th>Duration</th>
                <th className="w-12"></th>
              </tr>
            </thead>
            <tbody>
              {tracks.map((track, index) => {
                const isCurrentTrack = currentTrack?.id === track.id;
                return (
                  <tr
                    key={track.id}
                    className={`hover cursor-pointer ${isCurrentTrack ? 'bg-primary/10' : ''}`}
                    onClick={() => handlePlayTrack(index)}
                  >
                    <td>
                      {isCurrentTrack && isPlaying ? (
                        <span className="text-primary">▶</span>
                      ) : (
                        index + 1
                      )}
                    </td>
                    <td className="font-medium">{track.title}</td>
                    <td>{track.artist}</td>
                    <td>{track.album}</td>
                    <td>{formatDuration(track.duration)}</td>
                    <td>
                      <button
                        className="btn btn-ghost btn-xs btn-circle text-error"
                        onClick={(e) => {
                          e.stopPropagation();
                          void handleRemoveTrack(track.id);
                        }}
                        title="Remove from playlist"
                      >
                        ✕
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {showDeleteConfirm && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">Delete Playlist</h3>
            <p className="py-4">Are you sure you want to delete "{playlist.name}"?</p>
            <div className="modal-action">
              <button
                className="btn btn-error"
                onClick={() => void handleDelete()}
                disabled={deletePlaylist.isPending}
              >
                {deletePlaylist.isPending ? 'Deleting...' : 'Delete'}
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
