package auth

import (
	"github.com/google/uuid"
)

// NonceFunc is typing for functions that generate nonce(s).
type NonceFunc func(nonceSize int) ([]byte, error)

// UserID hides real user id implementation.
type UserID = *string

// NewUserID generates new random user identifier.
func NewUserID() UserID {
	id := uuid.New().String()
	return &id
}

// Auth is common interface for authentication implementation of helpers that can encode and decode a user identifier.
//go:generate mockery --name=Auth
type Auth interface {
	Decode(sequence []byte) (UserID, error)
	Encode(id UserID) ([]byte, error)
}
