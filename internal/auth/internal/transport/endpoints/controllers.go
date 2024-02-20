package endpoints

import (
	authservice "github.com/Antoha2/news/internal/auth/internal/services"

	"github.com/go-kit/kit/endpoint"
)

type AuthEndpoints struct {
	SignIn endpoint.Endpoint
	SignUp endpoint.Endpoint
	// SignUpUser  endpoint.Endpoint
	// DeleteUser  endpoint.Endpoint
	// UpdateUser  endpoint.Endpoint
	// ParseToken  endpoint.Endpoint
	// GetRoles    endpoint.Endpoint
}

func MakeAuthEndpoints(s authservice.AuthService) *AuthEndpoints {
	return &AuthEndpoints{
		SignIn: MakeLoginEndpoint(s),
		SignUp: MakeRegisterEndpoint(s),
		// SignUpUser:  MakeSignUpUserEndpoint(s),
		// DeleteUser:  MakeDeleteUserEndpoint(s),
		// UpdateUser:  MakeUpdateUserEndpoint(s),
		// ParseToken:  MakeParseTokenEndpoint(s),
		// GetRoles:    MakeGetRolesEndpoint(s),
	}
}
