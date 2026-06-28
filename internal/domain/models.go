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
	InviteCode   *string   `json:"invite_code" db:"invite_code"`
}

type Level struct {
	ID            int         `json:"id" db:"id"`
	Title         string      `json:"title" db:"title"`
	Color         string      `json:"color" db:"color"`
	Icon          string      `json:"icon" db:"icon"`
	OrderIndex    int         `json:"order_index" db:"order_index"`
	Steps         []LevelStep `json:"steps"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	BackgroundSrc string      `json:"background_src" db:"background_src"`
	PlanetImgSrc  string      `json:"planet_img_src" db:"planet_img_src"`
}

type LevelStepType string

const (
	LevelStepThoery      LevelStepType = "theory"
	LevelStepQuiz        LevelStepType = "quiz"
	LevelStepFinal       LevelStepType = "final"
	LevelStepDialog      LevelStepType = "dialog"
	LevelStepCategoriz   LevelStepType = "categorization"
	LevelStepToneListen  LevelStepType = "tone_listening"
	LevelStepPlanetClick LevelStepType = "planet_click"
	LevelStepPlanetMatch LevelStepType = "planet_matching"
)

type DialogStepItem struct {
	Speaker string `json:"speaker,omitempty"`
	Text    string `json:"text"`
	Emotion string `json:"emotion,omitempty"`
	Bg      string `json:"bg,omitempty"`
}

type LevelStep struct {
	ID          int              `json:"id" db:"id"`
	LevelID     int              `json:"level_id" db:"level_id"`
	Type        LevelStepType    `json:"type" db:"type"`
	Title       string           `json:"title" db:"title"`
	OrderIndex  int              `json:"order_index" db:"order_index"`
	Content     json.RawMessage  `json:"content" db:"content"`
	IsCompleted bool             `json:"is_completed" db:"is_completed"`
	Description string           `json:"description" db:"description"`
	Dialog      []DialogStepItem `json:"dialog,omitempty" db:"-"`
}

type StudentProgressInfo struct {
	ID             int       `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Coins          int       `json:"coins" db:"coins"`
	LastLevelTitle string    `json:"lastLevelTitle" db:"last_level_title"`
	CompletedSteps int       `json:"completedSteps" db:"completed_steps"`
	TotalSteps     int       `json:"totalSteps" db:"total_steps"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}
