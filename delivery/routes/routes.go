package routes

import (
	"be/delivery/controllers/auth"
	"be/delivery/controllers/clinic"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RoutesPath(e *echo.Echo, ac *auth.AuthController, cc *clinic.Controller) {
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}",
	}))
	/* no jwt */
	// User ====================================

	e.POST("/login", ac.Login())

	// clinic =================================

	e.POST("/clinic", cc.Create())

}
