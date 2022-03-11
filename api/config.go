package api

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupConfig(db_username, clientId, clientScret string) *oauth2.Config {

	// var urls = []string{"http://localhost:8080/google/callback, https://faliqadlan.cloud.okteto.net/google/callback"}

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
