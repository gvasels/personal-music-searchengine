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
  reorderPlaylistTracks,
  type GetPlaylistsParams,
  type CreatePlaylistData,
  type UpdatePlaylistData,
} from '../lib/api/playlists';
import type { PlaylistWithTracks, Track } from '@/types';

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

export function useReorderPlaylistTracks() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, trackIds }: { id: string; trackIds: string[] }) =>
      reorderPlaylistTracks(id, { trackIds }),
    onMutate: async ({ id, trackIds }) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: playlistKeys.detail(id) });

      // Snapshot the previous value
      const previousData = queryClient.getQueryData<PlaylistWithTracks>(playlistKeys.detail(id));

      // Optimistically update to the new order
      if (previousData) {
        const trackMap = new Map(previousData.tracks.map((t) => [t.id, t]));
        const reorderedTracks: Track[] = [];

        for (const trackId of trackIds) {
          const track = trackMap.get(trackId);
          if (track) {
            reorderedTracks.push(track);
          }
        }

        queryClient.setQueryData<PlaylistWithTracks>(playlistKeys.detail(id), {
          ...previousData,
          playlist: {
            ...previousData.playlist,
            trackIds,
          },
          tracks: reorderedTracks,
        });
      }

      return { previousData };
    },
    onError: (_err, { id }, context) => {
      // Rollback on error
      if (context?.previousData) {
        queryClient.setQueryData(playlistKeys.detail(id), context.previousData);
      }
    },
    onSettled: (_, __, { id }) => {
      // Always refetch after error or success
      void queryClient.invalidateQueries({ queryKey: playlistKeys.detail(id) });
    },
  });
}
