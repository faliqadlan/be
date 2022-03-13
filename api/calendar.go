package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func Calendar(file string, conf *oauth2.Config) string {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credential.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return err.Error()
	}
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return err.Error()
	}

	tokenInit, err := TokenFromFile(file, conf)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return err.Error()
	}

	srv, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tokenInit)))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
		return err.Error()
	}

	event := &calendar.Event{
		Summary:     "golang test",
		Location:    "golang3",
		Description: "test golang",
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(time.Hour).Format(time.RFC3339),
		},
	}
	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
		return err.Error()
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)

	return "success run calendar"
}
