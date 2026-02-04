/**
 * useFollow Hook - Global User Type Feature
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  followArtist,
  unfollowArtist,
  isFollowing,
  getFollowers,
  getFollowing,
  GetFollowersParams,
} from '../lib/api/follows';
import { artistProfileKeys } from './useArtistProfiles';

// Query key factory for follow relationships
export const followKeys = {
  all: ['follows'] as const,
  isFollowing: (artistId: string) => [...followKeys.all, 'isFollowing', artistId] as const,
  followers: (artistId: string, params?: GetFollowersParams) =>
    [...followKeys.all, 'followers', artistId, params] as const,
  following: (params?: GetFollowersParams) => [...followKeys.all, 'following', params] as const,
};

/**
 * Hook to check if current user is following an artist
 */
export function useIsFollowing(artistId: string) {
  return useQuery({
    queryKey: followKeys.isFollowing(artistId),
    queryFn: () => isFollowing(artistId),
    enabled: !!artistId,
  });
}

/**
 * Hook to get followers of an artist
 */
export function useFollowers(artistId: string, params?: GetFollowersParams) {
  return useQuery({
    queryKey: followKeys.followers(artistId, params),
    queryFn: () => getFollowers(artistId, params),
    enabled: !!artistId,
  });
}

/**
 * Hook to get artists the current user is following
 */
export function useFollowing(params?: GetFollowersParams) {
  return useQuery({
    queryKey: followKeys.following(params),
    queryFn: () => getFollowing(params),
  });
}

/**
 * Hook to follow an artist
 */
export function useFollowArtist() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (artistId: string) => followArtist(artistId),
    onSuccess: (_, artistId) => {
      // Update the isFollowing cache
      queryClient.setQueryData(followKeys.isFollowing(artistId), true);
      // Invalidate follower counts and following lists
      queryClient.invalidateQueries({ queryKey: followKeys.followers(artistId) });
      queryClient.invalidateQueries({ queryKey: followKeys.following() });
      // Invalidate artist profile to update follower count
      queryClient.invalidateQueries({ queryKey: artistProfileKeys.detail(artistId) });
    },
  });
}

/**
 * Hook to unfollow an artist
 */
export function useUnfollowArtist() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (artistId: string) => unfollowArtist(artistId),
    onSuccess: (_, artistId) => {
      // Update the isFollowing cache
      queryClient.setQueryData(followKeys.isFollowing(artistId), false);
      // Invalidate follower counts and following lists
      queryClient.invalidateQueries({ queryKey: followKeys.followers(artistId) });
      queryClient.invalidateQueries({ queryKey: followKeys.following() });
      // Invalidate artist profile to update follower count
      queryClient.invalidateQueries({ queryKey: artistProfileKeys.detail(artistId) });
    },
  });
}

/**
 * Combined hook for follow/unfollow functionality
 */
export function useFollowToggle(artistId: string) {
  const { data: isCurrentlyFollowing, isLoading } = useIsFollowing(artistId);
  const followMutation = useFollowArtist();
  const unfollowMutation = useUnfollowArtist();

  const toggle = () => {
    if (isCurrentlyFollowing) {
      unfollowMutation.mutate(artistId);
    } else {
      followMutation.mutate(artistId);
    }
  };

  return {
    isFollowing: isCurrentlyFollowing ?? false,
    isLoading,
    isToggling: followMutation.isPending || unfollowMutation.isPending,
    toggle,
    follow: () => followMutation.mutate(artistId),
    unfollow: () => unfollowMutation.mutate(artistId),
  };
}
