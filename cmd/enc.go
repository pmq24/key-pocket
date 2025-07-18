package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var encCmd = &cobra.Command{
	Use: "enc",
	Short: "Encrypt secret files",
	Run: runEnc,
}

func init() {
	rootCmd.AddCommand(encCmd)
}

func runEnc(cmd *cobra.Command, args []string) { 
	patterns := viper.GetStringSlice("secrets")

	var secretFiles []string
	for _, pattern := range patterns {
		glob := filepath.Join(Dir, pattern)
		matches, err := filepath.Glob(glob)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		secretFiles = append(secretFiles, matches...)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(Key)
	if err != nil {
		slog.Error("Corrupted key", "error", err)
		return
	}

	for _, file := range secretFiles {
		dir := filepath.Dir(file)
		
		encryptedFile := filepath.Join(dir, filepath.Base(file) + ".kpenc")
		
		content, err := os.ReadFile(file)
		if err != nil {
			slog.Error("Failed to read file", "file", file, "error", err)
			continue
		}
		
		encryptedContent, err := encryptAES256(content, keyBytes)
		if err != nil {
			slog.Error("Failed to encrypt file", "file", file, "error", err)
			continue
		}
		
		err = os.WriteFile(encryptedFile, encryptedContent, 0644)
		if err != nil {
			slog.Error("Failed to write encrypted file", "file", encryptedFile, "error", err)
			continue
		}
		
		slog.Info("Encrypted file", "source", file, "encrypted", encryptedFile)
	}
}

func encryptAES256(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}
