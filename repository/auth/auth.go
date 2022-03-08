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

func (ad *AuthDb) Login(username string, password string) (interface{}, error) {

	// check if clinic
	var clinic entities.Clinic
	var res = ad.db.Model(entities.Clinic{}).Where("user_name = ?", username).First(&clinic)
	if match := utils.CheckPasswordHash(password, clinic.Password); !match {
		return "clinic", errors.New("incorrect password")
	}

	if res.RowsAffected != 0 {
		return map[string]interface{}{
			"data": clinic,
			"type": "clinic",
		}, nil
	}

	// check if patient

	var patient entities.Patient

	res = ad.db.Model(entities.Patient{}).Where("user_name = ?", username).First(&patient)
	if match := utils.CheckPasswordHash(password, patient.Password); !match {
		return "patient", errors.New("incorrect password")
	}
	if res.RowsAffected != 0 {
		return map[string]interface{}{
			"data": patient,
			"type": "patient",
		}, nil
	}

	// check if doctor

	var doctor entities.Doctor

	res = ad.db.Model(entities.Doctor{}).Where("user_name = ?", username).First(&doctor)
	if match := utils.CheckPasswordHash(password, doctor.Password); !match {
		return "doctor", errors.New("incorrect password")
	}
	if res.RowsAffected != 0 {
		return map[string]interface{}{
			"data": doctor,
			"type": "doctor",
		}, nil
	}

	return "all", errors.New(gorm.ErrRecordNotFound.Error())
}
