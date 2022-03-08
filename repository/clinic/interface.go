package clinic

import "be/entities"

type Clinic interface {
	Create(clinicReq entities.Clinic) (entities.Clinic, error)
	Update(clinic_uid string, up entities.Clinic) (entities.Clinic, error)
	Delete(clinic_uid string) (entities.Clinic, error)
}
