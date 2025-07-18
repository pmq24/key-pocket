package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	KeyFile = "kp_key.%s"
	KeyAlg = "ALG"
	KeyVal = "KEY"
)

var newCmd = &cobra.Command{
	Use: "new",
	Short: "Create a new key",
	Run: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) {
	keyFilePath := filepath.Join(Dir, fmt.Sprintf(KeyFile, Profile))

	if _, err := os.Stat(keyFilePath); err == nil {
		slog.Error(fmt.Sprintf("Key file already exists for profile %s (%s)", Profile, keyFilePath))
		return
	}

	key := make([]byte, 32)
	rand.Read(key)
	encodedKey := base64.StdEncoding.EncodeToString(key)

	keyFileContent := fmt.Sprintf("%s=aes-256\n%s=%s\n", KeyAlg, KeyVal, encodedKey)
	
	err := os.WriteFile(keyFilePath, []byte(keyFileContent), 0600)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to write key file: %v", err))
		return
	}
}

