import { useQuery } from '@tanstack/react-query';
import { searchHello, getFeaturedTracks } from '../lib/api/helloSearch';

export const helloSearchKeys = {
  all: ['helloSearch'] as const,
  search: (query: string) => [...helloSearchKeys.all, 'search', query] as const,
  featured: () => [...helloSearchKeys.all, 'featured'] as const,
};

export function useHelloSearch(query: string) {
  return useQuery({
    queryKey: helloSearchKeys.search(query),
    queryFn: () => searchHello(query),
    enabled: query.length > 0,
  });
}

export function useFeaturedTracks() {
  return useQuery({
    queryKey: helloSearchKeys.featured(),
    queryFn: () => getFeaturedTracks(),
  });
}
