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

type completeStepInput struct {
	StepID int `json:"step_id"`
}

// CompleteStep godoc
// @Summary Update user progress for a specific step
// @Description Mark a step as completed for the authenticated user
// @Tags Progress
// @Param	input	body	completeStepInput	true	"Step ID"
// @Success	204 "Step marked as completed"
// @Failure	400	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router	/progress [post]
// @Security BearerAuth
func (h *ProgressHandler) CompleteStep(w http.ResponseWriter, r *http.Request) {
	var input completeStepInput

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

	w.WriteHeader(http.StatusNoContent)
}

// GetCompletedLevels godoc
// @Summary Get completed levels for the authenticated user
// @Description Retrieve a list of level IDs that the authenticated user has completed
// @Tags Progress
// @Success	200	{object}	[]int
// @Failure	500	{object}	ErrorResponse
// @Router	/progress/levels [get]
// @Security BearerAuth
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

// IsLevelCompleted godoc
// @Summary Check if a level is completed for the authenticated user
// @Description Determine if the authenticated user has completed a specific level
// @Tags Progress
// @Param	level_id	path	int	true	"Level ID"
// @Success	200	{object}	map[string]bool
// @Failure	400	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router	/progress/levels/{level_id} [get]
// @Security BearerAuth
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
