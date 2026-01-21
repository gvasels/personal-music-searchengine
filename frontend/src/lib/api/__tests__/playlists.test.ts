/**
 * Playlists API Tests - Wave 5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { apiClient } from '../client';
import {
  getPlaylists,
  getPlaylist,
  createPlaylist,
  updatePlaylist,
  deletePlaylist,
  addTrackToPlaylist,
  removeTrackFromPlaylist,
} from '../playlists';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('Playlists API (Wave 5)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getPlaylists', () => {
    it('should fetch paginated playlists', async () => {
      const mockResponse = {
        data: { items: [], total: 0, limit: 20, offset: 0 },
      };
      vi.mocked(apiClient.get).mockResolvedValue(mockResponse);

      const result = await getPlaylists({ limit: 20, offset: 0 });

      expect(apiClient.get).toHaveBeenCalledWith('/playlists', {
        params: { limit: 20, offset: 0 },
      });
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('getPlaylist', () => {
    it('should fetch a single playlist by ID', async () => {
      const mockPlaylist = {
        id: 'playlist-1',
        name: 'My Playlist',
        trackIds: [],
        trackCount: 0,
      };
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockPlaylist });

      const result = await getPlaylist('playlist-1');

      expect(apiClient.get).toHaveBeenCalledWith('/playlists/playlist-1');
      expect(result).toEqual(mockPlaylist);
    });
  });

  describe('createPlaylist', () => {
    it('should create a new playlist', async () => {
      const newPlaylist = { name: 'New Playlist', description: 'A test playlist' };
      const mockResponse = { id: 'playlist-1', ...newPlaylist, trackIds: [], trackCount: 0 };
      vi.mocked(apiClient.post).mockResolvedValue({ data: mockResponse });

      const result = await createPlaylist(newPlaylist);

      expect(apiClient.post).toHaveBeenCalledWith('/playlists', newPlaylist);
      expect(result).toEqual(mockResponse);
    });
  });

  describe('updatePlaylist', () => {
    it('should update an existing playlist', async () => {
      const updateData = { name: 'Updated Name' };
      const mockResponse = { id: 'playlist-1', name: 'Updated Name', trackIds: [], trackCount: 0 };
      vi.mocked(apiClient.put).mockResolvedValue({ data: mockResponse });

      const result = await updatePlaylist('playlist-1', updateData);

      expect(apiClient.put).toHaveBeenCalledWith('/playlists/playlist-1', updateData);
      expect(result).toEqual(mockResponse);
    });
  });

  describe('deletePlaylist', () => {
    it('should delete a playlist', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({ data: null });

      await deletePlaylist('playlist-1');

      expect(apiClient.delete).toHaveBeenCalledWith('/playlists/playlist-1');
    });
  });

  describe('addTrackToPlaylist', () => {
    it('should add a track to a playlist', async () => {
      const mockResponse = { id: 'playlist-1', trackIds: ['track-1'], trackCount: 1 };
      vi.mocked(apiClient.post).mockResolvedValue({ data: mockResponse });

      const result = await addTrackToPlaylist('playlist-1', 'track-1');

      expect(apiClient.post).toHaveBeenCalledWith('/playlists/playlist-1/tracks', {
        trackId: 'track-1',
      });
      expect(result).toEqual(mockResponse);
    });
  });

  describe('removeTrackFromPlaylist', () => {
    it('should remove a track from a playlist', async () => {
      const mockResponse = { id: 'playlist-1', trackIds: [], trackCount: 0 };
      vi.mocked(apiClient.delete).mockResolvedValue({ data: mockResponse });

      const result = await removeTrackFromPlaylist('playlist-1', 'track-1');

      expect(apiClient.delete).toHaveBeenCalledWith('/playlists/playlist-1/tracks/track-1');
      expect(result).toEqual(mockResponse);
    });
  });
});
