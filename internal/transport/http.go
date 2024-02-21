package transport

import (
	"context"
	"log/slog"
	"net/http"

	config "github.com/Antoha2/news/internal/config"
	"github.com/Antoha2/news/internal/service"
)

type Service interface {
	GetNews(ctx context.Context, param *service.SearchTerms) ([]*service.News, error)
	EditNews(ctx context.Context, id int, news *service.News) (*service.News, error)
	AddNews(ctx context.Context, news *service.News) (*service.News, error)
}

type AuthService interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	//ParseToken(token string) (int, error)
}

type apiImpl struct {
	authService AuthService
	service     Service
	cfg         *config.Config
	log         *slog.Logger
	server      *http.Server
}

// NewAPI
func NewApi(cfg *config.Config, log *slog.Logger, service Service, authService AuthService) *apiImpl {
	return &apiImpl{
		authService: authService,
		service:     service,
		log:         log,
		cfg:         cfg,
	}
}

type ServRegUser struct {
	FirstName string
	LastName  string
	Username  string
	Password  string
	Email     string
}
