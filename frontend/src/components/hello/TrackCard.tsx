/**
 * TrackCard Component
 *
 * Displays track information in a card format.
 * Shows title, artist, album, genre, year, and formatted duration.
 */
import { HelloTrack } from '../../lib/api/helloSearch';

interface TrackCardProps {
  track: HelloTrack;
}

export function TrackCard({ track }: TrackCardProps) {
  /**
   * Format duration from seconds to "M:SS" format.
   * Examples: 240 -> "4:00", 185 -> "3:05"
   */
  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${String(secs).padStart(2, '0')}`;
  };

  return (
    <div className="card bg-base-200">
      <div className="card-body">
        <h3 className="card-title">{track.title}</h3>
        <p>{track.artist}</p>
        <p>{track.album}</p>
        <div className="flex gap-2">
          <span className="badge">{track.genre}</span>
          <span>{track.year}</span>
          <span>{formatDuration(track.duration)}</span>
        </div>
      </div>
    </div>
  );
}
