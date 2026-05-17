package http

import (
	"chinese-game-backend/internal/service"
	"encoding/json"
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

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.services.SignUp(input.Username, input.Password)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(`{"error": "user_already_exists"}`))
				return
			}
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "user created"}`))
}
