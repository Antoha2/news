package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	models "github.com/Antoha2/news/internal/auth/internal/domain/models"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type AuthRepository interface {
	UserSaver(ctx context.Context, email string, passHash []byte) (uid int64, err error)
	UserProvider(ctx context.Context, email string) (models.User, error)
	App(ctx context.Context, appID int) (models.App, error)
}

type RepAuth struct {
	log *slog.Logger
	DB  *sqlx.DB
	AuthRepository
	TokenTTL time.Duration
}

func NewRepAuth(log *slog.Logger, dbx *sqlx.DB, tokenTTL time.Duration) *RepAuth {
	return &RepAuth{
		log:      log,
		DB:       dbx,
		TokenTTL: tokenTTL,
	}
}
