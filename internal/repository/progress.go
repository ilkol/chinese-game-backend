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

func (r *ProgressRepository) IsLevelCompleted(userID, levelID int) (bool, error) {
	query := `
		SELECT COUNT(*) = (
			SELECT COUNT(*) FROM level_steps WHERE level_id = $2
		)
		FROM user_progress up
		JOIN level_steps ls ON ls.id = up.step_id
		WHERE up.user_id = $1 
		AND ls.level_id = $2
		AND up.is_completed = true
	`
	var completed bool
	err := r.db.Get(&completed, query, userID, levelID)
	return completed, err
}

func (r *ProgressRepository) GetCompletedLevels(userID int) ([]int, error) {
	query := `
		SELECT ls.level_id
		FROM user_progress up
		JOIN level_steps ls ON ls.id = up.step_id
		WHERE up.user_id = $1 AND up.is_completed = true
		GROUP BY ls.level_id
		HAVING COUNT(*) = (
			SELECT COUNT(*) FROM level_steps WHERE level_id = ls.level_id
		)
	`
	var levelIDs []int
	err := r.db.Select(&levelIDs, query, userID)
	return levelIDs, err
}
