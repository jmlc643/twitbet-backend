package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/application/output"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type RegisterUseCase struct {
	userRepo 		repository.UserRepository
	hasher			port.PasswordHasher
	tokenService	port.TokenService
}

func NewRegisterUseCase(
	repo		repository.UserRepository,
	hasher		port.PasswordHasher,
	tokenSvc 	port.TokenService,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: 		repo,
		hasher: 		hasher,
		tokenService: 	tokenSvc,
	}
}

var defaultAvatars = []string{
	"/avatars/avatar1.png",
	"/avatars/avatar2.png",
	"/avatars/avatar3.png",
	"/avatars/avatar4.png",
}

func getRandomAvatar() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return defaultAvatars[r.Intn(len(defaultAvatars))]
}

func (uc *RegisterUseCase) Execute(ctx context.Context, in input.RegisterInput) (*output.AuthOutput, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, in.Email)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if existing != nil {
		return nil, apperror.ErrUserAlreadyExists
	}

	hashedPassword, err := uc.hasher.HashPassword(in.Password)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	avatar := getRandomAvatar()

	user := &entity.User{
		Username:     	in.Username,
		Email:        	in.Email,
		PasswordHash: 	hashedPassword,
		AvatarURL:		avatar,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.ErrInternal
	}

	token, err := uc.tokenService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	return &output.AuthOutput{
		Token: token,
		User: output.UserOutput{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}