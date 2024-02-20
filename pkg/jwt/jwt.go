package jwt

import (
	"time"

	"github.com/Antoha2/news/internal/config"
	"github.com/Antoha2/news/internal/lib/models"
	"github.com/golang-jwt/jwt/v5"
)

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
