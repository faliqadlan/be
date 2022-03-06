package main

import (
	"be/configs"
	"be/delivery/controllers/auth"
	"be/delivery/controllers/user"
	"be/delivery/routes"
	authRepo "be/repository/auth"
	userRepo "be/repository/user"
	"be/utils"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	config := configs.GetConfig()
	db := utils.InitDB(config)
	// log.Info(db)

	userRepo := userRepo.New(db)
	userCont := user.New(userRepo)

	authRepo := authRepo.New(db)
	authCont := auth.New(authRepo)

	e := echo.New()

	routes.RoutesPath(e, userCont, authCont)

	log.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))

}
