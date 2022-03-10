package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth2.Token)
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

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

func newOAuthClient(ctx context.Context, config *oauth2.Config) *http.Client {
	token, err := tokenFromFile("tok.json")
	if err != nil {
		log.Info(err)
		token = getTokenFromWeb(config)
	}

	return config.Client(ctx, token)
}

func InitConfig(client_id, client_secret string) *http.Client {

	var config = &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"mrClinic"},
	}

	ctx := context.Background()
	// if *debug {
	// 	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
	// 		Transport: &logTransport{http.DefaultTransport},
	// 	})
	// }

	var c = newOAuthClient(ctx, config)

	return c
}
