package configs

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type AppConfig struct {
	PORT        int
	DB          string
	DB_NAME     string
	DB_PORT     int
	DB_HOST     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_LOC      string
	S3_REGION   string
	S3_ID       string
	S3_SECRET   string
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

// type S3Config struct {
// 	S3_REGION string
// 	S3_ID     string
// 	S3_SECRET string
// }

// var s3Config *S3Config

// func initConfigS3() *S3Config {
// 	var exConfig S3Config

// 	exConfig.S3_REGION = os.Getenv("S3_REGION")
// 	exConfig.s
// }
