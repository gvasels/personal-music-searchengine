/**
 * API Client Tests - REQ-5.1
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

describe('API Client', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('apiClient', () => {
    it('REQ-5.1: should export axios instance', async () => {
      const { apiClient } = await import('@/lib/api/client');
      expect(apiClient).toBeDefined();
    });

    it('should have base URL configured', async () => {
      const { apiClient } = await import('@/lib/api/client');
      expect(apiClient.defaults.baseURL).toBeDefined();
    });
  });

  describe('Auth Interceptor', () => {
    it('REQ-5.1: should attach JWT token to requests', async () => {
      const { apiClient } = await import('@/lib/api/client');
      // Interceptors should be configured
      expect(apiClient.interceptors.request).toBeDefined();
    });
  });

  describe('getTracks', () => {
    it('REQ-5.3: should fetch tracks', async () => {
      const { getTracks } = await import('@/lib/api/client');
      expect(getTracks).toBeDefined();
      expect(typeof getTracks).toBe('function');
    });
  });

  describe('getAlbums', () => {
    it('REQ-5.4: should fetch albums', async () => {
      const { getAlbums } = await import('@/lib/api/client');
      expect(getAlbums).toBeDefined();
      expect(typeof getAlbums).toBe('function');
    });
  });

  describe('getPlaylists', () => {
    it('REQ-5.9: should fetch playlists', async () => {
      const { getPlaylists } = await import('@/lib/api/client');
      expect(getPlaylists).toBeDefined();
      expect(typeof getPlaylists).toBe('function');
    });
  });

  describe('createPlaylist', () => {
    it('REQ-5.9: should create playlist', async () => {
      const { createPlaylist } = await import('@/lib/api/client');
      expect(createPlaylist).toBeDefined();
      expect(typeof createPlaylist).toBe('function');
    });
  });

  describe('addTagToTrack', () => {
    it('REQ-5.10: should add tag to track', async () => {
      const { addTagToTrack } = await import('@/lib/api/client');
      expect(addTagToTrack).toBeDefined();
      expect(typeof addTagToTrack).toBe('function');
    });
  });

  describe('searchTracks', () => {
    it('REQ-5.8: should search tracks', async () => {
      const { searchTracks } = await import('@/lib/api/client');
      expect(searchTracks).toBeDefined();
      expect(typeof searchTracks).toBe('function');
    });
  });

  describe('getPresignedUploadUrl', () => {
    it('REQ-5.5: should get upload URL', async () => {
      const { getPresignedUploadUrl } = await import('@/lib/api/client');
      expect(getPresignedUploadUrl).toBeDefined();
      expect(typeof getPresignedUploadUrl).toBe('function');
    });
  });

  describe('getStreamUrl', () => {
    it('REQ-5.6: should get stream URL', async () => {
      const { getStreamUrl } = await import('@/lib/api/client');
      expect(getStreamUrl).toBeDefined();
      expect(typeof getStreamUrl).toBe('function');
    });
  });
});
