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

func (r *Repo) Create(doctorReq entities.Doctor) (entities.Doctor, error) {

	// check username

	type userNameCheck struct {
		UserNameC string
		UserNameP string
		UserNameD string
	}

	var checkUserName = r.db.Model(&entities.Doctor{}).Select("clinics.user_name as UserNameC, patients.user_name as UserNameP, doctors.user_name as UserNameD").Where("clinics.user_name = ? or patients.user_name = ? or doctors.user_name = ?", doctorReq.UserName, doctorReq.UserName, doctorReq.UserName).Joins("left join patients on 1=1 ").Joins("left join clinics on 1=1").Find(&userNameCheck{})

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
	doctorReq.Password, err = utils.HashPassword(doctorReq.Password)
	if err != nil {
		return entities.Doctor{}, errors.New("error in hash password")
	}
	doctorReq.Doctor_uid = uid

	if res := r.db.Model(&entities.Doctor{}).Create(&doctorReq); res.Error != nil {
		return entities.Doctor{}, res.Error
	}

	return doctorReq, nil
}
