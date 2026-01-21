import { Track } from '@/types';
import { usePlayerStore } from '@/lib/store/playerStore';

interface TrackListProps {
  tracks: Track[];
  isLoading?: boolean;
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function TrackList({ tracks, isLoading }: TrackListProps) {
  const { setQueue, currentTrack, isPlaying } = usePlayerStore();

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
            <td className="text-base-content/70">{track.artist}</td>
            <td className="text-base-content/70">{track.album}</td>
            <td className="tabular-nums">{formatDuration(track.duration)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
