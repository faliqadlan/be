package doctor

import "be/entities"

type Doctor interface {
	Create(req entities.Doctor) (entities.Doctor, error)
	Update(doctor_uid string, req entities.Doctor) (entities.Doctor, error)
	Delete(doctor_uid string) (entities.Doctor, error)
	GetProfile(doctor_uid, userName, email string) (ProfileResp, error)
	GetAll() (All, error)
}
