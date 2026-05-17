package http

import (
	"chinese-game-backend/internal/service"
	"net/http"

	"github.com/lib/pq"
)

type SignUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	services *service.UserService
}

func NewAuthHandler(services *service.UserService) *AuthHandler {
	return &AuthHandler{services: services}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input SignUpInput

	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.services.SignUp(input.Username, input.Password)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				errorJSON(w, http.StatusConflict, "user_already_exists")
				return
			}
		}

		errorJSON(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := readJSON(r, &input); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.services.SignIn(input.Username, input.Password)
	if err != nil {
		errorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
