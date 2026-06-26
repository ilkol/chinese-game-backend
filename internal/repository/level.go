package repository

import (
	"chinese-game-backend/internal/domain"
	"encoding/json"

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
		type stepWithDialog struct {
			domain.LevelStep
			DialogSteps *string `db:"dialog_steps"`
		}

		var rows []stepWithDialog
		err = r.db.Select(&rows, `
        SELECT 
            ls.*,
            sd.steps::text AS dialog_steps
        FROM level_steps ls
        LEFT JOIN step_dialogs sd ON sd.step_id = ls.id
        ORDER BY ls.order_index ASC
    `)
		if err != nil {
			return nil, err
		}

		stepMap := make(map[int][]domain.LevelStep)
		for _, row := range rows {
			step := row.LevelStep
			if row.DialogSteps != nil {
				json.Unmarshal([]byte(*row.DialogSteps), &step.Dialog)
			}
			stepMap[step.LevelID] = append(stepMap[step.LevelID], step)
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
			ls.*,
			COALESCE(p.is_completed, false) AS is_completed,
			sd.steps::text AS dialog_steps
		FROM level_steps ls
		LEFT JOIN user_progress p ON ls.id = p.step_id AND p.user_id = $2
		LEFT JOIN step_dialogs sd ON sd.step_id = ls.id
		WHERE ls.level_id = $1
		ORDER BY ls.order_index ASC
	`
	type stepWithDialog struct {
		domain.LevelStep
		DialogSteps *string `db:"dialog_steps"`
	}

	var rows []stepWithDialog
	err = r.db.Select(&rows, query, levelID, userID)
	if err != nil {
		return domain.Level{}, err
	}

	steps := make([]domain.LevelStep, 0, len(rows))
	for _, row := range rows {
		step := row.LevelStep
		if row.DialogSteps != nil {
			json.Unmarshal([]byte(*row.DialogSteps), &step.Dialog)
		}
		steps = append(steps, step)
	}

	level.Steps = steps
	return level, nil
}
