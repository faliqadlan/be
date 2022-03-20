package patient

import (
	"be/entities"
	"errors"
	"time"

	"gorm.io/datatypes"
)

type Req struct {
	UserName   string `json:"userName" form:"userName" validate:"required"`
	Email      string `json:"email" form:"email" validate:"required"`
	Password   string `json:"password" form:"password" validate:"required"`
	Nik        string `json:"nik" form:"nik" validate:"required"`
	Name       string `json:"name" form:"name"  validate:"required"`
	Image      string `json:"image" form:"image"`
	Gender     string `json:"gender" form:"gender"  validate:"required"`
	Address    string `json:"address" form:"address"  validate:"required"`
	PlaceBirth string `json:"placeBirth" form:"placeBirth"  validate:"required"`
	Dob        string `json:"dob" form:"dob"  validate:"required"`
	Job        string `json:"job" form:"job"  validate:"required"`
	Status     string `json:"status" form:"status"  validate:"required"`
	Religion   string `json:"religion" form:"religion"  validate:"required"`
}

func (r *Req) ToPatient() (*entities.Patient, error) {

	var layout = "02-01-2006"

	var dateConv, err = time.Parse(layout, r.Dob)
	if err != nil && r.Dob != "" {
		return &entities.Patient{}, errors.New("invalid date format")
	}
	// log.Info(dateConv)
	return &entities.Patient{
		UserName:   r.UserName,
		Email:      r.Email,
		Password:   r.Password,
		Nik:        r.Nik,
		Name:       r.Name,
		Image:      r.Image,
		Gender:     r.Gender,
		Address:    r.Address,
		PlaceBirth: r.PlaceBirth,
		Dob:        datatypes.Date(dateConv),
		Job:        r.Job,
		Status:     r.Status,
		Religion:   r.Religion,
	}, nil
}
