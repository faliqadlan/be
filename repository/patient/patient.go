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

func (r *Repo) Create(patientReq entities.Patient) (entities.Patient, error) {

	// check username

	type userNameCheck struct {
		UserNameC string
		UserNameP string
		UserNameD string
	}

	var checkUserName = r.db.Model(&entities.Patient{}).Select("clinics.user_name as UserNameC, patients.user_name as UserNameP, doctors.user_name as UserNameD").Where("clinics.user_name = ? or patients.user_name = ? or doctors.user_name = ?", patientReq.UserName, patientReq.UserName, patientReq.UserName).Joins("left join doctors on 1=1").Joins("left join clinics on 1=1").Find(&userNameCheck{})

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
	patientReq.Password, err = utils.HashPassword(patientReq.Password)
	if err != nil {
		return entities.Patient{}, errors.New("error in hash password")
	}
	patientReq.Patient_uid = uid

	if res := r.db.Model(&entities.Patient{}).Create(&patientReq); res.Error != nil {
		return entities.Patient{}, res.Error
	}

	return patientReq, nil
}
