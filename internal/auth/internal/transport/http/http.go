package transport

import (
	"context"
	"log/slog"
	"net/http"

	authService "github.com/Antoha2/news/internal/auth/internal/services"
)

type regUser struct {
	name  string `json:"name"`
	email string `json:"email"`
}

type apiImpl struct {
	authService authService.AuthService
	server      *http.Server
	log         *slog.Logger
	port        int
}

func NewApi(authService authService.AuthService, log *slog.Logger, port int) *apiImpl {
	return &apiImpl{
		authService: authService,
		log:         log,
		port:        port,
	}
}

func (wImpl *apiImpl) Stop() {

	if err := wImpl.server.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
