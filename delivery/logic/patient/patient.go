package patient

import (
	"be/utils"
	"errors"
)

type Logic struct{}

func New() *Logic {
	return &Logic{}
}

func (l *Logic) ValidationRequest(req Req) error {

	if err := utils.UserNameValid(req.UserName); err != nil && req.UserName != "" {
		return err
	}

	if err := utils.NikValid(req.Nik); err != nil && req.Nik != "" {
		return err
	}

	if err := utils.NameValid(req.Name); err != nil && req.Name != "" {
		return err
	}

	if _, ok := genders[req.Gender]; !ok && req.Gender != "" {
		return errors.New("invalid gender input")
	}

	if err := utils.AddressValid(req.Address); err != nil && req.Address != "" {
		return err
	}

	if _, ok := statueses[req.Status]; !ok && req.Status != "" {
		return errors.New("invalid status input")
	}

	if _, ok := religions[req.Religion]; !ok && req.Religion != "" {
		return errors.New("invalid religion input")
	}

	return nil
}

var genders = map[string]int{
	"pria":    0,
	"wanita":  1,
	"lainnya": 2,
}

var statueses = map[string]int{
	"belumKawin": 0,
	"kawin":      1,
	"ceraiHidup": 2,
	"ceraiMati":  3,
	"lainnya":    4,
}

var religions = map[string]int{
	"islam":     0,
	"kristen":   1,
	"katolik":   2,
	"budha":     3,
	"hindu":     4,
	"konghuchu": 5,
	"lainnya":   6,
}
