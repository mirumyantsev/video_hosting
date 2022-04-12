package vh

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	signingKey = "jD2@hSw2eGe7#HkU7fH@8kLe0#6GeD"
	tokenTTL   = 24 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateToken(username, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		username,
		password,
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(accessToken string) (NamePass, error) {
	var namepass NamePass
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return namepass, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return namepass, errors.New("token claims are not of type *tokenClaims")
	}

	namepass.Username = claims.Username
	namepass.PasswordHash = claims.Password
	return namepass, nil
}
