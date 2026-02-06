/**
 * useHelloSearch Hooks
 *
 * TanStack Query hooks for the Hello World search feature.
 * Provides reactive data fetching with caching and loading states.
 */
import { useQuery } from '@tanstack/react-query';
import {
  searchHelloTracks,
  getFeaturedHelloTracks,
  type HelloSearchResponse,
} from '../lib/api/helloSearch';

/**
 * Query key factory for hello-related queries.
 * Enables efficient cache invalidation and refetching.
 */
export const helloKeys = {
  all: ['hello'] as const,
  search: (q: string) => [...helloKeys.all, 'search', q] as const,
  featured: () => [...helloKeys.all, 'featured'] as const,
};

/**
 * Hook to search for hello tracks.
 * Query is disabled when the search string is empty.
 *
 * @param query - Search term
 * @returns TanStack Query result with items and total
 */
export function useHelloSearch(query: string) {
  return useQuery<HelloSearchResponse>({
    queryKey: helloKeys.search(query),
    queryFn: () => searchHelloTracks(query),
    enabled: query.length > 0,
  });
}

/**
 * Hook to fetch featured hello tracks.
 *
 * @returns TanStack Query result with items and total
 */
export function useHelloFeatured() {
  return useQuery<HelloSearchResponse>({
    queryKey: helloKeys.featured(),
    queryFn: () => getFeaturedHelloTracks(),
  });
}
