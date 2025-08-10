package encrypt

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"kp/cfg"
	"kp/enc"
	"kp/log"
)

var Cmd = &cobra.Command{
	Use:     "encrypt",
	Aliases: []string{"e"},
	Short:   "Encrypt the secrets",
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

		execute(opts{
			config:    config,
			encryptor: encryptor,
		})
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

type opts struct {
	config    *cfg.AppCfg
	encryptor enc.Encryptor
}

func execute(opts opts) {
	config := opts.config
	encryptor := opts.encryptor

	patterns := opts.config.GetSecrets()

	var secretFiles []string
	for _, pattern := range patterns {
		glob := filepath.Join(config.GetDir(), pattern)
		matches, err := filepath.Glob(glob)
		if err != nil {
			log.Errorf("Error globbing %s: %v. Skipping...", glob, err.Error())
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
			log.Errorf("Error reading %s: %v. Skipping...", file, err.Error())
			continue
		}

		encrypted := encryptor.Encrypt(original)

		err = os.WriteFile(encryptedFile, encrypted, 0644)
		if err != nil {
			log.Errorf("Failed to encrypt %s: %v. Skipping...", file, err)
			continue
		}
	}
}
