package auth

import (
	"be/configs"
	"be/entities"
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
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

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

	t.Run("delete patient", func(t *testing.T) {

		var mock1 = entities.Patient{UserName: "patient1", Email: "clinic@", Password: "clinic"}

		res, err := patient.New(db).Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := patient.New(db).Delete(res.Patient_uid); err != nil {
			log.Info(err)
			t.Fatal()
		}

		res1, err := r.Login(mock1.UserName, mock1.Password)
		assert.NotNil(t, err)
		log.Info(err)
		log.Info(res1["type"])
	})

	t.Run("incorrect password patient", func(t *testing.T) {

		var mock1 = entities.Patient{UserName: "patient3", Email: "clinic@", Password: "clinic"}

		_, err := patient.New(db).Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		res1, err := r.Login(mock1.UserName, "")
		assert.NotNil(t, err)
		log.Info(err)
		log.Info(res1["type"])
	})

	t.Run("success run login doctor", func(t *testing.T) {

		var mock2 = entities.Doctor{UserName: "doctor1", Email: "doctor@", Password: "doctor"}

		if _, err := doctor.New(db).Create(mock2); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res1, err1 = r.Login(mock2.UserName, mock2.Password)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		log.Info(res1["type"])
	})

	t.Run("delete doctor", func(t *testing.T) {

		var mock1 = entities.Doctor{UserName: "doctor2", Email: "doctor@", Password: "doctor"}

		res, err := doctor.New(db).Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := doctor.New(db).Delete(res.Doctor_uid); err != nil {
			log.Info(err)
			t.Fatal()
		}

		res1, err := r.Login(mock1.UserName, mock1.Password)
		assert.NotNil(t, err)
		log.Info(err)
		log.Info(res1["type"])
	})

	t.Run("incorrect password doctor", func(t *testing.T) {

		var mock2 = entities.Doctor{UserName: "doctor3", Email: "doctor@", Password: "doctor"}

		if _, err := doctor.New(db).Create(mock2); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res1, err1 = r.Login(mock2.UserName, "")
		assert.NotNil(t, err1)
		log.Info(err1)
		log.Info(res1["type"])
	})

	t.Run("fail run login", func(t *testing.T) {

		var res1, err1 = r.Login("", "")
		assert.NotNil(t, err1)
		log.Info(res1["type"])
		log.Info(err1)
	})
}
