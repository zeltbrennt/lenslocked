package models

import (
	"database/sql"
	"fmt"
)

type Session struct {
	ID               int
	UserID           int
	NewToken         string
	CurrentTokenHash string
}

type SessionService struct {
	DB *sql.DB
	TM TokenManager
}

func (ss *SessionService) Create(UserId int) (*Session, error) {
	token, hash, err := ss.TM.New()
	if err != nil {
		return nil, fmt.Errorf("creating Token: %w", err)
	}
	session := Session{
		UserID:           UserId,
		NewToken:         token,
		CurrentTokenHash: hash,
	}
	row := ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
			UPDATE
			SET token_hash = $2
		RETURNING id;`, session.UserID, session.CurrentTokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.TM.Hash(token)
	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.TM.Hash(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id, users.email
		FROM sessions
		JOIN users ON sessions.user_id = users.id 
		WHERE token_hash = $1`, tokenHash)
	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("getting session: %w", err)
	}

	return &user, nil
}
