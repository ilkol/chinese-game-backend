package http

import (
	"chinese-game-backend/internal/domain"
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type conextKey string

const (
	userContextKey conextKey = "user_id"
	roleContextKey conextKey = "user_role"
)

func (h *AuthHandler) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			errorJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			errorJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		token, err := jwt.Parse(headerParts[1], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(h.services.GetSecret()), nil
		})

		if err != nil || !token.Valid {
			errorJSON(w, http.StatusUnauthorized, "invalid_token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			errorJSON(w, http.StatusUnauthorized, "invalid_token")
			return
		}

		userID := int(claims["user_id"].(float64))
		userRole := claims["role"].(string)
		ctx := context.WithValue(r.Context(), userContextKey, userID)
		ctx = context.WithValue(r.Context(), roleContextKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *AuthHandler) CheckRole(allowedRoles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(roleContextKey).(domain.UserRole)
			if !ok {
				errorJSON(w, http.StatusForbidden, "role_not_found")
				return
			}

			if !slices.Contains(allowedRoles, userRole) {
				errorJSON(w, http.StatusForbidden, "insufficient_privileges")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
