package service

import "chinese-game-backend/internal/repository"

type ProgressService struct {
	repo *repository.ProgressRepository
}

func NewProgressService(repo *repository.ProgressRepository) *ProgressService {
	return &ProgressService{repo}
}

func (s *ProgressService) CompleteStep(userID, stepID int) error {
	return s.repo.CompleteStep(userID, stepID)
}

func (s *ProgressService) IsLevelCompleted(userID, levelID int) (bool, error) {
	return s.repo.IsLevelCompleted(userID, levelID)
}

func (s *ProgressService) GetCompletedLevels(userID int) ([]int, error) {
	return s.repo.GetCompletedLevels(userID)
}
