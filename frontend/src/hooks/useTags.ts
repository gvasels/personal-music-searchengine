/**
 * useTags Hook - Wave 5
 */
import { useQuery } from '@tanstack/react-query';
import { getTags, getTracksByTag, type GetTagsParams, type GetTracksByTagParams } from '../lib/api/tags';

export const tagKeys = {
  all: ['tags'] as const,
  lists: () => [...tagKeys.all, 'list'] as const,
  list: (params?: GetTagsParams) => [...tagKeys.lists(), params] as const,
  tracks: (tagName: string) => [...tagKeys.all, 'tracks', tagName] as const,
  trackList: (tagName: string, params?: GetTracksByTagParams) =>
    [...tagKeys.tracks(tagName), params] as const,
};

export function useTagsQuery(params?: GetTagsParams) {
  return useQuery({
    queryKey: tagKeys.list(params),
    queryFn: () => getTags(params),
  });
}

export function useTracksByTagQuery(tagName: string, params?: GetTracksByTagParams) {
  return useQuery({
    queryKey: tagKeys.trackList(tagName, params),
    queryFn: () => getTracksByTag(tagName, params),
    enabled: !!tagName,
  });
}
