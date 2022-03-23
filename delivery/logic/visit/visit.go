package visit

import "errors"

type Logic struct{}

func New() *Logic {
	return &Logic{}
}

func (l *Logic) ValidationRequest(req Req) error {
	if (Req{}) == req {
		return errors.New("data is empty")
	}

	if _, ok := statueses[req.Status]; !ok && req.Status != "" {
		return errors.New("invalid status input")
	}

	return nil
}

var statueses = map[string]int{
	"pending":   0,
	"ready":     1,
	"completed": 2,
	"cancelled": 3,
}
