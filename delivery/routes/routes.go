package routes

import (
	"be/delivery/controllers/auth"
	"be/delivery/controllers/doctor"
	"be/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RoutesPath(e *echo.Echo, ac *auth.AuthController, dc *doctor.Controller) {
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}",
	}))
	/* no jwt */
	// User ====================================

	e.POST("/login", ac.Login())

	// clinic =================================

	e.POST("/doctor", dc.Create())

	// with jwt

	var g = e.Group("", middlewares.JwtMiddleware())

	// doctor ===

	g.PUT("/doctor", dc.Update())
	g.DELETE("/doctor", dc.Delete())
	g.GET("/doctor/profile", dc.GetProfile())
	g.GET("/doctor/patient/all", dc.GetPatients())
	g.GET("/doctor/dashboard", dc.GetDashboard())

}
