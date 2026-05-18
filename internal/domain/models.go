package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

func (r UserRole) IsValid() error {
	switch r {
	case RoleStudent, RoleTeacher, RoleAdmin:
		return nil
	}
	return errors.New("invalid role")
}

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         UserRole  `json:"role" db:"role"`
	Coins        int       `json:"coins" db:"coins"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Level struct {
	ID        int             `json:"id" db:"id"`
	Title     string          `json:"title" db:"title"`
	Color     string          `json:"color" db:"color"`
	Icon      string          `json:"icon" db:"icon"`
	Blocks    json.RawMessage `json:"blocks" db:"blocks"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
