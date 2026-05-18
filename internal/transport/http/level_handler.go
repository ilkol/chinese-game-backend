package http

import (
	"chinese-game-backend/internal/service"
	"net/http"
)

type LevelHandler struct {
	service *service.LevelService
}

func NewLevelHandler(service *service.LevelService) *LevelHandler {
	return &LevelHandler{service}
}

func (h *LevelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	levels, err := h.service.GetAll()
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "failed to get levels")
		return
	}

	writeJSON(w, http.StatusOK, levels)
}
