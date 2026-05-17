package service

import (
	"chinese-game-backend/internal/domain"
	"chinese-game-backend/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) SignUp(username, password_hash string, role domain.UserRole) error {
	user := domain.User{
		Username:     username,
		PasswordHash: password_hash,
		Role:         role,
	}

	return s.repo.CreateUser(user)
}
