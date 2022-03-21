package visit

type Visit interface {
	ValidationRequest(req Req) error
}