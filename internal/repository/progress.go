package repository

import (
	"github.com/jmoiron/sqlx"
)

type ProgressRepository struct {
	db *sqlx.DB
}

func NewProgressRepository(db *sqlx.DB) *ProgressRepository {
	return &ProgressRepository{db}
}

func (r *ProgressRepository) CompleteStep(userID, stepID int) error {
	query := `
		INSERT INTO user_progress 
		(user_id, step_id, is_completed) VALUES ($1, $2, true)
		ON CONFLICT (user_id, step_id)
		DO UPDATE SET is_completed = true, updated_at = NOW()
	`
	_, err := r.db.Exec(query, userID, stepID)
	return err
}
