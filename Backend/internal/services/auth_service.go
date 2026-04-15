package service

import (
	"context"
	"errors"
	"fmt"
	"taskflow/internal/models"
	"taskflow/internal/repository"
	"taskflow/internal/utils"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) error {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("lookup error: %w", err)
	}

	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("password hashing error: %w", err)
	}

	user := &models.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("auth lookup error: %w", err)
	}
	invalidError := errors.New("invalid email or password")
	if user == nil {
		return "", invalidError
	}

	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return "", invalidError
	}

	expiry := 24 * time.Hour
	token, err := utils.GenerateToken(
		user.ID.String(),
		user.Email,
		s.jwtSecret,
		expiry,
	)
	if err != nil {
		return "", fmt.Errorf("token generation error: %w", err)
	}
	return token, nil
}
