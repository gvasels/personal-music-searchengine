/**
 * usePlaylists Hook - Wave 5
 */
import { useQuery } from '@tanstack/react-query';
import { getPlaylists, getPlaylist, type GetPlaylistsParams } from '../lib/api/playlists';

export const playlistKeys = {
  all: ['playlists'] as const,
  lists: () => [...playlistKeys.all, 'list'] as const,
  list: (params?: GetPlaylistsParams) => [...playlistKeys.lists(), params] as const,
  details: () => [...playlistKeys.all, 'detail'] as const,
  detail: (id: string) => [...playlistKeys.details(), id] as const,
};

export function usePlaylistsQuery(params?: GetPlaylistsParams) {
  return useQuery({
    queryKey: playlistKeys.list(params),
    queryFn: () => getPlaylists(params),
  });
}

export function usePlaylistQuery(id: string | undefined) {
  return useQuery({
    queryKey: playlistKeys.detail(id!),
    queryFn: () => getPlaylist(id!),
    enabled: !!id,
  });
}
