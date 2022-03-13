package google

import (
	"be/api"
	"be/configs"
	"be/delivery/controllers/templates"
	"be/repository/visit"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
)

type Controller struct {
	r    visit.Visit
	conf *oauth2.Config
}

func New(conf *oauth2.Config, r visit.Visit) *Controller {
	return &Controller{
		conf: conf,
		r:    r,
	}
}

func (cont *Controller) GoogleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		// var oauthState = api.GenerateStateOauthCookie(c.Response().Writer)

		// var oauthState, err = api.GenerateCookie(c)
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error get cookie "+ err.Error(), nil))
		// }

		// log.Info(oauthState)
		// var url = cont.conf.AuthCodeURL("randomstate", oauth2.AccessTypeOffline)
		// var url = cont.conf.AuthCodeURL( /* oauthState */ "randomstate")
		var url = cont.conf.AuthCodeURL("randomstate", oauth2.AccessTypeOffline)

		// log.Info(url)
		res := c.Redirect(http.StatusSeeOther, url)
		// log.Info(res)
		if res != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error Redirect to sign in "+res.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(nil, "success redirect to sign in", nil))
	}
}

func (cont *Controller) GoogleCalendar() echo.HandlerFunc {
	return func(c echo.Context) error {
		// var ctx = context.Background()
		// state := c.Request().URL.Query()["state"][0]
		// if state != "state" {
		// 	return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error cookie "+state, nil))
		// }

		// var code = c.FormValue("code")

		// _, token, err := api.GetUserDataFromGoogle(code, cont.conf)
		// if err != nil {
		// 	return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error get user profile info "+err.Error(), nil))
		// }

		// api.CacheToken(token, "./secret/token.json")

		// token, err = api.TokenFromFile("./secret/token.json")
		// if err != nil {
		// 	log.Warn(err)
		// }

		token, err := api.TokenFromFile("./secret/token.json", cont.conf)
		if err != nil {
			log.Warn(err)
		}

		// var conf = cont.conf.Client(context.Background(), tokenNew)

		response, err := http.Get(configs.OauthGoogleUrlAPI + token.AccessToken)
		if err != nil {
			log.Warn(err)
		}
		// defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Warn(err)
		}

		var user UserInfo

		if err := json.Unmarshal(contents, &user); err != nil {
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error convert to struct "+err.Error(), nil))
		}
		// api.CacheToken(token)
		log.Info(user)

		var res = api.Calendar(configs.TokenPath, cont.conf)

		return c.JSON(http.StatusOK, templates.Success(nil, "success run calendar", map[string]interface{}{
			"user":     user,
			"calendar": res,
		}))
	}
}

// func (cont *Controller) GoogleCalendarEventList() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var ctx = context.Background()
// 		// state := c.Request().URL.Query()["state"][0]
// 		// if state != "state" {
// 		// 	return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error cookie "+state, nil))
// 		// }

// 		var code = c.FormValue("code")

// 		_, token, err := api.GetUserDataFromGoogle(code, cont.conf)
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error get user profile info "+err.Error(), nil))
// 		}

// 		var client = cont.conf.Client(ctx, token)
// 		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
// 		if err != nil {
// 			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error in run calendar "+err.Error(), nil))
// 		}
// 		t := time.Now().Format(time.RFC3339)
// 		events, err := srv.Events.List("primary").s

// 		return c.JSON(http.StatusOK, templates.Success(nil, "success run calendar", nil))
// 	}
// }
