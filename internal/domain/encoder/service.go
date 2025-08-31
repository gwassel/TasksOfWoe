package encoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"
)

type Encoder struct {
	gcm cipher.AEAD
}

func New(key string) (*Encoder, error) {
	bytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher([]byte(bytes))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return &Encoder{gcm: gcm}, nil
}

func (e *Encoder) Encode(plaintext string) ([]byte, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	// fill in nonce with random bytes
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	// store both nonce and cyphertext
	return e.gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

func (e *Encoder) Decode(ciphertext []byte) (string, error) {
	if len(ciphertext) < e.gcm.NonceSize() {
		return "", errors.New("Invalid data encryption")
	}

	// split into nonce and ciphertext
	nonce, ciphertext := ciphertext[:e.gcm.NonceSize()], ciphertext[e.gcm.NonceSize():]
	// decode
	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
