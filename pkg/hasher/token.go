package hasher

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
)

type tokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateToken(namepass *auth.Namepass, signingKey string, tokenTTLHours int) (string, error) {
	tokenTTL := time.Duration(tokenTTLHours) * time.Hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		namepass.Username,
		namepass.PasswordHash,
	})
	return token.SignedString([]byte(signingKey))
}

func ParseToken(tokenContent, signingKey string) (*auth.Namepass, error) {
	var namepass auth.Namepass
	ok := false
	token, err := jwt.ParseWithClaims(tokenContent, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok = token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method.")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return &namepass, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return &namepass, errors.New("Token claims has wrong type.")
	}
	namepass.Username = claims.Username
	namepass.PasswordHash = claims.Password
	return &namepass, nil
}
