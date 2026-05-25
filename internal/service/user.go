package service

import (
	"chinese-game-backend/internal/domain"
	"chinese-game-backend/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{repo, jwtSecret}
}

func (s *UserService) GetSecret() string {
	return s.jwtSecret
}

func (s *UserService) SignUp(username, password, inviteCode string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	var user domain.User
	if inviteCode == "" {
		user := domain.User{
			Username:     username,
			PasswordHash: string(hash),
			Role:         domain.RoleStudent,
		}

		user, err = s.repo.CreateUser(user)
	} else {
		user, err = s.repo.CreateTeacher(username, string(hash), inviteCode)
	}

	if err != nil {
		return "", err
	}
	return s.generateJWTToken(user)
}

func (s *UserService) SignIn(username, password string) (string, error) {
	user, err := s.repo.GetUserByName(username)
	if err != nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidUserPassword
	}

	return s.generateJWTToken(user)
}

func (s *UserService) generateJWTToken(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"role":      user.Role,
		"expiredOn": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) JoinStudentToTeacher(studentID int, inviteCode string) error {
	return s.repo.JoinStudentToTeacher(studentID, inviteCode)
}

func (s *UserService) GetStudentByTeacher(teacherID int) ([]domain.StudentProgressInfo, error) {
	return s.repo.GetStudentByTeacher(teacherID)
}

func (s *UserService) GetInviteCode(teacherID int) (string, error) {
	return s.repo.GetInviteCode(teacherID)
}
