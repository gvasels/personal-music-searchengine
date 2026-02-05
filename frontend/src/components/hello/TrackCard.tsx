import type { HelloTrack } from '../../lib/api/helloSearch';

interface TrackCardProps {
  track: HelloTrack;
}

function formatDuration(seconds: number): string {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
}

export function TrackCard({ track }: TrackCardProps) {
  return (
    <div className="card bg-base-200 shadow-sm">
      <div className="card-body">
        <h3 className="card-title">{track.title}</h3>
        <p>{track.artist}</p>
        <p>{track.album}</p>
        {track.genre && <p>{track.genre}</p>}
        {track.year && <p>{track.year}</p>}
        <p>{formatDuration(track.duration)}</p>
      </div>
    </div>
  );
}
