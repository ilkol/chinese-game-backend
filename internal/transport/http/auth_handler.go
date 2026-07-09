package http

import (
	"chinese-game-backend/internal/service"
	"errors"
	"log"
	"net/http"

	"github.com/lib/pq"
)

type signUpInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	InviteCode string `json:"invite_code"`
}

type AuthHandler struct {
	services *service.UserService
}

func NewAuthHandler(services *service.UserService) *AuthHandler {
	return &AuthHandler{services: services}
}

type loginResponse struct {
	Token string `json:"token"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with student role
// @Tags auth
// @Param	input	body 		signUpInput true "User credentials"
// @Success	201		{object}	loginResponse
// @Failure	400		{object}	ErrorResponse
// @Failure	409		{object}	ErrorResponse
// @Failure	500		{object}	ErrorResponse
// @Router	/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input signUpInput

	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.services.SignUp(input.Username, input.Password, input.InviteCode)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				errorJSON(w, http.StatusConflict, "user_already_exists")
				return
			}
		}

		log.Println(err)
		errorJSON(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusCreated, loginResponse{Token: token})
}

// Login godoc
// @Summary Login user
// @Description Create a JWT token for user
// @Tags auth
// @Param	input	body 		signUpInput true "User credentials"
// @Success	200		{object}	loginResponse
// @Failure	400		{object}	ErrorResponse
// @Failure	401		{object}	ErrorResponse
// @Router	/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input signUpInput

	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.services.SignIn(input.Username, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserPassword) || errors.Is(err, service.ErrUserNotFound) {
			errorJSON(w, http.StatusUnauthorized, "invalid_credentials")
		} else {
			errorJSON(w, http.StatusInternalServerError, "internal_server_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}
