package visit

import (
	"be/entities"
	"errors"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
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

	var res = r.db.Unscoped().Model(&entities.Visit{}).Where("patient_uid = ?", patient_uid).Scan(&[]entities.Visit{})
	log.Info(res.Error)
	var uid string = patient_uid + "-" + strconv.Itoa(int(res.RowsAffected)+1)
	// log.Info(res.RowsAffected)
	req.Date = datatypes.Date(dateConv)

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


	return req, tx.Commit().Error

}

func (r *Repo) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {

	// check if there's appoinment

	var checkVisit entities.Visit

	if res := r.db.Model(&entities.Visit{}).Where("patient_uid = ? and status = 'pending'", patient_uid).Find(&checkVisit); res.Error != nil || res.RowsAffected != 0 {
		return entities.Visit{}, errors.New("there's another appoinment in pending")
	}

	var res = r.db.Unscoped().Model(&entities.Visit{}).Where("patient_uid = ?", patient_uid).Scan(&[]entities.Visit{})
	var uid string = patient_uid + "-" + strconv.Itoa(int(res.RowsAffected)+1)

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
		return entities.Visit{}, gorm.ErrRecordNotFound
	}
	resInit.ID = 0

	if res := tx.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		// log.Info(res.RowsAffected)
		tx.Rollback()
		return entities.Visit{}, gorm.ErrRecordNotFound
	}
	resInit.DeletedAt = gorm.DeletedAt{}

	if res := tx.Create(&resInit); res.Error != nil {
		tx.Rollback()
		return entities.Visit{}, res.Error
	}

	if req.MainDiagnose != "" {
		req.Status = "completed"
	}
	// log.Info(req)
	// log.Info(req.Event_uid)
	if res := tx.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Updates(entities.Visit{
		Event_uid:        req.Event_uid,
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
	}); res.Error != nil || res.RowsAffected == 0 {
		switch {
		case res.Error == nil:
			tx.Rollback()
			return entities.Visit{}, gorm.ErrRecordNotFound
		default:
			tx.Rollback()
			return entities.Visit{}, res.Error
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
		return entities.Visit{}, gorm.ErrRecordNotFound
	}

	if res := tx.Model(&entities.Visit{}).Where("visit_uid = ?", visit_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		log.Info(res.RowsAffected)
		tx.Rollback()
		return entities.Visit{}, gorm.ErrRecordNotFound
	}

	return resInit, tx.Commit().Error
}

func (r *Repo) GetVisitList(visit_uid string) (VisitCalendar, error) {
	var visits VisitCalendar

	if res := r.db.Model(&entities.Visit{}).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Joins("inner join doctors on visits.doctor_uid = doctors.doctor_uid").Where("visits.visit_uid = ?", visit_uid).Select("doctors.address as Address, complaint as Complaint, date_format(visits.date, '%d-%m-%Y') as Date, doctors.name as DoctorName, patients.name as PatientName, doctors.email as DoctorEmail, patients.email as PatientEmail, event_uid as Event_uid").Last(&visits); res.Error != nil || res.RowsAffected == 0 {
		return VisitCalendar{}, gorm.ErrRecordNotFound
	}

	return visits, nil
}

func (r *Repo) GetVisitsVer1(kind, uid, status, date, grouped string) (Visits, error) {

	switch kind {
	case "":
		kind = "visits.doctor_uid != '" + "uid" + "'"
	case "patient":
		kind = "patients.nik = '" + uid + "'"
	case "doctor":
		kind = "visits.doctor_uid = '" + uid + "'"
	case "visit":
		kind = "visits.visit_uid = '" + uid + "'"
	}

	switch status {
	case "":
		status = "and visits.status != '" + "status" + "'"
	default:
		status = "and visits.status = '" + status + "'"
	}

	switch date {
	case "":
		date = "and date(visits.date) != date('1000-01-01')"
	default:
		var layout = "02-01-2006"
		var dateConv, err = time.Parse(layout, date)
		if err != nil {
			return Visits{}, errors.New("error in time parse date")
		}
		date = "and date(visits.date) = date('" + dateConv.Format("2006-01-02") + "')"
	}
	// log.Info(date)
	switch grouped {
	case "":
		grouped = "id"
	case "patient":
		grouped = "patients.nik"
	case "doctor":
		grouped = "doctors.doctor_uid"
	}

	var visits Visits

	var condition = kind + status + date
	if res := r.db.Model(&entities.Visit{}).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Joins("inner join doctors on visits.doctor_uid = doctors.doctor_uid").Group(grouped).Where(condition).Order("date DESC, visits.updated_at DESC").Select("visit_uid as Visit_uid,  date_format(visits.date, '%d-%m-%Y') as Date, visits.status as Status, complaint as Complaint, main_diagnose as MainDiagnose, addition_diagnose as AdditionDiagnose, action as Action, recipe as Recipe, blood_pressure as BloodPressure, heart_rate as HeartRate, respiratory_rate as RespiratoryRate ,o2_saturate as O2Saturate, weight as Weight, height as Height, bmi as Bmi, visits.doctor_uid as Doctor_uid, doctors.name as DoctorName, doctors.address as DoctorAddress, visits.patient_uid as Patient_uid, patients.name as PatientName, patients.gender as Gender, patients.nik as Nik").Find(&visits.Visits); res.Error != nil {
		log.Info(res.Error)
		log.Info(condition)
		return Visits{}, res.Error
	}

	return visits, nil
}
