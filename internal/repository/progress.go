package repository

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type ProgressRepository struct {
	db *sqlx.DB
}

func NewProgressREpository(db *sqlx.DB) *ProgressRepository {
	return &ProgressRepository{db}
}

func (r *ProgressRepository) UpdateProgress(userID, levelID int, details json.RawMessage) error {
	query := `
		INSERT INTO progress 
		(user_id, level_id, details, is_completed) VALUES ($1, $1, $1, $1)
		ON CONFLICT (user_id, level_id)
		DO UPDATE SET 
			details = EXCLUDED.details, 
			is_completed = EXCLUDED.is_completed, 
			updated_at = NOW()
	`
	_, err := r.db.Exec(query, userID, levelID, details, false)
	return err
}
