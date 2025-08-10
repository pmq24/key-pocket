package enc

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

func NewAES256Encryptor(config *cfg.BaseCfg) (Encryptor, error) {
	keyFilePath := filepath.Join(config.GetDir(), fmt.Sprintf(cfg.KeyFileFormat, config.GetProfile()))

	keyBase64, err := os.ReadFile(keyFilePath)
	if err != nil {
	}

	key, err := base64.StdEncoding.DecodeString(string(keyBase64))
	if err != nil {
		return nil, ErrNewEncryptor{KeyFilePath: keyFilePath, Err: err}
	}

	_, err = aes.NewCipher(key)
	if err != nil {
		return nil, ErrNewEncryptor{KeyFilePath: keyFilePath, Err: err}
	}

	block, err := aes.NewCipher(key)
	gcm, err := goCipher.NewGCM(block)
	if err != nil {
		return nil, ErrNewEncryptor{KeyFilePath: keyFilePath, Err: err}
	}

	return &AES256Encryptor{gcm}, nil
}

func (c *AES256Encryptor) Encrypt(data []byte) []byte {
	nonce := make([]byte, c.gcm.NonceSize())
	rand.Read(nonce)
	return c.gcm.Seal(nonce, nonce, data, nil)
}

func (c *AES256Encryptor) Decrypt(data []byte) ([]byte, error) {
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

type AES256Encryptor struct {
	gcm goCipher.AEAD
}
