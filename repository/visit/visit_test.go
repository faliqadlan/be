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

	t.Run("error parsing date", func(t *testing.T) {
		var mock = entities.Doctor{UserName: shortuuid.New(), Email: "doctor@", Password: "doctor"}
		res, err := doctor.New(db).Create(mock)
		if err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Patient{UserName: shortuuid.New(), Email: "patient@", Password: "patient"}
		var res1, err1 = patient.New(db).Create(mock1)
		if err1 != nil {
			t.Log()
			t.Fatal()
		}

		var mock2 = entities.Visit{Complaint: "sick"}

		var _, err3 = r.Create(res.Doctor_uid, res1.Patient_uid, "dateNow", mock2)
		assert.NotNil(t, err3)
		// log.Info(err3)
	})

	t.Run("error duplicate", func(t *testing.T) {
		var mock = entities.Doctor{UserName: shortuuid.New(), Email: "doctor@", Password: "doctor"}
		res, err := doctor.New(db).Create(mock)
		if err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Patient{UserName: shortuuid.New(), Email: "patient@", Password: "patient"}
		var res1, err1 = patient.New(db).Create(mock1)
		if err1 != nil {
			t.Log()
			t.Fatal()
		}

		var layDate = "02-01-2006"

		var dateNow = time.Now().Local().Format(layDate)

		var mock2 = entities.Visit{Complaint: "sick"}

		if _, err := r.Create(res.Doctor_uid, res1.Patient_uid, dateNow, mock2); err != nil {
			log.Info(err)
			t.Fatal()
		}

		_, err = r.Create(res.Doctor_uid, res1.Patient_uid, dateNow, entities.Visit{ID: 1, Complaint: "sick new(type)"})
		assert.NotNil(t, err)
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

	t.Run("duplicate entry", func(t *testing.T) {
		var mock = entities.Doctor{UserName: shortuuid.New(), Email: "doctor@", Password: "doctor"}
		res, err := doctor.New(db).Create(mock)
		if err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Patient{UserName: shortuuid.New(), Email: "patient@", Password: "patient"}
		var res1, err1 = patient.New(db).Create(mock1)
		if err1 != nil {
			t.Log()
			t.Fatal()
		}

		var mock2 = entities.Visit{Complaint: "sick", Status: "ready"}

		if _, err := r.CreateVal(res.Doctor_uid, res1.Patient_uid, mock2); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var _, err3 = r.CreateVal(res.Doctor_uid, res1.Patient_uid, entities.Visit{ID: 1, Complaint: "sick new"})
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

		var res3, err3 = r.Update(res2.Visit_uid, entities.Visit{Status: "cancelled", Complaint: "update complaint", MainDiagnose: "update main diagnose", AdditionDiagnose: "update addition_diagnose", Action: "update action", Recipe: "update recipe", BloodPressure: "update blood_pressure", HeartRate: "update heart_rate", O2Saturate: "update o2_saturate", Weight: 100, Height: 100, Bmi: 100})
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		// log.Info(res3)
	})

	t.Run("invalid uid", func(t *testing.T) {
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
		_, err2 := r.CreateVal(res.Doctor_uid, res1.Patient_uid, mock2)
		if err2 != nil {
			t.Log()
			t.Fatal()
		}

		var _, err3 = r.Update(shortuuid.New(), entities.Visit{Status: "cancelled", Complaint: "update complaint", MainDiagnose: "update main diagnose", AdditionDiagnose: "update addition_diagnose", Action: "update action", Recipe: "update recipe", BloodPressure: "update blood_pressure", HeartRate: "update heart_rate", O2Saturate: "update o2_saturate", Weight: 100, Height: 100, Bmi: 100})
		assert.NotNil(t, err3)
		// log.Info(err3)
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

	t.Run("invalid uid", func(t *testing.T) {
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
		_, err2 := r.CreateVal(res.Doctor_uid, res1.Patient_uid, mock2)
		if err2 != nil {
			t.Log()
			t.Fatal()
		}

		var res3, err3 = r.Delete(shortuuid.New())
		assert.NotNil(t, err3)
		assert.Equal(t, false, res3.DeletedAt.Valid)
		// log.Info(res3)
	})

}

func TestGetVisitsList(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Visit{})

	t.Run("success", func(t *testing.T) {
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
		res1, err := r.Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain1"})
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res3, err3 = r.GetVisitList(res1.Visit_uid)
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		// log.Info(res3)
	})

	t.Run("invalid uid", func(t *testing.T) {
		var mock1 = entities.Doctor{UserName: "clinic2", Email: "clinic@", Password: "clinic"}

		var res, err = doctor.New(db).Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var mock2 = entities.Patient{UserName: "patient2", Email: "patient@", Password: "patient", Nik: "1", Name: "name 1"}

		var res2, err2 = patient.New(db).Create(mock2)
		if err2 != nil {
			log.Info(err2)
			t.Fatal()
		}

		var layDate = "02-01-2006"

		var dateNow = time.Now().Local().Format(layDate)
		_, err = r.Create(res.Doctor_uid, res2.Patient_uid, dateNow, entities.Visit{Complaint: "complain1"})
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var _, err3 = r.GetVisitList(shortuuid.New())
		assert.NotNil(t, err3)
		// log.Info(res3)
	})

}

func TestGetVisitsVer1(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Visit{})

	t.Run("case", func(t *testing.T) {
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

		if _, err := r.Create(res.Doctor_uid, res2.Patient_uid, time.Now().AddDate(0, 0, 1).Local().Format(layDate), entities.Visit{Complaint: "complain2"}); err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res3, err3 = r.GetVisitsVer1("", "", "", "", "")
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		// log.Info(res3)
		log.Info(len(res3.Visits))

		res3, err3 = r.GetVisitsVer1("doctor", res.Doctor_uid, "pending", "", "patient")
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		// log.Info(res3)
		log.Info(len(res3.Visits))

		res3, err3 = r.GetVisitsVer1("patient", res2.Patient_uid, "", time.Now().Format(layDate), "doctor")
		assert.Nil(t, err3)
		assert.NotNil(t, res3)
		// log.Info(res3)
		log.Info(len(res3.Visits))

		_, err3 = r.GetVisitsVer1("patient", res2.Patient_uid, "", "date error", "")
		assert.NotNil(t, err3)
	})

}
