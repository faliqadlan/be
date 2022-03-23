package patient

import (
	"be/utils"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type Logic struct{}

func New() *Logic {
	return &Logic{}
}

func (l *Logic) ValidationStruct(req Req) error {
	var v = validator.New()
	if err := v.Struct(req); err != nil {
		log.Warn(err)
		switch {
		case strings.Contains(err.Error(), "Nik"):
			err = errors.New("invalid nik")
		case strings.Contains(err.Error(), "Name"):
			err = errors.New("invalid name")
		case strings.Contains(err.Error(), "Gender"):
			err = errors.New("invalid gender")
		case strings.Contains(err.Error(), "Address"):
			err = errors.New("invalid address")
		case strings.Contains(err.Error(), "PlaceBirth"):
			err = errors.New("invalid place birth")
		case strings.Contains(err.Error(), "Dob"):
			err = errors.New("invalid date of birth")
		case strings.Contains(err.Error(), "Job"):
			err = errors.New("invalid job")
		case strings.Contains(err.Error(), "Status"):
			err = errors.New("invalid status")
		case strings.Contains(err.Error(), "Religion"):
			err = errors.New("invalid religion")
		default:
			err = errors.New("invalid input")
		}
		return err
	}
	return nil
}

func (l *Logic) ValidationRequest(req Req) error {

	if (Req{}) == req {
		return errors.New("data is empty")
	}

	if err := utils.UserNameValid(req.UserName); err != nil && req.UserName != "" {
		return err
	}

	if err := utils.EmailValid(req.Email); err != nil && req.Email != "" {
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
	"protestan": 7,
}
