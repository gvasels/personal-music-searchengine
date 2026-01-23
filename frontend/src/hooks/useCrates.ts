/**
 * useCrates Hook
 * Manages DJ crate operations with React Query
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getCrates,
  getCrate,
  createCrate,
  updateCrate,
  deleteCrate,
  addTracksToCrate,
  removeTracksFromCrate,
  reorderCrateTracks,
} from '@/lib/api/crates';
import { useFeatureGate } from './useFeatureFlags';

// Query key factory
export const crateKeys = {
  all: ['crates'] as const,
  lists: () => [...crateKeys.all, 'list'] as const,
  list: (filters?: Record<string, unknown>) => [...crateKeys.lists(), filters] as const,
  details: () => [...crateKeys.all, 'detail'] as const,
  detail: (id: string) => [...crateKeys.details(), id] as const,
};

export function useCrates() {
  const { isEnabled, showUpgrade } = useFeatureGate('CRATES');

  return {
    ...useQuery({
      queryKey: crateKeys.lists(),
      queryFn: getCrates,
      enabled: isEnabled,
      staleTime: 1000 * 60 * 5,
    }),
    isFeatureEnabled: isEnabled,
    showUpgrade,
  };
}

export function useCrate(id: string) {
  const { isEnabled } = useFeatureGate('CRATES');

  return useQuery({
    queryKey: crateKeys.detail(id),
    queryFn: () => getCrate(id),
    enabled: isEnabled && !!id,
    staleTime: 1000 * 60 * 5,
  });
}

export function useCreateCrate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createCrate,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: crateKeys.lists() });
    },
  });
}

export function useUpdateCrate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Parameters<typeof updateCrate>[1] }) =>
      updateCrate(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: crateKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: crateKeys.lists() });
    },
  });
}

export function useDeleteCrate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteCrate,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: crateKeys.lists() });
    },
  });
}

export function useAddTracksToCrate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      crateId,
      trackIds,
      position,
    }: {
      crateId: string;
      trackIds: string[];
      position?: number;
    }) => addTracksToCrate(crateId, trackIds, position),
    onSuccess: (_, { crateId }) => {
      queryClient.invalidateQueries({ queryKey: crateKeys.detail(crateId) });
      queryClient.invalidateQueries({ queryKey: crateKeys.lists() });
    },
  });
}

export function useRemoveTracksFromCrate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ crateId, trackIds }: { crateId: string; trackIds: string[] }) =>
      removeTracksFromCrate(crateId, trackIds),
    onSuccess: (_, { crateId }) => {
      queryClient.invalidateQueries({ queryKey: crateKeys.detail(crateId) });
      queryClient.invalidateQueries({ queryKey: crateKeys.lists() });
    },
  });
}

export function useReorderCrateTracks() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ crateId, trackIds }: { crateId: string; trackIds: string[] }) =>
      reorderCrateTracks(crateId, trackIds),
    onSuccess: (_, { crateId }) => {
      queryClient.invalidateQueries({ queryKey: crateKeys.detail(crateId) });
    },
  });
}
