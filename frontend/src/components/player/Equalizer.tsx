/**
 * Equalizer Component
 * 3-band EQ with presets and custom settings
 */
import { useState } from 'react';
import { useEQStore, useAllPresets } from '@/lib/store/eqStore';

interface EQSliderProps {
  label: string;
  frequency: string;
  value: number;
  onChange: (value: number) => void;
  disabled?: boolean;
}

function EQSlider({ label, frequency, value, onChange, disabled }: EQSliderProps) {
  return (
    <div className="flex flex-col items-center gap-1">
      <span className="text-xs font-medium">{label}</span>
      <span className="text-xs text-base-content/50">{frequency}</span>
      <div className="relative h-32 flex flex-col items-center justify-center">
        <input
          type="range"
          min={-12}
          max={12}
          step={1}
          value={value}
          onChange={(e) => onChange(Number(e.target.value))}
          disabled={disabled}
          className="range range-primary range-sm h-24"
          style={{
            writingMode: 'vertical-lr',
            direction: 'rtl',
          }}
          aria-label={`${label} ${value > 0 ? '+' : ''}${value}dB`}
        />
      </div>
      <span className="text-xs font-mono tabular-nums">
        {value > 0 ? '+' : ''}{value}
      </span>
    </div>
  );
}

interface EqualizerProps {
  compact?: boolean;
  className?: string;
}

export function Equalizer({ compact = false, className = '' }: EqualizerProps) {
  const {
    enabled,
    bands,
    selectedPreset,
    setEnabled,
    setBand,
    applyPreset,
    saveCustomPreset,
    deleteCustomPreset,
    reset,
  } = useEQStore();

  const allPresets = useAllPresets();
  const [showSaveDialog, setShowSaveDialog] = useState(false);
  const [newPresetName, setNewPresetName] = useState('');

  const handleSavePreset = () => {
    if (newPresetName.trim()) {
      saveCustomPreset(newPresetName.trim());
      setNewPresetName('');
      setShowSaveDialog(false);
    }
  };

  if (compact) {
    return (
      <div className={`flex items-center gap-2 ${className}`}>
        <label className="label cursor-pointer gap-2">
          <span className="label-text text-xs">EQ</span>
          <input
            type="checkbox"
            className="toggle toggle-xs toggle-primary"
            checked={enabled}
            onChange={(e) => setEnabled(e.target.checked)}
          />
        </label>
        {enabled && (
          <select
            className="select select-xs select-bordered"
            value={selectedPreset || ''}
            onChange={(e) => applyPreset(e.target.value)}
          >
            <option value="" disabled>Preset</option>
            {allPresets.map((preset) => (
              <option key={preset.name} value={preset.name}>
                {preset.name}
              </option>
            ))}
          </select>
        )}
      </div>
    );
  }

  return (
    <div className={`card bg-base-200 ${className}`}>
      <div className="card-body p-4">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <h3 className="font-semibold">Equalizer</h3>
            <label className="label cursor-pointer gap-2">
              <input
                type="checkbox"
                className="toggle toggle-sm toggle-primary"
                checked={enabled}
                onChange={(e) => setEnabled(e.target.checked)}
              />
            </label>
          </div>
          <button
            className="btn btn-ghost btn-xs"
            onClick={reset}
            disabled={!enabled}
          >
            Reset
          </button>
        </div>

        {/* Preset selector */}
        <div className="flex items-center gap-2 mb-4">
          <select
            className="select select-sm select-bordered flex-1"
            value={selectedPreset || ''}
            onChange={(e) => applyPreset(e.target.value)}
            disabled={!enabled}
          >
            <option value="" disabled>Select preset...</option>
            <optgroup label="Built-in">
              {allPresets.filter((p) => p.isBuiltIn).map((preset) => (
                <option key={preset.name} value={preset.name}>
                  {preset.name}
                </option>
              ))}
            </optgroup>
            {allPresets.some((p) => !p.isBuiltIn) && (
              <optgroup label="Custom">
                {allPresets.filter((p) => !p.isBuiltIn).map((preset) => (
                  <option key={preset.name} value={preset.name}>
                    {preset.name}
                  </option>
                ))}
              </optgroup>
            )}
          </select>
          <button
            className="btn btn-sm btn-outline"
            onClick={() => setShowSaveDialog(true)}
            disabled={!enabled}
          >
            Save
          </button>
          {selectedPreset && !allPresets.find((p) => p.name === selectedPreset)?.isBuiltIn && (
            <button
              className="btn btn-sm btn-ghost text-error"
              onClick={() => deleteCustomPreset(selectedPreset)}
              disabled={!enabled}
            >
              Delete
            </button>
          )}
        </div>

        {/* EQ Sliders */}
        <div className="flex justify-center gap-8">
          <EQSlider
            label="Bass"
            frequency="100Hz"
            value={bands.bass}
            onChange={(v) => setBand('bass', v)}
            disabled={!enabled}
          />
          <EQSlider
            label="Mid"
            frequency="1kHz"
            value={bands.mid}
            onChange={(v) => setBand('mid', v)}
            disabled={!enabled}
          />
          <EQSlider
            label="Treble"
            frequency="8kHz"
            value={bands.treble}
            onChange={(v) => setBand('treble', v)}
            disabled={!enabled}
          />
        </div>

        {/* Visual representation */}
        <div className="mt-4 h-16 bg-base-300 rounded-lg flex items-end justify-around px-4 pb-2">
          {Object.entries(bands).map(([band, value]) => (
            <div
              key={band}
              className="w-8 bg-primary rounded-t transition-all duration-150"
              style={{
                height: `${((value + 12) / 24) * 100}%`,
                opacity: enabled ? 1 : 0.3,
              }}
            />
          ))}
        </div>

        {!enabled && (
          <p className="text-center text-xs text-base-content/50 mt-2">
            Enable EQ to adjust audio settings
          </p>
        )}
      </div>

      {/* Save preset dialog */}
      {showSaveDialog && (
        <div className="modal modal-open">
          <div className="modal-box max-w-xs">
            <h3 className="font-bold text-lg">Save Preset</h3>
            <div className="py-4">
              <input
                type="text"
                placeholder="Preset name"
                className="input input-bordered w-full"
                value={newPresetName}
                onChange={(e) => setNewPresetName(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSavePreset()}
                autoFocus
              />
            </div>
            <div className="modal-action">
              <button className="btn btn-ghost" onClick={() => setShowSaveDialog(false)}>
                Cancel
              </button>
              <button
                className="btn btn-primary"
                onClick={handleSavePreset}
                disabled={!newPresetName.trim()}
              >
                Save
              </button>
            </div>
          </div>
          <div className="modal-backdrop" onClick={() => setShowSaveDialog(false)} />
        </div>
      )}
    </div>
  );
}

export default Equalizer;
