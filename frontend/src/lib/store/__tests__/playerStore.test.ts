/**
 * Player Store Tests - REQ-5.6, REQ-5.7
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { act, renderHook } from '@testing-library/react';
import { usePlayerStore } from '@/lib/store/playerStore';
import type { Track } from '@/types';

const mockTrack: Track = {
  id: 'track-1',
  title: 'Test Track',
  artist: 'Test Artist',
  album: 'Test Album',
  duration: 180,
  format: 'mp3',
  fileSize: 5000000,
  s3Key: 'tracks/1.mp3',
  tags: [],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};

const mockTrack2: Track = { ...mockTrack, id: 'track-2', title: 'Track Two' };

describe('Player Store', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Reset store to initial state
    usePlayerStore.setState({
      currentTrack: null,
      queue: [],
      queueIndex: 0,
      isPlaying: false,
      volume: 1,
      progress: 0,
      shuffle: false,
      repeat: 'off',
    });
  });

  describe('Initial State', () => {
    it('REQ-5.6: should have no current track', () => {
      const { result } = renderHook(() => usePlayerStore());
      expect(result.current.currentTrack).toBeNull();
    });

    it('should not be playing initially', () => {
      const { result } = renderHook(() => usePlayerStore());
      expect(result.current.isPlaying).toBe(false);
    });

    it('should have empty queue', () => {
      const { result } = renderHook(() => usePlayerStore());
      expect(result.current.queue).toEqual([]);
    });

    it('should have default volume at 1', () => {
      const { result } = renderHook(() => usePlayerStore());
      expect(result.current.volume).toBe(1);
    });
  });

  describe('setQueue', () => {
    it('REQ-5.7: should set queue and play first track', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.setQueue([mockTrack, mockTrack2], 0);
      });

      expect(result.current.queue).toHaveLength(2);
      expect(result.current.currentTrack?.id).toBe('track-1');
    });

    it('should set queue at specific index', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.setQueue([mockTrack, mockTrack2], 1);
      });

      expect(result.current.currentTrack?.id).toBe('track-2');
    });
  });

  describe('play/pause', () => {
    it('REQ-5.6: should toggle isPlaying', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.setQueue([mockTrack], 0);
        result.current.play();
      });

      expect(result.current.isPlaying).toBe(true);

      act(() => {
        result.current.pause();
      });

      expect(result.current.isPlaying).toBe(false);
    });
  });

  describe('next/previous', () => {
    it('REQ-5.6: should play next track', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.setQueue([mockTrack, mockTrack2], 0);
        result.current.next();
      });

      expect(result.current.currentTrack?.id).toBe('track-2');
    });

    it('should play previous track', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.setQueue([mockTrack, mockTrack2], 1);
        result.current.previous();
      });

      expect(result.current.currentTrack?.id).toBe('track-1');
    });
  });

  describe('volume', () => {
    it('REQ-5.6: should set volume', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.setVolume(0.5);
      });

      expect(result.current.volume).toBe(0.5);
    });

    it('should clamp volume to 0-1 range', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.setVolume(1.5);
      });

      expect(result.current.volume).toBeLessThanOrEqual(1);
    });
  });

  describe('shuffle', () => {
    it('REQ-5.7: should toggle shuffle', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.toggleShuffle();
      });

      expect(result.current.shuffle).toBe(true);
    });
  });

  describe('repeat', () => {
    it('REQ-5.7: should cycle repeat modes', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.cycleRepeat();
      });
      expect(result.current.repeat).toBe('all');

      act(() => {
        result.current.cycleRepeat();
      });
      expect(result.current.repeat).toBe('one');

      act(() => {
        result.current.cycleRepeat();
      });
      expect(result.current.repeat).toBe('off');
    });
  });

  describe('seek', () => {
    it('REQ-5.6: should update progress', () => {
      const { result } = renderHook(() => usePlayerStore());

      act(() => {
        result.current.setQueue([mockTrack], 0);
        result.current.seek(60);
      });

      expect(result.current.progress).toBe(60);
    });
  });
});
