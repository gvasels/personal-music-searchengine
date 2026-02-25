import { useQuery } from '@tanstack/react-query';
import { getArtistEvents, searchArtistEvents } from '../lib/api/events';

export const eventKeys = {
  all: ['events'] as const,
  artist: (artistName: string) => [...eventKeys.all, 'artist', artistName] as const,
  search: (query: string) => [...eventKeys.all, 'search', query] as const,
};

export function useArtistEvents(artistName: string | undefined) {
  return useQuery({
    queryKey: eventKeys.artist(artistName!),
    queryFn: () => getArtistEvents(artistName!),
    enabled: !!artistName,
  });
}

export function useSearchArtistEvents(query: string) {
  return useQuery({
    queryKey: eventKeys.search(query),
    queryFn: () => searchArtistEvents(query),
    enabled: query.length > 0,
  });
}
