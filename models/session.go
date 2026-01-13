package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/zeltbrennt/lenslocked/rand"
)

const minBytesPerToken = 32

type Session struct {
	ID               int
	UserID           int
	NewToken         string
	CurrentTokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(UserId int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	bytesPerToken = max(bytesPerToken, minBytesPerToken)
	token, err := rand.SessionToken(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:           UserId,
		NewToken:         token,
		CurrentTokenHash: ss.hash(token),
	}
	row := ss.DB.QueryRow(`INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2)
		RETURNING id;`, session.UserID, session.CurrentTokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
