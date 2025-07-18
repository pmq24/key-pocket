package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decCmd)
}

var decCmd = &cobra.Command{
  Use: "dec",
  Short: "Decrypt secret files",
	Run: runDec,
}

func runDec(cmd *cobra.Command, args []string) {
	encFiles, err := getEncryptedFiles()
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot retrieve encrypted files: %v", err))
		return
	}

	decryptor, err := createDecryptor()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create decryptor: %v", err))
		return
	}

	for _, encFile := range encFiles {
		encryptedContent, err := os.ReadFile(encFile)
		if err != nil {
			slog.Error("Failed to read encrypted file", "file", encFile, "error", err)
			continue
		}

		decryptedContent, err := decryptAES256(encryptedContent, *decryptor)
		if err != nil {
			slog.Error("Failed to decrypt file", "file", encFile, "error", err)
			continue
		}

		originalFilePath := encFile[:len(encFile) - len(".kpenc")]

		err = os.WriteFile(originalFilePath, decryptedContent, 0644)
		if err != nil {
			slog.Error("Failed to write decrypted file", "file", originalFilePath, "error", err)
			continue
		}

		slog.Info("Successfully decrypted file", "file", originalFilePath)
	}
}

func decryptAES256(data []byte, gcm cipher.AEAD) ([]byte, error) {
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func getEncryptedFiles() ([]string, error) {
	encFiles, err := filepath.Glob(filepath.Join(Dir, "*.kpenc"))
	if err != nil {
		return nil, err
	}

	return encFiles, nil
}

func createDecryptor() (*cipher.AEAD, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(Key)
	if err != nil {
		return nil, fmt.Errorf("Corrupted key: %v", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %v", err)
	}

	return &gcm, nil
}
