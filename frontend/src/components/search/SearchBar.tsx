import { useState, useRef, useEffect, useCallback } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { useAutocompleteQuery } from '@/hooks/useSearch';
import type { AutocompleteSuggestion } from '@/lib/api/search';

function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

export function SearchBar() {
  const navigate = useNavigate();
  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const inputRef = useRef<HTMLInputElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const debouncedQuery = useDebounce(query, 300);
  const { data, isLoading } = useAutocompleteQuery(debouncedQuery);

  const suggestions = data?.suggestions || [];

  // Group suggestions by type
  const groupedSuggestions = suggestions.reduce(
    (acc, suggestion) => {
      if (!acc[suggestion.type]) {
        acc[suggestion.type] = [];
      }
      acc[suggestion.type].push(suggestion);
      return acc;
    },
    {} as Record<string, AutocompleteSuggestion[]>
  );

  // Flatten for keyboard navigation
  const flatSuggestions = Object.values(groupedSuggestions).flat();

  const handleSelect = useCallback(
    (suggestion: AutocompleteSuggestion) => {
      setQuery('');
      setIsOpen(false);
      setSelectedIndex(-1);

      switch (suggestion.type) {
        case 'track':
          if (suggestion.trackId) {
            navigate({ to: '/tracks/$trackId', params: { trackId: suggestion.trackId } });
          }
          break;
        case 'album':
          // Search for the album
          navigate({ to: '/search', search: { q: suggestion.value } });
          break;
        case 'artist':
          // Search for the artist
          navigate({ to: '/search', search: { q: suggestion.value } });
          break;
      }
    },
    [navigate]
  );

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (!isOpen) {
        if (e.key === 'ArrowDown' && query.length >= 2) {
          setIsOpen(true);
        }
        return;
      }

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setSelectedIndex((prev) =>
            prev < flatSuggestions.length - 1 ? prev + 1 : prev
          );
          break;
        case 'ArrowUp':
          e.preventDefault();
          setSelectedIndex((prev) => (prev > 0 ? prev - 1 : -1));
          break;
        case 'Enter':
          e.preventDefault();
          if (selectedIndex >= 0 && flatSuggestions[selectedIndex]) {
            handleSelect(flatSuggestions[selectedIndex]);
          } else if (query.trim()) {
            // Navigate to search results page
            navigate({ to: '/search', search: { q: query.trim() } });
            setQuery('');
            setIsOpen(false);
          }
          break;
        case 'Escape':
          setIsOpen(false);
          setSelectedIndex(-1);
          inputRef.current?.blur();
          break;
      }
    },
    [isOpen, flatSuggestions, selectedIndex, query, handleSelect, navigate]
  );

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        !inputRef.current?.contains(event.target as Node)
      ) {
        setIsOpen(false);
        setSelectedIndex(-1);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Open dropdown when query changes and has results
  useEffect(() => {
    if (debouncedQuery.length >= 2 && suggestions.length > 0) {
      setIsOpen(true);
    } else if (debouncedQuery.length < 2) {
      setIsOpen(false);
    }
    setSelectedIndex(-1);
  }, [debouncedQuery, suggestions.length]);

  const renderSuggestionGroup = (type: string, items: AutocompleteSuggestion[]) => {
    const typeLabels: Record<string, string> = {
      track: 'Tracks',
      artist: 'Artists',
      album: 'Albums',
    };

    const startIndex = flatSuggestions.findIndex(
      (s) => s.type === type && s.value === items[0]?.value
    );

    return (
      <div key={type}>
        <div className="px-3 py-1.5 text-xs font-semibold text-base-content/60 uppercase tracking-wide">
          {typeLabels[type] || type}
        </div>
        {items.map((suggestion, idx) => {
          const globalIndex = startIndex + idx;
          const isSelected = globalIndex === selectedIndex;

          return (
            <button
              key={`${suggestion.type}-${suggestion.value}-${idx}`}
              type="button"
              className={`w-full px-3 py-2 text-left flex items-center gap-2 hover:bg-base-200 transition-colors text-base-content ${
                isSelected ? 'bg-base-200' : ''
              }`}
              onClick={() => handleSelect(suggestion)}
              onMouseEnter={() => setSelectedIndex(globalIndex)}
            >
              <span className="flex-1 truncate text-base-content">{suggestion.value}</span>
              <span className="text-xs text-base-content/50 capitalize">{suggestion.type}</span>
            </button>
          );
        })}
      </div>
    );
  };

  return (
    <div className="relative w-full max-w-md">
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          placeholder="Search tracks, artists, albums..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => {
            if (debouncedQuery.length >= 2 && suggestions.length > 0) {
              setIsOpen(true);
            }
          }}
          className="input input-bordered w-full pl-10 pr-4 bg-base-100 text-base-content"
          aria-label="Search"
          aria-expanded={isOpen}
          aria-haspopup="listbox"
          aria-autocomplete="list"
          role="combobox"
        />
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
        {isLoading && (
          <span className="absolute right-3 top-1/2 -translate-y-1/2 loading loading-spinner loading-sm" />
        )}
      </div>

      {isOpen && (
        <div
          ref={dropdownRef}
          className="absolute z-50 w-full mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-80 overflow-y-auto"
          role="listbox"
        >
          {suggestions.length === 0 && debouncedQuery.length >= 2 && !isLoading ? (
            <div className="px-3 py-4 text-center text-base-content/60">
              No results found for "{debouncedQuery}"
            </div>
          ) : (
            <>
              {['track', 'artist', 'album'].map(
                (type) =>
                  groupedSuggestions[type]?.length > 0 &&
                  renderSuggestionGroup(type, groupedSuggestions[type])
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
}
