import { create } from 'zustand';
import type { Track } from '@/types';
import { AudioService } from '../audio/audioService';
import { getStreamUrl } from '../api/client';

interface PlayerState {
  currentTrack: Track | null;
  queue: Track[];
  queueIndex: number;
  isPlaying: boolean;
  volume: number;
  progress: number;
  duration: number;
  shuffle: boolean;
  repeat: 'off' | 'all' | 'one';
  isLoading: boolean;
  _audioInitialized: boolean;
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
  _initAudio: () => void;
  _loadTrack: (track: Track) => Promise<void>;
}

type PlayerStore = PlayerState & PlayerActions;

export const usePlayerStore = create<PlayerStore>((set, get) => ({
  currentTrack: null,
  queue: [],
  queueIndex: 0,
  isPlaying: false,
  volume: 1,
  progress: 0,
  duration: 0,
  shuffle: false,
  repeat: 'off',
  isLoading: false,
  _audioInitialized: false,

  _initAudio: () => {
    const state = get();
    if (state._audioInitialized) return;

    const audioService = AudioService.getInstance();

    audioService.setCallbacks({
      onTimeUpdate: (currentTime, audioDuration) => {
        set({
          progress: currentTime,
          duration: audioDuration || 0
        });
      },
      onEnded: () => {
        const currentState = get();
        if (currentState.repeat === 'one') {
          // Repeat current track
          const audio = AudioService.getInstance();
          audio.seek(0);
          audio.play();
        } else {
          currentState.next();
        }
      },
      onError: (error) => {
        console.error('Audio playback error:', error);
        set({ isPlaying: false, isLoading: false });
      },
      onDurationChange: (audioDuration) => {
        set({ duration: audioDuration });
      },
      onCanPlay: () => {
        set({ isLoading: false });
      },
    });

    set({ _audioInitialized: true });
  },

  _loadTrack: async (track: Track) => {
    get()._initAudio();
    set({ isLoading: true });

    try {
      const { streamUrl } = await getStreamUrl(track.id);
      const audioService = AudioService.getInstance();
      audioService.load(streamUrl);
      await audioService.play();
      set({ isPlaying: true, isLoading: false, progress: 0 });
    } catch (error) {
      console.error('Failed to load track:', error);
      set({ isPlaying: false, isLoading: false });
    }
  },

  setQueue: (tracks: Track[], startIndex = 0) => {
    const track = tracks[startIndex] || null;
    set({
      queue: tracks,
      queueIndex: startIndex,
      currentTrack: track,
    });

    if (track) {
      get()._loadTrack(track);
    }
  },

  play: async () => {
    get()._initAudio();
    const state = get();

    if (!state.currentTrack && state.queue.length > 0) {
      // Start from beginning if no track is current
      const track = state.queue[0];
      set({ queueIndex: 0, currentTrack: track });
      await get()._loadTrack(track);
    } else {
      try {
        await AudioService.getInstance().play();
        set({ isPlaying: true });
      } catch (error) {
        console.error('Failed to play:', error);
      }
    }
  },

  pause: () => {
    AudioService.getInstance().pause();
    set({ isPlaying: false });
  },

  next: () => {
    const { queue, queueIndex, repeat, shuffle } = get();
    if (queue.length === 0) return;

    let nextIndex: number;

    if (shuffle) {
      // Pick random track (excluding current)
      const availableIndices = queue
        .map((_, i) => i)
        .filter((i) => i !== queueIndex);
      if (availableIndices.length === 0) {
        nextIndex = queueIndex;
      } else {
        nextIndex = availableIndices[Math.floor(Math.random() * availableIndices.length)];
      }
    } else {
      nextIndex = queueIndex + 1;

      if (nextIndex >= queue.length) {
        if (repeat === 'all') {
          nextIndex = 0;
        } else {
          // End of queue, stop playing
          AudioService.getInstance().pause();
          set({ isPlaying: false });
          return;
        }
      }
    }

    const nextTrack = queue[nextIndex];
    set({ queueIndex: nextIndex, currentTrack: nextTrack });
    get()._loadTrack(nextTrack);
  },

  previous: () => {
    const { queue, queueIndex, progress } = get();
    if (queue.length === 0) return;

    // If more than 3 seconds in, restart current track
    if (progress > 3) {
      AudioService.getInstance().seek(0);
      set({ progress: 0 });
      return;
    }

    if (queueIndex <= 0) return;

    const prevIndex = queueIndex - 1;
    const prevTrack = queue[prevIndex];
    set({ queueIndex: prevIndex, currentTrack: prevTrack });
    get()._loadTrack(prevTrack);
  },

  setVolume: (volume: number) => {
    const clampedVolume = Math.max(0, Math.min(1, volume));
    AudioService.getInstance().setVolume(clampedVolume);
    set({ volume: clampedVolume });
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

  seek: (position: number) => {
    AudioService.getInstance().seek(position);
    set({ progress: position });
  },
}));
