package api

import (
	"be/configs"
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func Calendar(file string, conf *oauth2.Config) string {
	ctx := context.Background()
	b, err := ioutil.ReadFile(configs.CredentialPath)
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

	var event *calendar.Event
	// for {

		var uid = uuid.New().ClockSequence()
		var uidS = strconv.Itoa(uid)
		log.Info(uidS)
		log.Info(len(uidS))
		event = &calendar.Event{
			Id:          uidS,
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
			log.Warn("Unable to create event. %v\n", err)
			return err.Error()
		}
		log.Info(strings.Contains(err.Error(), "id"))
	// 	if !strings.Contains(err.Error(), "id") {
	// 		break
	// 	}
	// 	log.Info(event.Id)
	// }

	fmt.Printf("Event created: %s\n", event.HtmlLink)
	log.Info(event.Id)
	log.Info(event.ICalUID)
	log.Info(event.Id)

	// findEvent, err := srv.Events.Get(calendarId, uid).Do()
	// if err != nil {
	// 	log.Info(err)
	// }
	// log.Info(findEvent.Id)
	// log.Info(findEvent.ICalUID)
	// log.Info(findEvent.Id == uid)
	// log.Info(findEvent)

	return "success run calendar"
}
