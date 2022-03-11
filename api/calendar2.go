package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const (
	cacheTokenDir  = "./tmp/google_tokens"
	credentialPath = "credentials.json"
)

func cacheToken(token *oauth2.Token) error {
	tokenByte, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cacheTokenDir, tokenByte, 0644)
}

func generateNewToken(config *oauth2.Config) (*oauth2.Token, error) {
	fmt.Println("fetching new google token")
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	fmt.Println(token)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	cacheToken(token)
	return token, nil

}

func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func Calendar() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("token.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tokenInit, err := generateNewToken(config)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	

	srv, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tokenInit)))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
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
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)

}
