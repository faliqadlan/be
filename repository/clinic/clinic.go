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

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? union all ?", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", clinicReq.UserName), r.db.Model(&entities.Clinic{}).Select("user_name").Where("user_name = ?", clinicReq.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", clinicReq.UserName)).Scan(&userNameCheck{})
	if checkUserName.RowsAffected != 0 {
		return entities.Clinic{}, errors.New("user name already exist")
	}
	var uid string
	for {
		uid = shortuuid.New()
		var find = entities.Clinic{}
		var res = r.db.Model(&entities.Clinic{}).Where("clinic_uid = ?", uid).Find(&find)
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

func (r *Repo) Update(clinic_uid string, up entities.Clinic) (entities.Clinic, error) {
	var tx = r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entities.Clinic{}, err
	}

	var resInit entities.Clinic

	if res := tx.Model(&entities.Clinic{}).Where("clinic_uid = ?", clinic_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
		tx.Rollback()
		return entities.Clinic{}, errors.New(gorm.ErrRecordNotFound.Error())
	}
	resInit.ID = 0
	if resInit.Clinic_uid != clinic_uid {
		tx.Rollback()
		return entities.Clinic{}, errors.New(gorm.ErrInvalidData.Error())
	}

	if res := tx.Model(&entities.Clinic{}).Where("clinic_uid = ?", clinic_uid).Delete(&resInit); res.RowsAffected == 0 {
		// log.Info(res.RowsAffected)
		tx.Rollback()
		return entities.Clinic{}, errors.New(gorm.ErrRecordNotFound.Error())
	}
	resInit.DeletedAt = gorm.DeletedAt{}

	if res := tx.Create(&resInit); res.Error != nil {
		tx.Rollback()
		return entities.Clinic{}, res.Error
	}

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? union all ?", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", up.UserName), r.db.Model(&entities.Clinic{}).Select("user_name").Where("user_name = ?", up.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", up.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		tx.Rollback()
		return entities.Clinic{}, errors.New("user name already exist")
	}

	if res := tx.Model(&entities.Clinic{}).Where("clinic_uid = ?", clinic_uid).Updates(entities.Clinic{UserName: up.UserName, Email: up.Email, Password: up.Password, DocterName: up.DocterName, ClinicName: up.ClinicName, Address: up.Address, OpenDay: up.OpenDay, CloseDay: up.CloseDay, Capacity: up.Capacity}); res.Error != nil {
		tx.Rollback()
		return entities.Clinic{}, res.Error
	}

	tx.Commit()

	if res := r.db.Model(&entities.Clinic{}).Where("clinic_uid = ?", clinic_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Clinic{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	return resInit, nil
}

func (r *Repo) Delete(clinic_uid string) (entities.Clinic, error) {
	var resInit entities.Clinic

	if res := r.db.Model(&entities.Clinic{}).Where("clinic_uid = ?", clinic_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Clinic{}, errors.New(gorm.ErrRecordNotFound.Error())
	}

	return resInit, nil
}
