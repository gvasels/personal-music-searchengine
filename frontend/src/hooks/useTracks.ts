/**
 * useTracks Hook - Wave 2
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getTracks,
  getTrack,
  updateTrack,
  deleteTrack,
  updateTrackVisibility,
  type GetTracksParams,
  type UpdateTrackData,
} from '../lib/api/tracks';
import type { TrackVisibility, Track } from '../types';

export const trackKeys = {
  all: ['tracks'] as const,
  lists: () => [...trackKeys.all, 'list'] as const,
  list: (params?: GetTracksParams) => [...trackKeys.lists(), params] as const,
  details: () => [...trackKeys.all, 'detail'] as const,
  detail: (id: string) => [...trackKeys.details(), id] as const,
};

export function useTracksQuery(params?: GetTracksParams) {
  return useQuery({
    queryKey: trackKeys.list(params),
    queryFn: () => getTracks(params),
  });
}

export function useTrackQuery(id: string | undefined) {
  return useQuery({
    queryKey: trackKeys.detail(id!),
    queryFn: () => getTrack(id!),
    enabled: !!id,
  });
}

export function useUpdateTrack() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTrackData }) => updateTrack(id, data),
    onSuccess: (data, variables) => {
      queryClient.setQueryData(trackKeys.detail(variables.id), data);
      void queryClient.invalidateQueries({ queryKey: trackKeys.lists() });
    },
  });
}

export function useDeleteTrack() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteTrack(id),
    onSuccess: (_, id) => {
      queryClient.removeQueries({ queryKey: trackKeys.detail(id) });
      void queryClient.invalidateQueries({ queryKey: trackKeys.lists() });
    },
  });
}

export function useUpdateTrackVisibility() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, visibility }: { id: string; visibility: TrackVisibility }) =>
      updateTrackVisibility(id, visibility),
    onSuccess: (data, variables) => {
      // Update the track in the cache with new visibility
      const existingTrack = queryClient.getQueryData<Track>(trackKeys.detail(variables.id));
      if (existingTrack) {
        queryClient.setQueryData(trackKeys.detail(variables.id), {
          ...existingTrack,
          visibility: data.visibility,
        });
      }
      void queryClient.invalidateQueries({ queryKey: trackKeys.lists() });
    },
  });
}
