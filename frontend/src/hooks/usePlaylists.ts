/**
 * usePlaylists Hook - Wave 5
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getPlaylists,
  getPlaylist,
  createPlaylist,
  updatePlaylist,
  deletePlaylist,
  addTrackToPlaylist,
  removeTrackFromPlaylist,
  type GetPlaylistsParams,
  type CreatePlaylistData,
  type UpdatePlaylistData,
} from '../lib/api/playlists';

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

export function useCreatePlaylist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreatePlaylistData) => createPlaylist(data),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: playlistKeys.lists() });
    },
  });
}

export function useUpdatePlaylist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePlaylistData }) => updatePlaylist(id, data),
    onSuccess: (_, { id }) => {
      void queryClient.invalidateQueries({ queryKey: playlistKeys.detail(id) });
      void queryClient.invalidateQueries({ queryKey: playlistKeys.lists() });
    },
  });
}

export function useDeletePlaylist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deletePlaylist(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: playlistKeys.lists() });
    },
  });
}

export function useAddTracksToPlaylist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, trackId }: { id: string; trackId: string }) => addTrackToPlaylist(id, trackId),
    onSuccess: (_, { id }) => {
      void queryClient.invalidateQueries({ queryKey: playlistKeys.detail(id) });
      void queryClient.invalidateQueries({ queryKey: playlistKeys.lists() });
    },
  });
}

export function useRemoveTracksFromPlaylist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, trackIds }: { id: string; trackIds: string[] }) =>
      Promise.all(trackIds.map((trackId) => removeTrackFromPlaylist(id, trackId))),
    onSuccess: (_, { id }) => {
      void queryClient.invalidateQueries({ queryKey: playlistKeys.detail(id) });
      void queryClient.invalidateQueries({ queryKey: playlistKeys.lists() });
    },
  });
}
