/**
 * useAlbums Hook - Wave 2
 */
import { useQuery } from '@tanstack/react-query';
import { getAlbums, getAlbum, type GetAlbumsParams } from '../lib/api/albums';

export const albumKeys = {
  all: ['albums'] as const,
  lists: () => [...albumKeys.all, 'list'] as const,
  list: (params?: GetAlbumsParams) => [...albumKeys.lists(), params] as const,
  details: () => [...albumKeys.all, 'detail'] as const,
  detail: (id: string) => [...albumKeys.details(), id] as const,
};

export function useAlbumsQuery(params?: GetAlbumsParams) {
  return useQuery({
    queryKey: albumKeys.list(params),
    queryFn: () => getAlbums(params),
  });
}

export function useAlbumQuery(id: string | undefined) {
  return useQuery({
    queryKey: albumKeys.detail(id!),
    queryFn: () => getAlbum(id!),
    enabled: !!id,
  });
}
