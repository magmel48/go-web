package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/magmel48/go-web/internal/config"
)

// CustomAuth is custom auth implementation.
type CustomAuth struct {
	algo      cipher.AEAD
	NonceFunc NonceFunc
}

// NewCustomAuth creates new CustomAuth instance.
func NewCustomAuth() (*CustomAuth, error) {
	aesBlock, err := aes.NewCipher([]byte(config.SecretKey)[:aes.BlockSize])
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	return &CustomAuth{algo: aesGCM, NonceFunc: DefaultNonceFunc}, nil
}

// Decode decodes encoded sequence (usually from user session)
// and returns user identifier if the input sequence is valid.
func (auth CustomAuth) Decode(sequence []byte) (UserID, error) {
	if len(sequence) == 0 {
		return nil, errors.New("wrong bytes sequence")
	}

	encoded := make([]byte, base64.RawStdEncoding.DecodedLen(len(sequence)))

	_, err := base64.RawStdEncoding.Decode(encoded, sequence)
	if err != nil {
		return nil, err
	}

	encrypted := encoded[:len(encoded) - auth.algo.NonceSize()]
	nonce := encoded[len(encoded) - auth.algo.NonceSize():]

	decrypted, err := auth.algo.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	result := string(decrypted)
	return &result, nil
}

// Encode encodes user identifier and puts iv into the end of result for further decoding.
func (auth CustomAuth) Encode(id UserID) ([]byte, error) {
	nonce, err := auth.NonceFunc(auth.algo.NonceSize())
	if err != nil {
		return nil, err
	}

	encrypted := auth.algo.Seal(nil, nonce, []byte(*id), nil)

	raw := append(encrypted, nonce...)
	serialized := make([]byte, base64.RawStdEncoding.EncodedLen(len(raw)))

	base64.RawStdEncoding.Encode(serialized, raw)

	return serialized, nil
}

// DefaultNonceFunc is default function for nonce retrieving.
func DefaultNonceFunc(nonceSize int) ([]byte, error) {
	nonce := make([]byte, nonceSize)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
