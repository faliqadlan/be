package google

import (
	"be/api"
	"be/delivery/controllers/templates"
	"be/repository/visit"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
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

func (cont *Controller) GoogleCalendar() echo.HandlerFunc {
	return func(c echo.Context) error {
		var ctx = context.Background()
		// state := c.Request().URL.Query()["state"][0]
		// if state != "state" {
		// 	return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error cookie "+state, nil))
		// }

		var code = c.FormValue("code")

		profile, token, err := api.GetUserDataFromGoogle(code, cont.conf)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error get user profile info "+err.Error(), nil))
		}

		var user UserInfo

		if err := json.Unmarshal(profile, &user); err != nil {
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error convert to struct "+err.Error(), nil))
		}

		log.Info(user)

		res, err := cont.r.GetVisitList(user.Email, "pending")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error in connect database "+err.Error(), nil))
		}

		var layout = "02-01-2006"
		dateConv, err := time.Parse(layout, res.Date)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error in parsing date "+err.Error(), nil))
		}

		var client = cont.conf.Client(ctx, token)
		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error in run calendar "+err.Error(), nil))
		}

		var eventInput = &calendar.Event{
			Summary:     "Apppoinment with " + res.DoctorName,
			Location:    res.Address,
			Description: res.Complaint,
			Start: &calendar.EventDateTime{
				DateTime: dateConv.Local().Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: dateConv.Add(24 * time.Hour).Local().Format(time.RFC3339),
			},
			Attendees: []*calendar.EventAttendee{
				{DisplayName: res.DoctorName, Email: res.DoctorEmail},
			},
			Reminders: &calendar.EventReminders{
				UseDefault: false,
				Overrides: []*calendar.EventReminder{
					{Method: "email", Minutes: 24 * 60},
					{Method: "email", Minutes: 2 * 60},
					{Method: "email", Minutes: 1 * 60},
					{Method: "email", Minutes: 30},
					{Method: "email", Minutes: 15},
				},
				ForceSendFields: []string{"UseDefault"},
				NullFields:      nil,
			},
		}

		event, err := srv.Events.Insert("primary", eventInput).Do()
		if err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error in run calendar "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(nil, "success run calendar", event))
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