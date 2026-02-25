/**
 * ArtistEventsSection Component - Show upcoming events for an artist
 */
import { useArtistEvents } from '../../hooks/useArtistEvents';
import { EventCard } from './EventCard';

interface ArtistEventsSectionProps {
  artistName: string;
}

export function ArtistEventsSection({ artistName }: ArtistEventsSectionProps) {
  const { data, isLoading, isError } = useArtistEvents(artistName);

  if (isLoading) {
    return <div className="text-center py-4">Loading events...</div>;
  }

  if (isError) {
    return <div className="text-center py-4 text-error">Failed to load events</div>;
  }

  if (!data?.events?.length) {
    return <div className="text-center py-4 text-base-content/60">No upcoming events</div>;
  }

  return (
    <div className="space-y-2">
      <h3 className="font-semibold">Upcoming Events</h3>
      <div className="space-y-2">
        {data.events.map((event) => (
          <EventCard key={event.id} event={event} />
        ))}
      </div>
    </div>
  );
}
