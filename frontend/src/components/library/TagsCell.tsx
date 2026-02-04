import { useState, useRef, useEffect } from 'react';
import { TagInput } from '@/components/tag/TagInput';

interface TagsCellProps {
  trackId: string;
  tags: string[];
  maxVisible?: number;
}

/**
 * TagsCell - Compact tag display with inline editing
 * Shows up to maxVisible tags as badges, with "+N more" for overflow
 * Click to expand into full TagInput for editing
 */
export function TagsCell({ trackId, tags, maxVisible = 3 }: TagsCellProps) {
  const [isEditing, setIsEditing] = useState(false);
  const cellRef = useRef<HTMLDivElement>(null);

  const visibleTags = tags.slice(0, maxVisible);
  const hiddenCount = Math.max(0, tags.length - maxVisible);

  // Close edit mode when clicking outside
  useEffect(() => {
    if (!isEditing) return;

    const handleClickOutside = (event: MouseEvent) => {
      if (cellRef.current && !cellRef.current.contains(event.target as Node)) {
        setIsEditing(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isEditing]);

  // Handle escape key
  useEffect(() => {
    if (!isEditing) return;

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setIsEditing(false);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isEditing]);

  if (isEditing) {
    return (
      <div ref={cellRef} className="min-w-48">
        <TagInput trackId={trackId} tags={tags} />
      </div>
    );
  }

  return (
    <div
      ref={cellRef}
      className="flex flex-wrap items-center gap-1 cursor-pointer group"
      onClick={() => setIsEditing(true)}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          setIsEditing(true);
        }
      }}
      tabIndex={0}
      role="button"
      aria-label={tags.length > 0 ? `Edit tags: ${tags.join(', ')}` : 'Add tags'}
    >
      {tags.length === 0 ? (
        <span className="text-base-content/40 text-sm group-hover:text-primary transition-colors">
          + Add tags
        </span>
      ) : (
        <>
          {visibleTags.map((tag) => (
            <span
              key={tag}
              className="badge badge-sm badge-outline group-hover:badge-primary transition-colors"
            >
              {tag}
            </span>
          ))}
          {hiddenCount > 0 && (
            <span className="badge badge-sm badge-ghost group-hover:badge-primary transition-colors">
              +{hiddenCount}
            </span>
          )}
        </>
      )}
    </div>
  );
}
