package calendar

import "google.golang.org/api/calendar/v3"

type Calendar interface {
	CreateEvent(visit_uid string) (*calendar.Event, error)
	InsertEvent(event *calendar.Event) error
}
