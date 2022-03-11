package api

import (
	"be/configs"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

	token.Expiry = token.Expiry.AddDate(1, 0,0 )


	tokenByte, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("tokenTest.json", tokenByte, 0644)
}
