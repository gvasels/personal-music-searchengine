/**
 * Equalizer Store
 * Manages 3-band EQ settings with presets and localStorage persistence
 */
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

export interface EQBands {
  bass: number;    // -12 to +12 dB (100Hz lowshelf)
  mid: number;     // -12 to +12 dB (1kHz peaking)
  treble: number;  // -12 to +12 dB (8kHz highshelf)
}

export interface EQPreset {
  name: string;
  bands: EQBands;
  isBuiltIn: boolean;
}

export interface EQState {
  enabled: boolean;
  bands: EQBands;
  selectedPreset: string | null;
  customPresets: EQPreset[];

  // Actions
  setEnabled: (enabled: boolean) => void;
  setBand: (band: keyof EQBands, value: number) => void;
  setBands: (bands: EQBands) => void;
  applyPreset: (presetName: string) => void;
  saveCustomPreset: (name: string) => void;
  deleteCustomPreset: (name: string) => void;
  reset: () => void;
}

// Built-in presets
export const EQ_PRESETS: EQPreset[] = [
  { name: 'Flat', bands: { bass: 0, mid: 0, treble: 0 }, isBuiltIn: true },
  { name: 'Bass Boost', bands: { bass: 6, mid: 0, treble: 0 }, isBuiltIn: true },
  { name: 'Bass Reducer', bands: { bass: -6, mid: 0, treble: 0 }, isBuiltIn: true },
  { name: 'Treble Boost', bands: { bass: 0, mid: 0, treble: 6 }, isBuiltIn: true },
  { name: 'Treble Reducer', bands: { bass: 0, mid: 0, treble: -6 }, isBuiltIn: true },
  { name: 'Vocal', bands: { bass: -2, mid: 4, treble: 2 }, isBuiltIn: true },
  { name: 'Rock', bands: { bass: 4, mid: -2, treble: 4 }, isBuiltIn: true },
  { name: 'Electronic', bands: { bass: 6, mid: 0, treble: 4 }, isBuiltIn: true },
  { name: 'Jazz', bands: { bass: 2, mid: 2, treble: -2 }, isBuiltIn: true },
  { name: 'Classical', bands: { bass: 0, mid: 0, treble: 4 }, isBuiltIn: true },
];

const DEFAULT_BANDS: EQBands = { bass: 0, mid: 0, treble: 0 };

export const useEQStore = create<EQState>()(
  persist(
    (set, get) => ({
      enabled: false,
      bands: { ...DEFAULT_BANDS },
      selectedPreset: 'Flat',
      customPresets: [],

      setEnabled: (enabled) => set({ enabled }),

      setBand: (band, value) => {
        const clampedValue = Math.max(-12, Math.min(12, value));
        set((state) => ({
          bands: { ...state.bands, [band]: clampedValue },
          selectedPreset: null, // Clear preset when manually adjusting
        }));
      },

      setBands: (bands) => {
        set({
          bands: {
            bass: Math.max(-12, Math.min(12, bands.bass)),
            mid: Math.max(-12, Math.min(12, bands.mid)),
            treble: Math.max(-12, Math.min(12, bands.treble)),
          },
          selectedPreset: null,
        });
      },

      applyPreset: (presetName) => {
        const allPresets = [...EQ_PRESETS, ...get().customPresets];
        const preset = allPresets.find((p) => p.name === presetName);
        if (preset) {
          set({
            bands: { ...preset.bands },
            selectedPreset: presetName,
          });
        }
      },

      saveCustomPreset: (name) => {
        const { bands, customPresets } = get();
        const existingIndex = customPresets.findIndex((p) => p.name === name);

        if (existingIndex >= 0) {
          // Update existing
          const updated = [...customPresets];
          updated[existingIndex] = { name, bands: { ...bands }, isBuiltIn: false };
          set({ customPresets: updated, selectedPreset: name });
        } else {
          // Add new
          set({
            customPresets: [...customPresets, { name, bands: { ...bands }, isBuiltIn: false }],
            selectedPreset: name,
          });
        }
      },

      deleteCustomPreset: (name) => {
        set((state) => ({
          customPresets: state.customPresets.filter((p) => p.name !== name),
          selectedPreset: state.selectedPreset === name ? null : state.selectedPreset,
        }));
      },

      reset: () => {
        set({
          bands: { ...DEFAULT_BANDS },
          selectedPreset: 'Flat',
        });
      },
    }),
    {
      name: 'eq-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
);

// Helper to get all presets (built-in + custom)
export function useAllPresets(): EQPreset[] {
  const customPresets = useEQStore((s) => s.customPresets);
  return [...EQ_PRESETS, ...customPresets];
}
