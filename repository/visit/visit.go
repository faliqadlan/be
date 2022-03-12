package visit

import (
	"be/entities"
	"errors"
	"time"

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

	if res := r.db.Model(&entities.Visit{}).Create(&req); res.Error != nil {
		return entities.Visit{}, res.Error
	}

	return req, nil

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
	if req.Height != 0 {
		req.Bmi = req.Weight / req.Height
	}

	if res := r.db.Model(&entities.Visit{}).Create(&req); res.Error != nil {
		return entities.Visit{}, res.Error
	}

	return req, nil
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
	if req.Height != 0 {
		req.Bmi = req.Weight / req.Height
	}
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

	return resInit, tx.Commit().Error
}

func (r *Repo) Delete(visit_uid string) (entities.Visit, error) {
	var resInit entities.Visit

	if res := r.db.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Visit{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	return resInit, nil
}

func (r *Repo) GetVisits(doctor_uid, status string) (Visits, error) {
	var visits Visits

	if res := r.db.Model(&entities.Visit{}).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Where("doctor_uid = ? and visits.status = ?", doctor_uid, status).Select("visit_uid as Visit_uid, name as Name, nik as Nik, gender as Gender, date_format(visits.date, '%d-%m-%Y') as Date, recipe as Recipe, main_diagnose as Diagnose").Find(&visits.Visits); res.Error != nil {
		return Visits{}, res.Error
	}

	return visits, nil
}

func (r *Repo) GetVisitList(visit_uid string) (VisitCalendar, error) {
	var visits VisitCalendar

	if res := r.db.Model(&entities.Visit{}).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Joins("inner join doctors on visits.doctor_uid = doctors.doctor_uid").Where("visits.visit_uid = ?", visit_uid).Select("doctors.address as Address, complaint as Complaint, date_format(visits.date, '%d-%m-%Y') as Date, doctors.name as DoctorName, patients.name as PatientName, doctors.email as DoctorEmail, patients.email as PatientEmail").Last(&visits); res.Error != nil || res.RowsAffected == 0 {
		return VisitCalendar{}, gorm.ErrRecordNotFound
	}

	return visits, nil
}
