/**
 * Keyboard Shortcuts Hook
 * Provides global keyboard shortcut handling for the application
 */
import { useEffect, useCallback, useState, useRef } from 'react';
import { usePlayerStore } from '@/lib/store/playerStore';
import { defaultShortcuts, matchesShortcut, type ShortcutConfig } from '@/lib/shortcuts/shortcuts';

interface UseKeyboardShortcutsOptions {
  enabled?: boolean;
  onShowShortcuts?: () => void;
  onEscape?: () => void;
  onFocusSearch?: () => void;
}

interface UseKeyboardShortcutsReturn {
  enabled: boolean;
  setEnabled: (enabled: boolean) => void;
  shortcuts: ShortcutConfig[];
}

// Volume adjustment step (0.1 = 10%)
const VOLUME_STEP = 0.1;

// Seek step in seconds
const SEEK_STEP = 10;

/**
 * Check if the event target is an input element that should not trigger shortcuts
 */
function isInputElement(target: EventTarget | null): boolean {
  if (!target || !(target instanceof HTMLElement)) return false;

  const tagName = target.tagName.toLowerCase();
  const isInput = tagName === 'input' || tagName === 'textarea' || tagName === 'select';
  const isContentEditable = target.isContentEditable;

  return isInput || isContentEditable;
}

/**
 * Hook for handling global keyboard shortcuts
 */
export function useKeyboardShortcuts(options: UseKeyboardShortcutsOptions = {}): UseKeyboardShortcutsReturn {
  const { enabled: initialEnabled = true, onShowShortcuts, onEscape, onFocusSearch } = options;

  const [enabled, setEnabled] = useState(initialEnabled);
  const previousVolumeRef = useRef<number>(1); // For mute/unmute

  // Get player actions
  const {
    isPlaying,
    volume,
    progress,
    duration,
    play,
    pause,
    next,
    previous,
    setVolume,
    toggleShuffle,
    cycleRepeat,
    seek,
  } = usePlayerStore();

  /**
   * Execute an action based on the shortcut action identifier
   */
  const executeAction = useCallback(
    (action: string) => {
      switch (action) {
        case 'togglePlayPause':
          if (isPlaying) {
            pause();
          } else {
            play();
          }
          break;

        case 'nextTrack':
          next();
          break;

        case 'previousTrack':
          previous();
          break;

        case 'seekForward':
          if (duration > 0) {
            seek(Math.min(progress + SEEK_STEP, duration));
          }
          break;

        case 'seekBackward':
          seek(Math.max(progress - SEEK_STEP, 0));
          break;

        case 'volumeUp':
          setVolume(Math.min(volume + VOLUME_STEP, 1));
          break;

        case 'volumeDown':
          setVolume(Math.max(volume - VOLUME_STEP, 0));
          break;

        case 'toggleMute':
          if (volume > 0) {
            previousVolumeRef.current = volume;
            setVolume(0);
          } else {
            setVolume(previousVolumeRef.current || 1);
          }
          break;

        case 'toggleShuffle':
          toggleShuffle();
          break;

        case 'cycleRepeat':
          cycleRepeat();
          break;

        case 'focusSearch':
          onFocusSearch?.();
          break;

        case 'escape':
          onEscape?.();
          break;

        case 'showShortcuts':
          onShowShortcuts?.();
          break;
      }
    },
    [isPlaying, volume, progress, duration, play, pause, next, previous, setVolume, toggleShuffle, cycleRepeat, seek, onFocusSearch, onEscape, onShowShortcuts]
  );

  /**
   * Handle keydown events
   */
  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      // Don't handle if shortcuts are disabled
      if (!enabled) return;

      // Don't handle if target is an input element
      if (isInputElement(event.target)) return;

      // Find matching shortcut
      for (const shortcut of defaultShortcuts) {
        if (matchesShortcut(event, shortcut)) {
          // Prevent default for matched shortcuts
          event.preventDefault();
          event.stopPropagation();

          executeAction(shortcut.action);
          return;
        }
      }
    },
    [enabled, executeAction]
  );

  // Set up global keydown listener
  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [handleKeyDown]);

  return {
    enabled,
    setEnabled,
    shortcuts: defaultShortcuts,
  };
}

/**
 * Helper hook to check if Mac
 */
export function useIsMac(): boolean {
  const [isMac, setIsMac] = useState(false);

  useEffect(() => {
    setIsMac(navigator.platform.toLowerCase().includes('mac'));
  }, []);

  return isMac;
}
