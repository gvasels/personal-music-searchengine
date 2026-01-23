import { useState, useRef, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { usePlaylistsQuery, playlistKeys } from '@/hooks/usePlaylists';
import { addTrackToPlaylist } from '@/lib/api/client';
import toast from 'react-hot-toast';

interface AddToPlaylistDropdownProps {
  trackId: string;
  onSuccess?: () => void;
}

export function AddToPlaylistDropdown({ trackId, onSuccess }: AddToPlaylistDropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [position, setPosition] = useState({ top: 0, left: 0 });
  const buttonRef = useRef<HTMLButtonElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();

  const { data: playlists, isLoading } = usePlaylistsQuery();

  const addMutation = useMutation({
    mutationFn: (playlistId: string) => addTrackToPlaylist(playlistId, trackId),
    onSuccess: (_, playlistId) => {
      const playlist = playlists?.items?.find(p => p.id === playlistId);
      toast.success(`Added to ${playlist?.name || 'playlist'}`);
      queryClient.invalidateQueries({ queryKey: playlistKeys.all });
      setIsOpen(false);
      onSuccess?.();
    },
    onError: () => {
      toast.error('Failed to add to playlist');
    },
  });

  // Update position when opening
  useEffect(() => {
    if (isOpen && buttonRef.current) {
      const rect = buttonRef.current.getBoundingClientRect();
      setPosition({
        top: rect.bottom + 4,
        left: rect.left - 180, // Position to the left so it doesn't go off screen
      });
    }
  }, [isOpen]);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        buttonRef.current &&
        !buttonRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    }
    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => document.removeEventListener('mousedown', handleClickOutside);
    }
  }, [isOpen]);

  const dropdownContent = isOpen && (
    <div
      ref={dropdownRef}
      className="bg-base-100 rounded-lg shadow-2xl border border-base-300"
      style={{
        position: 'fixed',
        top: position.top,
        left: Math.max(8, position.left), // Don't go off left edge
        zIndex: 99999,
        minWidth: '200px',
      }}
    >
      <div className="px-3 py-2 border-b border-base-300 bg-base-200 rounded-t-lg">
        <span className="font-semibold text-sm">Add to Playlist</span>
      </div>
      <ul className="py-1 max-h-60 overflow-y-auto">
        {isLoading ? (
          <li className="px-3 py-2 text-center">
            <span className="loading loading-spinner loading-sm" />
          </li>
        ) : playlists?.items && playlists.items.length > 0 ? (
          playlists.items.map((playlist) => (
            <li key={playlist.id}>
              <button
                className="w-full px-3 py-2 text-left hover:bg-base-200 flex justify-between items-center transition-colors"
                onClick={(e) => {
                  e.stopPropagation();
                  addMutation.mutate(playlist.id);
                }}
                disabled={addMutation.isPending}
              >
                <span className="truncate">{playlist.name}</span>
                {addMutation.isPending && (
                  <span className="loading loading-spinner loading-xs ml-2" />
                )}
              </button>
            </li>
          ))
        ) : (
          <li className="px-3 py-3 text-center text-base-content/60 text-sm">
            No playlists yet.
            <br />
            <span className="text-xs">Create one from the Playlists page.</span>
          </li>
        )}
      </ul>
    </div>
  );

  return (
    <>
      <button
        ref={buttonRef}
        className="btn btn-ghost btn-xs"
        onClick={(e) => {
          e.stopPropagation();
          setIsOpen(!isOpen);
        }}
        aria-label="Add to playlist"
        title="Add to playlist"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
        </svg>
      </button>
      {dropdownContent && createPortal(dropdownContent, document.body)}
    </>
  );
}
