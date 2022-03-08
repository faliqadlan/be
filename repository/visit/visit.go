package visit

import (
	"be/entities"
	"errors"
	"time"

	"github.com/lithammer/shortuuid"
	"gorm.io/datatypes"
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

func (r *Repo) Create(doctor_uid, patient_uid, date string, req entities.Visit) (entities.Visit, error) {

	var layout = "02-01-2006"

	var dateConv, err = time.Parse(layout, date)
	if err != nil {
		return entities.Visit{}, errors.New("error in time parse date")
	}

	req.Date = datatypes.Date(dateConv)

	var uid string
	for {
		uid = shortuuid.New()
		var find = entities.Visit{}
		var res = r.db.Model(&entities.Visit{}).Where("visit_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	req.Visit_uid = uid
	req.Doctor_uid = doctor_uid
	req.Patient_uid = patient_uid

	if res := r.db.Model(&entities.Visit{}).Create(&req); res.Error != nil {
		return entities.Visit{}, res.Error
	}

	return req, nil

}
