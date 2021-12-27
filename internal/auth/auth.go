package auth

import (
	"github.com/google/uuid"
)

type RandomNonceFunc func() ([]byte, error)

//go:generate mockery --name=Auth
type Auth interface {
	Decode(sequence []byte) (uuid.UUID, error)
	Encode(id uuid.UUID, nonceFunc RandomNonceFunc) (string, error)
}
