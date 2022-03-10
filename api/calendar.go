package api

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func TestApi(client_id, client_secret string) string {
	var ctx = context.Background()
	var client = InitConfig(client_id, client_secret)
	// conf := &oauth2.Config{
	// 	ClientID:     "https://console.developers.google.com/project/mrclinic/apiui/" + client_id,
	// 	ClientSecret: "https://console.developers.google.com/project/mrclinic/apiui/" + client_secret,
	// 	Scopes:       []string{urlshortener.UrlshortenerScope},
	// 	Endpoint:     google.Endpoint,
	// 	RedirectURL:  "https://oauth2.example.com/code",
	// }

	// url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	// fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// var code string
	// if _, err := fmt.Scan(&code); err != nil {
	// 	log.Info(err)
	// }

	// tok, err := conf.Exchange(ctx, code)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// client := conf.Client(ctx, tok)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Info(err)
	}

	event := &calendar.Event{
		Summary:     "Google I/O 2015",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A chance to hear more about Google's developer products.",
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Local().String(),
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().AddDate(0, 0, 7).Local().String(),
		},
		Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=2"},
		Attendees: []*calendar.EventAttendee{
			{Email: "wonganteng8@gmail.com"},
			{Email: "sbrin@example.com"},
		},
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		log.Info("Unable to create event. %v\n", err)
	}
	log.Info(event)

	return "test done"
}
