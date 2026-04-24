package impl

import (
	"crypto/sha256"
	"encoding/hex"
)

type AIAPIRepository struct{}

func NewAIAPIRepository() *AIAPIRepository {
	return &AIAPIRepository{}
}

func (r *AIAPIRepository) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
