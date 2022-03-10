package google

import (
	"be/api"
	"be/delivery/controllers/templates"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
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
		var state = api.GenerateStateOauthCookie(c.Response().Writer)

		var url = cont.conf.AuthCodeURL(state)

		res := c.Redirect(http.StatusSeeOther, url)

		if res != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error Redirect to sign in "+res.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(nil, "success redirect to sign in", nil))
	}
}

