package clinic

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

func (r *Repo) Create(clinicReq entities.Clinic) (entities.Clinic, error) {
	var clinicInit = entities.Clinic{}

	// check email

	var checkEmail = r.db.Model(&entities.Clinic{}).Where("email = ?", clinicReq.Email).Find(&clinicInit)

	if checkEmail.RowsAffected != 0 {
		return entities.Clinic{}, errors.New("email already exist")
	}

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Model(&entities.Clinic{}).Model(&entities.Patient{}).Model(&entities.Doctor{}).Where("clinics.user_name = ? or patients.user_name = ? or doctors.user_name = ?", clinicReq.UserName, clinicReq.UserName, clinicReq.UserName).Find(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		return entities.Clinic{}, errors.New("email already exist")
	}
	var uid string
	for {
		uid = shortuuid.New()
		var find = entities.Clinic{}
		var res = r.db.Model(&entities.Clinic{}).Where("patient_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	var err error
	clinicReq.Password, err = utils.HashPassword(clinicReq.Password)
	if err != nil {
		return entities.Clinic{}, errors.New("error in hash password")
	}
	clinicReq.Clinic_uid = uid

	if res := r.db.Model(&entities.Clinic{}).Create(&clinicReq); res.Error != nil {
		return entities.Clinic{}, res.Error
	}

	return clinicReq, nil
}
