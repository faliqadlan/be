package patient

import "be/delivery/controllers/patient"

type Patient interface {
	ValidationRegexPatient(req patient.Req) error
}
