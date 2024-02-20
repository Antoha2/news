package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Antoha2/news/internal/auth/internal/repository"
	AuthRepository "github.com/Antoha2/news/internal/auth/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
}

type servImpl struct {
	log *slog.Logger
	AuthService
	rep *AuthRepository.RepAuth
}

type ServRegUser struct {
	FirstName string
	LastName  string
	Username  string
	Password  string
	Email     string
}

func NewServAuth(log *slog.Logger, authRep *repository.RepAuth) *servImpl {
	return &servImpl{
		rep: authRep,
		log: log,
	}
}
