package doctor

import (
	"be/configs"
	"be/entities"
	"be/repository/patient"
	"be/repository/visit"
	"be/utils"
	"testing"
	"time"

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
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("success run Create", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

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

		var mock1 = entities.Doctor{UserName: "patient2", Email: "clinic@", Password: "clinic"}

		var _, err = r.Create(mock1)
		assert.NotNil(t, err)
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
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{Name: "clinic"}

		res, err = r.Update(res.Doctor_uid, mock1)
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

		var mock1 = entities.Doctor{UserName: "clinic2", Email: "clinic@", Password: "clinic"}

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
		var mock1 = entities.Doctor{UserName: "clinic3", Email: "clinic@", Password: "clinic"}

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
		var mock1 = entities.Doctor{UserName: "clinic3", Email: "clinic@", Password: "clinic"}

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
}

func TestGetPatients(t *testing.T) {
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

		var mock2 = entities.Patient{UserName: "patient1", Email: "patient@", Password: "patient", Nik: "1", Name: "name 1"}

		var res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, "05-05-2022", entities.Visit{Complaint: "complain1"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, "05-05-2022", entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock2 = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient", Nik: "1", Name: "name2"}

		res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, "05-05-2022", entities.Visit{Complaint: "complain1"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, "05-05-2022", entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock2 = entities.Patient{UserName: "patient3", Email: "patient@", Password: "patient", Nik: "3", Name: "name3"}

		res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, "05-05-2022", entities.Visit{Complaint: "complain1"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, "05-05-2022", entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res1, err1 = r.GetPatients(res.Doctor_uid)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		// log.Info(res1)
	})
}

func TestGetDashboard(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})
	db.AutoMigrate(&entities.Visit{})

	t.Run("success get Dashboard", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var mock2 = entities.Patient{UserName: "patient1", Email: "patient@", Password: "patient", Nik: "1", Name: "name 1"}

		var res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		var layDate = "02-01-2006"

		var dateNow = time.Now().Local().Format(layDate)

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain1"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock2 = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient", Nik: "1", Name: "name2"}

		res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain1", Status: "ready"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock2 = entities.Patient{UserName: "patient3", Email: "patient@", Password: "patient", Nik: "3", Name: "name3"}

		res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain1"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := visit.New(db).Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res1, err1 = r.GetDashboard(res.Doctor_uid)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		log.Info(res1)
		// res3, err2 := r.GetProfile(res.Doctor_uid)
		// log.Info(res3)
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

		mock1 = entities.Doctor{UserName: "doctor1", Email: "clinic@", Password: "clinic"}
		if _, err := r.Create(mock1); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Doctor{UserName: "doctor2", Email: "clinic@", Password: "clinic"}
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
