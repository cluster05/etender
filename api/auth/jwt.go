package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const jwtSecretKey = "SECRET_KEY"

func CreateJwtToken(userId int, username string) (string, error) {

	expiredAt := time.Now().Add(time.Hour * 24 * 7).Unix()

	claims := jwt.MapClaims{}
	claims["userId"] = userId
	claims["username"] = username
	claims["exp"] = expiredAt

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return to.SignedString([]byte(jwtSecretKey))

}
