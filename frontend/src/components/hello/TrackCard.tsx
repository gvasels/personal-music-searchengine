import type { HelloTrack } from '../../lib/api/helloSearch';

interface TrackCardProps {
  track: HelloTrack;
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function TrackCard({ track }: TrackCardProps) {
  const coverUrl =
    track.coverArtUrl ||
    `https://placehold.co/300x300/1a1a2e/e0e0e0?text=${encodeURIComponent(
      track.artist
        .split(' ')
        .map((w) => w[0])
        .join('')
    )}`;

  return (
    <div className="card bg-base-200 shadow-xl">
      <figure>
        <img
          src={coverUrl}
          alt={`${track.title} by ${track.artist}`}
          className="w-full h-48 object-cover"
        />
      </figure>
      <div className="card-body p-4">
        <h3 className="card-title text-base">{track.title}</h3>
        <p className="text-sm text-base-content/70">{track.artist}</p>
        <p className="text-xs text-base-content/50">{track.album}</p>
        <div className="card-actions justify-between items-center mt-2">
          <span className="badge badge-primary badge-sm">{track.genre}</span>
          <span className="text-xs text-base-content/60">
            {track.durationStr || formatDuration(track.duration)}
          </span>
        </div>
      </div>
    </div>
  );
}
