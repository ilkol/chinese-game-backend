package http

import (
	"chinese-game-backend/internal/domain"
	"chinese-game-backend/internal/repository"
	"chinese-game-backend/internal/service"
	"errors"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

type joinStudentInput struct {
	InviteCode string `json:"invite_code"`
}

// JoinStudentToTeacher godoc
// @Summary Join a student to a teacher using an invite code
// @Description Join a student to a teacher using an invite code
// @Tags User
// @Param	input	body	joinStudentInput	true	"Invite Code"
// @Success	204	"Student successfully joined the teacher"
// @Failure	400	{object}	ErrorResponse
// @Failure	404	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router	/user/join [post]
// @Security BearerAuth
func (h *UserHandler) JoinStudentToTeacher(w http.ResponseWriter, r *http.Request) {
	var input joinStudentInput

	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_input")
		return
	}

	studentID := r.Context().Value(userContextKey).(int)

	if err := h.service.JoinStudentToTeacher(studentID, input.InviteCode); err != nil {
		if errors.Is(err, repository.ErrInvalidInviteCode) {
			errorJSON(w, http.StatusNotFound, "invalid_teacher_code")
			return
		}
		errorJSON(w, http.StatusInternalServerError, "invalid_teacher_code")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type studentProgressResponse domain.StudentProgressInfo

// GetMyStudents godoc
// @Summary Get students associated with the authenticated teacher
// @Description Retrieve a list of students associated with the authenticated teacher
// @Tags User
// @Success	200	{array}		studentProgressResponse
// @Failure	500	{object}	ErrorResponse
// @Router	/teacher/students [get]
// @Security BearerAuth
func (h *UserHandler) GetMyStudents(w http.ResponseWriter, r *http.Request) {
	teacherID := r.Context().Value(userContextKey).(int)

	students, err := h.service.GetStudentByTeacher(teacherID)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "cannot_find_students")
		return
	}

	writeJSON(w, http.StatusOK, students)
}

// GetMyInviteCode godoc
// @Summary Get the invite code for the authenticated teacher
// @Description Retrieve the invite code for the authenticated teacher
// @Tags User
// @Success	200	{object}	string "Invite code"
// @Failure	500	{object}	ErrorResponse
// @Router	/teacher/invite-code [get]
// @Security BearerAuth
func (h *UserHandler) GetMyInviteCode(w http.ResponseWriter, r *http.Request) {
	teacherID := r.Context().Value(userContextKey).(int)

	code, err := h.service.GetInviteCode(teacherID)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "cannot_find_teacher")
	}

	writeJSON(w, http.StatusOK, code)
}
