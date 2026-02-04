/**
 * Settings Page
 * User preferences configuration
 */
import { usePreferencesStore, type TrackListColumn } from '../lib/store/preferencesStore';
import { useThemeStore } from '../lib/store/themeStore';

const COLUMN_OPTIONS: { value: TrackListColumn; label: string }[] = [
  { value: 'title', label: 'Title' },
  { value: 'artist', label: 'Artist' },
  { value: 'album', label: 'Album' },
  { value: 'duration', label: 'Duration' },
  { value: 'bpm', label: 'BPM' },
  { value: 'key', label: 'Key' },
  { value: 'genre', label: 'Genre' },
  { value: 'addedDate', label: 'Added Date' },
];

function SettingsPage() {
  const {
    sidebarVisible,
    trackListColumns,
    compactMode,
    showCoverArt,
    showUploadedBy,
    shortcutsEnabled,
    confirmBeforeDelete,
    autoPlayOnSelect,
    defaultVolume,
    crossfadeDuration,
    normalizeAudio,
    setSidebarVisible,
    toggleColumn,
    setCompactMode,
    setShowCoverArt,
    setShowUploadedBy,
    setShortcutsEnabled,
    setConfirmBeforeDelete,
    setAutoPlayOnSelect,
    setDefaultVolume,
    setCrossfadeDuration,
    setNormalizeAudio,
    resetToDefaults,
  } = usePreferencesStore();

  const { theme, toggleTheme } = useThemeStore();

  return (
    <div className="max-w-3xl mx-auto space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-base-content/60 mt-1">Customize your experience</p>
      </div>

      {/* Appearance */}
      <section className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">Appearance</h2>

          <div className="divider mt-0" />

          <div className="space-y-4">
            {/* Theme */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Dark Mode</span>
                  <span className="block text-sm text-base-content/60">
                    Use dark theme for the interface
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={theme === 'dark'}
                  onChange={toggleTheme}
                />
              </label>
            </div>

            {/* Sidebar */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Show Sidebar</span>
                  <span className="block text-sm text-base-content/60">
                    Display navigation sidebar on desktop
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={sidebarVisible}
                  onChange={(e) => setSidebarVisible(e.target.checked)}
                />
              </label>
            </div>

            {/* Compact Mode */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Compact Mode</span>
                  <span className="block text-sm text-base-content/60">
                    Use smaller spacing and text
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={compactMode}
                  onChange={(e) => setCompactMode(e.target.checked)}
                />
              </label>
            </div>

            {/* Show Cover Art */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Show Cover Art</span>
                  <span className="block text-sm text-base-content/60">
                    Display album art in track lists
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={showCoverArt}
                  onChange={(e) => setShowCoverArt(e.target.checked)}
                />
              </label>
            </div>

            {/* Show Uploaded By */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Show Uploaded By</span>
                  <span className="block text-sm text-base-content/60">
                    Display track uploader in track lists (admin/global users only)
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={showUploadedBy}
                  onChange={(e) => setShowUploadedBy(e.target.checked)}
                />
              </label>
            </div>
          </div>
        </div>
      </section>

      {/* Track List Columns */}
      <section className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">Track List Columns</h2>
          <p className="text-sm text-base-content/60">
            Choose which columns to display in track lists
          </p>

          <div className="divider mt-0" />

          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
            {COLUMN_OPTIONS.map((option) => (
              <label
                key={option.value}
                className={`
                  flex items-center gap-2 p-3 rounded-lg border cursor-pointer transition-colors
                  ${trackListColumns.includes(option.value)
                    ? 'border-primary bg-primary/10'
                    : 'border-base-300 hover:border-base-content/30'
                  }
                `}
              >
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm checkbox-primary"
                  checked={trackListColumns.includes(option.value)}
                  onChange={() => toggleColumn(option.value)}
                  disabled={trackListColumns.length === 1 && trackListColumns.includes(option.value)}
                />
                <span className="text-sm font-medium">{option.label}</span>
              </label>
            ))}
          </div>
        </div>
      </section>

      {/* Behavior */}
      <section className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">Behavior</h2>

          <div className="divider mt-0" />

          <div className="space-y-4">
            {/* Keyboard Shortcuts */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Keyboard Shortcuts</span>
                  <span className="block text-sm text-base-content/60">
                    Enable keyboard shortcuts (press ? to see all)
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={shortcutsEnabled}
                  onChange={(e) => setShortcutsEnabled(e.target.checked)}
                />
              </label>
            </div>

            {/* Confirm Before Delete */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Confirm Before Delete</span>
                  <span className="block text-sm text-base-content/60">
                    Ask for confirmation before deleting items
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={confirmBeforeDelete}
                  onChange={(e) => setConfirmBeforeDelete(e.target.checked)}
                />
              </label>
            </div>

            {/* Auto Play on Select */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Auto Play on Select</span>
                  <span className="block text-sm text-base-content/60">
                    Start playback when selecting a track
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={autoPlayOnSelect}
                  onChange={(e) => setAutoPlayOnSelect(e.target.checked)}
                />
              </label>
            </div>
          </div>
        </div>
      </section>

      {/* Audio */}
      <section className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">Audio</h2>

          <div className="divider mt-0" />

          <div className="space-y-6">
            {/* Default Volume */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">Default Volume</span>
                <span className="label-text-alt">{Math.round(defaultVolume * 100)}%</span>
              </label>
              <input
                type="range"
                className="range range-primary"
                min="0"
                max="100"
                value={defaultVolume * 100}
                onChange={(e) => setDefaultVolume(Number(e.target.value) / 100)}
              />
            </div>

            {/* Crossfade Duration */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">Crossfade Duration</span>
                <span className="label-text-alt">
                  {crossfadeDuration === 0 ? 'Off' : `${crossfadeDuration}s`}
                </span>
              </label>
              <input
                type="range"
                className="range range-primary"
                min="0"
                max="12"
                step="1"
                value={crossfadeDuration}
                onChange={(e) => setCrossfadeDuration(Number(e.target.value))}
              />
              <div className="flex justify-between text-xs text-base-content/50 px-1 mt-1">
                <span>Off</span>
                <span>3s</span>
                <span>6s</span>
                <span>9s</span>
                <span>12s</span>
              </div>
            </div>

            {/* Normalize Audio */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-4">
                <span className="label-text flex-1">
                  <span className="font-medium">Normalize Audio</span>
                  <span className="block text-sm text-base-content/60">
                    Balance volume levels across tracks
                  </span>
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={normalizeAudio}
                  onChange={(e) => setNormalizeAudio(e.target.checked)}
                />
              </label>
            </div>
          </div>
        </div>
      </section>

      {/* Reset */}
      <section className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">Reset</h2>

          <div className="divider mt-0" />

          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Reset to Defaults</p>
              <p className="text-sm text-base-content/60">
                Restore all settings to their default values
              </p>
            </div>
            <button
              className="btn btn-outline btn-error"
              onClick={() => {
                if (window.confirm('Are you sure you want to reset all settings to defaults?')) {
                  resetToDefaults();
                }
              }}
            >
              Reset All
            </button>
          </div>
        </div>
      </section>

      {/* Version info */}
      <div className="text-center text-sm text-base-content/50">
        <p>Personal Music Search Engine</p>
        <p>Settings are saved automatically to your browser</p>
      </div>
    </div>
  );
}

export default SettingsPage;
