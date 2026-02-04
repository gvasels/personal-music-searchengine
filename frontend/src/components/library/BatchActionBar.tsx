/**
 * Batch Action Bar
 * Shows batch operations when tracks are selected
 */
import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useSelection } from '@/hooks/useSelection';
import { getPlaylists, addTrackToPlaylist, addTagToTrack, deleteTrack } from '@/lib/api/client';
import toast from 'react-hot-toast';

interface BatchActionBarProps {
  /** Callback when selection is cleared */
  onClear?: () => void;
}

export function BatchActionBar({ onClear }: BatchActionBarProps) {
  const queryClient = useQueryClient();
  const { selectedIds, selectionCount, clearSelection } = useSelection();
  const [tagInput, setTagInput] = useState('');
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [progress, setProgress] = useState({ current: 0, total: 0 });

  // Fetch user's playlists for the dropdown
  const { data: playlistsData } = useQuery({
    queryKey: ['playlists'],
    queryFn: () => getPlaylists(),
    enabled: selectionCount > 0,
  });

  // Add to playlist mutation
  const addToPlaylistMutation = useMutation({
    mutationFn: async ({ playlistId, trackIds }: { playlistId: string; trackIds: string[] }) => {
      setIsProcessing(true);
      setProgress({ current: 0, total: trackIds.length });

      const results: { success: string[]; failed: string[] } = { success: [], failed: [] };

      for (let i = 0; i < trackIds.length; i++) {
        try {
          await addTrackToPlaylist(playlistId, trackIds[i]);
          results.success.push(trackIds[i]);
        } catch {
          results.failed.push(trackIds[i]);
        }
        setProgress({ current: i + 1, total: trackIds.length });
      }

      return results;
    },
    onSuccess: (results, { playlistId }) => {
      setIsProcessing(false);
      if (results.success.length > 0) {
        toast.success(`Added ${results.success.length} track(s) to playlist`);
        queryClient.invalidateQueries({ queryKey: ['playlist', playlistId] });
        clearSelection();
        onClear?.();
      }
      if (results.failed.length > 0) {
        toast.error(`Failed to add ${results.failed.length} track(s)`);
      }
    },
    onError: () => {
      setIsProcessing(false);
      toast.error('Failed to add tracks to playlist');
    },
  });

  // Add tags mutation
  const addTagsMutation = useMutation({
    mutationFn: async ({ tag, trackIds }: { tag: string; trackIds: string[] }) => {
      setIsProcessing(true);
      setProgress({ current: 0, total: trackIds.length });

      const results: { success: string[]; failed: string[] } = { success: [], failed: [] };

      for (let i = 0; i < trackIds.length; i++) {
        try {
          await addTagToTrack(trackIds[i], tag);
          results.success.push(trackIds[i]);
        } catch {
          results.failed.push(trackIds[i]);
        }
        setProgress({ current: i + 1, total: trackIds.length });
      }

      return results;
    },
    onSuccess: (results) => {
      setIsProcessing(false);
      setTagInput('');
      if (results.success.length > 0) {
        toast.success(`Tagged ${results.success.length} track(s)`);
        queryClient.invalidateQueries({ queryKey: ['tracks'] });
        clearSelection();
        onClear?.();
      }
      if (results.failed.length > 0) {
        toast.error(`Failed to tag ${results.failed.length} track(s)`);
      }
    },
    onError: () => {
      setIsProcessing(false);
      toast.error('Failed to add tags');
    },
  });

  // Delete tracks mutation
  const deleteTracksMutation = useMutation({
    mutationFn: async (trackIds: string[]) => {
      setIsProcessing(true);
      setProgress({ current: 0, total: trackIds.length });

      const results: { success: string[]; failed: string[] } = { success: [], failed: [] };

      for (let i = 0; i < trackIds.length; i++) {
        try {
          await deleteTrack(trackIds[i]);
          results.success.push(trackIds[i]);
        } catch {
          results.failed.push(trackIds[i]);
        }
        setProgress({ current: i + 1, total: trackIds.length });
      }

      return results;
    },
    onSuccess: (results) => {
      setIsProcessing(false);
      setShowDeleteConfirm(false);
      if (results.success.length > 0) {
        toast.success(`Deleted ${results.success.length} track(s)`);
        queryClient.invalidateQueries({ queryKey: ['tracks'] });
        clearSelection();
        onClear?.();
      }
      if (results.failed.length > 0) {
        toast.error(`Failed to delete ${results.failed.length} track(s)`);
      }
    },
    onError: () => {
      setIsProcessing(false);
      toast.error('Failed to delete tracks');
    },
  });

  const handleAddToPlaylist = (playlistId: string) => {
    const trackIds = Array.from(selectedIds);
    addToPlaylistMutation.mutate({ playlistId, trackIds });
  };

  const handleAddTag = () => {
    const tag = tagInput.trim().toLowerCase();
    if (!tag) return;

    const trackIds = Array.from(selectedIds);
    addTagsMutation.mutate({ tag, trackIds });
  };

  const handleDelete = () => {
    const trackIds = Array.from(selectedIds);
    deleteTracksMutation.mutate(trackIds);
  };

  const handleClear = () => {
    clearSelection();
    onClear?.();
  };

  // Don't render if nothing is selected
  if (selectionCount === 0) return null;

  return (
    <div className="fixed bottom-20 left-1/2 -translate-x-1/2 z-50 animate-in slide-in-from-bottom-4">
      <div className="bg-base-200 border border-base-300 rounded-lg shadow-xl p-3 flex items-center gap-3 flex-wrap">
        {/* Selection count */}
        <span className="font-medium whitespace-nowrap">
          {selectionCount} track{selectionCount !== 1 ? 's' : ''} selected
        </span>

        <div className="divider divider-horizontal m-0" />

        {/* Progress indicator */}
        {isProcessing && (
          <div className="flex items-center gap-2">
            <span className="loading loading-spinner loading-sm" />
            <span className="text-sm">
              {progress.current}/{progress.total}
            </span>
          </div>
        )}

        {!isProcessing && (
          <>
            {/* Add to Playlist dropdown */}
            <div className="dropdown dropdown-top">
              <label tabIndex={0} className="btn btn-sm btn-outline">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                Add to Playlist
              </label>
              <ul tabIndex={0} className="dropdown-content z-[1] menu p-2 shadow-lg bg-base-100 rounded-box w-52 max-h-48 overflow-y-auto mb-2">
                {playlistsData?.items.length === 0 && (
                  <li className="text-sm text-base-content/50 p-2">No playlists</li>
                )}
                {playlistsData?.items.map((playlist) => (
                  <li key={playlist.id}>
                    <button onClick={() => handleAddToPlaylist(playlist.id)}>
                      {playlist.name}
                    </button>
                  </li>
                ))}
              </ul>
            </div>

            {/* Add Tag */}
            <div className="join">
              <input
                type="text"
                placeholder="Add tag..."
                className="input input-bordered input-sm join-item w-24"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleAddTag()}
              />
              <button
                className="btn btn-sm btn-outline join-item"
                onClick={handleAddTag}
                disabled={!tagInput.trim()}
              >
                Tag
              </button>
            </div>

            {/* Delete */}
            {!showDeleteConfirm ? (
              <button
                className="btn btn-sm btn-error btn-outline"
                onClick={() => setShowDeleteConfirm(true)}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            ) : (
              <div className="flex items-center gap-2">
                <span className="text-sm text-error">Delete?</span>
                <button
                  className="btn btn-sm btn-error"
                  onClick={handleDelete}
                >
                  Yes
                </button>
                <button
                  className="btn btn-sm btn-ghost"
                  onClick={() => setShowDeleteConfirm(false)}
                >
                  No
                </button>
              </div>
            )}
          </>
        )}

        <div className="divider divider-horizontal m-0" />

        {/* Clear selection */}
        <button
          className="btn btn-sm btn-ghost"
          onClick={handleClear}
          disabled={isProcessing}
        >
          Clear
        </button>
      </div>
    </div>
  );
}

export default BatchActionBar;
