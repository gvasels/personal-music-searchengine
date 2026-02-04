/**
 * UserSearchForm - Admin Panel Component
 * Search input with debounce for finding users by email or display name
 */
import { useState, useRef } from 'react';

interface UserSearchFormProps {
  onSearch: (query: string) => void;
  isLoading?: boolean;
  placeholder?: string;
}

export function UserSearchForm({
  onSearch,
  isLoading = false,
  placeholder = 'Search by email or name...',
}: UserSearchFormProps) {
  const [query, setQuery] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);
    onSearch(value);
  };

  const handleClear = () => {
    setQuery('');
    onSearch('');
    inputRef.current?.focus();
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Trigger search on explicit submit
    onSearch(query);
  };

  return (
    <form onSubmit={handleSubmit} className="w-full">
      <div className="relative">
        {/* Search icon */}
        <svg
          className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>

        <input
          ref={inputRef}
          type="text"
          value={query}
          onChange={handleChange}
          placeholder={placeholder}
          className="input input-bordered w-full pl-10 pr-10 bg-base-100 text-base-content"
          aria-label="Search users"
        />

        {/* Loading spinner or clear button */}
        {isLoading ? (
          <span className="absolute right-3 top-1/2 -translate-y-1/2 loading loading-spinner loading-sm" />
        ) : query.length > 0 ? (
          <button
            type="button"
            onClick={handleClear}
            className="absolute right-3 top-1/2 -translate-y-1/2 btn btn-ghost btn-xs btn-circle"
            aria-label="Clear search"
          >
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        ) : null}
      </div>

      {/* Helper text */}
      <p className="text-sm text-base-content/60 mt-2">
        Search by email address or display name (minimum 1 character)
      </p>
    </form>
  );
}
