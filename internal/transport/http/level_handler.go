package http

import (
	"chinese-game-backend/internal/domain"
	"chinese-game-backend/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LevelHandler struct {
	levelService *service.LevelService
}

func NewLevelHandler(service *service.LevelService) *LevelHandler {
	return &LevelHandler{service}
}

// GetAll godoc
// @Summary Get all levels
// @Description Retrieve a list of all levels, optionally including their steps
// @Tags levels
// @Success	200	{array}		domain.Level
// @Failure	500	{object}	ErrorResponse
// @Router	/level [get]
func (h *LevelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	withSteps := r.URL.Query().Get("with_steps") == "true"

	levels, err := h.levelService.GetAll(withSteps)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "failed to get levels")
		return
	}

	writeJSON(w, http.StatusOK, levels)
}

// GetByID godoc
// @Summary Get level by ID
// @Description Retrieve a specific level by its ID, including its steps
// @Tags levels
// @Param	id	path	int	true	"Level ID"
// @Success	200	{object}	domain.Level
// @Failure	404	{object}	ErrorResponse
// @Router	/level/{id} [get]
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
		fmt.Print(err)
		errorJSON(w, http.StatusNotFound, "levevl not found")
		return
	}

	writeJSON(w, http.StatusOK, level)
}

// CreateStep godoc
// @Summary Create a new step for a level
// @Description Add a new step to a specific level
// @Tags levels
// @Param	id	path	int	true	"Level ID"
// @Param	step	body	domain.LevelStep	true	"Step data"
// @Success	201	{object}	domain.LevelStep
// @Failure	400	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router	/level/{id}/step [post]
func (h *LevelHandler) CreateStep(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	levelID, err := strconv.Atoi(idStr)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid level id")
		return
	}

	var input domain.LevelStep
	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid input")
		return
	}

	step, err := h.levelService.CreateStep(levelID, input)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "failed to create step")
		return
	}

	writeJSON(w, http.StatusCreated, step)
}

// UpdateStep godoc
// @Summary Update a step for a level
// @Description Update an existing step of a specific level
// @Tags levels
// @Param	step_id	path	int	true	"Step ID"
// @Param	id	path	int	true	"Level ID"
// @Param	step	body	domain.LevelStep	true	"Updated step data"
// @Success	200	{object}	map[string]string
// @Failure	400	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router	/level/{id}/step/{step_id} [put]
func (h *LevelHandler) UpdateStep(w http.ResponseWriter, r *http.Request) {
	stepID, err := strconv.Atoi(chi.URLParam(r, "step_id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid step id")
		return
	}

	var input domain.LevelStep
	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid input")
		return
	}

	if err := h.levelService.UpdateStep(stepID, input); err != nil {
		errorJSON(w, http.StatusInternalServerError, "failed to update step")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// DeleteStep godoc
// @Summary Delete a step for a level
// @Description Remove an existing step of a specific level
// @Tags levels
// @Param	step_id	path	int	true	"Step ID"
// @Success	200	{object}	map[string]string
// @Failure	400	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router	/level/{id}/step/{step_id} [delete]
func (h *LevelHandler) DeleteStep(w http.ResponseWriter, r *http.Request) {
	stepID, err := strconv.Atoi(chi.URLParam(r, "step_id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid step id")
		return
	}

	if err := h.levelService.DeleteStep(stepID); err != nil {
		errorJSON(w, http.StatusInternalServerError, "failed to delete step")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type upsertDialogInput struct {
	Steps json.RawMessage `json:"steps"`
}

// UpsertDialog godoc
// @Summary Upsert dialog for a step
// @Description Create or update the dialog for a specific step
// @Tags levels
// @Param	step_id	path	int	true	"Step ID"
// @Param	input	body	upsertDialogInput	true	"Dialog steps"
// @Success	200	{object}	map[string]string
// @Failure	400	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router	/step/{step_id}/dialog [put]
func (h *LevelHandler) UpsertDialog(w http.ResponseWriter, r *http.Request) {
	stepID, err := strconv.Atoi(chi.URLParam(r, "step_id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid step id")
		return
	}

	var input upsertDialogInput
	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid input")
		return
	}

	if err := h.levelService.UpsertDialog(stepID, input.Steps); err != nil {
		errorJSON(w, http.StatusInternalServerError, "failed to save dialog")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
