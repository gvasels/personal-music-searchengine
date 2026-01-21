import { useState, useCallback } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { addTagToTrack, removeTagFromTrack } from '@/lib/api/client';
import toast from 'react-hot-toast';

interface TagInputProps {
  trackId: string;
  tags: string[];
}

export function TagInput({ trackId, tags }: TagInputProps) {
  const [inputValue, setInputValue] = useState('');
  const [isAdding, setIsAdding] = useState(false);
  const queryClient = useQueryClient();

  const addMutation = useMutation({
    mutationFn: (tagName: string) => addTagToTrack(trackId, tagName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['track', trackId] });
      setInputValue('');
      toast.success('Tag added');
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to add tag');
    },
  });

  const removeMutation = useMutation({
    mutationFn: (tagName: string) => removeTagFromTrack(trackId, tagName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['track', trackId] });
      toast.success('Tag removed');
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to remove tag');
    },
  });

  const handleAddTag = useCallback(() => {
    const tagName = inputValue.trim().toLowerCase();
    if (tagName && !tags.includes(tagName)) {
      addMutation.mutate(tagName);
    }
    setIsAdding(false);
    setInputValue('');
  }, [inputValue, tags, addMutation]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === 'Enter') {
        e.preventDefault();
        handleAddTag();
      }
      if (e.key === 'Escape') {
        setIsAdding(false);
        setInputValue('');
      }
    },
    [handleAddTag]
  );

  return (
    <div className="space-y-2">
      <div className="flex flex-wrap gap-2">
        {tags.map((tag) => (
          <span key={tag} className="badge badge-outline gap-1">
            🏷️ {tag}
            <button
              className="btn btn-ghost btn-xs btn-circle"
              onClick={() => removeMutation.mutate(tag)}
              disabled={removeMutation.isPending}
              aria-label={`remove ${tag}`}
            >
              ✕
            </button>
          </span>
        ))}

        {isAdding ? (
          <div className="flex items-center gap-1">
            <input
              type="text"
              placeholder="Tag name..."
              className="input input-xs input-bordered w-24"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleKeyDown}
              onBlur={handleAddTag}
              autoFocus
              disabled={addMutation.isPending}
            />
            {addMutation.isPending && (
              <span className="loading loading-spinner loading-xs" />
            )}
          </div>
        ) : (
          <button
            className="badge badge-outline badge-ghost gap-1 cursor-pointer hover:badge-primary"
            onClick={() => setIsAdding(true)}
            aria-label="add tag"
          >
            + Add tag
          </button>
        )}
      </div>
    </div>
  );
}
