package endpoints

import (
	"context"

	"github.com/Antoha2/news/internal/auth/internal/services"
	"github.com/go-kit/kit/endpoint"
)

type RegisterRequest struct {
	// FirstName string `json:"firstname"`
	// LastName  string `json:"lastname"`
	// Username  string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserId int64 `json:"user_id"`
}

// const (
// 	roleAdmin = "admin"
// 	roleUser  = "user"
// 	roleDev   = "dev"
// )

func MakeRegisterEndpoint(s services.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(RegisterRequest)
		userId, err := s.RegisterNewUser(ctx, req.Email, req.Password)
		if err != nil {
			return nil, err
		}

		// inputUser := new(services.ServRegUser)
		// inputUser.Password = req.Password
		// inputUser.Email = req.Email
		// inputRoles := new(helper.UsersRoles)

		// inputUser.FirstName = req.FirstName
		// inputUser.LastName = req.LastName

		// inputUser.Username = req.Username

		// inputRoles.Roles = append(inputRoles.Roles, roleAdmin)
		// inputRoles.Roles = append(inputRoles.Roles, roleDev) // ?!? !!!!!!!!!!!!!!!!!!!!!!!!!!

		return RegisterResponse{UserId: userId}, nil
	}
}
