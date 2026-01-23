/**
 * Draggable Track Item for playlist reordering
 * Uses @dnd-kit/sortable for drag and drop functionality
 */
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import type { Track } from '@/types';

interface DraggableTrackItemProps {
  track: Track;
  index: number;
  isCurrentTrack: boolean;
  isPlaying: boolean;
  isSelected: boolean;
  onPlay: () => void;
  onRemove: () => void;
  onSelect?: (e: React.MouseEvent) => void;
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function DraggableTrackItem({
  track,
  index,
  isCurrentTrack,
  isPlaying,
  isSelected,
  onPlay,
  onRemove,
  onSelect,
}: DraggableTrackItemProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: track.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <tr
      ref={setNodeRef}
      style={style}
      className={`
        hover cursor-pointer
        ${isCurrentTrack ? 'bg-primary/10' : ''}
        ${isSelected ? 'bg-secondary/20' : ''}
        ${isDragging ? 'opacity-50 bg-base-300 shadow-lg z-50' : ''}
      `}
      onClick={(e) => {
        if (onSelect) {
          onSelect(e);
        } else {
          onPlay();
        }
      }}
    >
      {/* Drag handle */}
      <td className="w-8 p-0">
        <button
          className="btn btn-ghost btn-xs cursor-grab active:cursor-grabbing touch-none"
          {...attributes}
          {...listeners}
          onClick={(e) => e.stopPropagation()}
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-base-content/50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8h16M4 16h16" />
          </svg>
        </button>
      </td>

      {/* Index / Play indicator */}
      <td className="w-12">
        {isCurrentTrack && isPlaying ? (
          <span className="text-primary">▶</span>
        ) : (
          index + 1
        )}
      </td>

      {/* Track info */}
      <td className="font-medium">
        <div className="flex items-center gap-2">
          {track.coverArtUrl && (
            <img
              src={track.coverArtUrl}
              alt=""
              className="w-10 h-10 rounded object-cover"
            />
          )}
          <span>{track.title}</span>
        </div>
      </td>
      <td>{track.artist}</td>
      <td>{track.album}</td>
      <td>{formatDuration(track.duration)}</td>

      {/* Remove button */}
      <td className="w-12">
        <button
          className="btn btn-ghost btn-xs btn-circle text-error"
          onClick={(e) => {
            e.stopPropagation();
            onRemove();
          }}
          title="Remove from playlist"
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </td>
    </tr>
  );
}

export default DraggableTrackItem;
