package transport

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Antoha2/news/pkg/jwt"
	"github.com/Antoha2/news/pkg/logger/sl"
	"github.com/pkg/errors"
)

type ContextKey string

const (
	authorizationHeader = "Authorization"
	authErr             = "occurred error for Authorization user"
)

func (a *apiImpl) userIdentify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const op = "Authorization user"
		log := a.log.With(slog.String("op", op))

		header := r.Header.Get(authorizationHeader)
		if header == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("empty header", sl.Err(errors.New(authErr)))
			w.Write([]byte(authErr))
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {

			w.WriteHeader(http.StatusBadRequest)
			log.Error("wrong header", sl.Err(errors.New(authErr)))
			w.Write([]byte(authErr))
			return
		}

		userEmail, err := jwt.ParseToken(a.cfg, headerParts[1])
		if err != nil {

			log.Error("error ParseToken", sl.Err(errors.New(authErr)))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(authErr))
			return
		}

		if userEmail == "" {

			log.Error("no access", sl.Err(errors.New(authErr)))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(authErr))
			return
		}

		next.ServeHTTP(w, r)
	})
}
