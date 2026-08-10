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
	var avatar *string
	if out.User.AvatarURL != "" {
		avatar = &out.User.AvatarURL
	}
	return response.AuthResponse{
		Token: out.Token,
		User: response.UserResponse{
			ID:        out.User.ID,
			Username:  out.User.Username,
			Email:     out.User.Email,
			AvatarURL: avatar,
			CreatedAt: out.User.CreatedAt.Format(time.RFC3339),
		},
	}
}

func UserOutputToResponse(out *output.UserOutput) response.UserResponse {
	var avatar *string
	if out.AvatarURL != "" {
		avatar = &out.AvatarURL
	}
	return response.UserResponse{
		ID:        out.ID,
		Username:  out.Username,
		Email:     out.Email,
		AvatarURL: avatar,
		CreatedAt: out.CreatedAt.Format(time.RFC3339),
	}
}

func UpdateProfileRequestToInput(req request.UpdateProfileRequest, userID string) input.UpdateProfileInput {
	return input.UpdateProfileInput{
		UserID:    userID,
		Username:  req.Username,
		AvatarURL: req.AvatarURL,
	}
}

func VerifyAccountRequestToInput(req request.VerifyAccountRequest) input.VerifyAccountInput {
	return input.VerifyAccountInput{
		Email:   req.Email,
		OTPCode: req.OTPCode,
	}
}

func ForgotPasswordRequestToInput(req request.ForgotPasswordRequest) input.ForgotPasswordInput {
	return input.ForgotPasswordInput{
		Email: req.Email,
	}
}

func ResetPasswordRequestToInput(req request.ResetPasswordRequest) input.ResetPasswordInput {
	return input.ResetPasswordInput{
		Email:           req.Email,
		OTPCode:         req.OTPCode,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}
}

func ChangePasswordRequestToInput(req request.ChangePasswordRequest, userID string) input.ChangePasswordInput {
	return input.ChangePasswordInput{
		UserID:          userID,
		OldPassword:     req.OldPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}
}

func VerifyResetOtpRequestToInput(req request.VerifyResetOtpRequest) input.VerifyResetOtpInput {
	return input.VerifyResetOtpInput{
		Email:   req.Email,
		OTPCode: req.OTPCode,
	}
}
