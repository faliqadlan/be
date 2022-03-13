package utils

import (
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Warn(err)
	}
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	log.Warn(err)
	return err == nil
}

// func DecodePassword(password string) (string, error) {
// 	var pass, err = bcrypt.
// }
