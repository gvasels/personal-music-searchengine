/**
 * Track Card - Mobile-optimized track display
 * Card layout with swipe gestures for quick actions
 */
import { useState, useCallback, useRef } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { usePlayerStore } from '@/lib/store/playerStore';
import { deleteTrack } from '@/lib/api/client';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import type { Track } from '@/types';

interface TrackCardProps {
  track: Track;
  tracks: Track[];
  index: number;
  onDelete?: (trackId: string) => void;
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function TrackCard({ track, tracks, index, onDelete }: TrackCardProps) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { setQueue, currentTrack, isPlaying } = usePlayerStore();

  // Swipe state
  const [touchStart, setTouchStart] = useState<number | null>(null);
  const [touchEnd, setTouchEnd] = useState<number | null>(null);
  const [translateX, setTranslateX] = useState(0);
  const [showUndo, setShowUndo] = useState(false);
  const undoTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const isCurrentTrack = currentTrack?.id === track.id;

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: () => deleteTrack(track.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tracks'] });
      toast.success(`Deleted "${track.title}"`);
    },
    onError: () => {
      toast.error('Failed to delete track');
      setShowUndo(false);
    },
  });

  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    setTouchStart(e.targetTouches[0].clientX);
    setTouchEnd(e.targetTouches[0].clientX);
  }, []);

  const handleTouchMove = useCallback((e: React.TouchEvent) => {
    if (touchStart === null) return;
    const currentX = e.targetTouches[0].clientX;
    setTouchEnd(currentX);

    const diff = currentX - touchStart;
    // Limit swipe distance
    const clampedDiff = Math.max(-100, Math.min(100, diff));
    setTranslateX(clampedDiff);
  }, [touchStart]);

  const handleTouchEnd = useCallback(() => {
    if (touchStart === null || touchEnd === null) return;

    const diff = touchEnd - touchStart;
    const minSwipeDistance = 50;

    if (diff > minSwipeDistance) {
      // Swipe right - Play
      setQueue(tracks, index);
      toast.success(`Playing "${track.title}"`);
    } else if (diff < -minSwipeDistance) {
      // Swipe left - Show delete confirmation
      setShowUndo(true);
      // Auto-delete after 3 seconds if not undone
      undoTimeoutRef.current = setTimeout(() => {
        if (onDelete) {
          onDelete(track.id);
        } else {
          deleteMutation.mutate();
        }
        setShowUndo(false);
      }, 3000);
    }

    // Reset
    setTouchStart(null);
    setTouchEnd(null);
    setTranslateX(0);
  }, [touchStart, touchEnd, track, tracks, index, setQueue, deleteMutation, onDelete]);

  const handleUndo = useCallback(() => {
    if (undoTimeoutRef.current) {
      clearTimeout(undoTimeoutRef.current);
      undoTimeoutRef.current = null;
    }
    setShowUndo(false);
    toast.success('Delete cancelled');
  }, []);

  const handleCardClick = () => {
    if (!showUndo) {
      navigate({ to: '/tracks/$trackId', params: { trackId: track.id } });
    }
  };

  const handlePlayClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setQueue(tracks, index);
  };

  if (showUndo) {
    return (
      <div className="card bg-error/20 border border-error shadow-sm">
        <div className="card-body p-4 flex-row items-center justify-between">
          <span className="text-error font-medium">Deleting "{track.title}"...</span>
          <button className="btn btn-sm btn-ghost" onClick={handleUndo}>
            Undo
          </button>
        </div>
      </div>
    );
  }

  return (
    <div
      className={`card bg-base-100 shadow-sm transition-all duration-150 ${
        isCurrentTrack ? 'ring-2 ring-primary' : ''
      }`}
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
      onClick={handleCardClick}
      style={{
        transform: `translateX(${translateX}px)`,
        transition: touchStart === null ? 'transform 0.15s ease-out' : 'none',
      }}
    >
      {/* Swipe indicators (visible during swipe) */}
      {translateX !== 0 && (
        <>
          {/* Play indicator (swipe right) */}
          {translateX > 0 && (
            <div
              className="absolute left-0 top-0 bottom-0 bg-success/20 flex items-center pl-4"
              style={{ width: Math.abs(translateX) }}
            >
              <span className="text-success text-2xl">▶</span>
            </div>
          )}
          {/* Delete indicator (swipe left) */}
          {translateX < 0 && (
            <div
              className="absolute right-0 top-0 bottom-0 bg-error/20 flex items-center justify-end pr-4"
              style={{ width: Math.abs(translateX) }}
            >
              <svg className="w-6 h-6 text-error" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </div>
          )}
        </>
      )}

      <div className="card-body p-3 flex-row items-center gap-3">
        {/* Cover art */}
        <div className="avatar flex-shrink-0">
          <div className="w-14 h-14 rounded-lg bg-base-300">
            {track.coverArtUrl ? (
              <img src={track.coverArtUrl} alt="" className="object-cover" />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-2xl">
                🎵
              </div>
            )}
          </div>
        </div>

        {/* Track info */}
        <div className="flex-1 min-w-0">
          <h3 className={`font-medium truncate ${isCurrentTrack ? 'text-primary' : ''}`}>
            {track.title}
          </h3>
          <p className="text-sm text-base-content/60 truncate">{track.artist}</p>
          <div className="flex items-center gap-2 mt-1 text-xs text-base-content/50">
            <span>{formatDuration(track.duration)}</span>
            {track.bpm && (
              <>
                <span>•</span>
                <span>{track.bpm} BPM</span>
              </>
            )}
            {track.musicalKey && (
              <>
                <span>•</span>
                <span>{track.musicalKey}</span>
              </>
            )}
          </div>
        </div>

        {/* Play button */}
        <button
          className={`btn btn-circle btn-sm ${isCurrentTrack && isPlaying ? 'btn-primary' : 'btn-ghost'}`}
          onClick={handlePlayClick}
          aria-label={isCurrentTrack && isPlaying ? 'Pause' : 'Play'}
        >
          {isCurrentTrack && isPlaying ? (
            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
            </svg>
          ) : (
            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z" />
            </svg>
          )}
        </button>
      </div>

      {/* Swipe hint (shown on first card) */}
      {index === 0 && (
        <div className="px-3 pb-2 text-xs text-base-content/40 text-center">
          Swipe right to play, left to delete
        </div>
      )}
    </div>
  );
}

interface TrackCardListProps {
  tracks: Track[];
  isLoading?: boolean;
  onDelete?: (trackId: string) => void;
}

export function TrackCardList({ tracks, isLoading, onDelete }: TrackCardListProps) {
  if (isLoading) {
    return (
      <div className="space-y-3">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="card bg-base-100 shadow-sm">
            <div className="card-body p-3 flex-row items-center gap-3">
              <div className="skeleton w-14 h-14 rounded-lg" />
              <div className="flex-1 space-y-2">
                <div className="skeleton h-4 w-3/4" />
                <div className="skeleton h-3 w-1/2" />
              </div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (tracks.length === 0) {
    return (
      <div className="text-center py-12 text-base-content/60">
        No tracks found
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {tracks.map((track, index) => (
        <TrackCard
          key={track.id}
          track={track}
          tracks={tracks}
          index={index}
          onDelete={onDelete}
        />
      ))}
    </div>
  );
}

export default TrackCard;
