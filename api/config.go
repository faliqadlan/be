package api

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupConfig(clientId, clientScret string) *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientScret,
		RedirectURL:  "http://localhost:8080/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
