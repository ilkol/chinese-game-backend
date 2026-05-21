package service

import (
	"math/rand"
	"time"
)

const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func GenerateInviteCode(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
