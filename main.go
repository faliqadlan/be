package main

import (
	"be/api"
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
	"be/utils"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var awsS3Conf = utils.InitS3(config.S3_REGION, config.S3_ID, config.S3_SECRET)

	api.CreteCredentialJson(config.CLIENT_ID, config.PROJECT_ID, config.AUTH_URI, config.TOKEN_URI, config.Auth_provider_x509_cert_url, config.CLIENT_SECRET)

	var googleConf = googleApi.SetupConfig(config.DB_USERNAME, config.CLIENT_ID, config.CLIENT_SECRET)

	// api.CreateTokenJson(config.Access_token, config.Token_type, config.Refresh_token)

	var b, token = api.TokenInit(configs.CredentialPath, configs.TokenPath, googleConf)

	var srv = googleApi.InitCalendar(b, token)

	var authRepo = authRepo.New(db)
	var authCont = auth.New(authRepo)

	var doctorRepo = doctorRepo.New(db)
	var doctorCont = doctor.New(doctorRepo, awsS3Conf)

	var patientRepo = patientRepo.New(db)
	var patientCont = patient.New(patientRepo, awsS3Conf)

	var visitRepo = visitRepo.New(db)
	var calendar = calendar.New(visitRepo, srv)
	var visitCont = visit.New(visitRepo, calendar)

	var googleCont = google.New(googleConf, visitRepo)

	var e = echo.New()

	routes.RoutesPath(e, authCont, doctorCont, patientCont, visitCont, googleCont)

	log.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))

}
