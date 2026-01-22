/**
 * Tracks List Page - Wave 2
 * Updated with play button, tags column, and actions
 */
import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { useTracksQuery, useDeleteTrack } from '../../hooks/useTracks';
import { usePlayerStore } from '@/lib/store/playerStore';
import { TagsCell } from '@/components/library/TagsCell';
import { getDownloadUrl } from '@/lib/api/client';
import type { Track } from '../../types';

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
}

export default function TracksPage() {
  const navigate = useNavigate();
  const { data, isLoading, isError, error, refetch } = useTracksQuery();
  const { currentTrack, isPlaying, setQueue, pause } = usePlayerStore();
  const deleteTrackMutation = useDeleteTrack();
  const [trackToDelete, setTrackToDelete] = useState<Track | null>(null);

  const handleTrackClick = (track: Track) => {
    void navigate({
      to: '/tracks/$trackId',
      params: { trackId: track.id },
    });
  };

  const handlePlayClick = (e: React.MouseEvent, track: Track, index: number) => {
    e.stopPropagation(); // Don't navigate when clicking play

    if (currentTrack?.id === track.id) {
      // Toggle play/pause for current track
      if (isPlaying) {
        pause();
      } else {
        usePlayerStore.getState().play();
      }
    } else {
      // Play new track, queue all tracks starting from this one
      setQueue(data?.items || [], index);
    }
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    void navigate({
      to: '/tracks',
      search: { sortBy: e.target.value },
    });
  };

  const handleDownload = async (e: React.MouseEvent, track: Track) => {
    e.stopPropagation();
    try {
      const { downloadUrl, filename } = await getDownloadUrl(track.id);
      // Create a temporary link to trigger download
      const link = document.createElement('a');
      link.href = downloadUrl;
      link.download = filename || `${track.title}.${track.format}`;
      link.target = '_blank';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (err) {
      console.error('Failed to download track:', err);
    }
  };

  const handleDeleteClick = (e: React.MouseEvent, track: Track) => {
    e.stopPropagation();
    setTrackToDelete(track);
  };

  const handleConfirmDelete = async () => {
    if (trackToDelete) {
      try {
        await deleteTrackMutation.mutateAsync(trackToDelete.id);
        setTrackToDelete(null);
      } catch (err) {
        console.error('Failed to delete track:', err);
      }
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="alert alert-error">
        <span>Failed to load tracks: {error?.message}</span>
        <button className="btn btn-sm" onClick={() => void refetch()}>
          Retry
        </button>
      </div>
    );
  }

  if (!data?.items.length) {
    return (
      <div className="text-center py-12">
        <p className="text-base-content/60">No tracks found</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{data.total} tracks</h1>
        <label className="form-control w-auto">
          <select
            className="select select-bordered select-sm"
            onChange={handleSortChange}
            aria-label="Sort by"
          >
            <option value="">Sort by...</option>
            <option value="title">Title</option>
            <option value="artist">Artist</option>
            <option value="album">Album</option>
            <option value="createdAt">Date Added</option>
          </select>
        </label>
      </div>

      <div className="overflow-x-auto">
        <table className="table table-zebra">
          <thead>
            <tr>
              <th className="w-12"></th>
              <th>Title</th>
              <th>Artist</th>
              <th>Album</th>
              <th>Tags</th>
              <th>Duration</th>
              <th>Uploaded</th>
              <th className="w-24">Actions</th>
            </tr>
          </thead>
          <tbody>
            {data.items.map((track, index) => {
              const isCurrentTrack = currentTrack?.id === track.id;
              const isCurrentlyPlaying = isCurrentTrack && isPlaying;

              return (
                <tr
                  key={track.id}
                  className={`hover cursor-pointer ${isCurrentTrack ? 'bg-primary/10' : ''}`}
                  onClick={() => handleTrackClick(track)}
                >
                  <td className="w-12" onClick={(e) => e.stopPropagation()}>
                    <button
                      className={`btn btn-ghost btn-sm btn-circle ${isCurrentlyPlaying ? 'text-primary' : ''}`}
                      onClick={(e) => handlePlayClick(e, track, index)}
                      aria-label={isCurrentlyPlaying ? 'Pause' : 'Play'}
                    >
                      {isCurrentlyPlaying ? (
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                          <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                        </svg>
                      ) : isCurrentTrack ? (
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                          <path d="M8 5v14l11-7z" />
                        </svg>
                      ) : (
                        <svg className="w-5 h-5 opacity-50 group-hover:opacity-100" fill="currentColor" viewBox="0 0 24 24">
                          <path d="M8 5v14l11-7z" />
                        </svg>
                      )}
                    </button>
                  </td>
                  <td className="font-medium">{track.title}</td>
                  <td>{track.artist}</td>
                  <td>{track.album}</td>
                  <td onClick={(e) => e.stopPropagation()}>
                    <TagsCell trackId={track.id} tags={track.tags || []} maxVisible={2} />
                  </td>
                  <td>{formatDuration(track.duration)}</td>
                  <td className="text-sm text-base-content/60">{formatDate(track.createdAt)}</td>
                  <td onClick={(e) => e.stopPropagation()}>
                    <div className="flex gap-1">
                      <button
                        className="btn btn-ghost btn-xs"
                        onClick={(e) => handleDownload(e, track)}
                        aria-label="Download"
                        title="Download"
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                        </svg>
                      </button>
                      <button
                        className="btn btn-ghost btn-xs text-error"
                        onClick={(e) => handleDeleteClick(e, track)}
                        aria-label="Delete"
                        title="Delete"
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Delete Confirmation Modal */}
      {trackToDelete && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">Delete Track</h3>
            <p className="py-4">
              Are you sure you want to delete "{trackToDelete.title}" by {trackToDelete.artist}?
              This action cannot be undone.
            </p>
            <div className="modal-action">
              <button
                className="btn"
                onClick={() => setTrackToDelete(null)}
                disabled={deleteTrackMutation.isPending}
              >
                Cancel
              </button>
              <button
                className="btn btn-error"
                onClick={() => void handleConfirmDelete()}
                disabled={deleteTrackMutation.isPending}
              >
                {deleteTrackMutation.isPending ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : (
                  'Delete'
                )}
              </button>
            </div>
          </div>
          <div className="modal-backdrop" onClick={() => setTrackToDelete(null)} />
        </div>
      )}
    </div>
  );
}
