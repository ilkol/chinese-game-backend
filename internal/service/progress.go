package service

import "chinese-game-backend/internal/repository"

type ProgressService struct {
	repo *repository.ProgressRepository
}

func NewProgressRepository(repo *repository.ProgressRepository) *ProgressService {
	return &ProgressService{repo}
}

func (s *ProgressService) CompleteStep(userID, stepID int) error {
	return s.repo.CompleteStep(userID, stepID)
}
