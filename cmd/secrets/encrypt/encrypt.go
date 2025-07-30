package encrypt

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"kp/cfg"
	"kp/encryption"
	"kp/log"
)

var Cmd = &cobra.Command{
	Use:     "encrypt",
	Aliases: []string{"e"},
	Short:   "Encrypt the secrets",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cfg.ReadConfig()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		execute(config)
	},
}

func execute(config *cfg.Config) {
	global := config.GetGlobal()

	cipher, err := encryption.NewCipher(global)
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	patterns := config.GetStringSlice("secrets")
	var secretFiles []string
	for _, pattern := range patterns {
		glob := filepath.Join(global.Dir, pattern)
		matches, err := filepath.Glob(glob)
		if err != nil {
			log.Errorln(err.Error())
			continue
		}

		filteredMatches := []string{}

		for _, match := range matches {
			if !strings.HasSuffix(match, cfg.EncryptedExt) {
				filteredMatches = append(filteredMatches, match)
			}
		}

		secretFiles = append(secretFiles, filteredMatches...)
	}

	log.Verbosef("Detected %d secrets", len(secretFiles))

	for _, file := range secretFiles {
		log.Verbosef("Encrypting %s", file)
		dir := filepath.Dir(file)

		encryptedFile := filepath.Join(dir, filepath.Base(file)+cfg.EncryptedExt)

		original, err := os.ReadFile(file)
		if err != nil {
			log.Errorln(err.Error())
			continue
		}

		encrypted := cipher.Encrypt(original)

		err = os.WriteFile(encryptedFile, encrypted, 0644)
		if err != nil {
			log.Errorf("Failed to encrypt %s: %v", file, err)
			continue
		}
	}
}
