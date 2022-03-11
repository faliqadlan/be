package google

import (
	"be/api"
	"be/delivery/controllers/templates"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Controller struct {
	conf *oauth2.Config
}

func New(conf *oauth2.Config) *Controller {
	return &Controller{
		conf: conf,
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
		var url = cont.conf.AuthCodeURL( /* oauthState */ "randomstate")
		// log.Info(url)
		res := c.Redirect(http.StatusSeeOther, url)
		// log.Info(res)
		if res != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error Redirect to sign in "+res.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(nil, "success redirect to sign in", nil))
	}
}

func (cont *Controller) GoogleCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		state := c.Request().URL.Query()["state"][0]
		// log.Info(state)
		if state != "state" {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error cookie "+state, nil))
		}

		// var err =api.ReadCookie(c)
		// log.Info(err)
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error in callback "+err.Error(), nil))
		// }

		// code := c.Request().URL.Query()["code"][0]

		data, err := api.GetUserDataFromGoogle(c.FormValue("code"), cont.conf)
		if err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error in callback "+err.Error(), nil))
		}

		// log.Info(string(data))
		return c.JSON(http.StatusOK, templates.Success(nil, "success redirect to sign in", string(data)))
	}
}

func (cont *Controller) GoogleCalendar() echo.HandlerFunc {
	return func(c echo.Context) error {
		var ctx = context.Background()
		// state := c.Request().URL.Query()["state"][0]
		// if state != "state" {
		// 	return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error cookie "+state, nil))
		// }

		var code = c.FormValue("code")
		token, err := cont.conf.Exchange(context.Background(), code)
		if err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error in get token "+err.Error(), nil))
		}
		var client = cont.conf.Client(ctx, token)
		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error in run calendar "+err.Error(), nil))
		}

		var event = &calendar.Event{
			Summary:     "event from goolang",
			Location:    "golang",
			Description: "test insert event from goolang",
			Start: &calendar.EventDateTime{
				DateTime: time.Now().Local().Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: time.Now().Local().Format(time.RFC3339),
			},
		}

		event, err = srv.Events.Insert("primary", event).Do()
		if err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error in run calendar "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(nil, "success ru calendar", event))
	}
}
