package auth

import (
	"crypto/subtle"
	"errors"
	"xray_server/internal/domain"
)

type LoginUseCase struct {
	userRepo domain.UserRepository
	jwtSvc   domain.JWTService
}

func NewLoginUseCase(userRepo domain.UserRepository, jwtSvc domain.JWTService) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	Token     string
	ExpiresIn int64
}

func (uc *LoginUseCase) Execute(input LoginInput) (*LoginOutput, error) {
	user, err := uc.userRepo.FindByUsername(input.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(input.Password)) != 1 {
		return nil, errors.New("invalid credentials")
	}

	token, expiresIn, err := uc.jwtSvc.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token:     token,
		ExpiresIn: expiresIn,
	}, nil
}
