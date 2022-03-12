package api

import (
	"be/configs"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
)

func GetUserDataFromGoogle(code string, conf *oauth2.Config) ([]byte, *oauth2.Token, error) {
	// Use code to get token and get user info from Google.

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(configs.OauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, token, nil
}

func CacheToken(token *oauth2.Token) error {

	tokenByte, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("token.json", tokenByte, 0644)
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

func TokenInit(credentialPath, tokenPath string) ([]byte, *oauth2.Token) {
	b, err := ioutil.ReadFile(credentialPath)
	if err != nil {
		log.Warn(err)
	}

	token, err := TokenFromFile(tokenPath)
	if err != nil {
		log.Warn(err)
	}
	return b, token
}

type Credential struct {
	Client_uid                  string
	Project_id                  string
	Auth_uri                    string
	Token_uri                   string
	Auth_provider_x509_cert_url string
	Client_secret               string
	Redirect_uris               []string
}
type Web struct {
	Web Credential
}

func CreteCredentialJson(client_uid, project_id, auth_uri, token_uri, auth_provider_x509_cert_url, client_secret string) error {

	var redirect_uris = []string{
		"http://localhost:8080/google/callback", "https://faliqadlan.cloud.okteto.net/google/callback",
	}

	data := Web{
		Web: Credential{
			Client_uid:                  client_uid,
			Project_id:                  project_id,
			Auth_uri:                    auth_uri,
			Token_uri:                   token_uri,
			Auth_provider_x509_cert_url: auth_provider_x509_cert_url,
			Client_secret:               client_secret,
			Redirect_uris:               redirect_uris,
		},
	}

	file, err := json.MarshalIndent(data, "", "")
	if err != nil {
		log.Warn(err)
	}

	err = ioutil.WriteFile("credential.json", file, 0644)
	if err != nil {
		log.Warn(err)
	}
	return nil
}

type Token struct {
	Access_token  string
	Token_type    string
	Refresh_token string
	Expiry        string
}

func CreateTokenJson(access_token, token_type, refresh_token string) error {
	data := Token{
		Access_token:  access_token,
		Token_type:    token_type,
		Refresh_token: refresh_token,
		Expiry:        "2022-03-12T23:25:25.445233+07:00",
	}

	file, err := json.MarshalIndent(data, "", "")
	if err != nil {
		log.Warn(err)
	}

	err = ioutil.WriteFile("token.json", file, 0644)
	if err != nil {
		log.Warn(err)
	}
	return nil
}
