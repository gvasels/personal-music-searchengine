/**
 * Keyboard Shortcuts Configuration
 * Defines all keyboard shortcuts and their metadata
 */

export type ShortcutCategory = 'playback' | 'volume' | 'navigation' | 'modes' | 'ui';

export interface ShortcutConfig {
  key: string;
  code?: string; // KeyboardEvent.code for position-independent keys
  modifiers?: {
    ctrl?: boolean;
    meta?: boolean; // Cmd on Mac
    shift?: boolean;
    alt?: boolean;
  };
  description: string;
  category: ShortcutCategory;
  action: string; // Action identifier to dispatch
}

/**
 * Default keyboard shortcuts
 * Key is used for display, code for matching (when specified)
 */
export const defaultShortcuts: ShortcutConfig[] = [
  // Playback controls
  {
    key: 'Space',
    code: 'Space',
    description: 'Play / Pause',
    category: 'playback',
    action: 'togglePlayPause',
  },
  {
    key: '→',
    code: 'ArrowRight',
    description: 'Next track',
    category: 'playback',
    action: 'nextTrack',
  },
  {
    key: '←',
    code: 'ArrowLeft',
    description: 'Previous track / Restart',
    category: 'playback',
    action: 'previousTrack',
  },
  {
    key: 'Shift+→',
    code: 'ArrowRight',
    modifiers: { shift: true },
    description: 'Seek forward 10s',
    category: 'playback',
    action: 'seekForward',
  },
  {
    key: 'Shift+←',
    code: 'ArrowLeft',
    modifiers: { shift: true },
    description: 'Seek backward 10s',
    category: 'playback',
    action: 'seekBackward',
  },

  // Volume controls
  {
    key: '↑',
    code: 'ArrowUp',
    description: 'Volume up',
    category: 'volume',
    action: 'volumeUp',
  },
  {
    key: '↓',
    code: 'ArrowDown',
    description: 'Volume down',
    category: 'volume',
    action: 'volumeDown',
  },
  {
    key: 'M',
    description: 'Mute / Unmute',
    category: 'volume',
    action: 'toggleMute',
  },

  // Mode controls
  {
    key: 'S',
    description: 'Toggle shuffle',
    category: 'modes',
    action: 'toggleShuffle',
  },
  {
    key: 'R',
    description: 'Cycle repeat mode',
    category: 'modes',
    action: 'cycleRepeat',
  },

  // Navigation
  {
    key: '/',
    description: 'Focus search',
    category: 'navigation',
    action: 'focusSearch',
  },
  {
    key: 'Escape',
    code: 'Escape',
    description: 'Clear selection / Close modal',
    category: 'navigation',
    action: 'escape',
  },

  // UI
  {
    key: '?',
    modifiers: { shift: true },
    description: 'Show keyboard shortcuts',
    category: 'ui',
    action: 'showShortcuts',
  },
];

/**
 * Category labels for display
 */
export const categoryLabels: Record<ShortcutCategory, string> = {
  playback: 'Playback',
  volume: 'Volume',
  navigation: 'Navigation',
  modes: 'Modes',
  ui: 'Interface',
};

/**
 * Get shortcuts grouped by category
 */
export function getShortcutsByCategory(): Record<ShortcutCategory, ShortcutConfig[]> {
  const grouped: Record<ShortcutCategory, ShortcutConfig[]> = {
    playback: [],
    volume: [],
    navigation: [],
    modes: [],
    ui: [],
  };

  for (const shortcut of defaultShortcuts) {
    grouped[shortcut.category].push(shortcut);
  }

  return grouped;
}

/**
 * Format key for display (handles Mac vs Windows)
 */
export function formatKeyForDisplay(shortcut: ShortcutConfig, isMac: boolean): string {
  const parts: string[] = [];

  if (shortcut.modifiers?.ctrl) {
    parts.push(isMac ? '⌃' : 'Ctrl');
  }
  if (shortcut.modifiers?.meta) {
    parts.push(isMac ? '⌘' : 'Ctrl');
  }
  if (shortcut.modifiers?.alt) {
    parts.push(isMac ? '⌥' : 'Alt');
  }
  if (shortcut.modifiers?.shift) {
    parts.push(isMac ? '⇧' : 'Shift');
  }

  parts.push(shortcut.key);

  return parts.join(isMac ? '' : '+');
}

/**
 * Check if a keyboard event matches a shortcut config
 */
export function matchesShortcut(event: KeyboardEvent, shortcut: ShortcutConfig): boolean {
  // Check modifiers
  const modifiers = shortcut.modifiers || {};
  if (!!modifiers.ctrl !== event.ctrlKey) return false;
  if (!!modifiers.meta !== event.metaKey) return false;
  if (!!modifiers.shift !== event.shiftKey) return false;
  if (!!modifiers.alt !== event.altKey) return false;

  // Check key/code
  if (shortcut.code) {
    return event.code === shortcut.code;
  }

  // Fall back to key for letter keys
  return event.key.toUpperCase() === shortcut.key.toUpperCase();
}
