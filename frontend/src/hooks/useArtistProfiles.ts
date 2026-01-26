/**
 * useArtistProfiles Hook - Global User Type Feature
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  createArtistProfile,
  getArtistProfile,
  updateArtistProfile,
  deleteArtistProfile,
  listArtistProfiles,
  searchArtists,
  CreateArtistProfileData,
  UpdateArtistProfileData,
  ListArtistProfilesParams,
  SearchArtistsParams,
} from '../lib/api/artistProfiles';

// Query key factory for artist profiles
export const artistProfileKeys = {
  all: ['artistProfiles'] as const,
  lists: () => [...artistProfileKeys.all, 'list'] as const,
  list: (params?: ListArtistProfilesParams) => [...artistProfileKeys.lists(), params] as const,
  details: () => [...artistProfileKeys.all, 'detail'] as const,
  detail: (id: string) => [...artistProfileKeys.details(), id] as const,
  search: (params: SearchArtistsParams) => [...artistProfileKeys.all, 'search', params] as const,
};

/**
 * Hook to fetch a single artist profile
 */
export function useArtistProfile(profileId: string) {
  return useQuery({
    queryKey: artistProfileKeys.detail(profileId),
    queryFn: () => getArtistProfile(profileId),
    enabled: !!profileId,
  });
}

/**
 * Hook to fetch a list of artist profiles
 */
export function useArtistProfiles(params?: ListArtistProfilesParams) {
  return useQuery({
    queryKey: artistProfileKeys.list(params),
    queryFn: () => listArtistProfiles(params),
  });
}

/**
 * Hook to search artists
 */
export function useSearchArtists(params: SearchArtistsParams) {
  return useQuery({
    queryKey: artistProfileKeys.search(params),
    queryFn: () => searchArtists(params),
    enabled: params.query.length > 0,
  });
}

/**
 * Hook to create an artist profile
 */
export function useCreateArtistProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateArtistProfileData) => createArtistProfile(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: artistProfileKeys.lists() });
    },
  });
}

/**
 * Hook to update an artist profile
 */
export function useUpdateArtistProfile(profileId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateArtistProfileData) => updateArtistProfile(profileId, data),
    onSuccess: (updatedProfile) => {
      queryClient.setQueryData(artistProfileKeys.detail(profileId), updatedProfile);
      queryClient.invalidateQueries({ queryKey: artistProfileKeys.lists() });
    },
  });
}

/**
 * Hook to delete an artist profile
 */
export function useDeleteArtistProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (profileId: string) => deleteArtistProfile(profileId),
    onSuccess: (_, profileId) => {
      queryClient.removeQueries({ queryKey: artistProfileKeys.detail(profileId) });
      queryClient.invalidateQueries({ queryKey: artistProfileKeys.lists() });
    },
  });
}
