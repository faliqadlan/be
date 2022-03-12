package visit

import (
	"be/entities"
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/lithammer/shortuuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(doctor_uid, patient_uid, date string, req entities.Visit) (entities.Visit, error) {

	var layout = "02-01-2006"

	var dateConv, err = time.Parse(layout, date)
	if err != nil {
		return entities.Visit{}, errors.New("error in time parse date")
	}

	req.Date = datatypes.Date(dateConv)

	var uid string
	for {
		uid = shortuuid.New()
		var find = entities.Visit{}
		var res = r.db.Model(&entities.Visit{}).Where("visit_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	req.Visit_uid = uid
	req.Doctor_uid = doctor_uid
	req.Patient_uid = patient_uid

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entities.Visit{}, err
	}

	if res := tx.Model(&entities.Visit{}).Create(&req); res.Error != nil {
		tx.Rollback()
		return entities.Visit{}, res.Error
	}

	if req.Status == "pending" {

		var doctorInit entities.Doctor

		if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Find(&doctorInit); res.Error != nil {
			tx.Rollback()
			return entities.Visit{}, res.Error
		}

		if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Update("left_capacity", doctorInit.LeftCapacity-1); res.Error != nil || res.RowsAffected == 0 {
			tx.Rollback()
			return entities.Visit{}, gorm.ErrRecordNotFound
		}

	}

	return req, tx.Commit().Error

}

func (r *Repo) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {

	// check if there's appoinment

	var checkVisit entities.Visit

	if res := r.db.Model(&entities.Visit{}).Where("patient_uid = ? and status = 'pending'", patient_uid).Find(&checkVisit); res.Error != nil || res.RowsAffected != 0 {
		return entities.Visit{}, errors.New("there's another appoinment in pending")
	}

	var uid string
	for {
		uid = shortuuid.New()
		var find = entities.Visit{}
		var res = r.db.Model(&entities.Visit{}).Where("visit_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	req.Visit_uid = uid
	req.Doctor_uid = doctor_uid
	req.Patient_uid = patient_uid
	// if req.Height != 0 {
	// 	req.Bmi = req.Weight / req.Height
	// }

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entities.Visit{}, err
	}

	if res := tx.Model(&entities.Visit{}).Create(&req); res.Error != nil {
		tx.Rollback()
		return entities.Visit{}, res.Error
	}

	if req.Status == "pending" {

		var doctorInit entities.Doctor

		if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Find(&doctorInit); res.Error != nil {
			tx.Rollback()
			return entities.Visit{}, res.Error
		}

		if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Update("left_capacity", doctorInit.LeftCapacity-1); res.Error != nil || res.RowsAffected == 0 {
			tx.Rollback()
			return entities.Visit{}, gorm.ErrRecordNotFound
		}

	}

	return req, tx.Commit().Error
}

func (r *Repo) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entities.Visit{}, err
	}

	var resInit entities.Visit

	if res := tx.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
		tx.Rollback()
		return entities.Visit{}, errors.New(gorm.ErrRecordNotFound.Error())
	}
	resInit.ID = 0

	if res := tx.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Delete(&resInit); res.RowsAffected == 0 {
		// log.Info(res.RowsAffected)
		tx.Rollback()
		return entities.Visit{}, errors.New(gorm.ErrRecordNotFound.Error())
	}
	resInit.DeletedAt = gorm.DeletedAt{}

	if res := tx.Create(&resInit); res.Error != nil {
		tx.Rollback()
		return entities.Visit{}, res.Error
	}
	// if req.Height != 0 {
	// 	req.Bmi = req.Weight / req.Height
	// }
	if res := tx.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Updates(entities.Visit{
		Status:           req.Status,
		Complaint:        req.Complaint,
		MainDiagnose:     req.MainDiagnose,
		AdditionDiagnose: req.AdditionDiagnose,
		Action:           req.Action,
		Recipe:           req.Recipe,
		BloodPressure:    req.BloodPressure,
		HeartRate:        req.HeartRate,
		O2Saturate:       req.O2Saturate,
		Weight:           req.Weight,
		Height:           req.Height,
		Bmi:              req.Bmi,
	}); res.Error != nil {
		tx.Rollback()
		return entities.Visit{}, res.Error
	}

	if req.Status == "ready" || req.Status == "cancelled" {

		var doctorInit entities.Doctor

		if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", resInit.Doctor_uid).Find(&doctorInit); res.Error != nil {
			tx.Rollback()
			return entities.Visit{}, res.Error
		}

		if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", resInit.Doctor_uid).Update("left_capacity", doctorInit.LeftCapacity+1); res.Error != nil || res.RowsAffected == 0 {
			tx.Rollback()
			return entities.Visit{}, gorm.ErrRecordNotFound
		}

	}

	return resInit, tx.Commit().Error
}

