package patient

type Patient interface {
	ValidationRegexPatient(req Req) error
}
