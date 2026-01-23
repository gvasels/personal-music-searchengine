/**
 * useWaveform Hook
 * Fetches and caches waveform data for tracks
 */
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';

export interface WaveformData {
  version: number;
  sampleRate: number;
  duration: number;
  peaks: number[];
}

async function fetchWaveform(trackId: string): Promise<WaveformData | null> {
  try {
    const response = await apiClient.get<WaveformData>(`/tracks/${trackId}/waveform`);
    return response.data;
  } catch {
    // Waveform may not exist yet
    return null;
  }
}

export function useWaveform(trackId: string | undefined) {
  return useQuery({
    queryKey: ['waveform', trackId],
    queryFn: () => fetchWaveform(trackId!),
    enabled: !!trackId,
    staleTime: Infinity, // Waveform data doesn't change
    gcTime: 1000 * 60 * 30, // Keep in cache for 30 minutes
  });
}

/**
 * Generate mock waveform data for testing/demo
 */
export function generateMockWaveform(duration: number): WaveformData {
  const sampleRate = 100; // 100 samples per second
  const numSamples = Math.floor(duration * sampleRate);
  const peaks: number[] = [];

  for (let i = 0; i < numSamples; i++) {
    // Generate realistic-looking waveform with some patterns
    const t = i / numSamples;
    const base = 0.3 + Math.random() * 0.4;
    const beat = Math.sin(t * Math.PI * 50) * 0.2;
    const variation = Math.sin(t * Math.PI * 5) * 0.1;
    peaks.push(Math.min(1, Math.max(0.1, base + beat + variation)));
  }

  return {
    version: 1,
    sampleRate,
    duration,
    peaks,
  };
}
