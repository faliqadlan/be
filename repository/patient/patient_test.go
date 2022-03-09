package patient

import (
	"be/configs"
	"be/entities"
	"be/repository/doctor"
	"be/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("succress run Create", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "anonim", Email: "anonim@", Password: "anonim", Nik: "1"}

		var res, err = r.Create(mock1)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// log.Info(res)
	})

	t.Run("success handle username", func(t *testing.T) {
		var mock = entities.Doctor{UserName: "patient2", Email: "patient@", Password: "patient"}

		if _, err := doctor.New(db).Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Patient{UserName: "patient2", Email: "clinic@", Password: "clinic", Nik: "1"}

		var _, err = r.Create(mock1)
		assert.NotNil(t, err)
	})
}
