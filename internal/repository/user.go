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

var ErrInvalidInviteCode = errors.New("invalid invite code")

func (r *UserRepository) JoinStudentToTeacher(studentID int, inviteCode string) error {
	query := "SELECT * FROM users WHERE invite_code = $1 LIMIT 1"
	var teacher domain.User
	err := r.db.Get(&teacher, query, inviteCode)
	if err != nil {
		return ErrInvalidInviteCode
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

func (r *UserRepository) GetStudentByTeacher(teacherID int) ([]domain.StudentProgressInfo, error) {
	var result []domain.StudentProgressInfo

	query := `
		WITH student_list AS (
			SELECT u.id, u.username, u.coins 
			FROM users u
			JOIN teacher_students ts ON u.id = ts.student_id
			WHERE ts.teacher_id = $1
		),
		max_progress AS (
			SELECT DISTINCT ON (p.user_id)
				p.user_id,
				l.title as last_level_title,
				(SELECT COUNT(*) FROM user_progress up WHERE up.user_id = p.user_id AND up.step_id IN (SELECT id FROM level_steps ls WHERE ls.level_id = l.id)) as completed_steps,
				(SELECT COUNT(*) FROM level_steps ls WHERE ls.level_id = l.id) as total_steps,
				p.updated_at
			FROM user_progress p
			JOIN level_steps s ON p.step_id = s.id
			JOIN levels l ON s.level_id = l.id
			ORDER BY p.user_id, l.order_index DESC, p.updated_at DESC
		)
		SELECT 
			sl.*, 
			COALESCE(mp.last_level_title, 'Не начато') as last_level_title,
			COALESCE(mp.completed_steps, 0) as completed_steps,
			COALESCE(mp.total_steps, 0) as total_steps,
			mp.updated_at
		FROM student_list sl
		LEFT JOIN max_progress mp ON sl.id = mp.user_id
		ORDER BY sl.username ASC`

	err := r.db.Select(&result, query, teacherID)
	return result, err
}

func (r *UserRepository) GetInviteCode(teacherID int) (string, error) {
	query := "SELECT * FROM users WHERE id = $1 LIMIT 1"
	var user domain.User
	err := r.db.Get(&user, query, teacherID)
	return *user.InviteCode, err
}
