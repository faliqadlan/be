package api

import (
	"be/configs"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
)

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func GenerateCookie(c echo.Context) (string, error) {
	var expiration = time.Now().Add(time.Second)
	var b = make([]byte, 16)
	rand.Read(b)
	var state = base64.URLEncoding.EncodeToString(b)

	var cookie = new(http.Cookie)

	cookie.Name = "oauthState"
	cookie.Value = state
	cookie.Expires = expiration
	c.SetCookie(cookie)
	return state, c.String(http.StatusOK, "write cookie")
}

func ReadCookie(c echo.Context) error {
	cookie, err := c.Cookie("state")
	if err != nil {
		return err
	}
	log.Info(cookie.Name, cookie.Value)
	return c.String(http.StatusOK, "write cookie")
}

func GetUserDataFromGoogle(code string, conf *oauth2.Config) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(configs.OauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
