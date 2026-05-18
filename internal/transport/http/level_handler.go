package http

import (
	"chinese-game-backend/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LevelHandler struct {
	levelService *service.LevelService
	progressStep *service.ProgressService
}

func NewLevelHandler(service *service.LevelService, progressStep *service.ProgressService) *LevelHandler {
	return &LevelHandler{service, progressStep}
}

func (h *LevelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	levels, err := h.levelService.GetAll()
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

	level, err := h.levelService.GetByID(levelID, userID)
	if err != nil {
		errorJSON(w, http.StatusNotFound, "levevl not found")
		return
	}

	writeJSON(w, http.StatusOK, level)
}

func (h *LevelHandler) CompleteStep(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StepID int `json:"step_id"`
	}

	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid input")
		return
	}

	userID := r.Context().Value(userContextKey).(int)

	if err := h.progressStep.CompleteStep(userID, input.StepID); err != nil {
		log.Println(err)
		errorJSON(w, http.StatusInternalServerError, "failed to save progress")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
