package patient

import (
	"be/delivery/controllers/patient"
	"be/utils"
)

type Logic struct{}

func (l *Logic) ValidationRegexPatient(req patient.Req) error {
	if err := utils.UserNameValid(req.UserName); err != nil && req.UserName != "" {
		return err
	}
	if err := utils.NameValid(req.Name); err != nil && req.Name != "" {
		return err
	}
	if err := utils.AddressValid(req.Address); err != nil && req.Address != "" {
		return err
	}
	if err := utils.NikValid(req.Nik); err != nil && req.Nik != "" {
		return err
	}
	return nil
}
