import { useState, useEffect, useCallback } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createPlaylist } from '@/lib/api/client';
import toast from 'react-hot-toast';

interface CreatePlaylistModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CreatePlaylistModal({ isOpen, onClose }: CreatePlaylistModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState('');
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: createPlaylist,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['playlists'] });
      toast.success('Playlist created');
      handleClose();
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to create playlist');
    },
  });

  const handleClose = useCallback(() => {
    setName('');
    setDescription('');
    setError('');
    onClose();
  }, [onClose]);

  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      setError('');

      if (!name.trim()) {
        setError('Name is required');
        return;
      }

      if (name.trim().length < 3) {
        setError('Name must be at least 3 characters');
        return;
      }

      mutation.mutate({ name: name.trim(), description: description.trim() || undefined });
    },
    [name, description, mutation]
  );

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        handleClose();
      }
    };
    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, handleClose]);

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div
        data-testid="modal-backdrop"
        className="modal-backdrop"
        onClick={handleClose}
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        className="modal-box"
      >
        <h3 id="modal-title" className="font-bold text-lg">
          Create Playlist
        </h3>
        <form onSubmit={handleSubmit} className="py-4 space-y-4">
          <div className="form-control">
            <label className="label" htmlFor="playlist-name">
              <span className="label-text">Name</span>
            </label>
            <input
              id="playlist-name"
              type="text"
              placeholder="My Playlist"
              className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
            {error && <span className="label-text-alt text-error mt-1">{error}</span>}
          </div>
          <div className="form-control">
            <label className="label" htmlFor="playlist-description">
              <span className="label-text">Description</span>
            </label>
            <textarea
              id="playlist-description"
              placeholder="Optional description"
              className="textarea textarea-bordered w-full"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>
          <div className="modal-action">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={handleClose}
              aria-label="cancel"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={mutation.isPending}
              aria-label={mutation.isPending ? 'creating' : 'create'}
            >
              {mutation.isPending ? (
                <>
                  <span className="loading loading-spinner loading-sm" />
                  Creating...
                </>
              ) : (
                'Create'
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
