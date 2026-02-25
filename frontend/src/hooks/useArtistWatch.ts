import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { watchArtist, unwatchArtist, getWatchStatus, getWatchedArtists } from '../lib/api/artistWatch';

export const watchKeys = {
  all: ['artist-watch'] as const,
  status: (artistName: string) => [...watchKeys.all, 'status', artistName] as const,
  list: () => [...watchKeys.all, 'list'] as const,
};

export function useWatchStatus(artistName: string | undefined) {
  return useQuery({
    queryKey: watchKeys.status(artistName!),
    queryFn: () => getWatchStatus(artistName!),
    enabled: !!artistName,
  });
}

export function useWatchToggle(artistName: string) {
  const queryClient = useQueryClient();
  const { data: statusData, isLoading } = useWatchStatus(artistName);
  const isWatching = statusData?.watching ?? false;

  const watchMutation = useMutation({
    mutationFn: () => watchArtist(artistName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: watchKeys.status(artistName) });
      queryClient.invalidateQueries({ queryKey: watchKeys.list() });
    },
  });

  const unwatchMutation = useMutation({
    mutationFn: () => unwatchArtist(artistName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: watchKeys.status(artistName) });
      queryClient.invalidateQueries({ queryKey: watchKeys.list() });
    },
  });

  const toggle = async () => {
    if (isWatching) {
      await unwatchMutation.mutateAsync();
    } else {
      await watchMutation.mutateAsync();
    }
  };

  return {
    isWatching,
    isLoading,
    toggle,
  };
}

export function useWatchedArtists() {
  return useQuery({
    queryKey: watchKeys.list(),
    queryFn: () => getWatchedArtists(),
  });
}
