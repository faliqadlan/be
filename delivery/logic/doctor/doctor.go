package doctor

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
		case strings.Contains(err.Error(), "UserName"):
			err = errors.New("invalid userName")
		case strings.Contains(err.Error(), "Email"):
			err = errors.New("invalid email")
		case strings.Contains(err.Error(), "Password"):
			err = errors.New("invalid password")
		case strings.Contains(err.Error(), "Name"):
			err = errors.New("invalid name")
		case strings.Contains(err.Error(), "Address"):
			err = errors.New("invalid address")
		case strings.Contains(err.Error(), "Status"):
			err = errors.New("invalid status")
		case strings.Contains(err.Error(), "OpenDay"):
			err = errors.New("invalid open day")
		case strings.Contains(err.Error(), "CloseDay"):
			err = errors.New("invalid close day")
		case strings.Contains(err.Error(), "Capacity"):
			err = errors.New("invalid capacity ")
		default:
			err = errors.New("invalid input")
		}
		return err
	}
	return nil
}

func (l *Logic) ValidationRequest(req Req) error {

	if err := utils.UserNameValid(req.UserName); err != nil && req.UserName != "" {
		return err
	}

	if err := utils.EmailValid(req.Email); err != nil && req.Email != "" {
		return err
	}

	if err := utils.NameValid(req.Name); err != nil && req.Name != "" {
		return err
	}

	if err := utils.AddressValid(req.Address); err != nil && req.Address != "" {
		return err
	}

	if req.Capacity < 0 && req.Capacity != 0 {
		return errors.New("can't assign capacity below zero")
	}

	if _, ok := statueses[req.Status]; !ok && req.Status != "" {
		return errors.New("invalid status input")
	}

	if _, ok := days[req.OpenDay]; !ok && req.OpenDay != "" {
		return errors.New("invalid open day input")
	}

	if _, ok := days[req.CloseDay]; !ok && req.CloseDay != "" {
		return errors.New("invalid close day input")
	}

	return nil
}

var statueses = map[string]int{
	"available":   0,
	"unAvailable": 1,
}

var days = map[string]int{
	"senin":  0,
	"selasa": 1,
	"rabu":   2,
	"kamis":  3,
	"jumat":  4,
	"sabtu":  5,
	"minggu": 6,
}
