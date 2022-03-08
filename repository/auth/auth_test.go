package auth

import (
	"be/configs"
	"be/entities"
	"be/repository/clinic"
	"be/repository/doctor"
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
		log.Info(res1["type"])
		// log.Info(res1["data"].(entities.Clinic))
		// var data, _ = json.Marshal(res1["data"])
		// var testdata entities.Clinic
		// json.Unmarshal(data, &testdata)
		// log.Info(testdata)
	})

	t.Run("success run login patient", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "patient", Email: "clinic@", Password: "clinic"}

		var _, err = patient.New(db).Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}
		var res1, err1 = r.Login(mock1.UserName, mock1.Password)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		log.Info(res1["type"])
	})

	t.Run("success run login doctor", func(t *testing.T) {
		var mock1 = entities.Clinic{UserName: "clinic3", Email: "clinic@", Password: "clinic"}

		var res, err = clinic.New(db).Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var mock2 = entities.Doctor{Clinic_uid: res.Clinic_uid, UserName: "doctor1", Email: "doctor@", Password: "doctor"}

		if _, err := doctor.New(db).Create(mock2); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res1, err1 = r.Login(mock2.UserName, mock2.Password)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		log.Info(res1["type"])
		// log.Info(res1["data"].(entities.Clinic))
		// var data, _ = json.Marshal(res1["data"])
		// var testdata entities.Clinic
		// json.Unmarshal(data, &testdata)
		// log.Info(testdata)
	})

	t.Run("fail run login", func(t *testing.T) {

		var res1, err1 = r.Login("", "")
		assert.NotNil(t, err1)
		log.Info(res1["type"])

	})
}
