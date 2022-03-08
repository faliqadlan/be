package clinic

import "be/entities"

type Req struct {
	UserName   string `json:"userName" validate:"required"`
	Email      string `json:"email" form:"email" validate:"required"`
	Password   string `json:"password" form:"password" validate:"required"`
	DocterName string `json:"doctorName"`
	ClinicName string `json:"clinicName"`
	Address    string `json:"address"`
	OpenDay    string `json:"openDay"`
	CloseDay   string `json:"closeDay"`
	Capacity   int    `json:"capacity"`
}

func (r *Req) ToClinic() *entities.Clinic {
	return &entities.Clinic{
		UserName:   r.UserName,
		Email:      r.Email,
		Password:   r.Password,
		DocterName: r.DocterName,
		ClinicName: r.ClinicName,
		Address:    r.Address,
		OpenDay:    r.OpenDay,
		CloseDay:   r.CloseDay,
		Capacity:   r.Capacity,
	}
}
