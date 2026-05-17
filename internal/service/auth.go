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
	roles      map[string]int
	jWTSecret  []byte
}

func NewAuthService(repo repository.AuthRepository, JWTSecret string) *AuthService {
	return &AuthService{repository: repo, jWTSecret: []byte(JWTSecret)}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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

	return token.SignedString(s.jWTSecret)
}

func (s *AuthService) Register(ctx context.Context, email, password, role string) (domain.UserShort, error) {
	var us domain.UserShort

	if _, ok := s.roles[role]; !ok {
		return us, domain.ErrWrongUserRole
	}
	creatorRole, ok := ctx.Value("user_role").(string)
	if !ok || (creatorRole != "admin" && s.roles[role] >= s.roles[creatorRole]) {
		return us, domain.ErrWrongUserRole
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return us, err
	}

	return s.repository.CreateUserAuth(ctx, email, passwordHash, role)
}
