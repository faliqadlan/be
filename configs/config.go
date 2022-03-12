package configs

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type AppConfig struct {
	PORT                        int
	DB                          string
	DB_NAME                     string
	DB_PORT                     int
	DB_HOST                     string
	DB_USERNAME                 string
	DB_PASSWORD                 string
	DB_LOC                      string
	S3_REGION                   string
	S3_ID                       string
	S3_SECRET                   string
	CLIENT_ID                   string
	CLIENT_SECRET               string
	PROJECT_ID                  string
	AUTH_URI                    string
	TOKEN_URI                   string
	Auth_provider_x509_cert_url string
	Access_token                string
	Token_type                  string
	Refresh_token               string
}

var synchronizer = &sync.Mutex{}

var appConfig *AppConfig

func initConfig() *AppConfig {
	// if err := godotenv.Load("be-deploy.yaml"); err != nil {
	// 	log.Info(err)
	// }

	exConfig := AppConfig{}

	res, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		log.Warn(err)
	}

	exConfig.PORT = res
	exConfig.DB = os.Getenv("DB")
	exConfig.DB_NAME = os.Getenv("DB_NAME")
	res, err = strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		log.Warn(err)
	}

	exConfig.DB_PORT = res
	exConfig.DB_HOST = os.Getenv("DB_HOST")
	exConfig.DB_USERNAME = os.Getenv("DB_USERNAME")
	exConfig.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	exConfig.DB_LOC = os.Getenv("DB_LOC")
	exConfig.S3_REGION = os.Getenv("S3_REGION")
	exConfig.S3_ID = os.Getenv("S3_ID")
	exConfig.S3_SECRET = os.Getenv("S3_SECRET")
	exConfig.CLIENT_ID = os.Getenv("CLIENT_ID")
	exConfig.CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
	exConfig.PROJECT_ID = os.Getenv("PROJECT_ID")
	exConfig.AUTH_URI = os.Getenv("AUTH_URI")
	exConfig.TOKEN_URI = os.Getenv("TOKEN_URI")
	exConfig.Auth_provider_x509_cert_url = os.Getenv("auth_provider_x509_cert_url")
	exConfig.Access_token = os.Getenv("access_token")
	exConfig.Token_type = os.Getenv("token_type")
	exConfig.Refresh_token = os.Getenv("refresh_token")

	return &exConfig
}

func defaultConfig() *AppConfig {
	defaultConfig := AppConfig{}
	if err := godotenv.Load("local.env"); err != nil {
		log.Info(err)

		defaultConfig = AppConfig{PORT: 8080, DB: "mysql", DB_NAME: "crud_api_test", DB_PORT: 3306, DB_HOST: "localhost", DB_USERNAME: "root", DB_PASSWORD: "root", DB_LOC: "Local"}
		return &defaultConfig
	}

	defaultConfig = AppConfig{}

	res, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		log.Warn(err)
	}

	defaultConfig.PORT = res
	defaultConfig.DB = os.Getenv("DB")
	defaultConfig.DB_NAME = os.Getenv("DB_NAME")
	res, err = strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		log.Warn(err)
	}

	defaultConfig.DB_PORT = res
	defaultConfig.DB_HOST = os.Getenv("DB_HOST")
	defaultConfig.DB_USERNAME = os.Getenv("DB_USERNAME")
	defaultConfig.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	defaultConfig.DB_LOC = os.Getenv("DB_LOC")
	defaultConfig.S3_REGION = os.Getenv("S3_REGION")
	defaultConfig.S3_ID = os.Getenv("S3_ID")
	defaultConfig.S3_SECRET = os.Getenv("S3_SECRET")
	defaultConfig.CLIENT_ID = os.Getenv("CLIENT_ID")
	defaultConfig.CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
	defaultConfig.PROJECT_ID = os.Getenv("PROJECT_ID")
	defaultConfig.AUTH_URI = os.Getenv("AUTH_URI")
	defaultConfig.TOKEN_URI = os.Getenv("TOKEN_URI")
	defaultConfig.Auth_provider_x509_cert_url = os.Getenv("auth_provider_x509_cert_url")
	defaultConfig.Access_token = os.Getenv("access_token")
	defaultConfig.Token_type = os.Getenv("token_type")
	defaultConfig.Refresh_token = os.Getenv("refresh_token")

	return &defaultConfig
}

func GetConfig() *AppConfig {
	synchronizer.Lock()
	defer synchronizer.Unlock()
	appConfig = initConfig()
	// log.Info(appConfig, defaultConfig())
	if appConfig.DB == "" {
		appConfig = defaultConfig()
	}
	return appConfig
}
