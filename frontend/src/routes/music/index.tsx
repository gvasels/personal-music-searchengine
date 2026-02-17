import { useEffect } from 'react';
import { Link } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '../../hooks/useAuth';
import { getTracks } from '../../lib/api/client';
import { getLibraryStats } from '../../lib/api/stats';

export default function MusicIndexPage() {
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    document.title = 'Music Library - Music Search Engine';
  }, []);

  const { data: stats } = useQuery({
    queryKey: ['library', 'stats', 'own'],
    queryFn: () => getLibraryStats('own'),
    enabled: isAuthenticated,
  });

  const { data: tracksData } = useQuery({
    queryKey: ['tracks', 'recent'],
    queryFn: () => getTracks({ limit: 5 }),
    enabled: isAuthenticated,
  });

  const recentTracks = tracksData?.items || [];

  return (
    <main className="min-h-screen bg-base-200 p-4 md:p-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Music Library</h1>

        {stats && (
          <div className="stats stats-vertical md:stats-horizontal shadow bg-base-100 mb-8">
            <div className="stat">
              <div className="stat-title">Tracks</div>
              <div className="stat-value">{stats.totalTracks}</div>
            </div>
            <div className="stat">
              <div className="stat-title">Albums</div>
              <div className="stat-value">{stats.totalAlbums}</div>
            </div>
            <div className="stat">
              <div className="stat-title">Artists</div>
              <div className="stat-value">{stats.totalArtists}</div>
            </div>
          </div>
        )}

        <section>
          <h2 className="text-xl font-semibold mb-4">Recent Tracks</h2>
          {recentTracks.length > 0 ? (
            <ul className="space-y-2">
              {recentTracks.map((track) => (
                <li key={track.id} className="p-3 bg-base-100 rounded-lg">
                  <Link to="/music/tracks/$trackId" params={{ trackId: track.id }} className="font-medium">
                    {track.title}
                  </Link>
                  <span className="text-sm text-base-content/70 ml-2">{track.artist}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-base-content/70">No tracks yet.</p>
          )}
        </section>
      </div>
    </main>
  );
}
