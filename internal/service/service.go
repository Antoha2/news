package service

import (
	"context"
	"log/slog"

	"github.com/Antoha2/news/internal/config"
	"github.com/Antoha2/news/internal/repository"
)

type Repository interface {
	GetNews(ctx context.Context) ([]*repository.RepNews, error)
	AddNews(ctx context.Context, news *repository.RepNews) (int, error)
	// GetUsers(ctx context.Context, filter *repository.RepQueryFilter) ([]*repository.RepUser, error)
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

type News struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories []int  `json:"categories"`
}

type EditNewsFilter struct {
}
