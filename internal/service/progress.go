package service

import "chinese-game-backend/internal/repository"

type ProgressService struct {
	repo *repository.ProgressRepository
}

func NewProgressRepository(repo *repository.ProgressRepository) *ProgressService {
	return &ProgressService{repo}
}
