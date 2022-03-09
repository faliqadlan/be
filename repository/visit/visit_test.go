package visit

import (
	"be/configs"
	"be/entities"
	"be/repository/doctor"
	"be/repository/patient"
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

		var layDate = "02-01-2006"

		var dateNow = time.Now().Local().Format(layDate)

		var mock2 = entities.Visit{Complaint: "sick"}

		var res3, err3 = r.Create(res.Doctor_uid, res1.Patient_uid, dateNow, mock2)
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		// log.Info(res3)
	})

}

func TestCreateVal(t *testing.T) {
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

		var mock2 = entities.Visit{Complaint: "sick"}

		var res3, err3 = r.CreateVal(res.Doctor_uid, res1.Patient_uid, mock2)
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		// log.Info(res3)
	})

	t.Run("success handle pending", func(t *testing.T) {
		var mock = entities.Doctor{UserName: "doctor2", Email: "doctor@", Password: "doctor"}
		res, err := doctor.New(db).Create(mock)
		if err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient"}
		var res1, err1 = patient.New(db).Create(mock1)
		if err1 != nil {
			t.Log()
			t.Fatal()
		}

		var mock2 = entities.Visit{Complaint: "sick"}

		if _, err := r.CreateVal(res.Doctor_uid, res1.Patient_uid, mock2); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var _, err3 = r.CreateVal(res.Doctor_uid, res1.Patient_uid, mock2)
		assert.NotNil(t, err3)
		// log.Info(err3)
	})

}

func TestUpdate(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Visit{})

	t.Run("success run update", func(t *testing.T) {
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

		var mock2 = entities.Visit{Complaint: "sick"}
		res2, err2 := r.CreateVal(res.Doctor_uid, res1.Patient_uid, mock2)
		if err2 != nil {
			t.Log()
			t.Fatal()
		}

		var res3, err3 = r.Update(res2.Visit_uid, entities.Visit{Complaint: "very sick"})
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		// log.Info(res3)
	})

}

func TestDelete(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Visit{})

	t.Run("success run delete", func(t *testing.T) {
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

		var mock2 = entities.Visit{Complaint: "sick"}
		res2, err2 := r.CreateVal(res.Doctor_uid, res1.Patient_uid, mock2)
		if err2 != nil {
			t.Log()
			t.Fatal()
		}

		var res3, err3 = r.Delete(res2.Visit_uid)
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		assert.Equal(t, true, res3.DeletedAt.Valid)
		// log.Info(res3)
	})

}

func TestGetVisits(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Visit{})

	t.Run("success run update", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic1", Email: "clinic@", Password: "clinic"}

		var res, err = doctor.New(db).Create(mock1)
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

		if _, err := r.Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain1"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := r.Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain2", Status: "ready"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock2 = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient", Nik: "1", Name: "name2"}

		res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		if _, err := r.Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain1", Status: "ready"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := r.Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock2 = entities.Patient{UserName: "patient3", Email: "patient@", Password: "patient", Nik: "3", Name: "name3"}

		res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		if _, err := r.Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain1", Status: "cancelled"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		if _, err := r.Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var count int
		var res3, err3 = r.GetVisits(res.Doctor_uid, "pending")
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		count += len(res3.Visits)
		// log.Info(res3)

		res3, err3 = r.GetVisits(res.Doctor_uid, "ready")
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		count += len(res3.Visits)
		// log.Info(res3)

		res3, err3 = r.GetVisits(res.Doctor_uid, "cancelled")
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		count += len(res3.Visits)
		assert.Equal(t, 6, count)
		// log.Info(res3)
		// log.Info(count)
	})

}
