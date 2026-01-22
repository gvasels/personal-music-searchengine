/**
 * Home Page
 * Task 1.11 - Protected home page with library stats and recent tracks
 */

import { useEffect } from 'react';
import { useNavigate, Link } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '../hooks/useAuth';
import { getTracks, Track } from '../lib/api/client';

interface LibraryStats {
  totalTracks: number;
  totalAlbums: number;
  totalArtists: number;
  totalPlaylists: number;
  totalDuration: number;
}

async function fetchRecentTracks(): Promise<{ items: Track[] }> {
  return getTracks({ limit: 5 });
}

async function fetchLibraryStats(): Promise<LibraryStats> {
  // For now, derive stats from tracks data
  const { items: tracks } = await getTracks({ limit: 1000 });
  const albums = new Set(tracks.map((t) => t.album).filter(Boolean));
  const artists = new Set(tracks.map((t) => t.artist).filter(Boolean));

  return {
    totalTracks: tracks.length,
    totalAlbums: albums.size,
    totalArtists: artists.size,
    totalPlaylists: 0,
    totalDuration: tracks.reduce((sum, t) => sum + (t.duration || 0), 0),
  };
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

function formatNumber(num: number): string {
  return num.toLocaleString();
}

function HomePage() {
  const navigate = useNavigate();
  const { user, isLoading: isAuthLoading, isAuthenticated } = useAuth();

  useEffect(() => {
    document.title = 'Home - Music Search Engine';
  }, []);

  useEffect(() => {
    if (!isAuthLoading && !isAuthenticated) {
      navigate({ to: '/login', search: { redirect: '/' } });
    }
  }, [isAuthLoading, isAuthenticated, navigate]);

  const {
    data: tracksData,
    isLoading: isLoadingTracks,
    error: tracksError,
  } = useQuery({
    queryKey: ['tracks', 'recent'],
    queryFn: fetchRecentTracks,
    enabled: isAuthenticated,
  });

  const {
    data: stats,
    isLoading: isLoadingStats,
    error: statsError,
  } = useQuery({
    queryKey: ['library', 'stats'],
    queryFn: fetchLibraryStats,
    enabled: isAuthenticated,
  });

  if (isAuthLoading) {
    return (
      <main className="min-h-screen flex items-center justify-center bg-base-200">
        <span className="loading loading-spinner loading-lg" role="status" aria-label="Loading"></span>
      </main>
    );
  }

  if (!isAuthenticated || !user) {
    return null;
  }

  const displayName = user.name || user.email || 'there';
  const recentTracks = tracksData?.items || [];
  const isEmpty = !isLoadingStats && stats?.totalTracks === 0;

  return (
    <main id="main-content" className="min-h-screen bg-base-200 p-4 md:p-8">
      <div className="max-w-6xl mx-auto">
        {/* Welcome Header */}
        <header className="mb-8">
          <h1 className="text-3xl font-bold mb-2">Welcome back, {displayName}!</h1>
          <p className="text-base-content/70">Here's what's happening with your music library.</p>
        </header>

        {/* Library Stats */}
        <section aria-label="Library statistics" className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Library Stats</h2>
          {isLoadingStats ? (
            <div data-testid="stats-loading" className="flex gap-4">
              {[1, 2, 3].map((i) => (
                <div key={i} className="skeleton h-24 w-32"></div>
              ))}
            </div>
          ) : statsError ? (
            <div role="alert" className="alert alert-error">
              <span>Failed to load library stats</span>
            </div>
          ) : stats ? (
            <div className="stats stats-vertical md:stats-horizontal shadow bg-base-100">
              <div className="stat">
                <div className="stat-title">Tracks</div>
                <div className="stat-value">{formatNumber(stats.totalTracks)}</div>
              </div>
              <div className="stat">
                <div className="stat-title">Albums</div>
                <div className="stat-value">{formatNumber(stats.totalAlbums)}</div>
              </div>
              <div className="stat">
                <div className="stat-title">Artists</div>
                <div className="stat-value">{formatNumber(stats.totalArtists)}</div>
              </div>
            </div>
          ) : null}
        </section>

        {/* Quick Actions */}
        <section aria-labelledby="actions-heading" className="mb-8">
          <h2 id="actions-heading" className="text-xl font-semibold mb-4">Quick Actions</h2>
          <div className="flex flex-wrap gap-4">
            <button
              onClick={() => navigate({ to: '/upload' })}
              className="btn btn-primary"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM6.293 6.707a1 1 0 010-1.414l3-3a1 1 0 011.414 0l3 3a1 1 0 01-1.414 1.414L11 5.414V13a1 1 0 11-2 0V5.414L7.707 6.707a1 1 0 01-1.414 0z" clipRule="evenodd" />
              </svg>
              {isEmpty ? 'Upload Your First Track' : 'Upload Music'}
            </button>
            <Link to="/tracks" className="btn btn-outline">
              Browse Library
            </Link>
            <Link to="/playlists" className="btn btn-outline">
              Playlists
            </Link>
            <Link to="/search" className="btn btn-outline">
              Search
            </Link>
          </div>
        </section>

        {/* Recent Tracks */}
        <section aria-labelledby="recent-heading">
          <div className="flex items-center justify-between mb-4">
            <h2 id="recent-heading" className="text-xl font-semibold">Recent Tracks</h2>
            {recentTracks.length > 0 && (
              <Link to="/tracks" className="link link-primary text-sm">
                View all
              </Link>
            )}
          </div>

          {isLoadingTracks ? (
            <div data-testid="recent-tracks-loading" className="space-y-2">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="skeleton h-16 w-full"></div>
              ))}
            </div>
          ) : tracksError ? (
            <div role="alert" className="alert alert-error">
              <span>Failed to load recent tracks</span>
            </div>
          ) : recentTracks.length > 0 ? (
            <div className="bg-base-100 rounded-lg shadow overflow-hidden">
              <ul aria-label="Recent tracks" className="divide-y divide-base-200" role="list">
                {recentTracks.map((track) => (
                  <li
                    key={track.id}
                    className="p-4 hover:bg-base-200 transition-colors cursor-pointer"
                    onClick={() => navigate({ to: `/tracks/${track.id}` })}
                    role="listitem"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex-1 min-w-0">
                        <p className="font-medium truncate">{track.title}</p>
                        <p className="text-sm text-base-content/70 truncate">
                          {track.artist}
                          {track.album && ` • ${track.album}`}
                        </p>
                      </div>
                      <div className="ml-4 text-sm text-base-content/50">
                        {formatDuration(track.duration)}
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <div className="text-center py-12 bg-base-100 rounded-lg">
              <p className="text-base-content/70 mb-4">No tracks in your library yet.</p>
              <button
                onClick={() => navigate({ to: '/upload' })}
                className="btn btn-primary"
              >
                Upload Your First Track
              </button>
            </div>
          )}
        </section>
      </div>
    </main>
  );
}

export default HomePage;
