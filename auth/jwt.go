package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secure-secret") // Use environment variable in production

func GenerateJWT(userID string) (string, int64, error) {
	expTime := time.Now().Add(time.Hour * 12).Unix()
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    expTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(jwtSecret)
	return signedString, expTime, err
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}
