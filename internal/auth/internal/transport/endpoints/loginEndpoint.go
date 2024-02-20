package endpoints

import (
	"context"
	"log"

	authservice "github.com/Antoha2/news/internal/auth/internal/services"

	"github.com/go-kit/kit/endpoint"
)

type LoginRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
	Roles    []string
}

type LoginResponse struct {
	Token string `json:"token"`
}

func MakeLoginEndpoint(s authservice.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(LoginRequest)
		token, err := s.Login(ctx, req.Username, req.Password, 0)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return LoginResponse{Token: token}, nil
	}
}
