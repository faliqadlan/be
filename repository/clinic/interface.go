package clinic

import "be/entities"

type Clinic interface {
	Create(clinicReq entities.Clinic) (entities.Clinic, error)
}
