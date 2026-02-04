/**
 * Hot Cues API
 * Handles hot cue-related API calls
 */
import { apiClient } from './client';
import type { HotCue, TrackHotCuesResponse, HotCueColor } from '@/types';

// Get all hot cues for a track
export async function getTrackHotCues(trackId: string): Promise<TrackHotCuesResponse> {
  const response = await apiClient.get<TrackHotCuesResponse>(`/tracks/${trackId}/hotcues`);
  return response.data;
}

// Set a hot cue
export async function setHotCue(
  trackId: string,
  slot: number,
  data: {
    position: number;
    label?: string;
    color?: HotCueColor;
  }
): Promise<HotCue> {
  const response = await apiClient.put<HotCue>(`/tracks/${trackId}/hotcues/${slot}`, data);
  return response.data;
}

// Delete a hot cue
export async function deleteHotCue(trackId: string, slot: number): Promise<void> {
  await apiClient.delete(`/tracks/${trackId}/hotcues/${slot}`);
}

// Clear all hot cues for a track
export async function clearHotCues(trackId: string): Promise<void> {
  await apiClient.delete(`/tracks/${trackId}/hotcues`);
}
