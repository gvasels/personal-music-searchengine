/**
 * useArtists Hook - Wave 2
 */
import { useQuery } from '@tanstack/react-query';
import { getArtists, getArtist, type GetArtistsParams } from '../lib/api/artists';

export const artistKeys = {
  all: ['artists'] as const,
  lists: () => [...artistKeys.all, 'list'] as const,
  list: (params?: GetArtistsParams) => [...artistKeys.lists(), params] as const,
  details: () => [...artistKeys.all, 'detail'] as const,
  detail: (name: string) => [...artistKeys.details(), name] as const,
};

export function useArtistsQuery(params?: GetArtistsParams) {
  return useQuery({
    queryKey: artistKeys.list(params),
    queryFn: () => getArtists(params),
  });
}

export function useArtistQuery(name: string | undefined) {
  return useQuery({
    queryKey: artistKeys.detail(name!),
    queryFn: () => getArtist(name!),
    enabled: !!name,
  });
}
