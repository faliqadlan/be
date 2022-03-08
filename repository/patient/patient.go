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
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? union all ?", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", patientReq.UserName), r.db.Model(&entities.Clinic{}).Select("user_name").Where("user_name = ?", patientReq.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", patientReq.UserName)).Scan(&userNameCheck{})

	// log.Info(checkUserName.RowsAffected)
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
