package encoder

import (
	"crypto/aes"
	"crypto/cipher"
	"strings"
)

type Encoder struct {
	c cipher.Block
}

func (e *Encoder) Encode(plaintext string) []byte {
	// String lenght padding to make string size divisible by block size
	var length int
	if len(plaintext)%aes.BlockSize == 0 {
		length = len(plaintext)
	} else {
		length = len(plaintext) + aes.BlockSize - len(plaintext)%aes.BlockSize
	}

	in := make([]byte, length)
	copy(in, []byte(plaintext))
	out := make([]byte, length)
	e.c.Encrypt(out, in)

	return out
}

func (e *Encoder) Decode(ciphertext []byte) string {
	out := make([]byte, len(ciphertext))
	e.c.Decrypt(out, ciphertext)

	// Cut zeros added to the end of the strings before encoding
	return strings.Trim(string(out), string('\x00'))
}
