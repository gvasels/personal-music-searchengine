/**
 * Events API Tests - TDD Red Phase
 *
 * These tests define the contract for the events API client.
 * They MUST FAIL because events.ts does not exist yet.
 *
 * Contract critical:
 * - getArtistEvents sends GET to /v1/artists/{name}/events
 * - searchArtistEvents sends GET to /v1/events/search with q and limit params
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

// These imports will fail - events.ts does not exist yet
import { getArtistEvents, searchArtistEvents } from '../events';
import { apiClient } from '../client';
import type { ArtistEventsResponse, ArtistSearchResult } from '../../../types';

describe('Events API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockEvent = {
    id: 'evt-1',
    artistName: 'Daft Punk',
    title: 'Alive 2027 Tour',
    date: '2027-06-15T20:00:00Z',
    venue: 'Madison Square Garden',
    city: 'New York',
    region: 'NY',
    country: 'US',
    ticketUrl: 'https://tickets.example.com/daft-punk',
    status: 'scheduled' as const,
    source: 'ticketmaster',
  };

  const mockEventsResponse: ArtistEventsResponse = {
    artistName: 'Daft Punk',
    events: [mockEvent],
    totalCount: 1,
    source: 'ticketmaster',
  };

  const mockSearchResult: ArtistSearchResult = {
    name: 'Daft Punk',
    imageUrl: 'https://img.example.com/daft-punk.jpg',
    upcomingEvents: 3,
    source: 'ticketmaster',
  };

  describe('getArtistEvents', () => {
    it('should send GET to /v1/artists/{name}/events', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockEventsResponse });

      await getArtistEvents('Daft Punk');

      expect(apiClient.get).toHaveBeenCalledWith('/v1/artists/Daft Punk/events');
    });

    it('should return events data with artistName, events array, totalCount, and source', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockEventsResponse });

      const result = await getArtistEvents('Daft Punk');

      expect(result.artistName).toBe('Daft Punk');
      expect(result.events).toHaveLength(1);
      expect(result.totalCount).toBe(1);
      expect(result.source).toBe('ticketmaster');
    });

    it('should return event with all required fields', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockEventsResponse });

      const result = await getArtistEvents('Daft Punk');
      const event = result.events[0];

      expect(event.id).toBe('evt-1');
      expect(event.title).toBe('Alive 2027 Tour');
      expect(event.venue).toBe('Madison Square Garden');
      expect(event.city).toBe('New York');
      expect(event.country).toBe('US');
      expect(event.status).toBe('scheduled');
    });

    it('should return empty events for unknown artist', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { artistName: 'Unknown', events: [], totalCount: 0, source: 'ticketmaster' },
      });

      const result = await getArtistEvents('Unknown');

      expect(result.events).toEqual([]);
      expect(result.totalCount).toBe(0);
    });

    it('should throw error on API failure', async () => {
      vi.mocked(apiClient.get).mockRejectedValue(new Error('Service unavailable'));

      await expect(getArtistEvents('Daft Punk')).rejects.toThrow('Service unavailable');
    });
  });

  describe('searchArtistEvents', () => {
    it('should send GET to /v1/events/search with query params', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockSearchResult] },
      });

      await searchArtistEvents('Daft Punk');

      expect(apiClient.get).toHaveBeenCalledWith('/v1/events/search', {
        params: { q: 'Daft Punk', limit: 10 },
      });
    });

    it('should default limit to 10', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockSearchResult] },
      });

      await searchArtistEvents('Daft Punk');

      expect(apiClient.get).toHaveBeenCalledWith('/v1/events/search', {
        params: { q: 'Daft Punk', limit: 10 },
      });
    });

    it('should pass custom limit', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockSearchResult] },
      });

      await searchArtistEvents('Daft Punk', 5);

      expect(apiClient.get).toHaveBeenCalledWith('/v1/events/search', {
        params: { q: 'Daft Punk', limit: 5 },
      });
    });

    it('should return items array of ArtistSearchResult', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [mockSearchResult] },
      });

      const result = await searchArtistEvents('Daft Punk');

      expect(result.items).toHaveLength(1);
      expect(result.items[0].name).toBe('Daft Punk');
      expect(result.items[0].upcomingEvents).toBe(3);
      expect(result.items[0].source).toBe('ticketmaster');
    });

    it('should return empty items for no matches', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: { items: [] },
      });

      const result = await searchArtistEvents('nonexistent');

      expect(result.items).toEqual([]);
    });

    it('should throw error on API failure', async () => {
      vi.mocked(apiClient.get).mockRejectedValue(new Error('Bad request'));

      await expect(searchArtistEvents('Daft Punk')).rejects.toThrow('Bad request');
    });
  });
});
