import { describe, it, expect, vi, beforeEach } from 'vitest';
import { searchHello, getFeaturedTracks } from '../helloSearch';
import { apiClient } from '../client';

vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

const mockSearchResponse = {
  tracks: [
    { id: '1', title: 'Midnight Drift', artist: 'Luna Waves', album: 'Waveforms', genre: 'Electronic', year: 2024, duration: 240, durationStr: '4:00', coverArtUrl: '' },
  ],
  total: 1,
  query: 'luna',
};

describe('helloSearch API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('searchHello calls GET /v1/hello/search with query param', async () => {
    (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockSearchResponse });

    const result = await searchHello('luna');

    expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/search', { params: { q: 'luna' } });
    expect(result).toEqual(mockSearchResponse);
  });

  it('getFeaturedTracks calls GET /v1/hello/featured', async () => {
    const featuredResponse = { ...mockSearchResponse, query: '' };
    (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: featuredResponse });

    const result = await getFeaturedTracks();

    expect(apiClient.get).toHaveBeenCalledWith('/v1/hello/featured');
    expect(result).toEqual(featuredResponse);
  });

  it('searchHello response matches HelloSearchResponse type', async () => {
    (apiClient.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockSearchResponse });

    const result = await searchHello('test');

    expect(result).toHaveProperty('tracks');
    expect(result).toHaveProperty('total');
    expect(result).toHaveProperty('query');
    expect(Array.isArray(result.tracks)).toBe(true);
  });
});
