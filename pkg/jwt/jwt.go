package jwt

import (
	"errors"
	"time"

	"github.com/Antoha2/news/internal/config"
	"github.com/Antoha2/news/internal/lib/models"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func NewToken(cfg *config.Config, user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(cfg.TokenTTL).Unix()

	tokenString, err := token.SignedString([]byte(cfg.TokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(cfg *config.Config, accesToken string) (string, error) {

	token, err := jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signature method")
		}
		return []byte(cfg.TokenSecret), nil
	})
	if err != nil {
		return "", nil
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims not of type *tokenClaims")
	}
	return claims.Email, nil
}
