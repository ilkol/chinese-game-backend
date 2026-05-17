package repository

import (
	"chinese-game-backend/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(user domain.User) error {
	query := "INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, user.Username, user.PasswordHash, user.Role)
	return err
}

func (r *UserRepository) GetUserByName(username string) (domain.User, error) {
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1"
	var user domain.User
	err := r.db.Get(&user, query, username)
	return user, err
}
