/**
 * useSearch Hook - Wave 4
 */
import { useQuery } from '@tanstack/react-query';
import { searchTracks, searchAutocomplete, type SearchParams } from '../lib/api/search';

export const searchKeys = {
  all: ['search'] as const,
  results: (params: SearchParams) => [...searchKeys.all, 'results', params] as const,
  autocomplete: (query: string) => [...searchKeys.all, 'autocomplete', query] as const,
};

export function useSearchQuery(params: SearchParams) {
  return useQuery({
    queryKey: searchKeys.results(params),
    queryFn: () => searchTracks(params),
    enabled: !!params.query && params.query.length > 0,
  });
}

export function useAutocompleteQuery(query: string) {
  return useQuery({
    queryKey: searchKeys.autocomplete(query),
    queryFn: () => searchAutocomplete(query),
    enabled: query.length >= 3,
  });
}
