package calendar

import (
	"be/repository/visit"

	"google.golang.org/api/calendar/v3"
)

type Calendar interface {
	CreateEvent(res visit.VisitCalendar) (*calendar.Event, error)
	InsertEvent(event *calendar.Event) (*calendar.Event, error)
	UpdateEvent(event *calendar.Event, event_uid string) (*calendar.Event, error)
	DeleteEvent(event_uid string) error
}
