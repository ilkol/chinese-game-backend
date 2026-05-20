package repository

import (
	"chinese-game-backend/internal/domain"
	"errors"

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

func (r *UserRepository) CreateTeacher(username, passwordHash, inviteCode string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var inviteCodeExists bool

	err = tx.Get(&inviteCodeExists, "SELECT EXISTS(SELECT 1 FROM admin_invites WHERE code=$1 LIMIT 1)", inviteCode)
	if err != nil || !inviteCodeExists {
		return errors.New("invalid_invite_code")
	}

	_, err = tx.Exec("INSERT INTO users (username, password_hash, role, invite_code) VALUES ($1, $2, $3, $4)", username, passwordHash, domain.RoleTeacher, inviteCode)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE admin_invites SET used_at = NOW() WHERE code = $1", inviteCode)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetUserByName(username string) (domain.User, error) {
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1"
	var user domain.User
	err := r.db.Get(&user, query, username)
	return user, err
}
