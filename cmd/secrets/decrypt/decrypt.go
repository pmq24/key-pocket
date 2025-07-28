package decrypt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"kp/cfg"
	"kp/encryption"
	"kp/log"
)

var Cmd = &cobra.Command{
  Use: "decrypt",
	Aliases: []string{"d"},
  Short: "Decrypt the secrets",
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
	var encryptedFiles []string
	for _, pattern := range patterns {
		glob := filepath.Join(global.Dir, pattern + cfg.EncryptedExt)
		matches, err := filepath.Glob(glob)
		if err != nil {
			log.Errorln(err.Error())
			continue
		}

		encryptedFiles = append(encryptedFiles, matches...)
	}

	log.Verbosef("Detected %d encrypted files", len(encryptedFiles))

	for _, file := range encryptedFiles {
		log.Verboseln("Decrypting " + file)

		encryptedContent, err := os.ReadFile(file)
		if err != nil {
			log.Errorln(fmt.Sprintf("Error reading encrypted file %s: %v", file, err))
			continue
		}

		decryptedContent, err := cipher.Decrypt(encryptedContent)
		if err != nil {
			log.Errorln(fmt.Sprintf("Error decrypting %s: %v", file, err))
			continue
		}

		outputPath := file[:len(file)-len(cfg.EncryptedExt)]

		err = os.WriteFile(outputPath, decryptedContent, 0644)
		if err != nil {
			log.Errorln(fmt.Sprintf("Error writing decrypted file %s: %v", outputPath, err))
			continue
		}
	}
}
