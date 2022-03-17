package doctor

import "be/entities"

type Req struct {
	UserName string `json:"userName" form:"userName"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Name     string `json:"Name" form:"Name" validate:"required"`
	Image    string `json:"image" form:"image"`
	Address  string `json:"address" form:"address" validate:"required"`
	Status   string `json:"status" form:"status" validate:"required"`
	OpenDay  string `json:"openDay" form:"openDay" validate:"required"`
	CloseDay string `json:"closeDay" form:"closeDay" validate:"required"`
	Capacity int    `json:"capacity" form:"capacity" validate:"required"`
}

func (r *Req) ToDoctor() *entities.Doctor {
	return &entities.Doctor{
		UserName:     r.UserName,
		Email:        r.Email,
		Password:     r.Password,
		Name:         r.Name,
		Image:        r.Image,
		Address:      r.Address,
		Status:       r.Status,
		OpenDay:      r.OpenDay,
		CloseDay:     r.CloseDay,
		Capacity:     r.Capacity,
	}
}
