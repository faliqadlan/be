package patient

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

func (r *Repo) Create(req entities.Patient) (entities.Patient, error) {

	if req.Nik == "" && req.UserName == "" {
		return entities.Patient{}, errors.New("nik must filled")
	}

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		return entities.Patient{}, errors.New("user name already exist")
	}
	var uid string
	for {
		uid = shortuuid.New()
		var find = entities.Patient{}
		var res = r.db.Model(&entities.Patient{}).Where("patient_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	var err error
	req.Password, err = utils.HashPassword(req.Password)
	if err != nil {
		return entities.Patient{}, errors.New("error in hash password")
	}
	req.Patient_uid = uid

	if res := r.db.Model(&entities.Patient{}).Create(&req); res.Error != nil {
		return entities.Patient{}, res.Error
	}

	return req, nil
}

func (r *Repo) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {

	var resInit entities.Patient

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		return entities.Patient{}, errors.New("user name already exist")
	}

	if res := r.db.Model(&entities.Patient{}).Where("patient_uid = ?", patient_uid).Updates(entities.Patient{UserName: req.UserName, Email: req.Email, Password: req.Password, Name: req.Name, Image: req.Image, Gender: req.Gender, Address: req.Address, Status: req.Status, Religion: req.Religion}); res.Error != nil || res.RowsAffected == 0 {
		return entities.Patient{}, gorm.ErrRecordNotFound
	}

	// if res := r.db.Model(&entities.Doctor{}).Where("pa = ?", doctor_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
	// 	return entities.Doctor{}, errors.New(gorm.ErrRecordNotFound.Error())
	// }

	return resInit, nil
}

func (r *Repo) Delete(patient_uid string) (entities.Patient, error) {
	var resInit entities.Patient

	if res := r.db.Model(&entities.Patient{}).Where("patient_uid = ?", patient_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Patient{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	return resInit, nil
}

func (r *Repo) GetProfile(patient_uid string) (Profile, error) {
	var profileResp Profile

	if res := r.db.Model(&entities.Patient{}).Where("patient_uid = ?", patient_uid).Find(&profileResp); res.Error != nil || res.RowsAffected == 0 {
		return Profile{}, gorm.ErrRecordNotFound
	}

	return profileResp, nil
}

func (r *Repo) GetRecords(patient_uid string) (Records, error) {
	var resInit entities.Patient

	if res := r.db.Model(&entities.Patient{}).Where("patient_uid = ?", patient_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return Records{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	var records Records

	if res := r.db.Model(&entities.Visit{}).Joins("inner join patients on visits.patient_uid = patients.patient_uid").Where("patients.nik = ?", resInit.Nik).Find(&records.Records); res.Error != nil {
		return Records{}, res.Error
	}

	return records, nil
}
