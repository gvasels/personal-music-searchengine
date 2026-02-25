/**
 * EventCard Component - Display a single artist event
 */
import type { ArtistEvent } from '../../types';

interface EventCardProps {
  event: ArtistEvent;
}

export function EventCard({ event }: EventCardProps) {
  const date = new Date(event.date);
  const formattedDate = date.toLocaleDateString(undefined, {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });

  return (
    <div className="card bg-base-200 p-4">
      <div className="flex justify-between items-start">
        <div>
          <h4 className="font-semibold">{event.title}</h4>
          <p className="text-sm text-base-content/70">{event.venue}</p>
          <p className="text-sm text-base-content/60">
            {event.city}, {event.region} · {event.country}
          </p>
          <p className="text-xs text-base-content/50 mt-1">{formattedDate}</p>
        </div>
        <div className="flex flex-col items-end gap-1">
          {event.status !== 'scheduled' && (
            <span className={`badge badge-sm ${event.status === 'cancelled' ? 'badge-error' : 'badge-warning'}`}>
              {event.status === 'cancelled' ? 'Cancelled' : 'Postponed'}
            </span>
          )}
          {event.ticketUrl && (
            <a
              href={event.ticketUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="btn btn-xs btn-primary"
            >
              Tickets
            </a>
          )}
        </div>
      </div>
    </div>
  );
}
