package patient

type Patient interface {
	ValidationRequest(req Req) error
	ValidationStruct(req Req) error
}
