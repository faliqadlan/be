package doctor

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

func TestCreate(t *testing.T) {
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

	t.Run("success run Create", func(t *testing.T) {

		var mock1 = entities.Clinic{UserName: "anonim", Email: "anonim@", Password: "anonim"}

		var res1, err1 = clinic.New(db).Create(mock1)
		if err1 != nil {
			t.Log()
			t.Fatal()
		}

		var mock2 = entities.Doctor{Clinic_uid: res1.Clinic_uid, UserName: "doctor1", Email: "doctor1@", Password: "doctor"}

		var res, err = r.Create(mock2)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("success handle username", func(t *testing.T) {

		var mock = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient"}

		if _, err := patient.New(db).Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}
		var mock1 = entities.Clinic{UserName: "anonim1", Email: "anonim@", Password: "anonim"}

		var res1, err1 = clinic.New(db).Create(mock1)
		if err1 != nil {
			log.Info(err1)
			t.Fatal()
		}

		var mock2 = entities.Doctor{Clinic_uid: res1.Clinic_uid, UserName: "patient2", Email: "doctor1@", Password: "doctor"}

		var _, err = r.Create(mock2)
		assert.NotNil(t, err)
		// log.Info(err)
	})
}
