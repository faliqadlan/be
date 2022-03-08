package visit

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

func TestCreate(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Visit{})

	t.Run("success run create", func(t *testing.T) {
		var mock = entities.Doctor{UserName: "doctor1", Email: "doctor@", Password: "doctor"}
		res, err := doctor.New(db).Create(mock)
		if err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Patient{UserName: "patient1", Email: "patient@", Password: "patient"}
		var res1, err1 = patient.New(db).Create(mock1)
		if err1 != nil {
			t.Log()
			t.Fatal()
		}

		var date = "05-05-2022"

		var mock2 = entities.Visit{Complaint: "sick"}

		var res3, err3 = r.Create(res.Doctor_uid, res1.Patient_uid, date, mock2)
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		log.Info(res3)
	})

}
