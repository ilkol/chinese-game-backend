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

func (r *LevelRepository) CreateStep(levelID int, step domain.LevelStep) (domain.LevelStep, error) {
	query := `
		INSERT INTO level_steps (level_id, type, title, description, content, order_index)
		VALUES ($1, $2, $3, $4, $5::jsonb, (
			SELECT COALESCE(MAX(order_index), 0) + 1 FROM level_steps WHERE level_id = $1
		))
		RETURNING *
	`
	var created domain.LevelStep
	err := r.db.Get(&created, query,
		levelID, step.Type, step.Title, step.Description, []byte(step.Content),
	)
	return created, err
}

func (r *LevelRepository) UpdateStep(stepID int, step domain.LevelStep) error {
	query := `
		UPDATE level_steps
		SET type = $1, title = $2, description = $3, content = $4::jsonb, order_index = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(query,
		step.Type, step.Title, step.Description, []byte(step.Content), step.OrderIndex, stepID,
	)
	return err
}

func (r *LevelRepository) DeleteStep(stepID int) error {
	_, err := r.db.Exec("DELETE FROM level_steps WHERE id = $1", stepID)
	return err
}

func (r *LevelRepository) UpsertDialog(stepID int, steps json.RawMessage) error {
	query := `
		INSERT INTO step_dialogs (step_id, steps)
		VALUES ($1, $2)
		ON CONFLICT (step_id) DO UPDATE SET steps = $2, updated_at = NOW()
	`
	_, err := r.db.Exec(query, stepID, steps)
	return err
}

func (r *LevelRepository) GetStepByID(stepID int) (domain.LevelStep, error) {
	var step domain.LevelStep
	err := r.db.Get(&step, "SELECT * FROM level_steps WHERE id = $1", stepID)
	return step, err
}
