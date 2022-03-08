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
	var tx = r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entities.Doctor{}, err
	}

	var resInit entities.Doctor

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		tx.Rollback()
		return entities.Doctor{}, errors.New("user name already exist")
	}

	if res := tx.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Updates(entities.Doctor{UserName: req.UserName, Email: req.Email, Password: req.Password, Name: req.Name, Image: req.Image, Address: req.Address, Status: req.Status, OpenDay: req.OpenDay, CloseDay: req.CloseDay, Capacity: req.Capacity}); res.Error != nil {
		tx.Rollback()
		return entities.Doctor{}, res.Error
	}

	tx.Commit()

	if res := r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Doctor{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	return resInit, nil
}

func (r *Repo) Delete(doctor_uid string) (entities.Doctor, error) {
	var resInit entities.Doctor

	if res := r.db.Model(&entities.Doctor{}).Where("doctor_uid = ?", doctor_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Doctor{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	return resInit, nil
}

