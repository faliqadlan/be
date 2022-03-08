package doctor

import "be/entities"

type Doctor interface {
	Create(doctorReq entities.Doctor) (entities.Doctor, error)
}
