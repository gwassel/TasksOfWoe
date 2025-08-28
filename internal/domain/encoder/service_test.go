package encoder

import (
	"crypto/aes"
	"crypto/rand"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	plaintext := "Hello, world!"
	var enc Encoder

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Errorf(`Could not generate a key: %v`, err)
	}

	enc.c, err = aes.NewCipher([]byte(key))
	if err != nil {
		t.Errorf(`Could not create a cipher: %v`, err)
	}

	ciphertext := enc.Encode(plaintext)
	if err != nil {
		t.Errorf(`Could not encode message: %v`, err)
	}

	res := enc.Decode(ciphertext)
	if err != nil {
		t.Errorf(`Could not decode message: %v`, err)
	}

	if res != plaintext {
		t.Errorf(`Mismatch: Dec(Enc(%q)) = %q`, plaintext, res)
	}
}
