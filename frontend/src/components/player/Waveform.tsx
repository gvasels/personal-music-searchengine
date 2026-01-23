/**
 * Waveform Visualization Component
 * Canvas-based waveform display with progress overlay and seek functionality
 */
import { useRef, useEffect, useCallback, useState } from 'react';
import { useWaveform, generateMockWaveform, type WaveformData } from '@/hooks/useWaveform';

interface WaveformProps {
  trackId: string | undefined;
  duration: number;
  progress: number;
  onSeek: (position: number) => void;
  height?: number;
  className?: string;
}

export function Waveform({
  trackId,
  duration,
  progress,
  onSeek,
  height = 48,
  className = '',
}: WaveformProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [containerWidth, setContainerWidth] = useState(0);
  const [hoverPosition, setHoverPosition] = useState<number | null>(null);
  const [isHovering, setIsHovering] = useState(false);

  const { data: waveformData, isLoading } = useWaveform(trackId);

  // Use mock data if no real waveform available
  const effectiveWaveform: WaveformData | null = waveformData || (duration > 0 ? generateMockWaveform(duration) : null);

  // Track container width with ResizeObserver
  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        setContainerWidth(entry.contentRect.width);
      }
    });

    observer.observe(container);
    setContainerWidth(container.clientWidth);

    return () => observer.disconnect();
  }, []);

  // Draw waveform
  const drawWaveform = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas || !effectiveWaveform || containerWidth === 0) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const dpr = window.devicePixelRatio || 1;
    const width = containerWidth;

    // Set canvas size accounting for device pixel ratio
    canvas.width = width * dpr;
    canvas.height = height * dpr;
    canvas.style.width = `${width}px`;
    canvas.style.height = `${height}px`;
    ctx.scale(dpr, dpr);

    // Clear canvas
    ctx.clearRect(0, 0, width, height);

    const { peaks } = effectiveWaveform;
    const barWidth = Math.max(2, width / peaks.length);
    const barGap = 1;
    const effectiveBarWidth = barWidth - barGap;
    const progressPercent = duration > 0 ? progress / duration : 0;
    const progressX = width * progressPercent;

    // Draw each bar
    for (let i = 0; i < peaks.length; i++) {
      const x = (i / peaks.length) * width;
      const amplitude = peaks[i];
      const barHeight = amplitude * (height - 4);
      const y = (height - barHeight) / 2;

      // Determine color based on progress
      if (x < progressX) {
        ctx.fillStyle = 'hsl(var(--p))'; // Primary color (played)
      } else {
        ctx.fillStyle = 'hsl(var(--bc) / 0.3)'; // Base content with opacity (unplayed)
      }

      ctx.fillRect(x, y, effectiveBarWidth, barHeight);
    }

    // Draw hover indicator
    if (isHovering && hoverPosition !== null) {
      ctx.fillStyle = 'hsl(var(--p) / 0.5)';
      ctx.fillRect(hoverPosition, 0, 2, height);
    }
  }, [effectiveWaveform, containerWidth, height, progress, duration, isHovering, hoverPosition]);

  // Redraw on changes
  useEffect(() => {
    drawWaveform();
  }, [drawWaveform]);

  // Handle mouse/touch interactions
  const handleInteraction = useCallback(
    (clientX: number) => {
      const container = containerRef.current;
      if (!container || duration === 0) return;

      const rect = container.getBoundingClientRect();
      const x = clientX - rect.left;
      const percent = Math.max(0, Math.min(1, x / rect.width));
      return percent * duration;
    },
    [duration]
  );

  const handleMouseMove = useCallback(
    (e: React.MouseEvent) => {
      const container = containerRef.current;
      if (!container) return;

      const rect = container.getBoundingClientRect();
      setHoverPosition(e.clientX - rect.left);
    },
    []
  );

  const handleClick = useCallback(
    (e: React.MouseEvent) => {
      const position = handleInteraction(e.clientX);
      if (position !== undefined) {
        onSeek(position);
      }
    },
    [handleInteraction, onSeek]
  );

  const handleTouchEnd = useCallback(
    (e: React.TouchEvent) => {
      if (e.changedTouches.length > 0) {
        const position = handleInteraction(e.changedTouches[0].clientX);
        if (position !== undefined) {
          onSeek(position);
        }
      }
    },
    [handleInteraction, onSeek]
  );

  if (isLoading) {
    return (
      <div
        className={`flex items-center justify-center bg-base-200 rounded ${className}`}
        style={{ height }}
      >
        <span className="loading loading-dots loading-sm" />
      </div>
    );
  }

  return (
    <div
      ref={containerRef}
      className={`relative cursor-pointer select-none ${className}`}
      style={{ height }}
      onClick={handleClick}
      onMouseEnter={() => setIsHovering(true)}
      onMouseLeave={() => {
        setIsHovering(false);
        setHoverPosition(null);
      }}
      onMouseMove={handleMouseMove}
      onTouchEnd={handleTouchEnd}
      role="slider"
      aria-label="Waveform seek"
      aria-valuemin={0}
      aria-valuemax={duration}
      aria-valuenow={progress}
    >
      <canvas
        ref={canvasRef}
        className="w-full h-full rounded"
      />

      {/* Hover time tooltip */}
      {isHovering && hoverPosition !== null && duration > 0 && (
        <div
          className="absolute -top-8 transform -translate-x-1/2 bg-base-300 text-xs px-2 py-1 rounded pointer-events-none"
          style={{ left: hoverPosition }}
        >
          {formatTime((hoverPosition / containerWidth) * duration)}
        </div>
      )}
    </div>
  );
}

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export default Waveform;
