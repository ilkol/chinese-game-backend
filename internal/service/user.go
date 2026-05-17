package service

import (
	"chinese-game-backend/internal/domain"
	"chinese-game-backend/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{repo, jwtSecret}
}

func (s *UserService) SignUp(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         domain.RoleStudent,
	}

	return s.repo.CreateUser(user)
}

func (s *UserService) SignIn(username, password string) (string, error) {
	user, err := s.repo.GetUserByName(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"role":      user.Role,
		"expiredOn": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
