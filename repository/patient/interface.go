package patient

import "be/entities"

type Patient interface {
	Create(patientReq entities.Patient) (entities.Patient, error)
	Update(patient_uid string, req entities.Patient) (entities.Patient, error)
	Delete(patient_uid string) (entities.Patient, error)
	GetProfile(patient_uid, userName, email string) (Profile, error)
	GetAll() (All, error)
}
