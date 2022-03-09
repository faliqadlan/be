package doctor

import "be/entities"

type Req struct {
	UserName string `json:"userName" form:"userName" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Name     string `json:"Name" form:"Name"`
	Image    string `json:"image" form:"image"`
	Address  string `json:"address" form:"address"`
	Status   string `json:"status" form:"status"`
	OpenDay  string `json:"openDay" form:"openDay"`
	CloseDay string `json:"closeDay" form:"closeDay"`
	Capacity int    `json:"capacity" form:"capacity"`
}

func (r *Req) ToDoctor() *entities.Doctor {
	return &entities.Doctor{
		UserName: r.UserName,
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
		Image:    r.Image,
		Address:  r.Address,
		Status:   r.Status,
		OpenDay:  r.OpenDay,
		CloseDay: r.CloseDay,
		Capacity: r.Capacity,
	}
}

type ResponseFormat struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
