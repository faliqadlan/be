package utils

import (
	"errors"
	"net/mail"
	"regexp"
)

func UserNameRegex(s string) (string, error) {

	if regexp.MustCompile(`^[[:alnum:]]+$`).MatchString(s) {
		containNumber, err := regexp.Compile(`[[:digit:]]+`)
		if err != nil {
			return "error", err
		}
		containAphabet, err := regexp.Compile(`[[:alpha:]]+`)
		if err != nil {
			return "error", err
		}
		if containNumber.MatchString(s) && containAphabet.MatchString(s) {
			return "success", nil
		}
	}

	return "fail", nil
}

func NameRegex(s string) bool {

	return regexp.MustCompile(`^[[:alpha:]\s]+$`).MatchString(s)
}

func AddressRegex(s string) bool {

	return regexp.MustCompile(`^[a-zA-Z0-9/,.\s]+$`).MatchString(s)
}

func DigitRegex(s string) bool {

	return regexp.MustCompile(`^[[:digit:]]+$`).MatchString(s)
}

func NameValid(s string) error {
	if len(s) < 4 || len(s) > 255 {
		return errors.New("invalid length name")
	} else {
		if !NameRegex(s) {
			return errors.New("invalid name format")
		}
	}
	return nil
}

func AddressValid(s string) error {
	if len(s) < 15 || len(s) > 255 {
		return errors.New("invalid length address")
	} else {
		if !AddressRegex(s) {
			return errors.New("invalid address format")
		}
	}
	return nil
}

func UserNameValid(s string) error {
	if len(s) < 5 || len(s) > 255 {
		return errors.New("invalid length user name")
	} else {
		if s, err := UserNameRegex(s); s != "success" {
			if err != nil {
				return err
			}
			return errors.New("invalid user name format")
		}
	}
	return nil
}

func NikValid(s string) error {
	if len(s) != 16 {
		return errors.New("invalid length nik")
	} else {
		if !DigitRegex(s) {
			return errors.New("invalid nik format")
		}
	}
	return nil
}

func EmailValid(s string) error  {
	_, err := mail.ParseAddress(s)
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}