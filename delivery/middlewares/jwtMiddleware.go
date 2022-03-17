package middlewares

import (
	"be/configs"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GenerateToken(uid, kind string) (string, error) {
	if uid == "" {
		return "cannot Generate token", errors.New("uid is empyty")
	}

	codes := jwt.MapClaims{
		"uid":  uid,
		"kind": kind,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
		"auth": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, codes)
	// fmt.Println(token)
	return token.SignedString([]byte(configs.JWT_SECRET))
}

func ExtractTokenUid(e echo.Context) (uid string, kind string) {
	user := e.Get("user").(*jwt.Token) //convert to jwt token from interface
	if user.Valid {
		codes := user.Claims.(jwt.MapClaims)
		uid := codes["uid"].(string)
		kind := codes["kind"].(string)
		return uid, kind
	}
	return "", ""
}
