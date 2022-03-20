package doctor

type Doctor interface {
	ValidationRequest(req Req) error
	ValidationStruct(req Req) error
}
