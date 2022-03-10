package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func getTokenFromWeb(client_id, client_secret string) *oauth2.Token {

	var config = &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	log.Info(authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func TestApi(client_id, client_secret string) string {
	var ctx = context.Background()

	tok, err := tokenFromFile("token.json")

	if err != nil {
		tok = getTokenFromWeb(client_id, client_secret)
	}

	var config *oauth2.Config

	client := config.Client(ctx, tok)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Info(err)
	}

	log.Info(srv)

	return "test done"
}
