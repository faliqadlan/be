package calendar

import (
	"be/repository/visit"

	"google.golang.org/api/calendar/v3"
)

type Calendar interface {
	CreateEvent(res visit.VisitCalendar) (*calendar.Event, error)
	InsertEvent(event *calendar.Event) error
}
