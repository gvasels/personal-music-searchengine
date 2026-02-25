/**
 * EventCard Component Tests - TDD Red Phase
 *
 * These tests define the contract for the EventCard component.
 * They MUST FAIL because EventCard.tsx does not exist yet.
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { EventCard } from '../EventCard';
import type { ArtistEvent } from '../../../types';

describe('EventCard', () => {
  const mockEvent: ArtistEvent = {
    id: 'evt-1',
    artistName: 'Daft Punk',
    title: 'Alive 2027 Tour',
    date: '2027-06-15T20:00:00Z',
    venue: 'Madison Square Garden',
    city: 'New York',
    region: 'NY',
    country: 'US',
    ticketUrl: 'https://tickets.example.com/daft-punk',
    status: 'scheduled',
    source: 'ticketmaster',
  };

  it('should render event title', () => {
    render(<EventCard event={mockEvent} />);
    expect(screen.getByText('Alive 2027 Tour')).toBeInTheDocument();
  });

  it('should render venue name', () => {
    render(<EventCard event={mockEvent} />);
    expect(screen.getByText('Madison Square Garden')).toBeInTheDocument();
  });

  it('should render city and country', () => {
    render(<EventCard event={mockEvent} />);
    expect(screen.getByText(/New York/)).toBeInTheDocument();
    expect(screen.getByText(/US/)).toBeInTheDocument();
  });

  it('should render ticket link when ticketUrl present', () => {
    render(<EventCard event={mockEvent} />);
    const link = screen.getByRole('link', { name: /ticket/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute('href', 'https://tickets.example.com/daft-punk');
  });

  it('should open ticket link in new tab', () => {
    render(<EventCard event={mockEvent} />);
    const link = screen.getByRole('link', { name: /ticket/i });
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
  });

  it('should not render ticket link when ticketUrl is missing', () => {
    const eventNoTicket = { ...mockEvent, ticketUrl: undefined };
    render(<EventCard event={eventNoTicket} />);
    expect(screen.queryByRole('link', { name: /ticket/i })).not.toBeInTheDocument();
  });

  it('should show Cancelled badge when status is cancelled', () => {
    const cancelledEvent = { ...mockEvent, status: 'cancelled' as const };
    render(<EventCard event={cancelledEvent} />);
    expect(screen.getByText(/cancelled/i)).toBeInTheDocument();
  });

  it('should show Postponed badge when status is postponed', () => {
    const postponedEvent = { ...mockEvent, status: 'postponed' as const };
    render(<EventCard event={postponedEvent} />);
    expect(screen.getByText(/postponed/i)).toBeInTheDocument();
  });

  it('should not show status badge when status is scheduled', () => {
    render(<EventCard event={mockEvent} />);
    expect(screen.queryByText(/cancelled/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/postponed/i)).not.toBeInTheDocument();
  });
});
