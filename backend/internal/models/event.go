package models

import "time"

// Event represents a live music event or concert.
type Event struct {
	ID         string    `json:"id"`
	ArtistName string    `json:"artistName"`
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	Venue      string    `json:"venue"`
	City       string    `json:"city"`
	Region     string    `json:"region"`
	Country    string    `json:"country"`
	TicketURL  string    `json:"ticketUrl,omitempty"`
	Status     string    `json:"status"`
	Source     string    `json:"source"`
}

// EventResponse represents an event in API responses.
type EventResponse = Event

// ArtistSearchResult represents an artist found via event search.
type ArtistSearchResult struct {
	Name           string `json:"name"`
	ImageURL       string `json:"imageUrl,omitempty"`
	UpcomingEvents int    `json:"upcomingEvents"`
	Source         string `json:"source"`
}

// ArtistSearchResultResponse represents an artist search result in API responses.
type ArtistSearchResultResponse = ArtistSearchResult

// ArtistEventsResponse represents aggregated events for an artist in API responses.
type ArtistEventsResponse struct {
	ArtistName string  `json:"artistName"`
	Events     []Event `json:"events"`
	TotalCount int     `json:"totalCount"`
	Source     string  `json:"source"`
}
