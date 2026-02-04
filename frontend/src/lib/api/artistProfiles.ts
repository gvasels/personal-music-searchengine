/**
 * Artist Profiles API - Global User Type Feature
 */
import { apiClient } from './client';
import type { ArtistProfile } from '../../types';

export interface CreateArtistProfileData {
  displayName: string;
  bio?: string;
  location?: string;
  website?: string;
  socialLinks?: Record<string, string>;
}

export interface UpdateArtistProfileData {
  displayName?: string;
  bio?: string;
  avatarUrl?: string;
  headerImageUrl?: string;
  location?: string;
  website?: string;
  socialLinks?: Record<string, string>;
}

export interface ArtistProfilesResponse {
  items: ArtistProfile[];
  nextCursor?: string;
}

export async function createArtistProfile(data: CreateArtistProfileData): Promise<ArtistProfile> {
  const response = await apiClient.post<ArtistProfile>('/artists/entity', data);
  return response.data;
}

export async function getArtistProfile(profileId: string): Promise<ArtistProfile> {
  const response = await apiClient.get<ArtistProfile>(`/artists/entity/${profileId}`);
  return response.data;
}

export async function updateArtistProfile(
  profileId: string,
  data: UpdateArtistProfileData
): Promise<ArtistProfile> {
  const response = await apiClient.put<ArtistProfile>(`/artists/entity/${profileId}`, data);
  return response.data;
}

export async function deleteArtistProfile(profileId: string): Promise<void> {
  await apiClient.delete(`/artists/entity/${profileId}`);
}

export interface ListArtistProfilesParams {
  limit?: number;
  cursor?: string;
}

export async function listArtistProfiles(
  params?: ListArtistProfilesParams
): Promise<ArtistProfilesResponse> {
  const response = await apiClient.get<ArtistProfilesResponse>('/artists/entity', { params });
  return response.data;
}

export interface SearchArtistsParams {
  query: string;
  limit?: number;
}

export async function searchArtists(params: SearchArtistsParams): Promise<ArtistProfilesResponse> {
  const response = await apiClient.get<ArtistProfilesResponse>('/artists/entity/search', { params });
  return response.data;
}
