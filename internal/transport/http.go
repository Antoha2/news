package transport

import (
	"context"
	"log/slog"
	"net/http"

	config "github.com/Antoha2/news/internal/config"
	"github.com/Antoha2/news/internal/service"
)

type Service interface {
	GetNews(ctx context.Context) ([]*service.News, error)
	//	EditNews(ctx context.Context, id int) (*service.News, error)
	AddNews(ctx context.Context, news *service.News) (*service.News, error)
}

type apiImpl struct {
	service Service
	cfg     *config.Config
	log     *slog.Logger
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
