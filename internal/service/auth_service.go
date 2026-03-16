package service

import (
	"errors"
	"os"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/breamon/sinav-bilgi-sistemi/internal/utils"
)

type AuthService struct {
	userRepo *postgres.UserRepository
}

func NewAuthService(userRepo *postgres.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(fullName, email, password string) (*domain.User, string, error) {
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, "", errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*domain.User, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if err := utils.CheckPassword(password, user.PasswordHash); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Me(userID int64) (*domain.User, error) {
	return s.userRepo.GetByID(userID)
}
