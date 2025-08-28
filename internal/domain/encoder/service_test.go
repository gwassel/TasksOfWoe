package encoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	plaintext := `Hello, world!`
	var enc Encoder

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Errorf(`Could not generate a key: %v`, err)
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		t.Errorf(`Could not create a cipher: %v`, err)
	}

	enc.gcm, err = cipher.NewGCM(c)
	if err != nil {
		t.Errorf(`Could not create a GCM: %v`, err)
	}

	ciphertext, err := enc.Encode(plaintext)
	if err != nil {
		t.Errorf(`Could not encode message: %v`, err)
	}

	res, err := enc.Decode(ciphertext)
	if err != nil {
		t.Errorf(`Could not decode message: %v`, err)
	}

	if res != plaintext {
		t.Errorf(`Mismatch: Dec(Enc(%q)) = %q`, plaintext, res)
	}
}
