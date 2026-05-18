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

func (s *LevelService) GetAll() ([]domain.Level, error) {
	return s.repo.GetAll()
}

func (s *LevelService) GetByID(levelID, userID int) (domain.Level, error) {
	return s.repo.GetByID(levelID, userID)
}
