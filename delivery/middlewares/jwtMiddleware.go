package middlewares

import (
	"be/configs"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GenerateToken(uid string) (string, error) {

	codes := jwt.MapClaims{
		"uid":  uid,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
		"auth": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, codes)
	// fmt.Println(token)
	return token.SignedString([]byte(configs.JWT_SECRET))
}

func ExtractTokenUserUid(e echo.Context) string {
	user := e.Get("user").(*jwt.Token) //convert to jwt token from interface
	if user.Valid {
		codes := user.Claims.(jwt.MapClaims)
		uid := codes["uid"].(string)
		return uid
	}
	return ""
}

func ExtractTokenAdmin(e echo.Context) (result string) {
	user := e.Get("user").(*jwt.Token) //convert to jwt token from interface
	if user.Valid {
		codes := user.Claims.(jwt.MapClaims)
		result = codes["email"].(string)
		// result[1] = codes["password"].(string)
		return result
	}
	return ""
}
