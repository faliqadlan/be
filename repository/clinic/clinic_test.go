package clinic

import (
	"be/configs"
	"be/entities"
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
		var mock1 = entities.Clinic{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// log.Info(res)
	})

	t.Run("success handle username", func(t *testing.T) {
		var mock = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient"}

		if _, err := patient.New(db).Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Clinic{UserName: "patient2", Email: "clinic@", Password: "clinic"}

		var _, err = r.Create(mock1)
		assert.NotNil(t, err)
	})
}

func TestUpdate(t *testing.T) {
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

	t.Run("success update", func(t *testing.T) {
		var mock1 = entities.Clinic{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Clinic{ClinicName: "clinic"}

		res, err = r.Update(res.Clinic_uid, mock1)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// log.Info(res.ClinicName)
	})

	t.Run("success handle username", func(t *testing.T) {
		var mock = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient"}

		if _, err := patient.New(db).Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Clinic{UserName: "clinic2", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Clinic{ClinicName: "clinic", UserName: "patient2"}

		_, err = r.Update(res.Clinic_uid, mock1)
		assert.NotNil(t, err)
		// log.Info(err)
	})

	t.Run("error input uid", func(t *testing.T) {
		var mock1 = entities.Clinic{UserName: "clinic3", Email: "clinic@", Password: "clinic"}

		var _, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Clinic{ClinicName: "clinic"}

		_, err = r.Update("", mock1)
		assert.NotNil(t, err)
		// log.Info(err)
	})

}

func TestDelete(t *testing.T) {
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

	t.Run("success delete", func(t *testing.T) {
		var mock1 = entities.Clinic{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Clinic{ClinicName: "clinic"}

		res, err = r.Delete(res.Clinic_uid)
		assert.Nil(t, err)
		assert.Equal(t, true, res.DeletedAt.Valid)
		// log.Info(res.ClinicName)
	})

	t.Run("error input uid", func(t *testing.T) {
		var mock1 = entities.Clinic{UserName: "clinic3", Email: "clinic@", Password: "clinic"}

		var _, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Clinic{ClinicName: "clinic"}

		_, err = r.Delete("")
		assert.NotNil(t, err)
		// log.Info(err)
	})
}
