package doctor

import "be/entities"

type Req struct {
	UserName string `json:"userName" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Name     string `json:"Name"`
	Image    string `json:"image"`
	Address  string `json:"address"`
	Status   string `json:"status"`
	OpenDay  string `json:"openDay"`
	CloseDay string `json:"closeDay"`
	Capacity int    `json:"capacity"`
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
