package patient

import "be/entities"

type Patient interface {
	Create(patientReq entities.Patient) (entities.Patient, error)
}
