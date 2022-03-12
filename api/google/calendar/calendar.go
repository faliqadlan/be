package calendar

import (
	"be/repository/visit"
	"time"

	"google.golang.org/api/calendar/v3"
)

type CalendarConf struct {
	r   visit.Visit
	srv *calendar.Service
}

func New(r visit.Visit, srv *calendar.Service) *CalendarConf {
	return &CalendarConf{
		r:   r,
		srv: srv,
	}
}

func (cal *CalendarConf) CreateEvent(visit_uid string) (*calendar.Event, error) {
	res, err := cal.r.GetVisitList(visit_uid)
	if err != nil {
		return &calendar.Event{}, err
	}
	var layout = "02-01-2006"
	dateConv, err := time.Parse(layout, res.Date)
	if err != nil {
		return &calendar.Event{}, err
	}

	var event = &calendar.Event{
		Summary:     "Apppoinment with " + res.DoctorName,
		Location:    res.Address,
		Description: res.Complaint,
		Start: &calendar.EventDateTime{
			DateTime: dateConv.Local().Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: dateConv.Add(24 * time.Hour).Local().Format(time.RFC3339),
		},
		Attendees: []*calendar.EventAttendee{
			{DisplayName: res.DoctorName, Email: res.DoctorEmail},
			{DisplayName: res.PatientName, Email: res.PatientEmail},
		},
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides: []*calendar.EventReminder{
				{Method: "email", Minutes: 24 * 60},
				{Method: "email", Minutes: 2 * 60},
				{Method: "email", Minutes: 1 * 60},
				{Method: "email", Minutes: 30},
				{Method: "email", Minutes: 15},
			},
			ForceSendFields: []string{"UseDefault"},
			NullFields:      nil,
		},
	}
	return event, nil
}

func (cal *CalendarConf) InsertEvent(event *calendar.Event) error {

	_, err := cal.srv.Events.Insert("primary", event).Do()
	if err != nil {
		return err
	}
	return nil
}
