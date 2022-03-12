package google

import (
	"context"

	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func SetupConfig(db_username, clientId, clientScret string) *oauth2.Config {

	// var url = urls[0]
	var conf = &oauth2.Config{}
	if db_username != "root" {
		conf = &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientScret,
			RedirectURL:  "https://faliqadlan.cloud.okteto.net/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/calendar",
			},
			Endpoint: google.Endpoint,
		}
	} else {
		conf = &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientScret,
			RedirectURL:  "http://localhost:8080/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/calendar",
			},
			Endpoint: google.Endpoint,
		}
	}

	return conf
}

func InitCalendar(b []byte, token *oauth2.Token) *calendar.Service {
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		log.Warn(err)
	}
	return srv
}
