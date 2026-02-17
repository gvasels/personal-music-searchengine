import { Link } from '@tanstack/react-router';
import { useSimilarTracks } from '../../hooks/useSimilarTracks';

interface SimilarTracksProps {
  trackId: string;
}

export function SimilarTracks({ trackId }: SimilarTracksProps) {
  const { data, isLoading } = useSimilarTracks(trackId);

  if (isLoading) {
    return <div className="text-center py-4">Loading similar tracks...</div>;
  }

  if (!data?.similar?.length) {
    return <div className="text-center py-4 text-base-content/60">No similar tracks found</div>;
  }

  return (
    <div className="space-y-2">
      <h3 className="font-semibold">Similar Tracks</h3>
      <ul className="space-y-1">
        {data.similar.map((item) => (
          <li key={item.track.id} className="flex items-center gap-3 p-2 bg-base-200 rounded hover:bg-base-300 transition-colors">
            {item.track.coverArtUrl ? (
              <img src={item.track.coverArtUrl} alt="" className="w-10 h-10 rounded object-cover" />
            ) : (
              <div className="w-10 h-10 rounded bg-base-300 flex items-center justify-center text-base-content/40">♪</div>
            )}
            <div className="flex-1 min-w-0">
              <Link to="/tracks/$trackId" params={{ trackId: item.track.id }} className="text-sm font-medium truncate block hover:underline">
                {item.track.title}
              </Link>
              <span className="text-xs text-base-content/60 truncate block">{item.track.artist}</span>
            </div>
            <div className="flex items-center gap-2 text-xs text-base-content/60 shrink-0">
              {item.track.bpm ? <span>{item.track.bpm} BPM</span> : null}
              {item.track.keyCamelot ? <span>{item.track.keyCamelot}</span> : null}
              <span className="badge badge-sm badge-ghost">{((1 - item.score) * 100).toFixed(0)}%</span>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
