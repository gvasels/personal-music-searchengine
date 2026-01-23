/**
 * useHotCues Hook
 * Manages hot cue operations with React Query
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getTrackHotCues,
  setHotCue,
  deleteHotCue,
  clearHotCues,
} from '@/lib/api/hotcues';
import { useFeatureGate } from './useFeatureFlags';
import type { HotCueColor } from '@/types';

// Query key factory
export const hotCueKeys = {
  all: ['hotcues'] as const,
  track: (trackId: string) => [...hotCueKeys.all, trackId] as const,
};

export function useHotCues(trackId: string | undefined) {
  const { isEnabled, showUpgrade } = useFeatureGate('HOT_CUES');

  const query = useQuery({
    queryKey: hotCueKeys.track(trackId || ''),
    queryFn: () => getTrackHotCues(trackId!),
    enabled: isEnabled && !!trackId,
    staleTime: 1000 * 60 * 5,
  });

  return {
    ...query,
    isFeatureEnabled: isEnabled,
    showUpgrade,
  };
}

export function useSetHotCue() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      trackId,
      slot,
      position,
      label,
      color,
    }: {
      trackId: string;
      slot: number;
      position: number;
      label?: string;
      color?: HotCueColor;
    }) => setHotCue(trackId, slot, { position, label, color }),
    onSuccess: (_, { trackId }) => {
      queryClient.invalidateQueries({ queryKey: hotCueKeys.track(trackId) });
    },
  });
}

export function useDeleteHotCue() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ trackId, slot }: { trackId: string; slot: number }) =>
      deleteHotCue(trackId, slot),
    onSuccess: (_, { trackId }) => {
      queryClient.invalidateQueries({ queryKey: hotCueKeys.track(trackId) });
    },
  });
}

export function useClearHotCues() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: clearHotCues,
    onSuccess: (_, trackId) => {
      queryClient.invalidateQueries({ queryKey: hotCueKeys.track(trackId) });
    },
  });
}
