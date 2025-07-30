package encryption

import (
	"crypto/aes"
	goCipher "crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"kp/cfg"
)

func NewCipher(global cfg.Global) (*Cipher, error) {
	keyFilePath := filepath.Join(global.Dir, fmt.Sprintf(cfg.KeyFileFormat, global.Profile))

	keyBase64, err := os.ReadFile(keyFilePath)
	if err != nil {
	}

	key, err := base64.StdEncoding.DecodeString(string(keyBase64))
	if err != nil {
		return nil, ErrNewCipher{KeyFilePath: keyFilePath, Err: err}
	}

	_, err = aes.NewCipher(key)
	if err != nil {
		return nil, ErrNewCipher{KeyFilePath: keyFilePath, Err: err}
	}

	block, err := aes.NewCipher(key)
	gcm, err := goCipher.NewGCM(block)
	if err != nil {
		return nil, ErrNewCipher{KeyFilePath: keyFilePath, Err: err}
	}

	return &Cipher{gcm}, nil
}

func (c *Cipher) Encrypt(data []byte) []byte {
	nonce := make([]byte, c.gcm.NonceSize())
	rand.Read(nonce)
	return c.gcm.Seal(nonce, nonce, data, nil)
}

func (c *Cipher) Decrypt(data []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()

	if nonceSize > len(data) {
		return nil, ErrDecrypt{Err: fmt.Errorf("encrypted data too short")}
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	original, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecrypt{Err: err}
	}

	return original, nil
}

type Cipher struct {
	gcm goCipher.AEAD
}
