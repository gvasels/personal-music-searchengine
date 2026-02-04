import { useState } from 'react';
import { usePlayerStore } from '@/lib/store/playerStore';
import { useAuth } from '@/hooks/useAuth';
import { Waveform } from './Waveform';
import { Equalizer } from './Equalizer';

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function PlayerBar() {
  const [showEQ, setShowEQ] = useState(false);
  const { isAuthenticated } = useAuth();
  const {
    currentTrack,
    isPlaying,
    volume,
    progress,
    shuffle,
    repeat,
    play,
    pause,
    next,
    previous,
    seek,
    setVolume,
    toggleShuffle,
    cycleRepeat,
  } = usePlayerStore();

  // Don't show player bar when not authenticated
  if (!isAuthenticated) {
    return null;
  }

  return (
    <div
      data-testid="player-bar"
      className="fixed bottom-0 left-0 right-0 h-24 bg-base-100 border-t border-base-300 flex items-center px-4 z-50"
    >
      {/* Track Info */}
      <div className="flex items-center gap-3 w-1/4">
        {currentTrack && (
          <>
            {currentTrack.coverArtUrl ? (
              <img
                src={currentTrack.coverArtUrl}
                alt={currentTrack.album}
                className="w-14 h-14 rounded object-cover"
              />
            ) : (
              <div className="w-14 h-14 bg-base-300 rounded flex items-center justify-center">
                🎵
              </div>
            )}
            <div className="min-w-0">
              <div className="font-medium truncate">{currentTrack.title}</div>
              <div className="text-sm text-base-content/70 truncate">
                {currentTrack.artist}
              </div>
            </div>
          </>
        )}
      </div>

      {/* Controls */}
      <div className="flex flex-col items-center flex-1">
        <div className="flex items-center gap-4 player-controls">
          <button
            onClick={toggleShuffle}
            className={`btn btn-ghost btn-sm ${shuffle ? 'text-primary' : ''}`}
            aria-label="shuffle"
          >
            🔀
          </button>
          <button
            onClick={previous}
            className="btn btn-ghost btn-circle"
            aria-label="previous"
          >
            ⏮️
          </button>
          <button
            onClick={isPlaying ? pause : play}
            className="btn btn-primary btn-circle"
            aria-label={isPlaying ? 'pause' : 'play'}
          >
            {isPlaying ? '⏸️' : '▶️'}
          </button>
          <button
            onClick={next}
            className="btn btn-ghost btn-circle"
            aria-label="next"
          >
            ⏭️
          </button>
          <button
            onClick={cycleRepeat}
            className={`btn btn-ghost btn-sm ${repeat !== 'off' ? 'text-primary' : ''}`}
            aria-label="repeat"
          >
            {repeat === 'one' ? '🔂' : '🔁'}
          </button>
        </div>
        <div className="flex items-center gap-2 w-full max-w-md mt-1">
          <span className="text-xs tabular-nums">{formatTime(progress)}</span>
          <Waveform
            trackId={currentTrack?.id}
            duration={currentTrack?.duration || 0}
            progress={progress}
            onSeek={seek}
            height={32}
            className="flex-1"
          />
          <span className="text-xs tabular-nums">
            {formatTime(currentTrack?.duration || 0)}
          </span>
        </div>
      </div>

      {/* Volume & EQ */}
      <div className="flex items-center gap-2 w-1/4 justify-end">
        <Equalizer compact className="mr-2" />
        <button
          className={`btn btn-ghost btn-sm ${showEQ ? 'text-primary' : ''}`}
          onClick={() => setShowEQ(!showEQ)}
          aria-label="equalizer settings"
          title="Equalizer"
        >
          ⚙️
        </button>
        <button className="btn btn-ghost btn-sm" aria-label="volume">
          {volume === 0 ? '🔇' : volume < 0.5 ? '🔉' : '🔊'}
        </button>
        <input
          type="range"
          min={0}
          max={1}
          step={0.01}
          value={volume}
          onChange={(e) => setVolume(Number(e.target.value))}
          className="range range-xs w-24"
          aria-label="volume"
        />
      </div>

      {/* EQ Panel */}
      {showEQ && (
        <div className="absolute bottom-24 right-4 z-50">
          <div className="relative">
            <button
              className="absolute -top-2 -right-2 btn btn-circle btn-xs btn-ghost"
              onClick={() => setShowEQ(false)}
              aria-label="close equalizer"
            >
              ✕
            </button>
            <Equalizer className="shadow-xl" />
          </div>
        </div>
      )}
    </div>
  );
}
