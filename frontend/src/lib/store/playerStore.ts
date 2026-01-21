import { create } from 'zustand';
import type { Track } from '@/types';

interface PlayerState {
  currentTrack: Track | null;
  queue: Track[];
  queueIndex: number;
  isPlaying: boolean;
  volume: number;
  progress: number;
  shuffle: boolean;
  repeat: 'off' | 'all' | 'one';
}

interface PlayerActions {
  setQueue: (tracks: Track[], startIndex?: number) => void;
  play: () => void;
  pause: () => void;
  next: () => void;
  previous: () => void;
  setVolume: (volume: number) => void;
  toggleShuffle: () => void;
  cycleRepeat: () => void;
  seek: (position: number) => void;
}

type PlayerStore = PlayerState & PlayerActions;

export const usePlayerStore = create<PlayerStore>((set, get) => ({
  currentTrack: null,
  queue: [],
  queueIndex: 0,
  isPlaying: false,
  volume: 1,
  progress: 0,
  shuffle: false,
  repeat: 'off',

  setQueue: (tracks: Track[], startIndex = 0) => {
    set({
      queue: tracks,
      queueIndex: startIndex,
      currentTrack: tracks[startIndex] || null,
    });
  },

  play: () => set({ isPlaying: true }),
  pause: () => set({ isPlaying: false }),

  next: () => {
    const { queue, queueIndex, repeat } = get();
    if (queue.length === 0) return;

    let nextIndex = queueIndex + 1;
    if (nextIndex >= queue.length) {
      if (repeat === 'all') {
        nextIndex = 0;
      } else {
        return;
      }
    }

    set({
      queueIndex: nextIndex,
      currentTrack: queue[nextIndex],
      progress: 0,
    });
  },

  previous: () => {
    const { queue, queueIndex } = get();
    if (queue.length === 0 || queueIndex <= 0) return;

    const prevIndex = queueIndex - 1;
    set({
      queueIndex: prevIndex,
      currentTrack: queue[prevIndex],
      progress: 0,
    });
  },

  setVolume: (volume: number) => {
    set({ volume: Math.max(0, Math.min(1, volume)) });
  },

  toggleShuffle: () => set((state) => ({ shuffle: !state.shuffle })),

  cycleRepeat: () => {
    set((state) => {
      const cycle: Record<'off' | 'all' | 'one', 'off' | 'all' | 'one'> = {
        off: 'all',
        all: 'one',
        one: 'off',
      };
      return { repeat: cycle[state.repeat] };
    });
  },

  seek: (position: number) => set({ progress: position }),
}));
