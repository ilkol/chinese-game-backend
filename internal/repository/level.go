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

func (r *LevelRepository) GetAll() ([]domain.Level, error) {
	query := "SELECT * FROM levels"
	var levels []domain.Level
	err := r.db.Select(&levels, query)
	if err != nil {
		return nil, err
	}
	return levels, nil
}
