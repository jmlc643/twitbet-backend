package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type RegisterUseCase struct {
	userRepo     repository.UserRepository
	hasher       port.PasswordHasher
	tokenService port.TokenService
	otpRepo      port.OTPRepository
	emailService port.EmailService
}

func NewRegisterUseCase(
	repo repository.UserRepository,
	hasher port.PasswordHasher,
	tokenSvc port.TokenService,
	otpRepo port.OTPRepository,
	emailService port.EmailService,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:     repo,
		hasher:       hasher,
		tokenService: tokenSvc,
		otpRepo:      otpRepo,
		emailService: emailService,
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

func generateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

func (uc *RegisterUseCase) Execute(ctx context.Context, in input.RegisterInput) error {
	existing, err := uc.userRepo.FindByEmail(ctx, in.Email)
	if err != nil {
		return apperror.ErrInternal
	}
	if existing != nil {
		return apperror.ErrUserAlreadyExists
	}

	hashedPassword, err := uc.hasher.HashPassword(in.Password)
	if err != nil {
		return apperror.ErrInternal
	}

	avatar := getRandomAvatar()

	user := &entity.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: hashedPassword,
		AvatarURL:    avatar,
		IsVerified:   false,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return apperror.ErrInternal
	}

	otpCode := generateOTP()
	otpKey := fmt.Sprintf("verify:%s", user.Email)
	
	if err := uc.otpRepo.SaveOTP(ctx, otpKey, otpCode, 10*time.Minute); err != nil {
		fmt.Printf("Error saving OTP to Redis: %v\n", err)
		return apperror.ErrInternal
	}

	go func() {
		if err := uc.emailService.SendVerificationEmail(context.Background(), user.Email, otpCode); err != nil {
			fmt.Printf("Error sending verification email to %s: %v\n", user.Email, err)
		}
	}()

	return nil
}
