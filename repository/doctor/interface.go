package doctor

import "be/entities"

type Clinic interface {
	Create(req entities.Doctor) (entities.Doctor, error)
	Update(doctor_uid string, req entities.Doctor) (entities.Doctor, error)
	Delete(doctor_uid string) (entities.Doctor, error)
}
