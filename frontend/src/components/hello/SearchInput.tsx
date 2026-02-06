/**
 * SearchInput Component
 *
 * A controlled text input for search functionality.
 * Uses DaisyUI input classes for consistent styling.
 */

interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

export function SearchInput({
  value,
  onChange,
  placeholder = 'Search tracks...',
}: SearchInputProps) {
  return (
    <input
      type="text"
      className="input input-bordered w-full"
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder={placeholder}
    />
  );
}
