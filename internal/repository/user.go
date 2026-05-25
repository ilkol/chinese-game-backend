package repository

import (
	"chinese-game-backend/internal/domain"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(inputUser domain.User) (domain.User, error) {
	var user domain.User
	query := "INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING *"
	err := r.db.QueryRowx(query, inputUser.Username, inputUser.PasswordHash, inputUser.Role).StructScan(&user)
	return user, err
}

func (r *UserRepository) CreateTeacher(username, passwordHash, inviteCode string) (domain.User, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback()

	var inviteCodeExists bool

	err = tx.Get(&inviteCodeExists, "SELECT EXISTS(SELECT 1 FROM admin_invites WHERE code=$1 LIMIT 1)", inviteCode)
	if err != nil || !inviteCodeExists {
		return domain.User{}, errors.New("invalid_invite_code")
	}

	quey := "INSERT INTO users (username, password_hash, role, invite_code) VALUES ($1, $2, $3, $4) RETURNING *"
	var user domain.User
	err = tx.QueryRowx(quey, username, passwordHash, domain.RoleTeacher, inviteCode).StructScan(&user)
	if err != nil {
		return domain.User{}, err
	}
	_, err = tx.Exec("UPDATE admin_invites SET used_at = NOW() WHERE code = $1", inviteCode)

	if err != nil {
		return domain.User{}, err
	}

	err = tx.Commit()
	return user, err
}

func (r *UserRepository) GetUserByName(username string) (domain.User, error) {
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1"
	var user domain.User
	err := r.db.Get(&user, query, username)
	return user, err
}

func (r *UserRepository) JoinStudentToTeacher(studentID int, inviteCode string) error {
	query := "SELECT * FROM users WHERE invite_code = $1 LIMIT 1"
	var teacher domain.User
	err := r.db.Get(&teacher, query, inviteCode)
	if err != nil {
		return nil
	}

	query = `
		INSERT INTO teacher_students
		(teacher_id, student_id)
		VALUES
		($1, $2)
		ON CONFLICT (teacher_id, student_id) DO NOTHING
	`
	_, err = r.db.Exec(query, teacher.ID, studentID)
	return err
}

func (r *UserRepository) GetStudentByTeacher(teacherID int) ([]domain.User, error) {
	query := `
		SELECT u.* 
		FROM users u
		JOIN teacher_students ts ON u.id = ts.student_id
		WHERE ts.teacher_id = $1
		ORDER BY u.username ASC
	`
	students := make([]domain.User, 0)
	err := r.db.Select(&students, query, teacherID)
	fmt.Println(students, err)
	return students, err
}

func (r *UserRepository) GetInviteCode(teacherID int) (string, error) {
	query := "SELECT * FROM users WHERE id = $1 LIMIT 1"
	var user domain.User
	err := r.db.Get(&user, query, teacherID)
	return *user.InviteCode, err
}
