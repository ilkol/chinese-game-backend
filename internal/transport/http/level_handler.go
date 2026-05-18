package http

import (
	"chinese-game-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *LevelHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	levelID, err := strconv.Atoi(idStr)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid level id")
		return
	}
	userID := r.Context().Value(userContextKey).(int)

	level, err := h.service.GetByID(levelID, userID)
	if err != nil {
		errorJSON(w, http.StatusNotFound, "levevl not found")
		return
	}

	writeJSON(w, http.StatusOK, level)
}
