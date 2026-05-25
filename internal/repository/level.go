package repository

import (
	"chinese-game-backend/internal/domain"

	"github.com/jmoiron/sqlx"
)

type LevelRepository struct {
	db *sqlx.DB
}

func NewLevelRepository(db *sqlx.DB) *LevelRepository {
	return &LevelRepository{db}
}

func (r *LevelRepository) GetAll(withSteps bool) ([]domain.Level, error) {
	query := "SELECT * FROM levels ORDER BY order_index ASC"
	levels := make([]domain.Level, 0)
	err := r.db.Select(&levels, query)
	if err != nil {
		return nil, err
	}

	if withSteps {
		var steps []domain.LevelStep
		err = r.db.Select(&steps, "SELECT * FROM level_steps ORDER BY order_index ASC")
		if err != nil {
			return nil, err
		}

		stepMap := make(map[int][]domain.LevelStep)
		for _, s := range steps {
			stepMap[s.LevelID] = append(stepMap[s.LevelID], s)
		}

		for i := range levels {
			levels[i].Steps = stepMap[levels[i].ID]
		}
	}

	return levels, nil
}

func (r *LevelRepository) GetByID(levelID, userID int) (domain.Level, error) {
	query := "SELECT * FROM levels WHERE id=$1"
	var level domain.Level
	err := r.db.Get(&level, query, levelID)
	if err != nil {
		return domain.Level{}, err
	}

	query = `
		SELECT 
			s.*,
			COALESCE(p.is_completed, false) as is_completed
		FROM level_steps s
		LEFT JOIN user_progress p ON s.id = p.step_id AND p.user_id = $2
		WHERE level_id=$1
		ORDER BY s.order_index ASC
	`
	var steps []domain.LevelStep
	err = r.db.Select(&steps, query, levelID, userID)
	if err != nil {
		return domain.Level{}, err
	}

	level.Steps = steps
	return level, nil
}
