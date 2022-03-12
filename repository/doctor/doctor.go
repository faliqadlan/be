package doctor

import (
	"be/entities"
	"be/utils"
	"errors"

	"github.com/lithammer/shortuuid"
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

func (r *Repo) Create(req entities.Doctor) (entities.Doctor, error) {

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		return entities.Doctor{}, errors.New("user name already exist")
	}
	var uid string
	for {
		uid = shortuuid.New()
		var find = entities.Doctor{}
		var res = r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	var err error
	req.Password, err = utils.HashPassword(req.Password)
	if err != nil {
		return entities.Doctor{}, errors.New("error in hash password")
	}
	req.Doctor_uid = uid

	if res := r.db.Model(&entities.Doctor{}).Create(&req); res.Error != nil {
		return entities.Doctor{}, res.Error
	}

	return req, nil
}

func (r *Repo) Update(doctor_uid string, req entities.Doctor) (entities.Doctor, error) {

	var resInit entities.Doctor

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		return entities.Doctor{}, errors.New("user name already exist")
	}

	if res := r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Updates(entities.Doctor{UserName: req.UserName, Email: req.Email, Password: req.Password, Name: req.Name, Image: req.Image, Address: req.Address, Status: req.Status, OpenDay: req.OpenDay, CloseDay: req.CloseDay, Capacity: req.Capacity}); res.Error != nil || res.RowsAffected == 0 {
		return entities.Doctor{}, gorm.ErrRecordNotFound
	}

	// if res := r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
	// 	return entities.Doctor{}, errors.New(gorm.ErrRecordNotFound.Error())
	// }

	return resInit, nil
}

func (r *Repo) Delete(doctor_uid string) (entities.Doctor, error) {
	var resInit entities.Doctor

	if res := r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Doctor{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	return resInit, nil
}

func (r *Repo) GetProfile(doctor_uid string) (ProfileResp, error) {

	var profileResp ProfileResp

	var query = "doctor_uid as Doctor_uid, user_name as UserName, email as Email, name as Name, image as Image, address as Address, status as Status, open_day as OpenDay, close_day as CloseDay, capacity as Capacity, "

	var sub = " capacity - (select count(*) from visits where visits.doctor_uid = ? and visits.status = 'pending' and visits.deleted_at is Null) as LeftCapacity"

	query = query + sub
	// log.Info(query)
	if res := r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Select(query, doctor_uid).Find(&profileResp); res.Error != nil || res.RowsAffected == 0 {
		return ProfileResp{}, gorm.ErrRecordNotFound
	}

	return profileResp, nil
}

func (r *Repo) GetPatients(doctor_uid string) (PatientsResp, error) {
	var patientsResp PatientsResp
	// var patient PatientResp

	if res := r.db.Model(&entities.Visit{}).Where("doctor_uid = ?", doctor_uid).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Group("nik").Select("patients.patient_uid as Patient_uid, patients.name as Name, gender as Gender, nik as Nik, count(*) as TotalVisit").Find(&patientsResp.Patients); res.Error != nil {
		return PatientsResp{}, res.Error
	}

	return patientsResp, nil
}

func (r *Repo) GetDashboard(doctor_uid string) (Dashboard, error) {
	var dashResp Dashboard

	if res := r.db.Model(&entities.Visit{}).Where("doctor_uid = ?", doctor_uid).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Group("nik").Select("count(nik) as TotalPatient").Find(&dashResp.TotalPatient); res.Error != nil {
		return Dashboard{}, res.Error
	}

	if res := r.db.Model(&entities.Visit{}).Where("doctor_uid = ? and date(created_at) = date(curdate())", doctor_uid).Select("count(*) as TotalVisitDay").Find(&dashResp.TotalVisitDay); res.Error != nil {
		return Dashboard{}, res.Error
	}

	if res := r.db.Model(&entities.Visit{}).Where("doctor_uid = ? and date(created_at) = date(curdate()) and status = 'pending'", doctor_uid).Select("count(*) as TotalAppointment").Find(&dashResp.TotalAppointment); res.Error != nil {
		return Dashboard{}, res.Error
	}

	// if res := r.db.Model(&entities.Visit{}).Where("doctor_uid = ?", doctor_uid).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Select("patients.patient_uid as Patient_uid, patients.name as Name, patients.gender as Gender, patients.nik as Nik, visits.status as Status ").Find(&dashResp.Visits); res.Error != nil {
	// 	return Dashboard{}, res.Error
	// }

	return dashResp, nil
}

func (r *Repo) GetAll() (All, error) {
	var all All

	if res := r.db.Model(&entities.Doctor{}).Select("doctor_uid as Doctor_uid, name as Name, image as Image, address as Address, status as Status").Find(&all.Doctors); res.Error != nil {
		return All{}, res.Error
	}

	return all, nil
}
