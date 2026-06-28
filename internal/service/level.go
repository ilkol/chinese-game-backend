package service

import (
	"chinese-game-backend/internal/domain"
	"chinese-game-backend/internal/repository"
	"encoding/json"
)

type LevelService struct {
	repo *repository.LevelRepository
}

func NewLevelService(repo *repository.LevelRepository) *LevelService {
	return &LevelService{repo}
}

func (s *LevelService) GetAll(withSteps bool) ([]domain.Level, error) {
	return s.repo.GetAll(withSteps)
}

func (s *LevelService) GetByID(levelID, userID int) (domain.Level, error) {
	return s.repo.GetByID(levelID, userID)
}

func (s *LevelService) CreateStep(levelID int, step domain.LevelStep) (domain.LevelStep, error) {
	return s.repo.CreateStep(levelID, step)
}

func (s *LevelService) UpdateStep(stepID int, step domain.LevelStep) error {
	return s.repo.UpdateStep(stepID, step)
}

func (s *LevelService) DeleteStep(stepID int) error {
	return s.repo.DeleteStep(stepID)
}

func (s *LevelService) UpsertDialog(stepID int, steps json.RawMessage) error {
	return s.repo.UpsertDialog(stepID, steps)
}
