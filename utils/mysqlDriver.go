package utils

import (
	"be/configs"
	"be/entities"
	"fmt"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&entities.Patient{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Visit{})
}

func InitDB(config *configs.AppConfig) *gorm.DB {

	connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=%v",
		config.DB_USERNAME,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
		config.DB_LOC,
	)
	fmt.Println(connectionString)
	DB, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Info("error in connect database ", err)
		panic(err)
	}

	autoMigrate(DB)
	return DB
}
