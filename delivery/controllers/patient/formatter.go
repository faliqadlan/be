package patient

import (
	"be/entities"
	"time"

	"gorm.io/datatypes"
)

type Req struct {
	UserName   string `json:"userName" form:"userName" validate:"required"`
	Email      string `json:"email" form:"email" validate:"required"`
	Password   string `json:"password" form:"password" validate:"required"`
	Nik        string `json:"nik" form:"nik" validate:"required"`
	Name       string `json:"name" form:"name"`
	Image      string `json:"image" form:"image"`
	Gender     string `json:"gender" form:"gender"`
	Address    string `json:"address" form:"address"`
	PlaceBirth string `json:"placeBirth" form:"placeBirth"`
	Dob        string `json:"dob" form:"dob"`
	Job        string `json:"job" form:"job"`
	Status     string `json:"status" form:"status"`
	Religion   string `json:"religion" form:"religion"`
}

func (r *Req) ToPatient() *entities.Patient {

	var layout = "02-01-2006"

	var dateConv, err = time.Parse(layout, r.Dob)
	if err != nil {
		return &entities.Patient{}
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
	}
}

type ResponseFormat struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
