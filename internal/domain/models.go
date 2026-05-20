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
	InviteCode   string    `json:"invite_code" db:"invite_code"`
}

type Level struct {
	ID         int         `json:"id" db:"id"`
	Title      string      `json:"title" db:"title"`
	Color      string      `json:"color" db:"color"`
	Icon       string      `json:"icon" db:"icon"`
	OrderIndex int         `json:"order_index" db:"order_index"`
	Steps      []LevelStep `json:"steps"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
}

type LevelStepType string

const (
	LevelStepThoery LevelStepType = "theory"
	LevelStepQuiz   LevelStepType = "quiz"
	LevelStepFinal  LevelStepType = "final"
)

type LevelStep struct {
	ID          int             `json:"id" db:"id"`
	LevelID     int             `json:"level_id" db:"level_id"`
	Type        LevelStepType   `json:"type" db:"type"`
	Title       string          `json:"title" db:"title"`
	OrderIndex  int             `json:"order_index" db:"order_index"`
	Content     json.RawMessage `json:"content" db:"content"`
	IsCompleted bool            `json:"is_completed" db:"is_completed"`
}
