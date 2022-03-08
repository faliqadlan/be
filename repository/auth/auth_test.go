package auth

import (
	"be/configs"
	"be/entities"
	"be/repository/clinic"
	"be/repository/patient"
	"be/utils"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Clinic{})
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Clinic{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("success run login clinic", func(t *testing.T) {
		var mock1 = entities.Clinic{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var _, err = clinic.New(db).Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}
		var res1, err1 = r.Login(mock1.UserName, mock1.Password)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		log.Info(res1)
	})

	t.Run("success run login patient", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var _, err = patient.New(db).Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}
		var res1, err1 = r.Login(mock1.UserName, mock1.Password)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		log.Info(res1)
	})

	// t.Run("fail run login", func(t *testing.T) {
	// 	mockLogin := entities.User{Email: "anonim@456", Password: "anonim456"}
	// 	_, err := repo.Login(mockLogin)
	// 	assert.NotNil(t, err)
	// })

}
