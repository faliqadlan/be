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

	var res = ad.db.Unscoped().Model(entities.Patient{}).Where("user_name = ?", userName).First(&patient)
	if res.RowsAffected != 0 {
		if match := utils.CheckPasswordHash(password, patient.Password); !match {
			return map[string]interface{}{"type": "patient"}, errors.New("incorrect password")
		}
	}

	if res.RowsAffected != 0 {

		if patient.DeletedAt.Valid {
			return map[string]interface{}{
				"type": "patient",
			}, errors.New("account is deleted")
		}

		return map[string]interface{}{
			"data": patient.Patient_uid,
			"type": "patient",
		}, nil
	}

	// check if doctor

	var doctor entities.Doctor

	res = ad.db.Unscoped().Model(entities.Doctor{}).Where("user_name = ?", userName).First(&doctor)
	if res.RowsAffected != 0 {
		if match := utils.CheckPasswordHash(password, doctor.Password); !match {
			return map[string]interface{}{"type": "doctor"}, errors.New("incorrect password")
		}
	}

	if res.RowsAffected != 0 {

		if doctor.DeletedAt.Valid {
			return map[string]interface{}{
				"type": "doctor",
			}, errors.New("account is deleted")
		}

		return map[string]interface{}{
			"data": doctor.Doctor_uid,
			"type": "doctor",
		}, nil
	}

	return map[string]interface{}{"type": "all"}, gorm.ErrRecordNotFound
}
