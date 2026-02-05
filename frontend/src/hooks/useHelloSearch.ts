import { useQuery } from '@tanstack/react-query';
import { searchHelloTracks, getFeaturedTracks } from '@/lib/api/helloSearch';

export const helloKeys = {
  all: ['hello'] as const,
  search: (query: string) => ['hello', 'search', query] as const,
  featured: (limit?: number) =>
    limit !== undefined
      ? (['hello', 'featured', limit] as const)
      : (['hello', 'featured'] as const),
};

export function useHelloSearch(query: string) {
  return useQuery({
    queryKey: helloKeys.search(query),
    queryFn: () => searchHelloTracks(query),
    enabled: query.length > 0,
  });
}

export function useFeaturedTracks(limit?: number) {
  return useQuery({
    queryKey: helloKeys.featured(limit),
    queryFn: () => getFeaturedTracks(limit),
  });
}
