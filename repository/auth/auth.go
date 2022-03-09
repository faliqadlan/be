package auth

import (
	"be/entities"
	"be/utils"
	"errors"

	"gorm.io/gorm"
)

type AuthDb struct {
	db *gorm.DB
}

func New(db *gorm.DB) *AuthDb {
	return &AuthDb{
		db: db,
	}
}

func (ad *AuthDb) Login(userName string, password string) (map[string]interface{}, error) {

	// check if patient

	var patient entities.Patient

	var res = ad.db.Model(entities.Patient{}).Where("user_name = ?", userName).First(&patient)
	if res.RowsAffected != 0 {
		if match := utils.CheckPasswordHash(password, patient.Password); !match {
			return map[string]interface{}{"type": "patient"}, errors.New("incorrect password")
		}
	}
	if res.RowsAffected != 0 {
		return map[string]interface{}{
			"data": patient.Patient_uid,
			"type": "patient",
		}, nil
	}

	// check if doctor

	var doctor entities.Doctor

	res = ad.db.Model(entities.Doctor{}).Where("user_name = ?", userName).First(&doctor)
	if res.RowsAffected != 0 {
		if match := utils.CheckPasswordHash(password, doctor.Password); !match {
			return map[string]interface{}{"type": "doctor"}, errors.New("incorrect password")
		}
	}
	if res.RowsAffected != 0 {
		return map[string]interface{}{
			"data": doctor.Doctor_uid,
			"type": "doctor",
		}, nil
	}

	return map[string]interface{}{"type": "all"}, errors.New(gorm.ErrRecordNotFound.Error())
}
