package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/gommon/log"
)

type Token struct {
	Client_uid                  string
	Project_id                  string
	Auth_uri                    string
	Token_uri                   string
	Auth_provider_x509_cert_url string
	Client_secret               string
}
type Web struct {
	Web Token
}

func CreteTokenJson(client_uid, project_id, auth_uri, token_uri, auth_provider_x509_cert_url, client_secret string) error {

	data := Web{
		Web: Token{
			Client_uid:                  client_uid,
			Project_id:                  project_id,
			Auth_uri:                    auth_uri,
			Token_uri:                   token_uri,
			Auth_provider_x509_cert_url: auth_provider_x509_cert_url,
			Client_secret:               client_secret,
		},
	}

	file, err := json.MarshalIndent(data, "", "")
	if err != nil {
		log.Info(err)
	}

	res := ioutil.WriteFile("tok.json", file, 0644)
	return res
}
