package main

import (
	"be/api"
	"be/configs"
	"be/delivery/controllers/auth"
	"be/delivery/controllers/doctor"
	"be/delivery/controllers/patient"
	"be/delivery/controllers/visit"
	"be/delivery/routes"
	authRepo "be/repository/auth"
	doctorRepo "be/repository/doctor"
	patientRepo "be/repository/patient"
	visitRepo "be/repository/visit"
	"be/utils"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var awsS3Conf = utils.InitS3(config.S3_REGION, config.S3_ID, config.S3_SECRET)

	var authRepo = authRepo.New(db)
	var authCont = auth.New(authRepo)

	var doctorRepo = doctorRepo.New(db)
	var doctorCont = doctor.New(doctorRepo, awsS3Conf)

	var patientRepo = patientRepo.New(db)
	var patientCont = patient.New(patientRepo, awsS3Conf)

	var visitRepo = visitRepo.New(db)
	var visitCont = visit.New(visitRepo)

	api.Calendar()

	var e = echo.New()

	routes.RoutesPath(e, authCont, doctorCont, patientCont, visitCont)

	log.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))

}
