package visit

import "be/entities"

type Visit interface {
	CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error)
	Update(visit_uid string, req entities.Visit) (entities.Visit, error)
}