func (r *Repo) Delete(visit_uid string) (entities.Visit, error) {

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entities.Visit{}, err
	}

	var resInit entities.Visit

	if res := tx.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
		tx.Rollback()
		return entities.Visit{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	if res := tx.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		log.Info(res.RowsAffected)
		tx.Rollback()
		return entities.Visit{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	var doctorInit entities.Doctor

	if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", resInit.Doctor_uid).Find(&doctorInit); res.Error != nil {
		tx.Rollback()
		return entities.Visit{}, res.Error
	}

	if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", resInit.Doctor_uid).Update("left_capacity", doctorInit.LeftCapacity+1); res.Error != nil || res.RowsAffected == 0 {
		log.Info(res.RowsAffected)
		tx.Rollback()
		return entities.Visit{}, gorm.ErrRecordNotFound
	}

	return resInit, tx.Commit().Error
}

func (r *Repo) GetVisitList(email, status string) (VisitCalendar, error) {
	var visits VisitCalendar

	if res := r.db.Model(&entities.Visit{}).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Joins("inner join doctors on visits.doctor_uid = doctors.doctor_uid").Where("patients.email = ? and visits.status = ?", email, status).Select("doctors.address as Address, complaint as Complaint, date_format(visits.date, '%d-%m-%Y') as Date, doctors.name as DoctorName, patients.name as PatientName, doctors.email as DoctorEmail").Last(&visits); res.Error != nil || res.RowsAffected == 0 {
		return VisitCalendar{}, gorm.ErrRecordNotFound
	}

	return visits, nil
}

func (r *Repo) GetVisitsVer1(kind, uid, status, signStatus string) (Visits, error) {

	switch kind {
	case "patient":
		kind = "patients.nik"
	case "doctor":
		kind = "visits.doctor_uid"
	}

	switch signStatus {
	case "equal":
		signStatus = " = "
	case "notequal":
		signStatus = " != "
	}

	var visits Visits

	var condition = kind + " = '" + uid + "' and visits.status" + signStatus + "'" + status + "'"
	if res := r.db.Model(&entities.Visit{}).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Joins("inner join doctors on visits.doctor_uid = doctors.doctor_uid").Where(condition).Select("visit_uid as Visit_uid,  date_format(visits.date, '%d-%m-%Y') as Date, visits.status as Status, complaint as Complaint, main_diagnose as MainDiagnose, addition_diagnose as AdditionDiagnose, action as Action, recipe as Recipe, blood_pressure as BloodPressure, heart_rate as HeartRate, o2_saturate as O2Saturate, weight as Weight, height as Height, bmi as Bmi, visits.doctor_uid as Doctor_uid, doctors.name as DoctorName, doctors.address as DoctorAddress, visits.patient_uid as Patient_uid, patients.name as PatientName, patients.gender as Gender, patients.nik as Nik").Find(&visits.Visits); res.Error != nil {
		return Visits{}, res.Error
	}

	return visits, nil
}
