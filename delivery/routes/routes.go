package routes

import (
	"be/delivery/controllers/auth"
	"be/delivery/controllers/doctor"
	"be/delivery/controllers/patient"
	"be/delivery/controllers/visit"
	"be/delivery/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RoutesPath(e *echo.Echo, ac *auth.AuthController, dc *doctor.Controller, pc *patient.Controller, vc *visit.Controller) {
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}",
	}))
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

	// with jwt

	var g = e.Group("")

	g.Use(middlewares.JwtMiddleware())

	// doctor =================================

	g.PUT("/doctor", dc.Update())
	g.DELETE("/doctor", dc.Delete())
	g.GET("/doctor/profile", dc.GetProfile())
	g.GET("/doctor/patient/all", dc.GetPatients())
	g.GET("/doctor/dashboard", dc.GetDashboard())
	g.GET("/doctor/all", dc.GetAll())

	// patient ===================================

	g.PUT("/patient", pc.Update())
	g.DELETE("/patient", pc.Delete())
	g.GET("/patient/profile", pc.GetProfile())
	g.GET("/patient/record", pc.GetRecords())
	g.GET("/patient/history", pc.GetHistories())
	g.GET("/patient/appontment", pc.GetAppointMent())

	// visit

	g.POST("/visit", vc.Create())
	g.PUT("/visit/:visit_uid", vc.Update())

}
