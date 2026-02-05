interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;
}

function formatDuration(seconds: number): string {
  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${minutes}:${secs.toString().padStart(2, '0')}`;
}

export function TrackCard({ track }: { track: HelloTrack }) {
  return (
    <div className="card bg-base-200">
      <div className="card-body">
        <h3 className="card-title text-primary">{track.title}</h3>
        <p>{track.artist}</p>
        <p>{track.album}</p>
        <span className="badge">{track.genre}</span>
        <span>{track.year}</span>
        <span>{formatDuration(track.duration)}</span>
      </div>
    </div>
  );
}
