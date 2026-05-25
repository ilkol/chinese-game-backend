package service

import (
	"chinese-game-backend/internal/domain"
	"chinese-game-backend/internal/repository"
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
