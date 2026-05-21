package http

import (
	"chinese-game-backend/internal/service"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) JoinStudentToTeacher(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InviteCode string `json:"invite_code"`
	}

	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_input")
		return
	}

	studentID := r.Context().Value(userContextKey).(int)

	if err := h.service.JoinStudentToTeacher(studentID, input.InviteCode); err != nil {
		errorJSON(w, http.StatusNotFound, "invalid_teacher_code")
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "joined"})
}
