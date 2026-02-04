/**
 * VisibilitySelector Component - Select playlist visibility
 */
import type { PlaylistVisibility } from '../../types';

interface VisibilitySelectorProps {
  value: PlaylistVisibility;
  onChange: (visibility: PlaylistVisibility) => void;
  disabled?: boolean;
  size?: 'sm' | 'md' | 'lg';
}

const VISIBILITY_OPTIONS: { value: PlaylistVisibility; label: string; description: string; icon: string }[] = [
  {
    value: 'private',
    label: 'Private',
    description: 'Only you can see this playlist',
    icon: '🔒',
  },
  {
    value: 'unlisted',
    label: 'Unlisted',
    description: 'Anyone with the link can see',
    icon: '🔗',
  },
  {
    value: 'public',
    label: 'Public',
    description: 'Visible to everyone',
    icon: '🌐',
  },
];

export function VisibilitySelector({
  value,
  onChange,
  disabled = false,
  size = 'md',
}: VisibilitySelectorProps) {
  const sizeClasses = {
    sm: 'select-sm',
    md: '',
    lg: 'select-lg',
  };

  return (
    <div className="form-control">
      <select
        className={`select select-bordered ${sizeClasses[size]}`}
        value={value}
        onChange={(e) => onChange(e.target.value as PlaylistVisibility)}
        disabled={disabled}
      >
        {VISIBILITY_OPTIONS.map((option) => (
          <option key={option.value} value={option.value}>
            {option.icon} {option.label}
          </option>
        ))}
      </select>
    </div>
  );
}

interface VisibilityBadgeProps {
  visibility: PlaylistVisibility;
  size?: 'sm' | 'md' | 'lg';
}

export function VisibilityBadge({ visibility, size = 'md' }: VisibilityBadgeProps) {
  const option = VISIBILITY_OPTIONS.find((o) => o.value === visibility);
  if (!option) return null;

  const sizeClasses = {
    sm: 'badge-sm',
    md: '',
    lg: 'badge-lg',
  };

  const colorClasses: Record<PlaylistVisibility, string> = {
    private: 'badge-ghost',
    unlisted: 'badge-warning',
    public: 'badge-success',
  };

  return (
    <span className={`badge ${sizeClasses[size]} ${colorClasses[visibility]}`}>
      {option.icon} {option.label}
    </span>
  );
}

interface VisibilityRadioGroupProps {
  value: PlaylistVisibility;
  onChange: (visibility: PlaylistVisibility) => void;
  disabled?: boolean;
}

export function VisibilityRadioGroup({ value, onChange, disabled = false }: VisibilityRadioGroupProps) {
  return (
    <div className="space-y-2">
      {VISIBILITY_OPTIONS.map((option) => (
        <label
          key={option.value}
          className={`flex items-start gap-3 p-3 rounded-lg border cursor-pointer transition-colors ${
            value === option.value
              ? 'border-primary bg-primary/5'
              : 'border-base-200 hover:border-base-300'
          } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
        >
          <input
            type="radio"
            className="radio radio-primary mt-1"
            checked={value === option.value}
            onChange={() => onChange(option.value)}
            disabled={disabled}
          />
          <div>
            <div className="font-medium">
              {option.icon} {option.label}
            </div>
            <div className="text-sm text-base-content/60">{option.description}</div>
          </div>
        </label>
      ))}
    </div>
  );
}
