package doctor

import (
	"be/configs"
	"be/entities"
	"be/repository/patient"
	"be/utils"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/lithammer/shortuuid"
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

	t.Run("success run Create", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		log.Info(res)
	})

	t.Run("success handle username", func(t *testing.T) {
		var mock = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient"}

		if _, err := patient.New(db).Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Doctor{UserName: "patient2", Email: "clinic@", Password: "clinic"}

		var _, err = r.Create(mock1)
		assert.NotNil(t, err)
	})

	t.Run("success handle email", func(t *testing.T) {
		var mock = entities.Doctor{UserName: "patient5", Email: "patient@", Password: "patient"}

		if _, err := r.Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Doctor{UserName: "patient1", Email: "patient@", Password: "clinic"}

		var _, err = r.Create(mock1)
		assert.NotNil(t, err)
		// log.Info(err)
	})

	t.Run("success handle enum", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "patient1", Email: "doctor@", Password: "clinic", CloseDay: "daksndka"}

		var _, err = r.Create(mock1)
		assert.NotNil(t, err)
		log.Info(err)
	})
}

func TestUpdate(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("success update", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic", LeftCapacity: 5, Capacity: 15}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{Name: "clinic", Capacity: 10}

		res, err = r.Update(res.Doctor_uid, mock1)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// log.Info(res.ClinicName)
	})

	t.Run("handle invalid input capacity", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1n", Email: shortuuid.New(), Password: "clinic", LeftCapacity: 5, Capacity: 15}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{Name: "clinic", Capacity: 5}

		_, err = r.Update(res.Doctor_uid, mock1)
		assert.NotNil(t, err)
		// log.Info(res.ClinicName)
	})

	t.Run("success handle username", func(t *testing.T) {
		var mock = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient"}

		if _, err := patient.New(db).Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Doctor{UserName: "clinic2", Email: shortuuid.New(), Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{Name: "clinic", UserName: "patient2"}

		_, err = r.Update(res.Doctor_uid, mock1)
		assert.NotNil(t, err)
		// log.Info(err)
	})

	t.Run("error input uid", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic3", Email: shortuuid.New(), Password: "clinic"}

		var _, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{Name: "clinic"}

		_, err = r.Update("", mock1)
		assert.NotNil(t, err)
		// log.Info(err)
	})

	t.Run("error input enum", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: shortuuid.New(), Email: shortuuid.New(), Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{Name: "clinic", OpenDay: "ndjandsjka"}

		_, err = r.Update(res.Doctor_uid, mock1)
		assert.NotNil(t, err)
		// log.Info(err)
	})

}

func TestDelete(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("success delete", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		res, err = r.Delete(res.Doctor_uid)
		assert.Nil(t, err)
		assert.Equal(t, true, res.DeletedAt.Valid)
		// log.Info(res.ClinicName)
	})

	t.Run("error input uid", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic3", Email: shortuuid.New(), Password: "clinic"}

		var _, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		_, err = r.Delete("")
		assert.NotNil(t, err)
		// log.Info(err)
	})
}

func TestGetProfile(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})
	db.AutoMigrate(&entities.Visit{})

	t.Run("success get profile", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res1, err1 = r.GetProfile(res.Doctor_uid)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		// log.Info(res1)
	})

	t.Run("input uid", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic2", Email: shortuuid.New(), Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var _, err1 = r.GetProfile(res.Status)
		assert.NotNil(t, err1)
		// log.Info(res1)
	})
}

func TestGetAll(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("success get all", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		if _, err := r.Create(mock1); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{UserName: "doctor1", Email: shortuuid.New(), Password: "clinic"}
		if _, err := r.Create(mock1); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{UserName: "doctor2", Email: shortuuid.New(), Password: "clinic"}
		if _, err := r.Create(mock1); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res1, err1 = r.GetAll()
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		// log.Info(res1)
	})
}
