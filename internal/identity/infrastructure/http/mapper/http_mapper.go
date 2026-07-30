package mapper

import (
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/application/output"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/dto/response"
)

func RegisterRequestToInput(req request.RegisterRequest) input.RegisterInput {
	return input.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
}

func LoginRequestToInput(req request.LoginRequest) input.LoginInput {
	return input.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

func AuthOutputToResponse(out *output.AuthOutput) response.AuthResponse {
	return response.AuthResponse{
		Token: out.Token,
		User: response.UserResponse{
			ID:        out.User.ID,
			Username:  out.User.Username,
			Email:     out.User.Email,
			AvatarURL: out.User.AvatarURL,
			CreatedAt: out.User.CreatedAt.Format(time.RFC3339),
		},
	}
}

func UserOutputToResponse(out *output.UserOutput) response.UserResponse {
	return response.UserResponse{
		ID:        out.ID,
		Username:  out.Username,
		Email:     out.Email,
		AvatarURL: out.AvatarURL,
		CreatedAt: out.CreatedAt.Format(time.RFC3339),
	}
}
