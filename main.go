package main

import (
	"be/configs"
	"be/delivery/controllers/auth"
	"be/delivery/controllers/clinic"
	"be/delivery/routes"
	authRepo "be/repository/auth"
	clinicRepo "be/repository/clinic"
	"be/utils"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	config := configs.GetConfig()
	db := utils.InitDB(config)
	// log.Info(db)

	authRepo := authRepo.New(db)
	authCont := auth.New(authRepo)

	clinicRepo := clinicRepo.New(db)
	clinicCont := clinic.New(clinicRepo)

	e := echo.New()

	routes.RoutesPath(e, authCont, clinicCont)

	log.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))

}
