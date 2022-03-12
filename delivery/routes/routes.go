package routes

import (
	"be/delivery/controllers/auth"
	"be/delivery/controllers/doctor"
	"be/delivery/controllers/google"
	"be/delivery/controllers/patient"
	"be/delivery/controllers/visit"
	"be/delivery/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RoutesPath(e *echo.Echo, ac *auth.AuthController, dc *doctor.Controller, pc *patient.Controller, vc *visit.Controller, gc* google.Controller) {
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}",
	}))
	// e.AcquireContext().Cookies()
	/* no jwt */

	e.GET("/test", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
	})

	// login ====================================

	e.POST("/login", ac.Login())

	// doctor =================================

	e.POST("/doctor", dc.Create())

	// patient

	e.POST("/patient", pc.Create())

	// google

	e.GET("/google/login", gc.GoogleLogin())
	e.GET("google/callback", gc.GoogleCalendar())
	

	// with jwt

	var g = e.Group("")

	g.Use(middlewares.JwtMiddleware())

	// doctor =================================

	g.PUT("/doctor", dc.Update())
	g.DELETE("/doctor", dc.Delete())
	g.GET("/doctor/profile", dc.GetProfile())
	g.GET("/doctor/all", dc.GetAll())

	// patient ===================================

	g.PUT("/patient", pc.Update())
	g.DELETE("/patient", pc.Delete())
	g.GET("/patient/profile", pc.GetProfile())

	// visit

	g.POST("/visit", vc.Create())
	g.PUT("/visit/:visit_uid", vc.Update())
	g.DELETE("/visit/:visit_uid", vc.Delete())
	g.GET("/visit", vc.GetVisits())

}
