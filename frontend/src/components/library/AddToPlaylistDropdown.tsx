import { useState, useRef, useEffect } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { usePlaylistsQuery, playlistKeys } from '@/hooks/usePlaylists';
import { addTrackToPlaylist } from '@/lib/api/client';

interface AddToPlaylistDropdownProps {
  trackId: string;
  onSuccess?: () => void;
}

export function AddToPlaylistDropdown({ trackId, onSuccess }: AddToPlaylistDropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();

  const { data: playlists, isLoading } = usePlaylistsQuery();

  const addMutation = useMutation({
    mutationFn: (playlistId: string) => addTrackToPlaylist(playlistId, trackId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: playlistKeys.all });
      setIsOpen(false);
      onSuccess?.();
    },
  });

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className="dropdown dropdown-end" ref={dropdownRef}>
      <button
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

      {isOpen && (
        <ul className="dropdown-content menu bg-base-200 rounded-box z-50 w-52 p-2 shadow-lg">
          {isLoading ? (
            <li className="p-2 text-center">
              <span className="loading loading-spinner loading-sm" />
            </li>
          ) : playlists?.items && playlists.items.length > 0 ? (
            playlists.items.map((playlist) => (
              <li key={playlist.id}>
                <button
                  className="text-left"
                  onClick={(e) => {
                    e.stopPropagation();
                    addMutation.mutate(playlist.id);
                  }}
                  disabled={addMutation.isPending}
                >
                  {playlist.name}
                  {addMutation.isPending && (
                    <span className="loading loading-spinner loading-xs" />
                  )}
                </button>
              </li>
            ))
          ) : (
            <li className="p-2 text-center text-base-content/60">
              No playlists yet
            </li>
          )}
        </ul>
      )}
    </div>
  );
}
