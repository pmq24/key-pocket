package decrypt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"kp/cfg"
	"kp/enc"
	"kp/log"
)

var Cmd = &cobra.Command{
	Use:     "decrypt",
	Aliases: []string{"d"},
	Short:   "Decrypt the secrets",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := buildConfig(cmd)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		encryptor, err := enc.NewAES256Encryptor(&config.BaseCfg)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		opts := executeOpts{
			config:    config,
			encryptor: encryptor,
		}
		execute(opts)
	},
}

func buildConfig(cmd *cobra.Command) (*cfg.AppCfg, error) {
	pflags := cmd.Root().PersistentFlags()
	cfgOpts := cfg.NewCfgOpts{
		Dir:     pflags.Lookup("dir").Value.String(),
		Profile: pflags.Lookup("profile").Value.String(),
	}
	return cfg.NewAppCfg(cfgOpts)
}

type executeOpts struct {
	config    *cfg.AppCfg
	encryptor enc.Encryptor
}

func execute(opts executeOpts) {
	patterns := opts.config.GetSecrets()

	var encryptedFiles []string
	for _, pattern := range patterns {
		p := pattern + cfg.EncryptedExt
		glob := filepath.Join(opts.config.GetDir(), p)
		matches, err := filepath.Glob(glob)
		if err != nil {
			log.Errorf("Error globbing %s: %v. Skipping...", p, err.Error())
			continue
		}

		encryptedFiles = append(encryptedFiles, matches...)
	}

	log.Verbosef("Detected %d encrypted files", len(encryptedFiles))

	for _, file := range encryptedFiles {
		log.Verboseln("Decrypting " + file)

		encryptedContent, err := os.ReadFile(file)
		if err != nil {
			log.Errorln(fmt.Sprintf("Error reading encrypted file %s: %v. Skipping...", file, err))
			continue
		}

		decryptedContent, err := opts.encryptor.Decrypt(encryptedContent)
		if err != nil {
			log.Errorln(fmt.Sprintf("Error decrypting %s: %v", file, err))
			continue
		}

		outputPath := file[:len(file)-len(cfg.EncryptedExt)]

		err = os.WriteFile(outputPath, decryptedContent, 0644)
		if err != nil {
			log.Errorln(fmt.Sprintf("Error writing decrypted file %s: %v. Skipping...", outputPath, err))
			continue
		}
	}
}
