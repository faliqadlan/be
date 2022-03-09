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
