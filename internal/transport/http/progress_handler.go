package http

import (
	"chinese-game-backend/internal/service"
	"log"
	"net/http"
	"strconv"
)

type ProgressHandler struct {
	progressService *service.ProgressService
}

func NewProgressHandler(progressService *service.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressService}
}

func (h *ProgressHandler) CompleteStep(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StepID int `json:"step_id"`
	}

	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid input")
		return
	}

	userID := r.Context().Value(userContextKey).(int)

	if err := h.progressService.CompleteStep(userID, input.StepID); err != nil {
		log.Println(err)
		errorJSON(w, http.StatusInternalServerError, "failed to save progress")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *ProgressHandler) GetCompletedLevels(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userContextKey).(int)
	levelIDs, err := h.progressService.GetCompletedLevels(userID)
	if err != nil {
		log.Println(err)
		errorJSON(w, http.StatusInternalServerError, "failed to get completed levels")
		return
	}

	writeJSON(w, http.StatusOK, levelIDs)
}

func (h *ProgressHandler) IsLevelCompleted(w http.ResponseWriter, r *http.Request) {
	levelIDStr := r.URL.Query().Get("level_id")
	levelID, err := strconv.Atoi(levelIDStr)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid level_id")
		return
	}

	userID := r.Context().Value(userContextKey).(int)
	isCompleted, err := h.progressService.IsLevelCompleted(userID, levelID)
	if err != nil {
		log.Println(err)
		errorJSON(w, http.StatusInternalServerError, "failed to check level completion")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"is_completed": isCompleted})
}

func (h *ProgressHandler) GetCompletedLevelsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userContextKey).(int)
	levelIDs, err := h.progressService.GetCompletedLevels(userID)
	if err != nil {
		log.Println(err)
		errorJSON(w, http.StatusInternalServerError, "failed to get completed levels")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"level_ids": levelIDs})
}

func (h *ProgressHandler) IsLevelCompletedHandler(w http.ResponseWriter, r *http.Request) {
	levelIDStr := r.URL.Query().Get("level_id")
	levelID, err := strconv.Atoi(levelIDStr)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid level_id")
		return
	}
	isCompleted, err := h.progressService.IsLevelCompleted(r.Context().Value(userContextKey).(int), levelID)
	if err != nil {
		log.Println(err)
		errorJSON(w, http.StatusInternalServerError, "failed to check level completion")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"is_completed": isCompleted})
}
