/**
 * User Preferences Store
 * Manages user preferences with localStorage persistence
 */
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

// Track list column options
export type TrackListColumn = 'title' | 'artist' | 'album' | 'duration' | 'bpm' | 'key' | 'genre' | 'addedDate';

export interface PreferencesState {
  // Display preferences
  sidebarVisible: boolean;
  trackListColumns: TrackListColumn[];
  compactMode: boolean;
  showCoverArt: boolean;
  showUploadedBy: boolean; // Show "Uploaded By" column in track list (for admin/global users)

  // Behavior preferences
  shortcutsEnabled: boolean;
  confirmBeforeDelete: boolean;
  autoPlayOnSelect: boolean;

  // Audio preferences
  defaultVolume: number;
  crossfadeDuration: number;
  normalizeAudio: boolean;

  // Actions
  setSidebarVisible: (visible: boolean) => void;
  toggleSidebar: () => void;
  setTrackListColumns: (columns: TrackListColumn[]) => void;
  toggleColumn: (column: TrackListColumn) => void;
  setCompactMode: (compact: boolean) => void;
  setShowCoverArt: (show: boolean) => void;
  setShowUploadedBy: (show: boolean) => void;
  setShortcutsEnabled: (enabled: boolean) => void;
  setConfirmBeforeDelete: (confirm: boolean) => void;
  setAutoPlayOnSelect: (autoPlay: boolean) => void;
  setDefaultVolume: (volume: number) => void;
  setCrossfadeDuration: (duration: number) => void;
  setNormalizeAudio: (normalize: boolean) => void;
  resetToDefaults: () => void;
}

const DEFAULT_COLUMNS: TrackListColumn[] = ['title', 'artist', 'album', 'duration'];

const defaultPreferences = {
  sidebarVisible: true,
  trackListColumns: DEFAULT_COLUMNS,
  compactMode: false,
  showCoverArt: true,
  showUploadedBy: false, // Hidden by default, toggle in settings
  shortcutsEnabled: true,
  confirmBeforeDelete: true,
  autoPlayOnSelect: true,
  defaultVolume: 0.8,
  crossfadeDuration: 0,
  normalizeAudio: false,
};

export const usePreferencesStore = create<PreferencesState>()(
  persist(
    (set, get) => ({
      // Initial state
      ...defaultPreferences,

      // Display actions
      setSidebarVisible: (visible) => set({ sidebarVisible: visible }),
      toggleSidebar: () => set((state) => ({ sidebarVisible: !state.sidebarVisible })),

      setTrackListColumns: (columns) => set({ trackListColumns: columns }),
      toggleColumn: (column) => {
        const current = get().trackListColumns;
        if (current.includes(column)) {
          // Don't allow removing the last column
          if (current.length > 1) {
            set({ trackListColumns: current.filter((c) => c !== column) });
          }
        } else {
          set({ trackListColumns: [...current, column] });
        }
      },

      setCompactMode: (compact) => set({ compactMode: compact }),
      setShowCoverArt: (show) => set({ showCoverArt: show }),
      setShowUploadedBy: (show) => set({ showUploadedBy: show }),

      // Behavior actions
      setShortcutsEnabled: (enabled) => set({ shortcutsEnabled: enabled }),
      setConfirmBeforeDelete: (confirm) => set({ confirmBeforeDelete: confirm }),
      setAutoPlayOnSelect: (autoPlay) => set({ autoPlayOnSelect: autoPlay }),

      // Audio actions
      setDefaultVolume: (volume) => set({ defaultVolume: Math.max(0, Math.min(1, volume)) }),
      setCrossfadeDuration: (duration) => set({ crossfadeDuration: Math.max(0, Math.min(12, duration)) }),
      setNormalizeAudio: (normalize) => set({ normalizeAudio: normalize }),

      // Reset
      resetToDefaults: () => set(defaultPreferences),
    }),
    {
      name: 'user-preferences',
      storage: createJSONStorage(() => localStorage),
      version: 1,
      migrate: (persistedState, version) => {
        // Handle migration for future schema changes
        if (version === 0) {
          // Migration from version 0 to 1
          return { ...defaultPreferences, ...(persistedState as Partial<PreferencesState>) };
        }
        return persistedState as PreferencesState;
      },
    }
  )
);

// Selector hooks for specific preferences
export const useSidebarVisible = () => usePreferencesStore((s) => s.sidebarVisible);
export const useShortcutsEnabled = () => usePreferencesStore((s) => s.shortcutsEnabled);
export const useTrackListColumns = () => usePreferencesStore((s) => s.trackListColumns);
export const useCompactMode = () => usePreferencesStore((s) => s.compactMode);
export const useShowUploadedBy = () => usePreferencesStore((s) => s.showUploadedBy);
