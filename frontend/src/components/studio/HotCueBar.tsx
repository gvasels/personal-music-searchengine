/**
 * HotCueBar Component
 * Displays 8 hot cue buttons with set/jump/delete functionality
 */
import { useHotCues, useSetHotCue, useDeleteHotCue } from '@/hooks/useHotCues';
import type { HotCue, HotCueColor } from '@/types';

// Default colors for hot cue slots
const DEFAULT_COLORS: HotCueColor[] = [
  '#FF0000', // Red
  '#FF8C00', // Orange
  '#FFFF00', // Yellow
  '#00FF00', // Green
  '#00FFFF', // Cyan
  '#0000FF', // Blue
  '#800080', // Purple
  '#FF69B4', // Pink
];

interface HotCueBarProps {
  trackId: string | undefined;
  currentPosition: number; // Current playback position in seconds
  onSeek?: (position: number) => void;
}

export function HotCueBar({ trackId, currentPosition, onSeek }: HotCueBarProps) {
  const { data, isLoading, isFeatureEnabled, showUpgrade } = useHotCues(trackId);
  const setHotCue = useSetHotCue();
  const deleteHotCue = useDeleteHotCue();

  if (!trackId) {
    return null;
  }

  if (!isFeatureEnabled && showUpgrade) {
    return (
      <div className="flex items-center gap-2 p-2 bg-base-200 rounded">
        <span className="text-xs text-base-content/50">Hot Cues</span>
        <a href="/subscription" className="btn btn-xs btn-outline">
          Upgrade
        </a>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="flex items-center gap-1">
        {Array.from({ length: 8 }).map((_, i) => (
          <div
            key={i}
            className="w-8 h-8 bg-base-300 animate-pulse rounded"
          />
        ))}
      </div>
    );
  }

  // Build hot cue map from response
  const hotCueMap = new Map<number, HotCue>();
  data?.hotCues.forEach((cue) => {
    hotCueMap.set(cue.slot, cue);
  });

  const handleClick = (slot: number) => {
    const existingCue = hotCueMap.get(slot);

    if (existingCue) {
      // Jump to cue position
      onSeek?.(existingCue.position);
    } else {
      // Set new cue at current position
      setHotCue.mutate({
        trackId: trackId!,
        slot,
        position: currentPosition,
        color: DEFAULT_COLORS[slot - 1],
      });
    }
  };

  const handleRightClick = (slot: number, e: React.MouseEvent) => {
    e.preventDefault();
    const existingCue = hotCueMap.get(slot);
    if (existingCue) {
      deleteHotCue.mutate({ trackId: trackId!, slot });
    }
  };

  const formatPosition = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="flex items-center gap-1">
      <span className="text-xs text-base-content/50 mr-2">Cues</span>
      {Array.from({ length: 8 }, (_, i) => i + 1).map((slot) => {
        const cue = hotCueMap.get(slot);
        const color = cue?.color || DEFAULT_COLORS[slot - 1];

        return (
          <button
            key={slot}
            className={`w-8 h-8 rounded font-bold text-xs transition-all ${
              cue
                ? 'text-white shadow-md'
                : 'bg-base-300 text-base-content/30 hover:bg-base-200'
            }`}
            style={cue ? { backgroundColor: color } : undefined}
            onClick={() => handleClick(slot)}
            onContextMenu={(e) => handleRightClick(slot, e)}
            title={
              cue
                ? `${cue.label || `Cue ${slot}`} - ${formatPosition(cue.position)}\nRight-click to delete`
                : `Set cue ${slot} at ${formatPosition(currentPosition)}`
            }
          >
            {slot}
          </button>
        );
      })}
    </div>
  );
}
