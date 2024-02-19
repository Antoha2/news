package transport

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Antoha2/news/internal/config"
)

type Service interface {
	GetUsers(ctx context.Context, filter *service.QueryUsersFilter) ([]*service.User, error)
	GetUser(ctx context.Context, id int) (*service.User, error)
}

type apiImpl struct {
	cfg     *config.Config
	log     *slog.Logger
	service Service
	server  *http.Server
}

// NewAPI
func NewApi(cfg *config.Config, log *slog.Logger, service Service) *apiImpl {
	return &apiImpl{
		service: service,
		log:     log,
		cfg:     cfg,
	}
}
