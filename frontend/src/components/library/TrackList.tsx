import { Track } from '@/types';
import { useNavigate } from '@tanstack/react-router';
import { usePlayerStore } from '@/lib/store/playerStore';
import { getDownloadUrl } from '@/lib/api/client';
import { AddToPlaylistDropdown } from './AddToPlaylistDropdown';

interface TrackListProps {
  tracks: Track[];
  isLoading?: boolean;
  showDownload?: boolean;
  showAddedDate?: boolean;
  showAddToPlaylist?: boolean;
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
}

export function TrackList({ tracks, isLoading, showDownload = false, showAddedDate = false, showAddToPlaylist = false }: TrackListProps) {
  const navigate = useNavigate();
  const { setQueue, currentTrack, isPlaying } = usePlayerStore();

  const handleArtistClick = (e: React.MouseEvent, artist: string) => {
    e.stopPropagation();
    navigate({ to: '/artists/$artistName', params: { artistName: artist } });
  };

  const handleAlbumClick = (e: React.MouseEvent, album: string) => {
    e.stopPropagation();
    navigate({ to: '/search', search: { q: album } });
  };

  const handleDownload = async (e: React.MouseEvent, track: Track) => {
    e.stopPropagation();
    try {
      const { downloadUrl, fileName } = await getDownloadUrl(track.id);
      const link = document.createElement('a');
      link.href = downloadUrl;
      link.download = fileName || `${track.title}.${track.format}`;
      link.target = '_blank';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (err) {
      console.error('Failed to download track:', err);
    }
  };

  if (isLoading) {
    return (
      <div data-testid="track-list-skeleton" className="space-y-2">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="skeleton h-12 w-full" />
        ))}
      </div>
    );
  }

  if (tracks.length === 0) {
    return (
      <div className="text-center py-12 text-base-content/60">
        No tracks found
      </div>
    );
  }

  const handleTrackClick = (_track: Track, index: number) => {
    setQueue(tracks, index);
  };

  return (
    <table role="table" className="table w-full">
      <thead>
        <tr>
          <th className="w-12">#</th>
          <th>Title</th>
          <th>Artist</th>
          <th>Album</th>
          <th className="w-20">Duration</th>
          {showAddedDate && <th>Added</th>}
          {(showDownload || showAddToPlaylist) && <th className="w-24">Actions</th>}
        </tr>
      </thead>
      <tbody>
        {tracks.map((track, index) => (
          <tr
            key={track.id}
            role="row"
            className={`hover:bg-base-200 cursor-pointer ${
              currentTrack?.id === track.id ? 'bg-primary/10' : ''
            }`}
            onClick={() => handleTrackClick(track, index)}
          >
            <td>
              {currentTrack?.id === track.id && isPlaying ? (
                <span className="text-primary">▶</span>
              ) : (
                index + 1
              )}
            </td>
            <td className="font-medium">{track.title}</td>
            <td onClick={(e) => e.stopPropagation()}>
              <button
                type="button"
                className="text-base-content/70 hover:text-primary hover:underline cursor-pointer"
                onClick={(e) => handleArtistClick(e, track.artist)}
              >
                {track.artist}
              </button>
            </td>
            <td onClick={(e) => e.stopPropagation()}>
              <button
                type="button"
                className="text-base-content/70 hover:text-primary hover:underline cursor-pointer"
                onClick={(e) => handleAlbumClick(e, track.album)}
              >
                {track.album}
              </button>
            </td>
            <td className="tabular-nums">{formatDuration(track.duration)}</td>
            {showAddedDate && (
              <td className="text-sm text-base-content/60">{formatDate(track.createdAt)}</td>
            )}
            {(showDownload || showAddToPlaylist) && (
              <td onClick={(e) => e.stopPropagation()}>
                <div className="flex gap-1">
                  {showAddToPlaylist && (
                    <AddToPlaylistDropdown trackId={track.id} />
                  )}
                  {showDownload && (
                    <button
                      className="btn btn-ghost btn-xs"
                      onClick={(e) => handleDownload(e, track)}
                      aria-label="Download"
                      title="Download for offline"
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                    </button>
                  )}
                </div>
              </td>
            )}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
