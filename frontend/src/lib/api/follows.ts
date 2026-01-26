/**
 * Follow API - Global User Type Feature
 */
import { apiClient } from './client';
import type { ArtistProfile } from '../../types';

export interface FollowersResponse {
  items: ArtistProfile[];
  nextCursor?: string;
  totalCount: number;
}

export interface FollowingResponse {
  items: ArtistProfile[];
  nextCursor?: string;
  totalCount: number;
}

export interface IsFollowingResponse {
  isFollowing: boolean;
}

export async function followArtist(artistId: string): Promise<void> {
  await apiClient.post(`/artists/entity/${artistId}/follow`);
}

export async function unfollowArtist(artistId: string): Promise<void> {
  await apiClient.delete(`/artists/entity/${artistId}/follow`);
}

export async function isFollowing(artistId: string): Promise<boolean> {
  const response = await apiClient.get<IsFollowingResponse>(
    `/artists/entity/${artistId}/following`
  );
  return response.data.isFollowing;
}

export interface GetFollowersParams {
  limit?: number;
  cursor?: string;
}

export async function getFollowers(
  artistId: string,
  params?: GetFollowersParams
): Promise<FollowersResponse> {
  const response = await apiClient.get<FollowersResponse>(
    `/artists/entity/${artistId}/followers`,
    { params }
  );
  return response.data;
}

export async function getFollowing(params?: GetFollowersParams): Promise<FollowingResponse> {
  const response = await apiClient.get<FollowingResponse>('/users/me/following', { params });
  return response.data;
}
