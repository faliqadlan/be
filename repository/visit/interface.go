package visit

import "be/entities"

type Visit interface {
	CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error)
	Update(visit_uid string, req entities.Visit) (entities.Visit, error)
	Delete(visit_uid string) (entities.Visit, error)
	GetVisitsVer1(kind, uid, status, sign_status string) (Visits, error)
	GetVisitList(email, status string) (VisitCalendar, error)
}
