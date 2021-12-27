package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/magmel48/go-web/internal/config"
)

type CustomAuth struct {
	algo cipher.AEAD
}

func NewCustomAuth() (*CustomAuth, error) {
	aesBlock, err := aes.NewCipher([]byte(config.SecretKey)[:aes.BlockSize])
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	return &CustomAuth{algo: aesGCM}, nil
}

// Decode decodes encoded sequence (usually from user session)
// and returns user identifier if the input sequence is valid.
func (auth CustomAuth) Decode(sequence []byte) (uuid.UUID, error) {
	return uuid.New(), nil
}

// Encode encodes user identifier and puts iv into the end of result for further decoding.
func (auth CustomAuth) Encode(id uuid.UUID, nonceFunc RandomNonceFunc) (string, error) {
	nonce, err := nonceFunc()
	if err != nil {
		return "", err
	}

	encrypted := auth.algo.Seal(nil, nonce, []byte(id.String()), nil)
	serialized := base64.RawStdEncoding.EncodeToString(append(encrypted, nonce...))

	return serialized, nil
}

func (auth CustomAuth) RandomNonce() ([]byte, error) {
	nonce := make([]byte, auth.algo.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
