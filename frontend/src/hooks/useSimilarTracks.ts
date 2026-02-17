import { useQuery } from '@tanstack/react-query';
import { getSimilarTracks } from '../lib/api/tracks';

export const similarTracksKeys = {
  all: ['similarTracks'] as const,
  detail: (id: string) => [...similarTracksKeys.all, id] as const,
};

export function useSimilarTracks(trackId: string | undefined) {
  return useQuery({
    queryKey: similarTracksKeys.detail(trackId!),
    queryFn: () => getSimilarTracks(trackId!),
    enabled: !!trackId,
  });
}
