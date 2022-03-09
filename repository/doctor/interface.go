package doctor

import "be/entities"

type Doctor interface {
	Create(req entities.Doctor) (entities.Doctor, error)
	Update(doctor_uid string, req entities.Doctor) (entities.Doctor, error)
	Delete(doctor_uid string) (entities.Doctor, error)
	GetProfile(doctor_uid string) (ProfileResp, error)
	GetPatients(doctor_uid string) (PatientsResp, error)
	GetDashboard(doctor_uid string) (Dashboard, error)
	GetAll() (All, error)
}
