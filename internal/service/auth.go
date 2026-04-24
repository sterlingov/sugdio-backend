package service

import (
	"context"
	"time"

	"sugdio/internal/domain"
	"sugdio/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repository repository.AuthRepository
	JWTSecret  []byte
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repository: repo}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repository.GetByEmailAuth(ctx, email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(s.JWTSecret)
}
