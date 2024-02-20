package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Antoha2/news/internal/config"
	"github.com/Antoha2/news/internal/lib/models"
	"github.com/Antoha2/news/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Repository interface {
	GetNews(ctx context.Context, pNews *repository.SearchTerms) ([]*repository.RepNews, error)
	AddNews(ctx context.Context, news *repository.RepNews) (*repository.RepNews, error)
	EditNews(ctx context.Context, id int, news *repository.RepNews) (*repository.RepNews, error)
}

type AuthRepository interface {
	UserSaver(ctx context.Context, email string, passHash []byte) (uid int64, err error)
	UserProvider(ctx context.Context, email string) (models.User, error)
}

type ServImpl struct {
	rep *repository.RepImpl
	cfg *config.Config
	log *slog.Logger

	//*repository.RepImpl
}

func NewServ(
	cfg *config.Config,
	log *slog.Logger,
	rep *repository.RepImpl,
) *ServImpl {
	return &ServImpl{
		rep: rep,
		log: log,
		cfg: cfg,
	}
}

type authServImpl struct {
	rep *repository.RepImpl
	cfg *config.Config
	log *slog.Logger
}

// // Login implements transport.AuthService.
// func (*authServImpl) Login(ctx context.Context, email string, password string, appID int) (token string, err error) {
// 	panic("unimplemented")
// }

// // RegisterNewUser implements transport.AuthService.
// func (*authServImpl) RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error) {
// 	panic("unimplemented")
// }

func NewServAuth(
	cfg *config.Config,
	log *slog.Logger,
	rep *repository.RepImpl,
) *authServImpl {
	return &authServImpl{
		cfg: cfg,
		rep: rep,
		log: log,
	}
}

type News struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories []int  `json:"categories"`
}

type SearchTerms struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// type regUser struct {
// 	name  string `json:"name"`
// 	email string `json:"email"`
// }

type LoginRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
	Roles    []string
}

type RegisterRequest struct {
	// FirstName string `json:"firstname"`
	// LastName  string `json:"lastname"`
	// Username  string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
