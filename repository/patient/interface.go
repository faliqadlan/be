package patient

import "be/entities"

type Patient interface {
	Create(patientReq entities.Patient) (entities.Patient, error)
	Update(patient_uid string, req entities.Patient) (entities.Patient, error)
	Delete(patient_uid string) (entities.Patient, error)
	GetProfile(patient_uid string) (Profile, error)
	GetRecords(patient_uid string) (Records, error)
	GetHistories(patient_uid string) (Histories, error)
	GetAppointMent(Patient_uid string) (Apppoinment, error)
}
