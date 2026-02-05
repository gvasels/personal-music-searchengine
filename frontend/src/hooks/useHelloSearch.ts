import { useQuery } from '@tanstack/react-query';
import { searchHelloTracks, getFeaturedHelloTracks } from '../lib/api/helloSearch';

export const helloKeys = {
  all: ['hello'] as const,
  search: (query: string) => ['hello', 'search', query] as const,
  featured: () => ['hello', 'featured'] as const,
};

export function useHelloSearch(query: string) {
  return useQuery({
    queryKey: helloKeys.search(query),
    queryFn: () => searchHelloTracks(query),
    enabled: query.length > 0,
  });
}

export function useHelloFeatured() {
  return useQuery({
    queryKey: helloKeys.featured(),
    queryFn: () => getFeaturedHelloTracks(),
  });
}
