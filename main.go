package main

import (
	"be/api"
	"be/api/aws"
	"be/api/aws/s3"
	googleApi "be/api/google"
	"be/api/google/calendar"
	"be/configs"
	"be/delivery/controllers/auth"
	"be/delivery/controllers/doctor"
	"be/delivery/controllers/google"
	"be/delivery/controllers/patient"
	"be/delivery/controllers/visit"
	"be/delivery/routes"
	authRepo "be/repository/auth"
	doctorRepo "be/repository/doctor"
	patientRepo "be/repository/patient"
	visitRepo "be/repository/visit"
	logicDoctor "be/delivery/logic/doctor"

	"be/utils"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var awsS3Conf = aws.InitS3(config.S3_REGION, config.S3_ID, config.S3_SECRET)

	var googleConf = googleApi.SetupConfig(config.DB_USERNAME, config.CLIENT_ID, config.CLIENT_SECRET)

	var awsS3 = s3.New(awsS3Conf)

	var b, token = api.TokenInit(configs.CredentialPath, configs.TokenPath, googleConf)

	var srv = googleApi.InitCalendar(b, token)

	var authRepo = authRepo.New(db)
	var authCont = auth.New(authRepo)

	var doctorRepo = doctorRepo.New(db)
	var doctorLogic = logicDoctor.New()
	var doctorCont = doctor.New(doctorRepo, awsS3, doctorLogic)

	var patientRepo = patientRepo.New(db)
	var patientCont = patient.New(patientRepo, awsS3)

	var visitRepo = visitRepo.New(db)
	var calendar = calendar.New(visitRepo, srv)
	var visitCont = visit.New(visitRepo, calendar)

	var googleCont = google.New(googleConf, visitRepo)

	var e = echo.New()

	routes.RoutesPath(e, authCont, doctorCont, patientCont, visitCont, googleCont)

	log.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))

}
