package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/zeltbrennt/lenslocked/rand"
)

const minBytesPerToken = 32

type TokenManager struct {
	BytesPerToken int
}

func (tm TokenManager) New() (token, tokenHash string, err error) {
	bytesPerToken := max(tm.BytesPerToken, minBytesPerToken)
	token, err = rand.RandomString(bytesPerToken)
	if err != nil {
		return "", "", fmt.Errorf("create: %w", err)
	}
	tokenHash = tm.Hash(token)
	return token, tokenHash, err
}

func (tm TokenManager) Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
